package main

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/bootstrap"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-auth/app"
)

func main() {
	bootstrap.Load()
	app.Run()
}
