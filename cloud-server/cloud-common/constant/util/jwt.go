package util

import "time"

const (
	JwtTokenExpire   = time.Duration(86400 * 30) //token过期时间
	JwtTokenKey      = "SinKingCloud-Key"        //token加密key(必须为16位)
	JwtTokenName     = "sinking-token"           //用户token名称
	JwtTokenNameType = "sinking-type"            //jwt类型
)
