package route

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/http"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// Init 初始化server
func Init() {
	host, port := global.App.ServerAddr()
	addr := host + ":" + strconv.Itoa(port)
	r := http.NewServer(addr, global.App.IsDebug()).Handle(func(engine *sinking_web.Engine) {
		loadApp(engine)
	})
	_ = r.Start()
	//等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	_ = <-sigChan
	_ = r.Stop()
}
