package rpc

import (
	"encoding/json"
	"errors"
	"service-cron/app/enum/task_concurrent"
	"service-cron/app/enum/task_status"
	"service-cron/app/enum/task_type"
	"service-cron/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/cron"
)

type taskService struct {
}

// RegisterRpc 注册rpc任务
func (*taskService) RegisterRpc(params json.RawMessage) (interface{}, error) {
	var param *cron.RegisterRpcTaskParam
	err := json.Unmarshal(params, &param)
	if err != nil {
		return false, err
	}
	if param == nil {
		return false, errors.New("参数错误")
	}
	if param.Id <= 0 || param.Name == "" || param.Spec == "" || param.Service == "" || param.Function == "" {
		return false, errors.New("参数错误")
	}
	var concurrent, status int
	if param.Concurrent {
		concurrent = task_concurrent.Enable
	} else {
		concurrent = task_concurrent.Disable
	}
	if param.Status {
		status = task_status.Running
	} else {
		status = task_status.Stop
	}
	err = service.Task.Register(param.Id, task_type.Rpc, param.Name, param.Spec, param.Service+"."+param.Function, concurrent, status)
	if err != nil {
		return false, err
	}
	return true, nil
}

// RegisterScript 注册script任务
func (*taskService) RegisterScript(params json.RawMessage) (interface{}, error) {
	var param *cron.RegisterScriptTaskParam
	err := json.Unmarshal(params, &param)
	if err != nil {
		return false, err
	}
	if param == nil {
		return false, errors.New("参数错误")
	}
	if param.Id <= 0 || param.Name == "" || param.Spec == "" || param.Script == "" {
		return false, errors.New("参数错误")
	}
	var concurrent, status int
	if param.Concurrent {
		concurrent = task_concurrent.Enable
	} else {
		concurrent = task_concurrent.Disable
	}
	if param.Status {
		status = task_status.Running
	} else {
		status = task_status.Stop
	}
	err = service.Task.Register(param.Id, task_type.Script, param.Name, param.Spec, param.Script, concurrent, status)
	if err != nil {
		return false, err
	}
	return true, nil
}
