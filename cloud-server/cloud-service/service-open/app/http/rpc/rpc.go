package rpc

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/http/rpc/oauth"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/http/rpc/user_app"
)

var (
	Oauth   = oauth.ControllerOauth{}
	UserApp = user_app.ControllerUserApp{}
)
