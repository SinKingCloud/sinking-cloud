package route

import (
	http2 "net/http"
	"os"
	"os/signal"
	"service-user/app/constant"
	"strconv"
	"syscall"

	constant2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/middleware"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/http"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// Init 初始化server
func Init() {
	host, port := global.App.ServerAddr()
	addr := host + ":" + strconv.Itoa(port)
	r := http.NewServer(addr, global.App.IsDebug()).
		BeforeHandle(func(r *http2.Request, w http2.ResponseWriter) bool {
			originService := r.Header.Get(constant2.RpcOriginService)
			if originService == constant2.ServiceGateway {
				r.URL.Path = constant2.ServiceUser + "/" + r.URL.Path
			}
			return true
		}).
		Handle(func(engine *sinking_web.Engine) {
			engine.Use(context.HandleFunc(middleware.InterceptIllegalRequest))
			loadApp(engine)
		}).
		ErrorHandle(&sinking_web.ErrorHandel{
			NotFound: func(c *sinking_web.Context) {
				c.Set(constant.NeedProxy, true)
			},
			Fail: func(c *sinking_web.Context, code int, message string) {
				ctx := context.New(c)
				ctx.ErrorWithData("请求失败", message)
			},
		})
	_ = r.Start()
	//等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	_ = <-sigChan
	_ = r.Stop()
}
