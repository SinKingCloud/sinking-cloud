package proxy

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_proxy/proxy"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_proxy/proxy_area"
	"gorm.io/gorm"
	"strconv"
)

type ControllerProxy struct {
}

// Extract 提取代理
func (ControllerProxy) Extract(c *sinking_server.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		OrderByField string `form:"order_by_field" json:"order_by_field" default:"id" validate:"oneof=id extract_num" label:"排序字段"`
		OrderByType  string `form:"order_by_type" json:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		Province     string `json:"province" validate:"omitempty,max=50" label:"归属省"`
		City         string `json:"city" validate:"omitempty,max=50" label:"归属市"`
		Type         string `json:"type" default:"socks5" validate:"omitempty,max=30" label:"代理类型"`
		Echo         string `json:"echo" default:"json" validate:"oneof=txt json" label:"输出类型"`
	}
	var message string
	var data []*model.AppProxy
	form := &Form{}
	if ok, msg := c.ValidatorGet(form); !ok {
		message = msg
	} else {
		where := make(map[string]interface{})
		if form.Province != "" {
			where["province"] = form.Province
		}
		if form.City != "" {
			where["city"] = form.City
		}
		area, err := proxy_area.GetIns().Find(where)
		if err != nil || area == nil {
			message = "获取地区信息失败"
		} else {
			list, _, e := proxy.GetIns().Select(map[string]string{"area_id": strconv.Itoa(area.Id)}, form.OrderByField, form.OrderByType, pageInfo.Page, pageInfo.PageSize)
			if e != nil {
				message = "提取代理列表失败"
			} else {
				//更改提取次数
				go func() {
					var ids []int
					for _, v := range list {
						ids = append(ids, v.Id)
					}
					_ = proxy.GetIns().UpdateByIds(ids, map[string]interface{}{
						"extract_num": gorm.Expr("`extract_num`+1"),
					})
				}()
				message = "提取代理列表成功"
				data = list
			}
		}
	}
	if data != nil {
		if form.Echo == "txt" {
			txt := ""
			for _, v := range data {
				txt += v.Ip + ":" + strconv.Itoa(v.Port) + " " + v.Auth + "\n"
			}
			c.String(200, txt)
		} else {
			c.SuccessWithData(message, data)
		}
	} else {
		if form.Echo == "txt" {
			c.String(500, message)
		} else {
			c.Error(message)
		}
	}
}
