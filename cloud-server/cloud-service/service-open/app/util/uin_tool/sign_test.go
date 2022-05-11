package uin_tool

import (
	"fmt"
	"testing"
)

const (
	uin   = "1178710004"
	sKey  = "Mh8A8iSNkN"
	psKey = "WvKr9DFbhqortgXJYxK7lOFvoNVlIOgnBr4OSqGevaM_"
)

func Test_Gtk(t *testing.T) {
	sign := NewSign(uin, sKey, psKey)
	fmt.Println(sign.getGTK("0123"))
}

func Test_VipSign(t *testing.T) {
	sign := NewSign(uin, sKey, psKey)
	sign.VipSign()
}

func Test_BigVipSign(t *testing.T) {
	sign := NewSign(uin, sKey, psKey)
	sign.BigVipSign()
}

func Test_YellowDiaSin(t *testing.T) {
	sign := NewSign(uin, sKey, psKey)
	sign.YellowDiaSin()
}
