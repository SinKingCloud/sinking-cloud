package main

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/bootstrap"
)

func main() {
	bootstrap.Load()
	app.Run()
}
