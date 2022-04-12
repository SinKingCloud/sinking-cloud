package web_uin

// Setting 功能设置
var Setting = map[string]string{
	"mobile":       "1",
	"pc":           "1",
	"online":       "1",
	"zone":         "1",
	"medal":        "1",
	"ws":           "1",
	"micro_show":   "1",
	"safe":         "1",
	"game":         "1",
	"sport":        "1",
	"mobile_clock": "1",
	"zone_clock":   "1",
	"login":        "1",
	"friend":       "1",
}

// JobMap 映射表
var JobMap = map[string]map[string]int{
	"speed": {
		"mobile":       0,
		"pc":           1,
		"online":       2,
		"zone":         3,
		"medal":        4,
		"ws":           5,
		"micro_show":   6,
		"safe":         7,
		"game":         8,
		"sport":        9,
		"mobile_clock": 10,
		"zone_clock":   11,
		"login":        12,
		"friend":       13,
	},
	"sign": {},
}

// Project 项目列表
var Project = map[string]map[string]map[string]string{
	"speed": {
		"mobile": {
			"name":     "手机QQ",
			"describe": "QQ手机版连续在线满6小时，最高可累积1天",
			"logo":     "https://cdn.open-api.cn/app/uin/img/sq.png",
			"edit":     "true",
		},
		"pc": {
			"name":     "电脑QQ",
			"describe": "电脑QQ在线2小时加速0.5天",
			"logo":     "https://cdn.open-api.cn/app/uin/img/pc.png",
			"edit":     "true",
		},
		"online": {
			"name":     "非隐身在线",
			"describe": "非隐身时长满2小时，可加速0.2天(与电脑捆绑)",
			"logo":     "https://cdn.open-api.cn/app/uin/img/ys.png",
			"edit":     "false",
		},
		"zone": {
			"name":     "空间访客",
			"describe": "QQ空间今日访客数超过10人，可累积0.5天",
			"logo":     "https://cdn.open-api.cn/app/uin/img/kj.png",
			"edit":     "true",
		},
		"medal": {
			"name":     "勋章墙",
			"describe": "点亮QQ勋章墙图标，加速0.2天(与电脑捆绑)",
			"logo":     "https://cdn.open-api.cn/app/uin/img/xz.png",
			"edit":     "false",
		},
		"ws": {
			"name":     "微视观看",
			"describe": "微视看视频满30分钟，加速0.5天",
			"logo":     "https://cdn.open-api.cn/app/uin/img/ws.png",
			"edit":     "true",
		},
		"micro_show": {
			"name":     "厘米秀活跃",
			"describe": "厘米秀今日活跃达到70，加速0.2天",
			"logo":     "https://cdn.open-api.cn/app/uin/img/cm.png",
			"edit":     "true",
		},
		"safe": {
			"name":     "电脑管家",
			"describe": "电脑管家日登录达30分钟，加速0.2天",
			"logo":     "https://cdn.open-api.cn/app/uin/img/gj.png",
			"edit":     "true",
		},
		"game": {
			"name":     "QQ游戏",
			"describe": "玩10分钟游戏加速0.2天",
			"logo":     "https://cdn.open-api.cn/app/uin/img/sy.png",
			"edit":     "true",
		},
		"sport": {
			"name":     "QQ运动",
			"describe": "步行满5000步或跑步满1公里，加速0.2天（月限8次）",
			"logo":     "https://cdn.open-api.cn/app/uin/img/yd.png",
			"edit":     "true",
		},
		"mobile_clock": {
			"name":     "登陆打卡",
			"describe": "手机QQ签到打卡加速0.5天",
			"logo":     "https://cdn.open-api.cn/app/uin/img/dk.png",
			"edit":     "true",
		},
		"zone_clock": {
			"name":     "发布说说",
			"describe": "发布说说签到打卡加速0.5天(系统自动删除发布的说说)",
			"logo":     "https://cdn.open-api.cn/app/uin/img/ss.png",
			"edit":     "true",
		},
		"login": {
			"name":     "连续登录",
			"describe": "连续登录手机QQ加速0.5天(和手机QQ捆绑)",
			"logo":     "https://cdn.open-api.cn/app/uin/img/dl.png",
			"edit":     "false",
		},
		"friend": {
			"name":     "好友加速",
			"describe": "每日添加一位好友加速0.5天(系统自动加速,无需添加好友)",
			"logo":     "https://cdn.open-api.cn/app/uin/img/hy.png",
			"edit":     "true",
		},
	},
	"sign": {},
}
