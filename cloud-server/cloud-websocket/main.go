package main

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/bootstrap"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-websocket/app"
)

func main() {
	bootstrap.Load()
	app.Run()
}
