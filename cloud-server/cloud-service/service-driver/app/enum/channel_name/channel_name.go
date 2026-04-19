package channel_name

import "service-driver/app/service"

// Map 映射数据
func Map() map[int64]string {
	data, err := service.Channel.GetIdNameMap(false)
	if err != nil {
		return map[int64]string{}
	}
	return data
}
