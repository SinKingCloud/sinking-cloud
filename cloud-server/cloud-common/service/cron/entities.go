package cron

// RegisterRpcTaskParam 注册rpc类型任务参数
type RegisterRpcTaskParam struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Spec       string `json:"spec"`
	Service    string `json:"service"`
	Function   string `json:"function"`
	Concurrent bool   `json:"concurrent"`
	Status     bool   `json:"status"`
}

// RegisterScriptTaskParam 注册脚本类型任务参数
type RegisterScriptTaskParam struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Spec       string `json:"spec"`
	Script     string `json:"script"`
	Concurrent bool   `json:"concurrent"`
	Status     bool   `json:"status"`
}
