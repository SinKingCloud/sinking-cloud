package route

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"net/http/pprof"
)

// loadDebug 加载debug
func loadDebug(s *sinking_web.Engine) {
	r := s.Group("/debug/" + util.Conf.GetString("register.tokenValue"))
	r.GET("/index", func(c *sinking_web.Context) {
		pprof.Index(c.Writer, c.Request)
	})
	r.GET("/profile", func(c *sinking_web.Context) {
		pprof.Profile(c.Writer, c.Request)
	})
	r.GET("/cmdline", func(c *sinking_web.Context) {
		pprof.Cmdline(c.Writer, c.Request)
	})
	r.GET("/symbol", func(c *sinking_web.Context) {
		pprof.Symbol(c.Writer, c.Request)
	})
	r.GET("/trace", func(c *sinking_web.Context) {
		pprof.Trace(c.Writer, c.Request)
	})
	r.GET("/goroutine", func(c *sinking_web.Context) {
		pprof.Handler("goroutine").ServeHTTP(c.Writer, c.Request)
	})
	r.GET("/allocs", func(c *sinking_web.Context) {
		pprof.Handler("allocs").ServeHTTP(c.Writer, c.Request)
	})
	r.GET("/block", func(c *sinking_web.Context) {
		pprof.Handler("block").ServeHTTP(c.Writer, c.Request)
	})
	r.GET("/heap", func(c *sinking_web.Context) {
		pprof.Handler("heap").ServeHTTP(c.Writer, c.Request)
	})
	r.GET("/mutex", func(c *sinking_web.Context) {
		pprof.Handler("mutex").ServeHTTP(c.Writer, c.Request)
	})
	r.GET("/threadcreate", func(c *sinking_web.Context) {
		pprof.Handler("threadcreate").ServeHTTP(c.Writer, c.Request)
	})
	r.GET("/threadcreate", func(c *sinking_web.Context) {
		pprof.Handler("threadcreate").ServeHTTP(c.Writer, c.Request)
	})
}
