package api

import (
	"service-driver/app/http/controller/api/file"
	"service-driver/app/http/controller/api/upload"
)

type Controller struct {
	File   file.Controller
	Upload upload.Controller
}
