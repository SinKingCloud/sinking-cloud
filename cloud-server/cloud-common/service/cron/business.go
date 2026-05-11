package cron

import (
	"errors"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
)

// RegisterRpcTask 注册rpc类型定时任务
func (s *Service) RegisterRpcTask(id int64, name string, spec string, service string, function string, concurrent bool, status bool) error {
	if id <= 0 || name == "" || spec == "" || service == "" || function == "" {
		return errors.New("参数错误")
	}
	param := &RegisterRpcTaskParam{
		Id:         id,
		Name:       name,
		Spec:       spec,
		Service:    service,
		Function:   function,
		Concurrent: concurrent,
		Status:     status,
	}
	return s.registerWithRetry("RegisterRpcTask", param)
}

// RegisterScriptTask 注册script类型定时任务
func (s *Service) RegisterScriptTask(id int64, name string, spec string, script string, concurrent bool, status bool) error {
	if id <= 0 || name == "" || spec == "" || script == "" {
		return errors.New("参数错误")
	}
	param := &RegisterScriptTaskParam{
		Id:         id,
		Name:       name,
		Spec:       spec,
		Script:     script,
		Concurrent: concurrent,
		Status:     status,
	}
	return s.registerWithRetry("RegisterScriptTask", param)
}

// registerWithRetry 注册任务
func (*Service) registerWithRetry(function string, param interface{}) error {
	var lastErr error
	for i := 0; ; i++ {
		var result bool
		err := rpc.InvokeFunc(constant.ServiceCron, function, param, &result)
		if err == nil && result {
			return nil
		}
		lastErr = err
		if i >= 10 {
			break
		}
		time.Sleep(5 * time.Second)
	}
	if lastErr == nil {
		return errors.New("注册任务失败")
	}
	return errors.Join(errors.New("注册任务失败"), lastErr)
}
