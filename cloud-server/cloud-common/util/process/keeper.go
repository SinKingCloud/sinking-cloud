package process

import (
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Keeper struct {
	Cmd           []string
	RetryInterval time.Duration
	stopChan      chan struct{}
	wg            sync.WaitGroup
	mu            sync.Mutex
	process       *os.Process
	running       bool
}

func NewKeeper(cmd []string) *Keeper {
	return &Keeper{
		Cmd:           cmd,
		RetryInterval: 5 * time.Second,
	}
}

// Start 启动守护进程，如果已启动则不再重复启动
func (pk *Keeper) Start() {
	pk.mu.Lock()
	defer pk.mu.Unlock()
	// 如果已经启动，则直接返回
	if pk.running {
		return
	}
	// 标记为运行中，并为stopChan创建一个新的通道
	pk.running = true
	pk.stopChan = make(chan struct{})
	pk.wg.Add(1)
	go pk.mainLoop()
}

// Stop 停止守护进程，等待goroutine结束
func (pk *Keeper) Stop() {
	pk.mu.Lock()
	// 如果没有运行，则直接返回
	if !pk.running {
		pk.mu.Unlock()
		return
	}
	// 关闭stopChan通知 goroutine 退出，并修改状态标识
	close(pk.stopChan)
	pk.running = false
	pk.mu.Unlock()
	pk.wg.Wait()
}

// Restart 停止后重新启动
func (pk *Keeper) Restart() {
	pk.Stop()
	time.Sleep(500 * time.Millisecond)
	pk.Start()
}

// mainLoop 主循环，用于运行和重启进程
func (pk *Keeper) mainLoop() {
	defer pk.wg.Done()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	for {
		select {
		case <-pk.stopChan:
			pk.killProcess()
			return
		case <-sigChan:
			// 收到信号后直接退出当前循环，结束守护进程
			pk.killProcess()
			return
		default:
			pk.runProcess()
		}
	}
}

// runProcess 启动具体进程，并等待其退出；进程退出后等待一段重试间隔
func (pk *Keeper) runProcess() {
	cmd := exec.Command(pk.Cmd[0], pk.Cmd[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 启动前更新process状态（注意：cmd.Process在Start后才会生效）
	pk.mu.Lock()
	pk.process = nil // 清空旧进程引用
	pk.mu.Unlock()

	if err := cmd.Start(); err != nil {
		time.Sleep(pk.RetryInterval)
		return
	}

	// 更新进程引用
	pk.mu.Lock()
	pk.process = cmd.Process
	pk.mu.Unlock()

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case <-done:
	case <-pk.stopChan:
		pk.killProcess()
		return
	}
	time.Sleep(pk.RetryInterval)
}

// killProcess 杀死正在运行的进程
func (pk *Keeper) killProcess() {
	pk.mu.Lock()
	defer pk.mu.Unlock()

	if pk.process != nil {
		_ = pk.process.Kill()
	}
}
