package web_uin

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/util/oicq"
)

type JobName string

const (
	Mobile      JobName = "mobile"
	Pc          JobName = "pc"
	Ws          JobName = "ws"
	MobileClock JobName = "mobile_clock"
	ZoneClock   JobName = "zone_clock"
	Login       JobName = "login"
	Friend      JobName = "friend"
	Word        JobName = "word"
	QQMusic     JobName = "qq_music"
	Vip         JobName = "vip"
	BigVip      JobName = "big_vip"
	Hz          JobName = "hz"
	Lz          JobName = "lz"
	Ld          JobName = "ld"
	GameCenter  JobName = "game_center"
	QQGame      JobName = "qq_game"
	Sweet       JobName = "sweet"
	Video       JobName = "video"
	WeiYun      JobName = "weiyun"
)

type JobType string

const (
	Speed JobType = "speed"
	Sign  JobType = "sign"
)

// Setting 功能设置
var Setting = map[JobName]string{
	Mobile:      "1",
	Pc:          "0",
	Ws:          "1",
	MobileClock: "1",
	ZoneClock:   "1",
	Login:       "1",
	Friend:      "1",
	Word:        "1",
	QQMusic:     "1",
	Vip:         "1",
	BigVip:      "1",
	Hz:          "1",
	Lz:          "1",
	Ld:          "1",
	GameCenter:  "1",
	QQGame:      "1",
	Sweet:       "1",
	Video:       "1",
	WeiYun:      "1",
}

// JobMap 映射表
var JobMap = map[JobType]map[JobName]oicq.Job{
	Speed: {
		Mobile:      oicq.None,
		Pc:          oicq.Pc,
		Ws:          oicq.Ws,
		MobileClock: oicq.MobileClock,
		ZoneClock:   oicq.ZoneClock,
		Login:       oicq.None,
		Friend:      oicq.Friend,
		QQMusic:     oicq.QQMusic,
		Word:        oicq.Word,
	},
	Sign: {
		Vip:        oicq.Vip,
		BigVip:     oicq.BigVip,
		Hz:         oicq.Hz,
		Lz:         oicq.Lz,
		Ld:         oicq.Ld,
		GameCenter: oicq.GameCenter,
		QQGame:     oicq.QQGame,
		Sweet:      oicq.Sweet,
		Video:      oicq.Video,
		WeiYun:     oicq.WeiYun,
	},
}

