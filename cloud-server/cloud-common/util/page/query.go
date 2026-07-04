package page

// Query 分页和排序查询参数。
type Query struct {
	Page         int
	PageSize     int
	OrderByField string
	OrderByType  string
}
