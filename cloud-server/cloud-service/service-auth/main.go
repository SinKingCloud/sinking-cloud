package main

import (
	"service-auth/app"
	"service-auth/bootstrap"
)

func main() {
	bootstrap.Load()
	app.Run()
}
