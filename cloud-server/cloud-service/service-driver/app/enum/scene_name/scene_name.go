package scene_name

import "service-driver/app/service"

// Map 映射数据
func Map() map[int64]string {
	data, err := service.Scene.GetIdNameMap(false)
	if err != nil {
		return map[int64]string{}
	}
	return data
}
