package systems

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service/system_log"
)

// List 获取数据
func (ControllerFile) List(c *sinking_server.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		OrderByField    string `form:"order_by_field" json:"order_by_field" default:"id" validate:"oneof=id hash ext create_time update_time" label:"排序字段"`
		OrderByType     string `form:"order_by_type" json:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		Type            string `form:"type" json:"type" default:"" validate:"omitempty,numeric" label:"所属类别"`
		Hash            string `form:"hash" json:"hash" default:"" validate:"omitempty,len=32" label:"文件指纹"`
		Ext             string `form:"ext" json:"ext" default:"" validate:"omitempty,max=10" label:"文件类型"`
		CreateTimeStart string `form:"create_time_start" json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd   string `form:"create_time_end" json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
		UpdateTimeStart string `form:"update_time_start" json:"update_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新起始时间"`
		UpdateTimeEnd   string `form:"update_time_end" json:"update_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新结束时间"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	where := make(map[string]string)
	if form.Type != "" {
		where["type"] = form.Type
	}
	if form.Hash != "" {
		where["hash"] = form.Hash
	}
	if form.Ext != "" {
		where["ext"] = form.Ext
	}
	if form.CreateTimeStart != "" {
		where["create_time_start"] = form.CreateTimeStart
	}
	if form.CreateTimeEnd != "" {
		where["create_time_end"] = form.CreateTimeEnd
	}
	if form.UpdateTimeStart != "" {
		where["update_time_start"] = form.UpdateTimeStart
	}
	if form.UpdateTimeEnd != "" {
		where["update_time_end"] = form.UpdateTimeEnd
	}
	data, total, err := service.SystemFile.Select(where, form.OrderByField, form.OrderByType, pageInfo.Page, pageInfo.PageSize)
	if err != nil {
		c.Error("获取附件数据失败")
	} else {
		system_log.GetIns().CreateByQueue(c, system_log.EventShow, "查看附件", "查看附件列表")
		c.SuccessWithData("获取附件数据成功", page.NewPage(total, pageInfo.Page, pageInfo.PageSize, data))
	}
}

// Delete 删除信息
func (ControllerFile) Delete(c *sinking_server.Context) {
	type Form struct {
		Ids []int `form:"ids" json:"ids" default:"" validate:"required,min=1,max=1000,unique" label:"ID列表"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	paths := service.SystemFile.SelectPathByIds(form.Ids)
	err := service.SystemFile.DeleteByIds(form.Ids)
	if err == nil {
		go service.SystemFile.DeleteFiles(paths)
		service.SystemLog.CreateByQueue(c, system_log.EventDelete, "删除附件", "删除系统附件")
		c.Success("删除数据成功")
	} else {
		c.Error("删除数据失败")
	}
}
