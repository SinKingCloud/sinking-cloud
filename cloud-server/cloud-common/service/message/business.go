package message

import (
	"errors"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
)

// Create 创建用户消息
func (*Service) Create(userId int64, title string, content string) error {
	var result bool
	err := rpc.InvokeFunc(constant.ServiceUser, "CreateMessage", &CreateMessageParam{
		UserId:  userId,
		Title:   title,
		Content: content,
	}, &result)
	if err != nil {
		return err
	}
	if !result {
		return errors.New("创建消息失败")
	}
	return nil
}

// Creates 批量创建用户消息
func (*Service) Creates(userIds []int64, title string, content string) error {
	var result bool
	err := rpc.InvokeFunc(constant.ServiceUser, "CreateMessages", &CreateMessagesParam{
		UserIds: userIds,
		Title:   title,
		Content: content,
	}, &result)
	if err != nil {
		return err
	}
	if !result {
		return errors.New("创建消息失败")
	}
	return nil
}
