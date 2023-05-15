package main

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/bootstrap"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		if err := http.ListenAndServe("0.0.0.0:6060", nil); err != nil {
			log.Println(err)
		}
	}()
	bootstrap.Load()
	app.Run()
}
