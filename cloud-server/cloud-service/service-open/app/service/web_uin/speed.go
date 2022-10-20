package web_uin

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/util/oicq"
)

// Setting 功能设置
var Setting = map[string]string{
	"mobile":       "1",
	"pc":           "1",
	"online":       "1",
	"ws":           "1",
	"mobile_clock": "1",
	"zone_clock":   "1",
	"login":        "1",
	"friend":       "1",
	"word":         "1",
	"vip":          "1",
	"big_vip":      "1",
	"hz":           "1",
	"lz":           "1",
	"ld":           "1",
	"game_center":  "1",
	"yy":           "1",
}

// JobMap 映射表
var JobMap = map[string]map[string]oicq.Job{
	"speed": {
		"mobile":       oicq.None,
		"pc":           oicq.Pc,
		"online":       oicq.None,
		"ws":           oicq.Ws,
		"mobile_clock": oicq.MobileClock,
		"zone_clock":   oicq.ZoneClock,
		"login":        oicq.None,
		"friend":       oicq.Friend,
		"word":         oicq.Word,
	},
	"sign": {
		"vip":         oicq.Vip,
		"big_vip":     oicq.BigVip,
		"hz":          oicq.Hz,
		"lz":          oicq.Lz,
		"ld":          oicq.Ld,
		"game_center": oicq.GameCenter,
		"yy":          oicq.Music,
	},
}

// Project 项目列表
var Project = map[string]map[string]map[string]string{
	"speed": {
		"mobile": {
			"sort":     "1",
			"type":     "speed",
			"name":     "手机在线",
			"describe": "手机在线满6小时加速1天",
			"logo":     "https://cdn.yunhuoke.cn/app/uin/img/sq.png",
			"edit":     "true",
		},
		"pc": {
			"sort":     "2",
			"type":     "speed",
			"name":     "电脑在线",
			"describe": "电脑在线2小时加速0.5天",
			"logo":     "https://cdn.yunhuoke.cn/app/uin/img/pc.png",
			"edit":     "true",
		},
		"online": {
			"sort":     "3",
			"type":     "speed",
			"name":     "非隐身在线",
			"describe": "非隐身时长满2小时加速0.2天",
			"logo":     "https://cdn.yunhuoke.cn/app/uin/img/ys.png",
			"edit":     "false",
		},
		"ws": {
			"sort":     "6",
			"type":     "speed",
			"name":     "微视观看",
			"describe": "微视看视频满30分钟，加速0.5天",
			"logo":     "https://cdn.yunhuoke.cn/app/uin/img/ws.png",
			"edit":     "true",
		},
		"word": {
			"sort":     "7",
			"type":     "speed",
			"name":     "小世界加速",
			"describe": "每日发布动态加速0.5天",
			"logo":     "https://cdn.yunhuoke.cn/app/uin/img/xsj.png",
			"edit":     "true",
		},
		"mobile_clock": {
			"sort":     "11",
			"type":     "speed",
			"name":     "登陆打卡",
			"describe": "手机签到打卡加速0.5天",
			"logo":     "https://cdn.yunhuoke.cn/app/uin/img/dk.png",
			"edit":     "true",
		},
		"zone_clock": {
			"sort":     "12",
			"type":     "speed",
			"name":     "发布说说",
			"describe": "发布说说签到打卡加速0.5天",
			"logo":     "https://cdn.yunhuoke.cn/app/uin/img/ss.png",
			"edit":     "true",
		},
		"login": {
			"sort":     "13",
			"type":     "speed",
			"name":     "连续登录",
			"describe": "连续登录手机加速0.5天",
			"logo":     "https://cdn.yunhuoke.cn/app/uin/img/dl.png",
			"edit":     "false",
		},
		"friend": {
			"sort":     "14",
			"type":     "speed",
			"name":     "好友加速",
			"describe": "每日添加一位好友加速0.5天",
			"logo":     "https://cdn.yunhuoke.cn/app/uin/img/hy.png",
			"edit":     "true",
		},
	},
	"sign": {
		"vip": {
			"sort":     "15",
			"type":     "sign",
			"name":     "会员成长",
			"describe": "SVIP6级年费以上用户，每月SVIP成长值840点左右",
			"logo":     "https://cdn.yunhuoke.cn/app/uin/img/vip.png",
			"edit":     "true",
		},
		"big_vip": {
			"sort":     "16",
			"type":     "sign",
			"name":     "大会员成长",
			"describe": "大会员签到/访客/点赞/师徒/分享，每日25点成长值",
			"logo":     "https://cdn.yunhuoke.cn/app/uin/img/big_vip.png",
			"edit":     "true",
		},
		"hz": {
			"sort":     "17",
			"type":     "sign",
			"name":     "黄钻成长",
			"describe": "黄钻每日完成三次打卡任务，每日获得15点成长值",
			"logo":     "https://cdn.yunhuoke.cn/app/uin/img/hz.png",
			"edit":     "true",
		},
		"lz": {
			"sort":     "18",
			"type":     "sign",
			"name":     "蓝钻成长",
			"describe": "蓝钻每日领取两次成长值4+4，每日获得8点成长值",
			"logo":     "https://cdn.yunhuoke.cn/app/uin/img/lz.png",
			"edit":     "true",
		},
		"ld": {
			"sort":     "19",
			"type":     "sign",
			"name":     "乐斗签到",
			"describe": "自动完成QQ空间Q宠大乐斗每日签到、领取礼包等",
			"logo":     "https://cdn.yunhuoke.cn/app/uin/img/ld.png",
			"edit":     "true",
		},
		"game_center": {
			"sort":     "20",
			"type":     "sign",
			"name":     "游戏中心",
			"describe": "每日完成手Q游戏中心日常任务获得游戏币和活跃值",
			"logo":     "https://cdn.yunhuoke.cn/app/uin/img/game_center.png",
			"edit":     "true",
		},
		"yy": {
			"sort":     "21",
			"type":     "sign",
			"name":     "音乐签到",
			"describe": "每日完成QQ音乐签到和刷歌时长",
			"logo":     "https://cdn.yunhuoke.cn/app/uin/img/yq.png",
			"edit":     "true",
		},
	},
}
