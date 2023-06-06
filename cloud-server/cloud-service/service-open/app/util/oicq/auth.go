package oicq

import (
	"errors"
)

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

// AllowDeviceLock 扫码授权设备锁
func (o *Oicq) AllowDeviceLock(verifyUrl string) error {
	if !o.Mobile.Online.Load() {
		return ErrNotLogin
	}
	token, err := o.Mobile.GetQrCodeAuthToken(verifyUrl)
	if err != nil || token == "" {
		return errors.New("解析授权token")
	}
	err = o.Mobile.QrCodeAuth(token)
	if err != nil {
		return errors.New("授权登录失败")
	}
	return nil
}
