package main

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/bootstrap"
)

func main() {
	bootstrap.Load()
	app.Run()
}
