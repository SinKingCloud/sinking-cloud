package aliyun

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/aliyun"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"strconv"
	"sync"
)

// Service 单例对象
type Service struct {
}

// obj 单例对象
var (
	obj  *Service
	once sync.Once
)

// GetIns 获取单例
func GetIns() *Service {
	once.Do(func() {
		obj = &Service{}
	})
	return obj
}

// FaceInit 初始化人脸识别
func (s *Service) FaceInit(userId int, cardType string, cardNo string, cardName string, jumpUrl string, callbackUrl string) (string, string, error) {
	type Form struct {
		UserId   int    `form:"user_id" json:"user_id"`
		CardNo   string `form:"card_no" json:"card_no"`
		CardName string `form:"card_name" json:"card_name"`
	}
	form := &Form{
		UserId:   userId,
		CardNo:   cardNo,
		CardName: cardName,
	}
	configs := service.GetSystemConfigServiceIns().GetGroupByCache(config.AliyunGroup, false)
	configs2 := service.GetSystemConfigServiceIns().GetGroupByCache(config.SystemGroup, false)
	if jumpUrl == "" {
		jumpUrl = configs2[config.SystemUrl]
	}
	if callbackUrl == "" {
		callbackUrl = "/auth/callback/face"
	}
	id, url, err := aliyun.NewFace(configs[config.AliyunAccessKey], configs[config.AliyunAccessSecret]).
		Init(strconv.FormatInt(str.GetSnowWorkIns().GetId(), 10),
			configs[config.AliyunAuthSceneId],
			"FACE",
			cardType,
			cardNo,
			cardName,
			jumpUrl,
			configs2[config.SystemGateway]+callbackUrl,
		)
	if err != nil {
		return "", "", err
	}
	util.Cache.SetWithExpire("Face_"+id, func() string {
		marshal, err := json.Marshal(&form)
		if err != nil {
			return ""
		}
		return string(marshal)
	}(), 3600*1000)
	return id, url, nil
}

// FaceCheck 判断人脸识别
func (s *Service) FaceCheck(id string) (res bool, userId int, cardNo string, cardName string) {
	configs := service.GetSystemConfigServiceIns().GetGroupByCache(config.AliyunGroup, false)
	res, err := aliyun.NewFace(configs[config.AliyunAccessKey], configs[config.AliyunAccessSecret]).
		Query(configs[config.AliyunAuthSceneId], id)
	if err != nil {
		return false, 0, "", ""
	}
	if res {
		data, exists := util.Cache.Get("Face_" + id)
		if exists {
			type Form struct {
				UserId   int    `form:"user_id" json:"user_id"`
				CardNo   string `form:"card_no" json:"card_no"`
				CardName string `form:"card_name" json:"card_name"`
			}
			form := &Form{}
			err = json.Unmarshal([]byte(data), &form)
			if err == nil {
				return true, form.UserId, form.CardNo, form.CardName
			}
		}
	}
	return res, 0, "", ""
}
