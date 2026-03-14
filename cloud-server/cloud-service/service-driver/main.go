package main

import (
	"service-driver/app"
	"service-driver/bootstrap"
)

func main() {
	bootstrap.Load()
	app.Run()
}
