package main

import (
	"service-gateway/app"
	"service-gateway/bootstrap"
)

func main() {
	bootstrap.Load()
	app.Run()
}
