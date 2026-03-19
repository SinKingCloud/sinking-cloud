package service

type AuthService interface {
	VerifyUserToken(token string) error  //验证网站用户token有效性
	VerifyAdminToken(token string) error //验证后台用户token有效性
}
