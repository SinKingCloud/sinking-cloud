package main

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-spider/app"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-spider/bootstrap"
)

func main() {
	bootstrap.Load()
	app.Run()
}
