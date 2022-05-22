package uin_tool

import (
	"fmt"
	"testing"
)

const (
	uin      = "2013654855"
	sKey     = "Mm8JB7bAtT"
	psKey    = "VgD7BWoO6dEm4vXWzpTg1yAdL1MjhW6yWTJVklk33ac_"
	superKey = "apBIw0Tb1*iKGsYH9XWEHm4OcQfThfey-tB9xJxV0*E_"
)

func Test_Gtk(t *testing.T) {
	sign := NewSignTool(uin, sKey, psKey, superKey)
	fmt.Println(sign.getGTK("0123"))
}

func Test_Token(t *testing.T) {
	sign := NewSignTool(uin, sKey, psKey, superKey)
	fmt.Println(sign.getToken2(superKey))
}

func Test_VipSign(t *testing.T) {
	sign := NewSignTool(uin, sKey, psKey, superKey)
	sign.VipSign()
}

func Test_BigVipSign(t *testing.T) {
	sign := NewSignTool(uin, sKey, psKey, superKey)
	sign.BigVipSign()
}

func Test_YellowDiaSin(t *testing.T) {
	sign := NewSignTool(uin, sKey, psKey, superKey)
	sign.YellowDiaSin()
}

func Test_BlueDiaSin(t *testing.T) {
	sign := NewSignTool(uin, sKey, psKey, superKey)
	sign.BlueDiaSin()
}

func Test_DaLeDouSign(t *testing.T) {
	sign := NewSignTool(uin, sKey, psKey, superKey)
	sign.DaleDouSign()
}

func Test_GameCenterSign(t *testing.T) {
	sign := NewSignTool(uin, sKey, psKey, superKey)
	sign.GameCenterSign()
}
