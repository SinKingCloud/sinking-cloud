package role

type SelectRole struct {
	Name            string
	Status          string
	CreateTimeStart string
	CreateTimeEnd   string
	UpdateTimeStart string
	UpdateTimeEnd   string
}

// UpdateRole 角色更新
type UpdateRole struct {
	Name   interface{}
	Status interface{}
}
