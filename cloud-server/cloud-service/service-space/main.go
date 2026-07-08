package main

import (
	"service-space/app"
	"service-space/bootstrap"
)

func main() {
	bootstrap.Load()
	app.Run()
}
