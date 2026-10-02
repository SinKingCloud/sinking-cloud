package user

import (
	"errors"
	"fmt"
	"log"
	"service-admin/app/constant"
	"service-admin/app/model"
	"service-admin/app/queue/user_log"
	userLogRepository "service-admin/app/repository/user_log"
	"strconv"
	"strings"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/user/user_log_type"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/ip"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// SelectLog 查询用户操作日志
func (s *service) SelectLog(where *userLogRepository.SelectUserLog, queryPage *page.Query) (*page.Result[*model.UserLog], error) {
	if where != nil {
		if where.UserId != nil && *where.UserId <= 0 {
			return nil, errors.New("用户ID格式错误")
		}
		if where.Type != nil {
			if _, ok := user_log_type.Map()[*where.Type]; !ok {
				return nil, errors.New("类型参数不合法")
			}
		}
	}
	return s.repositoryUserLog.Select(where, queryPage)
}

// CreateLog 异步创建用户操作日志，成功表示消息已入队。
func (s *service) CreateLog(id int64, userId int64, requestId string, requestIp string, types int, title string, content string) error {
	if id <= 0 {
		return errors.New("日志ID格式错误")
	}
	if userId <= 0 {
		return errors.New("用户ID格式错误")
	}
	if _, ok := user_log_type.Map()[types]; !ok {
		return errors.New("类型参数不合法")
	}
	if title == "" {
		return errors.New("日志标题不能为空")
	}
	return s.logQueue.SendUnicastTask(&user_log.Task{
		Log: &model.UserLog{
			Id:         id,
			UserId:     userId,
			RequestId:  requestId,
			RequestIp:  requestIp,
			Type:       types,
			Title:      title,
			Content:    content,
			CreateTime: str.DateTime(time.Now()),
		},
	})
}

// consumeLogTask 消费用户操作日志消息。
func (s *service) consumeLogTask(task *user_log.Task) error {
	if task == nil || task.Log == nil {
		log.Println("用户操作日志消息为空")
		return nil
	}
	data := task.Log
	if _, ok := user_log_type.Map()[data.Type]; !ok || data.Id <= 0 || data.UserId <= 0 || data.Title == "" || time.Time(data.CreateTime).IsZero() {
		log.Println("用户操作日志消息参数不合法")
		return nil
	}
	// 与删除用户共用锁，避免日志写入和清理并发执行。
	lockKey := constant.LockUserOperate + strconv.FormatInt(data.UserId, 10)
	if !s.cache.Lock(lockKey, constant.LockTimeUserOperate) {
		return errors.New("系统繁忙,请稍候重试")
	}
	defer s.cache.UnLock(lockKey)
	data.Location = "未知"
	ipInfo, err := ip.Query(data.RequestIp)
	if err == nil && ipInfo != nil {
		data.Location = strings.ReplaceAll(fmt.Sprintf("%s %s %s %s", ipInfo.Country, ipInfo.Province, ipInfo.City, ipInfo.Isp), "  ", "")
	}
	return s.repositoryUserLog.Create(data)
}
