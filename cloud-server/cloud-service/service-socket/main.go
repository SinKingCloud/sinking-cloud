package main

import (
	"service-socket/app"
	"service-socket/bootstrap"
)

func main() {
	bootstrap.Load()
	app.Run()
}
