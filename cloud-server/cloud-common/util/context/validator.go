package context

import (
	"strings"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/validator"
)

// Validator 结构体验证
func (c *Context) Validator(data interface{}) (bool, string) {
	return validator.Check(data)
}

// ValidatorAll 参数绑定
func (c *Context) ValidatorAll(data interface{}) (bool, string) {
	if c.BindAll(data) != nil {
		return false, "参数绑定失败"
	}
	return validator.Check(data)
}

// ValidatorJson json参数验证
func (c *Context) ValidatorJson(data interface{}) (bool, string) {
	if c.BindJson(data) != nil {
		return false, "json参数绑定失败"
	}
	return validator.Check(data)
}

// ValidatorPost post参数验证
func (c *Context) ValidatorPost(data interface{}) (bool, string) {
	if c.BindForm(data) != nil {
		return false, "post参数绑定失败"
	}
	return validator.Check(data)
}

// ValidatorGet get参数验证
func (c *Context) ValidatorGet(data interface{}) (bool, string) {
	if c.BindQuery(data) != nil {
		return false, "get参数绑定失败"
	}
	return validator.Check(data)
}

// ValidatorPath path参数验证
func (c *Context) ValidatorPath(data interface{}) (bool, string) {
	if c.BindParam(data) != nil {
		return false, "path参数绑定失败"
	}
	return validator.Check(data)
}

// ValidatePage 验证分页和排序参数，并返回分页查询对象。
// 可选参数按类型匹配：int 参数依次表示 defaultPage、defaultPageSize；string 参数依次表示 defaultField、defaultType、allowField。
// 例如 ValidatePage("id", "desc", "id,create_time") 表示默认按 id 倒序，并且只允许 id、create_time 排序。
// defaultField 和 defaultType 默认为空；allowField 为空字符串或 "*" 时不限制排序字段；请求中的排序字段不在 allowField 内时，会回退到 defaultField。
func (c *Context) ValidatePage(args ...interface{}) *page.Query {
	defaultPage := 1
	defaultPageSize := 20
	defaultField := ""
	defaultType := ""
	allowField := defaultField
	intIndex := 0
	stringIndex := 0
	for _, arg := range args {
		switch value := arg.(type) {
		case int:
			if intIndex == 0 {
				defaultPage = value
			}
			if intIndex == 1 {
				defaultPageSize = value
			}
			intIndex++
		case string:
			if value == "" {
				continue
			}
			if stringIndex == 0 {
				defaultField = value
				allowField = value
			}
			if stringIndex == 1 {
				defaultType = strings.ToLower(value)
			}
			if stringIndex == 2 {
				allowField = value
			}
			stringIndex++
		}
	}
	if defaultPage <= 0 {
		defaultPage = 1
	}
	if defaultPageSize <= 0 || defaultPageSize > 1000 {
		defaultPageSize = 20
	}
	if defaultType != "" && defaultType != "desc" && defaultType != "asc" {
		defaultType = "desc"
	}
	var pageInfo struct {
		Page     int    `json:"page" label:"页码编号"`
		PageSize int    `json:"page_size" label:"每页数量"`
		Field    string `json:"order_by_field" label:"排序字段"`
		Type     string `json:"order_by_type" label:"排序类型"`
	}
	if c.BindAll(&pageInfo) != nil {
		return &page.Query{
			Page:         defaultPage,
			PageSize:     defaultPageSize,
			OrderByField: defaultField,
			OrderByType:  defaultType,
		}
	}
	query := &page.Query{
		Page:         pageInfo.Page,
		PageSize:     pageInfo.PageSize,
		OrderByField: pageInfo.Field,
		OrderByType:  strings.ToLower(pageInfo.Type),
	}
	if query.Page <= 0 {
		query.Page = defaultPage
	}
	if query.PageSize <= 0 || query.PageSize > 1000 {
		query.PageSize = defaultPageSize
	}
	if query.OrderByField == "" {
		query.OrderByField = defaultField
	}
	if query.OrderByType == "" {
		query.OrderByType = defaultType
	}
	if query.OrderByType != "" && query.OrderByType != "desc" && query.OrderByType != "asc" {
		query.OrderByType = defaultType
	}
	if allowField == "" || allowField == "*" {
		return query
	}
	fields := strings.Fields(strings.ReplaceAll(allowField, ",", " "))
	for _, v := range fields {
		if v == query.OrderByField {
			return query
		}
	}
	query.OrderByField = defaultField
	return query
}
