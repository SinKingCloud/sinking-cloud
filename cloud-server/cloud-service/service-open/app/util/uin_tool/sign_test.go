package uin_tool

import (
	"fmt"
	"testing"
)

const (
	uin      = "1178710004"
	proxy    = "117.68.107.99:1088"
	sKey     = "MOsO3udek1"
	psKey    = "9VccAD0z2MNshIZ2Ej3YyXOt5UeuifCgZ-4HhG8SHoI_"
	superKey = ""
)

func Test_Gtk(t *testing.T) {
	sign := NewSignTool(uin, proxy, sKey, psKey, superKey)
	fmt.Println(sign.getGTK("0123"))
}

func Test_Token(t *testing.T) {
	sign := NewSignTool(uin, proxy, sKey, psKey, superKey)
	fmt.Println(sign.getToken2(superKey))
}

func Test_VipSign(t *testing.T) {
	sign := NewSignTool(uin, proxy, sKey, psKey, superKey)
	sign.VipSign()
}

func Test_BigVipSign(t *testing.T) {
	sign := NewSignTool(uin, proxy, sKey, psKey, superKey)
	sign.BigVipSign()
}

func Test_YellowDiaSin(t *testing.T) {
	sign := NewSignTool(uin, proxy, sKey, psKey, superKey)
	sign.YellowDiaSin()
}

func Test_BlueDiaSin(t *testing.T) {
	sign := NewSignTool(uin, proxy, sKey, psKey, superKey)
	sign.BlueDiaSin()
}

func Test_DaLeDouSign(t *testing.T) {
	sign := NewSignTool(uin, proxy, sKey, psKey, superKey)
	sign.DaleDouSign()
}

func Test_GameCenterSign(t *testing.T) {
	sign := NewSignTool(uin, proxy, sKey, psKey, superKey)
	sign.GameCenterSign()
}

func Test_QQGamSign(t *testing.T) {
	sign := NewSignTool(uin, proxy, sKey, psKey, superKey)
	sign.QQGamSign()
}

func Test_WeiYunSign(t *testing.T) {
	sign := NewSignTool(uin, proxy, sKey, psKey, superKey)
	sign.WeiYunSign()
}

func Test_QQMusicSign(t *testing.T) {
	sign := NewSignTool(uin, proxy, sKey, psKey, superKey)
	sign.QQMusicSign()
}

func Test_VideoSign(t *testing.T) {
	sign := NewSignTool(uin, proxy, sKey, psKey, superKey)
	sign.VideoSign()
}
