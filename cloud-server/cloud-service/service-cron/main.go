package main

import (
	"service-cron/app"
	"service-cron/bootstrap"
)

func main() {
	bootstrap.Load()
	app.Run()
}
