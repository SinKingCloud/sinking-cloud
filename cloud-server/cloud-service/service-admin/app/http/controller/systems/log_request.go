package systems

import (
	"errors"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service/system_log"
)

// List 列表
func (ControllerLogRequest) List(c *sinking_server.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		OrderByField     string `form:"order_by_field" json:"order_by_field" default:"id" validate:"oneof=id request_id request_time" label:"排序字段"`
		OrderByType      string `form:"order_by_type" json:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		RequestId        string `form:"request_id" json:"request_id" default:"" validate:"omitempty,len=36" label:"请求ID"`
		RequestIp        string `form:"request_ip" json:"request_ip" default:"" validate:"omitempty,max=20" label:"请求IP"`
		RequestUrl       string `form:"request_url" json:"request_url" default:"" validate:"omitempty,max=200" label:"请求地址"`
		RequestMethod    string `form:"request_method" json:"request_method" default:"" validate:"omitempty,oneof=GET POST PUT UPDATE" label:"请求类型"`
		RequestHeader    string `form:"request_header" json:"request_header" default:"" validate:"omitempty" label:"请求头"`
		RequestBody      string `form:"request_body" json:"request_body" default:"" validate:"omitempty" label:"请求体"`
		ResponseCode     string `form:"response_code" json:"response_code" default:"" validate:"omitempty,numeric" label:"响应码"`
		ResponseHeader   string `form:"response_header" json:"response_header" default:"" validate:"omitempty" label:"响应头"`
		ResponseBody     string `form:"response_body" json:"response_body" default:"" validate:"omitempty" label:"响应体"`
		ResponseTimeMax  string `form:"response_time_max" json:"response_time_max" default:"" validate:"omitempty,numeric" label:"最大响应时间"`
		ResponseTimeMin  string `form:"response_time_min" json:"response_time_min" default:"" validate:"omitempty,numeric" label:"最小响应时间"`
		GroupName        string `form:"group_name" json:"group_name" default:"" validate:"omitempty" label:"服务组名"`
		AppName          string `form:"app_name" json:"app_name" default:"" validate:"omitempty" label:"服务名"`
		AppAddr          string `form:"app_addr" json:"app_addr" default:"" validate:"omitempty" label:"服务地址"`
		RequestTimeStart string `form:"request_time_start" json:"request_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		RequestTimeEnd   string `form:"request_time_end" json:"request_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	where := make(map[string]string)
	if form.RequestId != "" {
		where["request_id"] = form.RequestId
	}
	if form.RequestIp != "" {
		where["request_ip"] = form.RequestIp
	}
	if form.RequestUrl != "" {
		where["request_url"] = form.RequestUrl
	}
	if form.RequestMethod != "" {
		where["request_method"] = form.RequestMethod
	}
	if form.RequestHeader != "" {
		where["request_header"] = form.RequestHeader
	}
	if form.RequestBody != "" {
		where["request_body"] = form.RequestBody
	}
	if form.ResponseCode != "" {
		where["response_code"] = form.ResponseCode
	}
	if form.ResponseHeader != "" {
		where["response_header"] = form.ResponseHeader
	}
	if form.ResponseBody != "" {
		where["response_body"] = form.ResponseBody
	}
	if form.ResponseTimeMax != "" {
		where["response_time_max"] = form.ResponseTimeMax
	}
	if form.ResponseTimeMin != "" {
		where["response_time_min"] = form.ResponseTimeMin
	}
	if form.GroupName != "" {
		where["group_name"] = form.GroupName
	}
	if form.AppName != "" {
		where["app_name"] = form.AppName
	}
	if form.AppAddr != "" {
		where["app_addr"] = form.AppAddr
	}
	if form.RequestTimeStart != "" {
		where["request_time_start"] = form.RequestTimeStart
	}
	if form.RequestTimeEnd != "" {
		where["request_time_end"] = form.RequestTimeEnd
	}
	data, total, err := service.SystemLogRequest.Select(where, form.OrderByField, form.OrderByType, pageInfo.Page, pageInfo.PageSize)
	if err == nil {
		service.SystemLog.CreateByQueue(c, system_log.EventShow, "查看日志", "查看请求日志")
		c.SuccessWithData("获取请求数据成功", page.NewPage(total, pageInfo.Page, pageInfo.PageSize, data))
		return
	}
	c.Error("获取数据失败")
}

// Info 请求记录详情
func (ControllerLogRequest) Info(c *sinking_server.Context) {
	type Form struct {
		Id        string `form:"id" json:"id" default:"" validate:"omitempty,max=36" label:"记录ID"`
		RequestId string `form:"request_id" json:"request_id" default:"" validate:"omitempty,len=36" label:"请求ID"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	var err error
	var data *model.SystemLogRequest
	if form.Id != "" {
		data, err = service.SystemLogRequest.FindById(form.Id)
	} else if form.RequestId != "" {
		data, err = service.SystemLogRequest.FindByRequestId(form.RequestId)
	} else {
		err = errors.New("field error")
	}
	if err == nil && data != nil {
		service.SystemLog.CreateByQueue(c, system_log.EventShow, "查看日志", "查看详细请求日志")
		c.SuccessWithData("获取详情数据成功", data)
	} else {
		c.Error("获取详情数据失败")
	}
}

// Delete 删除请求记录
func (ControllerLogRequest) Delete(c *sinking_server.Context) {
	type Form struct {
		RequestTimeStart string `form:"request_time_start" json:"request_time_start" default:"" validate:"required" label:"创建起始时间"`
		RequestTimeEnd   string `form:"request_time_end" json:"request_time_end" default:"" validate:"required" label:"创建结束时间"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	err := service.SystemLogRequest.DeleteByRequestTime(form.RequestTimeStart, form.RequestTimeEnd)
	if err == nil {
		service.SystemLog.CreateByQueue(c, system_log.EventDelete, "清理日志", "清理请求日志")
		c.Success("删除数据成功")
	} else {
		c.Error("删除数据失败")
	}
}
