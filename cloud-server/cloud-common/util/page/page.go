package page

type Page struct {
	Total    int64       `json:"total"`     //记录总数
	Page     int         `json:"page"`      //分页页码
	PageSize int         `json:"page_size"` //分页容量
	Data     interface{} `json:"data"`      //记录详情
}

func NewPage(total int64, page int, pageSize int, data interface{}) *Page {
	return &Page{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Data:     data,
	}
}
