package ai

import (
	"fmt"
	"os"
)

// Accurate 通用文字识别
func (a *ai) Accurate() {
	post := map[string]string{
		"detect_direction":  "true",
		"vertexes_location": "true",
	}
	bytes, _ := os.ReadFile("img.png")
	file := map[string][]byte{
		"image": bytes,
	}
	body, err := a.SendFormDataRequest("/weapp/rest/2.0/ocr/v1/accurate", &post, &file)
	fmt.Println(string(body), err)
}
