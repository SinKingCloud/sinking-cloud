package auth

import (
	"encoding/json"
	"errors"
	"service-auth/app/constant"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/captcha"
	"github.com/wenlng/go-captcha/v2/slide"
)

// GetCaptcha 获取滑块验证码
func (s *service) GetCaptcha(key string) (map[string]interface{}, error) {
	result, masterImage, tileImage, height, width, err := captcha.GenSlide()
	if err != nil {
		return nil, err
	}
	bytes, err := json.Marshal(result)
	if err != nil {
		return nil, errors.New("获取验证码失败")
	}
	set := s.cache.SetWithExpire(constant.CacheNameWithCaptcha+key, string(bytes), constant.CacheTimeWithCaptcha)
	if !set {
		return nil, errors.New("存储验证码失败")
	}
	ret := map[string]interface{}{
		"key":          key,
		"image_base64": masterImage,
		"width":        width,
		"height":       height,
		"tile_base64":  tileImage,
		"tile_width":   result.Width,
		"tile_height":  result.Height,
		"tile_x":       result.DX,
		"tile_y":       result.DY,
	}
	return ret, nil
}

// CheckCaptcha 判断滑块验证码是否正确
func (s *service) CheckCaptcha(key string, x int, y int) bool {
	value, exists := s.cache.Get(constant.CacheNameWithCaptcha + key)
	s.cache.Delete(constant.CacheNameWithCaptcha + key)
	if !exists {
		return false
	}
	var dct *slide.Block
	if err := json.Unmarshal([]byte(value), &dct); err != nil {
		return false
	}
	return captcha.CheckSlide(x, y, dct)
}
