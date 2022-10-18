package oicq

import "errors"

// GetCookies 根据域名获取sKey pSKey
func (o *Oicq) GetCookies(domain string) (string, string) {
	return o.Mobile.GetCookiesWithDomain(domain)
}

// QrLogin 扫码授权登陆
func (o *Oicq) QrLogin(pk string) error {
	if !o.Mobile.Online.Load() {
		return ErrNotLogin
	}
	_, err := o.Mobile.QrCodeAllow(pk)
	if err != nil {
		return errors.New("授权登录失败")
	}
	return nil
}
