package util

import "time"

const (
	JwtTokenExpire         = time.Duration(86400 * 30) //token过期时间
	JwtTokenKey            = "SinKingCloudJwtTokenKey" //token加密key
	JwtTokenNameWebUser    = "sinking-web-token"       //网站用户token名称
	JwtTokenNameSystemUser = "sinking-system-token"    //系统用户token名称
)
