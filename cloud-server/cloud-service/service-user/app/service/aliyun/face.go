package aliyun

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/aliyun"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-user/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-user/app/service/user"
	"strconv"
)

// FaceInit 初始化人脸识别
func FaceInit(userId int, cardType string, cardNo string, cardName string) (string, string, error) {
	type Form struct {
		UserId   int    `json:"user_id"`
		CardNo   string `json:"card_no"`
		CardName string `json:"card_name"`
	}
	form := &Form{
		UserId:   userId,
		CardNo:   cardNo,
		CardName: cardName,
	}
	configs := service.GetSystemConfigServiceIns().GetByCache(false)
	id, url, err := aliyun.NewFace(configs[config.AliyunAccessKey], configs[config.AliyunAccessSecret]).
		Init(strconv.FormatInt(str.GetSnowWorkIns().GetId(), 10),
			configs[config.AliyunAuthSceneId],
			"FACE",
			cardType,
			cardNo,
			cardName,
			configs[config.SystemUrl],
			configs[config.SystemGateway]+"/auth/callback/face",
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
func FaceCheck(id string) bool {
	configs := service.GetSystemConfigServiceIns().GetByCache(false)
	res, err := aliyun.NewFace(configs[config.AliyunAccessKey], configs[config.AliyunAccessSecret]).
		Query(configs[config.AliyunAuthSceneId], id)
	if err != nil {
		return false
	}
	if res {
		data, exists := util.Cache.Get("Face_" + id)
		if exists {
			type Form struct {
				UserId   int    `json:"user_id"`
				CardNo   string `json:"card_no"`
				CardName string `json:"card_name"`
			}
			form := &Form{}
			err := json.Unmarshal([]byte(data), &form)
			if err == nil {
				//写入信息
				service.GetWebUserConfSrvIns().SetConfigs(form.UserId, map[string]string{
					constant.UserConfigPersonCardNo:   form.CardNo,
					constant.UserConfigPersonCardName: form.CardName,
				})
				//更新用户信息
				_ = user.GetIns().UpdateById(form.UserId, map[string]interface{}{
					"real_auth": "1",
				})
			}
		}
	}
	return res
}
