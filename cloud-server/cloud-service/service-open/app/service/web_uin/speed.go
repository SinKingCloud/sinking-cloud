package web_uin

// Project 项目列表
var Project = map[string]map[string]map[string]string{
	"speed": {
		"mobile": {
			"name":     "手机QQ",
			"describe": "QQ手机版连续在线满6小时，最高可累积1天",
			"edit":     "true",
		},
		"pc": {
			"name":     "电脑QQ",
			"describe": "电脑QQ在线2小时加速0.5天",
			"edit":     "true",
		},
		"online": {
			"name":     "非隐身在线",
			"describe": "非隐身时长满2小时，可加速0.2天",
			"edit":     "false",
		},
		"zone": {
			"name":     "空间访客",
			"describe": "QQ空间今日访客数超过10人，可累积0.5天",
			"edit":     "true",
		},
		"medal": {
			"name":     "勋章墙",
			"describe": "点亮QQ勋章墙图标，加速0.2天",
			"edit":     "false",
		},
		"ws": {
			"name":     "微视观看",
			"describe": "微视看视频满30分钟，加速0.5天",
			"edit":     "true",
		},
		"micro_show": {
			"name":     "厘米秀活跃",
			"describe": "厘米秀今日活跃达到70，加速0.2天",
			"edit":     "true",
		},
		"safe": {
			"name":     "电脑管家",
			"describe": "电脑管家日登录达30分钟，加速0.2天",
			"edit":     "true",
		},
		"game": {
			"name":     "QQ游戏",
			"describe": "手Q→动态→游戏→启动并玩10分钟游戏 ＋0.2天",
			"edit":     "true",
		},
		"sport": {
			"name":     "QQ运动",
			"describe": "步行满5000步或跑步满1公里，加速0.2天（月限8次）",
			"edit":     "true",
		},
		"mobile_clock": {
			"name":     "手Q签到",
			"describe": "手机QQ签到打卡加速0.5天",
			"edit":     "true",
		},
		"zone_clock": {
			"name":     "空间签到",
			"describe": "空间签到打卡加速0.5天",
			"edit":     "true",
		},
	},
}
