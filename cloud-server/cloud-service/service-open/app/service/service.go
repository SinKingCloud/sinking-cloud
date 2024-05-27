package service

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_app"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_auth"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_auth/auth_app"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_auth/auth_log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_auth/auth_version"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_auth/user_auth"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_oauth"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_oauth/oauth"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_oauth/oauth_log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_oauth/oauth_user"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_plugin"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_plugin/plugin"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_plugin/user_plugin"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_user"
)

var (
	WebApp = web_app.GetIns()

	WebAuthPackage = web_auth.Package
	WebAuthShop    = web_auth.Shop
	AuthApp        = auth_app.GetIns()
	AuthLog        = auth_log.GetIns()
	AuthVersion    = auth_version.GetIns()
	UserAuth       = user_auth.GetIns()

	WebOauth  = web_oauth.GetIns()
	Oauth     = oauth.GetIns()
	OauthLog  = oauth_log.GetIns()
	OauthUser = oauth_user.GetIns()

	WebPlugin  = web_plugin.GetIns()
	Plugin     = plugin.GetIns()
	UserPlugin = user_plugin.GetIns()

	WebUser = web_user.GetIns()
)