// Project 项目列表
var Project = map[JobType]map[JobName]map[string]string{
	Speed: {
		Mobile: {
			"sort":     "1",
			"type":     "speed",
			"name":     "手机在线",
			"describe": "手机在线满6小时加速1天",
			"logo":     "https://cdn.clwl.online/app/uin/img/sq.png",
			"edit":     "true",
		},
		Pc: {
			"sort":     "2",
			"type":     "speed",
			"name":     "电脑在线",
			"describe": "电脑在线2小时加速0.5天",
			"logo":     "https://cdn.clwl.online/app/uin/img/pc.png",
			"edit":     "true",
		},
		Ws: {
			"sort":     "6",
			"type":     "speed",
			"name":     "微视观看",
			"describe": "微视看视频30分钟加速0.5天",
			"logo":     "https://cdn.clwl.online/app/uin/img/ws.png",
			"edit":     "true",
		},
		Word: {
			"sort":     "7",
			"type":     "speed",
			"name":     "发小世界",
			"describe": "每日发布动态加速0.5天",
			"logo":     "https://cdn.clwl.online/app/uin/img/xsj.png",
			"edit":     "true",
		},
		MobileClock: {
			"sort":     "11",
			"type":     "speed",
			"name":     "登陆打卡",
			"describe": "手机签到打卡加速0.5天",
			"logo":     "https://cdn.clwl.online/app/uin/img/dk.png",
			"edit":     "true",
		},
		ZoneClock: {
			"sort":     "12",
			"type":     "speed",
			"name":     "发布说说",
			"describe": "发布说说签到打卡加速0.5天",
			"logo":     "https://cdn.clwl.online/app/uin/img/ss.png",
			"edit":     "true",
		},
		Login: {
			"sort":     "13",
			"type":     "speed",
			"name":     "连续登录",
			"describe": "连续登录手机加速0.5天",
			"logo":     "https://cdn.clwl.online/app/uin/img/dl.png",
			"edit":     "false",
		},
		Friend: {
			"sort":     "14",
			"type":     "speed",
			"name":     "好友加速",
			"describe": "每日添加一位好友加速0.5天",
			"logo":     "https://cdn.clwl.online/app/uin/img/hy.png",
			"edit":     "true",
		},
		QQMusic: {
			"sort":     "15",
			"type":     "speed",
			"name":     "音乐加速",
			"describe": "每日听歌30分钟加速0.5天",
			"logo":     "https://cdn.clwl.online/app/uin/img/yq.png",
			"edit":     "true",
		},
	},
	Sign: {
		Vip: {
			"sort":     "15",
			"type":     "sign",
			"name":     "会员成长",
			"describe": "SVIP6级年费以上用户，每月SVIP成长值840点左右",
			"logo":     "https://cdn.clwl.online/app/uin/img/vip.png",
			"edit":     "true",
		},
		BigVip: {
			"sort":     "16",
			"type":     "sign",
			"name":     "大会员成长",
			"describe": "大会员签到/访客/点赞/师徒/分享，每日25点成长值",
			"logo":     "https://cdn.clwl.online/app/uin/img/big_vip.png",
			"edit":     "true",
		},
		Hz: {
			"sort":     "17",
			"type":     "sign",
			"name":     "黄钻成长",
			"describe": "黄钻每日完成三次打卡任务，每日获得15点成长值",
			"logo":     "https://cdn.clwl.online/app/uin/img/hz.png",
			"edit":     "true",
		},
		Lz: {
			"sort":     "18",
			"type":     "sign",
			"name":     "蓝钻成长",
			"describe": "蓝钻每日领取两次成长值4+4，每日获得8点成长值",
			"logo":     "https://cdn.clwl.online/app/uin/img/lz.png",
			"edit":     "true",
		},
		Ld: {
			"sort":     "19",
			"type":     "sign",
			"name":     "乐斗签到",
			"describe": "自动完成QQ空间Q宠大乐斗每日签到、领取礼包等",
			"logo":     "https://cdn.clwl.online/app/uin/img/ld.png",
			"edit":     "true",
		},
		GameCenter: {
			"sort":     "20",
			"type":     "sign",
			"name":     "游戏中心",
			"describe": "每日完成手Q游戏中心日常任务获得游戏币和活跃值",
			"logo":     "https://cdn.clwl.online/app/uin/img/game_center.png",
			"edit":     "true",
		},
		QQGame: {
			"sort":     "22",
			"type":     "sign",
			"name":     "QQ游戏",
			"describe": "QQ游戏,道聚城,精品页游签到领取礼包",
			"logo":     "https://cdn.clwl.online/app/uin/img/dr.png",
			"edit":     "true",
		},
		Sweet: {
			"sort":     "23",
			"type":     "sign",
			"name":     "情侣空间",
			"describe": "每日情侣空间签到发表状态",
			"logo":     "https://cdn.clwl.online/app/uin/img/ql.png",
			"edit":     "true",
		},
		Video: {
			"sort":     "24",
			"type":     "sign",
			"name":     "腾讯视频",
			"describe": "腾讯视频会员每日签到领取V力值",
			"logo":     "https://cdn.clwl.online/app/uin/img/sp.png",
			"edit":     "true",
		},
		WeiYun: {
			"sort":     "25",
			"type":     "sign",
			"name":     "微云签到",
			"describe": "微云每日签到领取积分",
			"logo":     "https://cdn.clwl.online/app/uin/img/xy.png",
			"edit":     "true",
		},
	},
}
