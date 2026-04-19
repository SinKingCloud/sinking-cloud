package route

import (
	"service-driver/app/http/controller"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadApi 加载API
func loadApi(s *sinking_web.Engine) {
	upload := s.Group("/upload")
	upload.ANY("/simple", context.HandleFunc(controller.Api.Upload.Simple))
	upload.ANY("/chunk", context.HandleFunc(controller.Api.Upload.Chunk))

	file := s.Group("/preview")
	file.ANY("/*", context.HandleFunc(controller.Api.File.Preview))
}
