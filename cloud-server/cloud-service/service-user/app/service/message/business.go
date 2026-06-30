package message

import (
	"errors"
	"service-user/app/enum/message_status"
	"service-user/app/model"
	repositoryMessage "service-user/app/repository/message"
	"strconv"
	"strings"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/socket"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// Create 创建用户消息
func (s *service) Create(userId int64, title string, content string) error {
	if userId <= 0 {
		return errors.New("用户信息错误")
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return errors.New("消息标题不能为空")
	}
	if err := s.repositoryMessage.Create(&model.Message{
		Id:      str.GetSnowWorkIns().GetId(),
		UserId:  userId,
		Title:   title,
		Content: content,
		Status:  message_status.Unread,
	}); err != nil {
		return err
	}
	s.notifyMessageRefresh([]int64{userId})
	return nil
}

// Creates 批量创建用户消息
func (s *service) Creates(userIds []int64, title string, content string) error {
	if len(userIds) == 0 {
		return errors.New("用户列表不能为空")
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return errors.New("消息标题不能为空")
	}
	messages := make([]*model.Message, 0, len(userIds))
	notifyUserIds := make([]int64, 0, len(userIds))
	userMap := make(map[int64]struct{}, len(userIds))
	for _, userId := range userIds {
		if userId <= 0 {
			continue
		}
		if _, ok := userMap[userId]; ok {
			continue
		}
		userMap[userId] = struct{}{}
		notifyUserIds = append(notifyUserIds, userId)
		messages = append(messages, &model.Message{
			Id:      str.GetSnowWorkIns().GetId(),
			UserId:  userId,
			Title:   title,
			Content: content,
			Status:  message_status.Unread,
		})
	}
	if len(messages) == 0 {
		return errors.New("用户列表不能为空")
	}
	if err := s.repositoryMessage.Creates(messages); err != nil {
		return err
	}
	s.notifyMessageRefresh(notifyUserIds)
	return nil
}

// Read 设置用户消息为已读
func (s *service) Read(userId int64, ids []int64) error {
	if userId <= 0 {
		return errors.New("用户信息错误")
	}
	if len(ids) == 0 {
		return errors.New("消息列表不能为空")
	}
	readIds := make([]int64, 0, len(ids))
	idMap := make(map[int64]struct{}, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := idMap[id]; ok {
			continue
		}
		idMap[id] = struct{}{}
		readIds = append(readIds, id)
	}
	if len(readIds) == 0 {
		return errors.New("消息列表不能为空")
	}
	if err := s.repositoryMessage.ReadByIdsAndUserId(readIds, userId, str.DateTime(time.Now())); err != nil {
		return err
	}
	s.notifyMessageRefresh([]int64{userId})
	return nil
}

// Delete 删除用户消息
func (s *service) Delete(userId int64, ids []int64) error {
	if userId <= 0 {
		return errors.New("用户信息错误")
	}
	if len(ids) == 0 {
		return errors.New("消息列表不能为空")
	}
	deleteIds := make([]int64, 0, len(ids))
	idMap := make(map[int64]struct{}, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := idMap[id]; ok {
			continue
		}
		idMap[id] = struct{}{}
		deleteIds = append(deleteIds, id)
	}
	if len(deleteIds) == 0 {
		return errors.New("消息列表不能为空")
	}
	if err := s.repositoryMessage.DeleteByIdsAndUserId(deleteIds, userId); err != nil {
		return err
	}
	s.notifyMessageRefresh([]int64{userId})
	return nil
}

// Clear 清空用户消息
func (s *service) Clear(userId int64) error {
	if userId <= 0 {
		return errors.New("用户信息错误")
	}
	if err := s.repositoryMessage.DeleteByUserId(userId); err != nil {
		return err
	}
	s.notifyMessageRefresh([]int64{userId})
	return nil
}

// SelectByUserId 查询用户消息列表
func (s *service) SelectByUserId(userId int64, where *repositoryMessage.SelectMessage, orderByField string, orderByType string, page int, pageSize int) ([]*Message, int64, error) {
	if userId <= 0 {
		return nil, 0, errors.New("用户信息错误")
	}
	if where == nil {
		where = &repositoryMessage.SelectMessage{}
	}
	where.UserId = userId
	if where.Status != "" {
		status, err := strconv.Atoi(where.Status)
		if err != nil {
			return nil, 0, errors.New("消息状态参数错误")
		}
		if _, ok := message_status.Map()[status]; !ok {
			return nil, 0, errors.New("消息状态参数不合法")
		}
	}
	list, total, err := s.repositoryMessage.Select(where, orderByField, orderByType, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*Message, 0, len(list))
	for _, item := range list {
		if item == nil {
			continue
		}
		result = append(result, &Message{
			Id:         item.Id,
			Title:      item.Title,
			Status:     item.Status,
			ReadTime:   item.ReadTime,
			CreateTime: item.CreateTime,
			UpdateTime: item.UpdateTime,
		})
	}
	return result, total, nil
}

// GetInfoByUserId 查询用户消息详情，未读消息会自动标记为已读。
func (s *service) GetInfoByUserId(userId int64, id int64) (*Message, error) {
	if userId <= 0 || id <= 0 {
		return nil, errors.New("参数错误")
	}
	info, err := s.repositoryMessage.FindByIdAndUserId(id, userId)
	if err != nil {
		return nil, err
	}
	if info == nil || info.Id <= 0 {
		return nil, errors.New("消息不存在")
	}
	if info.Status == message_status.Unread {
		readTime := str.DateTime(time.Now())
		err = s.repositoryMessage.UpdateByIdAndUserId(info.Id, userId, &repositoryMessage.UpdateMessage{
			Status:   message_status.Read,
			ReadTime: readTime,
		})
		if err != nil {
			return nil, err
		}
		info.Status = message_status.Read
		info.ReadTime = readTime
		info.UpdateTime = readTime
		s.notifyMessageRefresh([]int64{userId})
	}
	return &Message{
		Id:         info.Id,
		Title:      info.Title,
		Content:    info.Content,
		Status:     info.Status,
		ReadTime:   info.ReadTime,
		CreateTime: info.CreateTime,
		UpdateTime: info.UpdateTime,
	}, nil
}

// notifyMessageRefresh 通知前端重新拉取消息列表；提醒失败不影响主业务。
func (s *service) notifyMessageRefresh(userIds []int64) {
	if s.socketService == nil || len(userIds) == 0 {
		return
	}
	ids := make([]int64, 0, len(userIds))
	idMap := make(map[int64]struct{}, len(userIds))
	for _, userId := range userIds {
		if userId <= 0 {
			continue
		}
		if _, ok := idMap[userId]; ok {
			continue
		}
		idMap[userId] = struct{}{}
		ids = append(ids, userId)
	}
	if len(ids) == 0 {
		return
	}
	message := &socket.Message{
		Id:   strconv.FormatInt(str.GetSnowWorkIns().GetId(), 10),
		Type: socket.MessageRefresh,
		Data: map[string]interface{}{
			"refresh": true,
		},
		Time: time.Now().UnixMilli(),
	}
	_ = s.socketService.SendMessage(message, ids)
}
