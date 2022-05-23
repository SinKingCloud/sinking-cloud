package uin_tool

import (
	"fmt"
	"testing"
)

const (
	uin      = "2013654855"
	proxy    = "117.68.107.99:1088"
	sKey     = "M0VrS1nJJC"
	psKey    = "rrApGYyZnLHRG8E1CzHYKXEeHYJzAGg61hIpOoJBVTw_"
	superKey = "h2tdNoPz3IoeH*uMQ*hn*Vrb9Ny56ByT0-5dEJoxXRI_"
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
