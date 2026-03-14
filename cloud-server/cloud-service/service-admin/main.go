package main

import (
	"service-admin/app"
	"service-admin/bootstrap"
)

func main() {
	bootstrap.Load()
	app.Run()
}
