package main

import (
	"service-pay/app"
	"service-pay/bootstrap"
)

func main() {
	bootstrap.Load()
	app.Run()
}
