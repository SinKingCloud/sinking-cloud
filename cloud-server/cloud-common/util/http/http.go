package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"net/http"
	"sync"
	"time"
)

// Server HTTP服务器
type Server struct {
	engine   *sinking_web.Engine
	server   *http.Server
	running  bool
	mutex    sync.RWMutex
	handlers []func(engine *sinking_web.Engine)
}

// NewServer 创建新的HTTP服务器实例
func NewServer() *Server {
	return &Server{
		running:  false,
		handlers: make([]func(engine *sinking_web.Engine), 0),
	}
}

// init 初始化服务器
func (s *Server) init() {
	if s.engine != nil {
		return
	}
	debugMode := false
	if global.Config != nil {
		debugMode = global.Config.GetString(constant.ServerMode) == "debug"
	}
	sinking_web.SetDebugMode(debugMode)
	s.engine = sinking_web.Default()
	s.engine.SetErrorHandle(&sinking_web.ErrorHandel{
		NotFound: func(c *sinking_web.Context) {
			c.JSON(404, sinking_web.H{"code": 404, "message": "请求资源不存在"})
		},
		Fail: func(c *sinking_web.Context, code int, message string) {
			c.JSON(500, sinking_web.H{"code": code, "message": message})
		},
	})
	for _, handler := range s.handlers {
		handler(s.engine)
	}
}

// Handle 添加处理器函数
func (s *Server) Handle(handler func(engine *sinking_web.Engine)) *Server {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.handlers = append(s.handlers, handler)
	if s.engine != nil {
		handler(s.engine)
	}
	return s
}

// Start 启动服务器
func (s *Server) Start() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.running {
		return fmt.Errorf("服务器已经在运行中")
	}
	s.init()
	host := "0.0.0.0"
	port := "5678"
	if global.Config != nil {
		if h := global.Config.GetString(constant.ServerHost); h != "" {
			host = h
		}
		if p := global.Config.GetString(constant.ServerPort); p != "" {
			port = p
		}
	}
	addr := host + ":" + port
	s.server = &http.Server{
		Addr:    addr,
		Handler: s.engine,
	}
	s.running = true
	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.mutex.Lock()
			s.running = false
			s.mutex.Unlock()
		}
	}()
	return nil
}

// Stop 停止服务器
func (s *Server) Stop() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if !s.running {
		return fmt.Errorf("服务器未运行")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	s.running = false
	s.server = nil
	return nil
}

// Restart 重启服务器
func (s *Server) Restart() error {
	if s.running {
		if err := s.Stop(); err != nil {
			return fmt.Errorf("停止服务器失败: %v", err)
		}
	}
	time.Sleep(500 * time.Millisecond)
	if err := s.Start(); err != nil {
		return fmt.Errorf("启动服务器失败: %v", err)
	}
	return nil
}
