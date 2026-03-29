package str

import (
	"math/rand"
	"time"
)

// NickTool 网名生成工具
type NickTool struct {
	r *rand.Rand
}

// NewNickTool 实例化（带随机种子，避免每次生成一样）
func NewNickTool() *NickTool {
	return &NickTool{
		r: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// 前缀：风格词
func (n *NickTool) getPrefix() []string {
	return []string{
		"温柔", "安静", "快乐", "星星", "月亮", "晚风", "橘子", "草莓",
		"汽水", "奶糖", "芝士", "云朵", "深海", "森林", "夏日", "冬日",
		"心动", "微光", "甜柚", "白茶", "清欢", "拾光", "晚风", "雾岛",
		"懒懒", "软软", "暖暖", "皮皮", "呆呆", "萌萌", "小小", "落落",
		"芝士", "泡芙", "蜜桃", "椰果", "冰柠", "薄荷", "蓝莓", "荔枝",
	}
}

// 后缀：名词/状态词
func (n *NickTool) getSuffix() []string {
	return []string{
		"邮递员", "收藏家", "便利店", "小星球", "贩卖机", "观察员",
		"不打烊", "睡着了", "在逃中", "限定款", "爱好者", "小队长",
		"日记", "晚风", "日落", "星河", "月亮", "岛屿", "山海", "烟火",
		"小熊", "兔子", "猫咪", "小狗", "狐狸", "企鹅", "小鹿", "仓鼠",
		"泡泡", "星星", "云朵", "汽水", "面包", "饼干", "布丁", "奶茶",
	}
}

// 英文前缀（简约风）
func (n *NickTool) getEnPrefix() []string {
	return []string{
		"Lovely", "Sweet", "Happy", "Sunny", "Moon", "Star",
		"Gentle", "Cute", "Cool", "Easy", "Free", "Light",
	}
}

// GetRand 生成自然好看的随机网名（主方法）
func (n *NickTool) GetRand() string {
	// 随机风格 1~4
	style := n.r.Intn(4) + 1

	prefix := n.getPrefix()
	suffix := n.getSuffix()
	enPrefix := n.getEnPrefix()

	p := prefix[n.r.Intn(len(prefix))]
	s := suffix[n.r.Intn(len(suffix))]
	en := enPrefix[n.r.Intn(len(enPrefix))]

	switch style {
	case 1:
		// 风格1：前缀+后缀（文艺）
		return p + s
	case 2:
		// 风格2：前缀+的+后缀（可爱）
		return p + "的" + s
	case 3:
		// 风格3：英文+中文（洋气）
		return en + p
	case 4:
		// 风格4：简约两字
		return p[:2]
	default:
		return p + s
	}
}
