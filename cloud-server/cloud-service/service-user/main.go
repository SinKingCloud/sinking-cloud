package main

import (
	"service-user/app"
	"service-user/bootstrap"
)

func main() {
	bootstrap.Load()
	app.Run()
}
