package cron

import (
	"errors"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
)

// RegisterRpcTask 注册rpc类型定时任务
func (*Service) RegisterRpcTask(id int64, name string, spec string, service string, function string, concurrent bool, status bool) error {
	if id <= 0 || name == "" || spec == "" || service == "" || function == "" {
		return errors.New("参数错误")
	}
	var result bool
	err := rpc.InvokeFunc(constant.ServiceCron, "RegisterRpcTask", &RegisterRpcTaskParam{
		Id:         id,
		Name:       name,
		Spec:       spec,
		Service:    service,
		Function:   function,
		Concurrent: concurrent,
		Status:     status,
	}, &result)
	if err != nil || !result {
		return errors.Join(errors.New("注册任务失败"), err)
	}
	return nil
}

// RegisterScriptTask 注册script类型定时任务
func (*Service) RegisterScriptTask(id int64, name string, spec string, script string, concurrent bool, status bool) error {
	if id <= 0 || name == "" || spec == "" || script == "" {
		return errors.New("参数错误")
	}
	var result bool
	err := rpc.InvokeFunc(constant.ServiceCron, "RegisterScriptTask", &RegisterScriptTaskParam{
		Id:         id,
		Name:       name,
		Spec:       spec,
		Script:     script,
		Concurrent: concurrent,
		Status:     status,
	}, &result)
	if err != nil || !result {
		return errors.Join(errors.New("注册任务失败"), err)
	}
	return nil
}
