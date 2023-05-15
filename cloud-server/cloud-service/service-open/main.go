package main

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/bootstrap"
	_ "net/http/pprof"
)

func main() {
	bootstrap.Load()
	app.Run()
}
