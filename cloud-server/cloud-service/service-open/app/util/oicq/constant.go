package oicq

import "errors"

var (
	ErrNotLogin   = errors.New("account is not login")
	ErrLoginError = errors.New("account is login error")
)

type Job int //挂机项目

const (
	None        Job = iota //空占位
	Pc                     //Pc加速
	Ws                     //微视加速
	Word                   //小世界加速
	Game                   //游戏加速
	MobileClock            //手Q打卡
	ZoneClock              //空间打卡
	Friend                 //加好友加速

	Vip        //会员签到
	BigVip     //大会员签到
	Hz         //黄钻签到
	Lz         //蓝钻签到
	Ld         //大乐斗签到
	GameCenter //游戏中心签到
	Music      //音乐签到
	QQGame     //QQ游戏签到
	Sweet      //情侣空间签到
	Video      //腾讯视频签到
	WeiYun     //微云签到
)
