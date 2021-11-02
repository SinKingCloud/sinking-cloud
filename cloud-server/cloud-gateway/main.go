package main

import (
	"cloud-gateway/app"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/bootstrap"
)

func main() {
	bootstrap.Load()
	app.Run()
}
