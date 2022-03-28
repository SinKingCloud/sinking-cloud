package speed

import "encoding/json"

// toJson 转json
func (*speed) toJson(data interface{}) string {
	str, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(str)
}
