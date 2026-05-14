package menu

type SelectMenu struct {
	ParentId        string
	Code            string
	Name            string
	Icon            string
	Type            string
	Status          string
	CreateTimeStart string
	CreateTimeEnd   string
	UpdateTimeStart string
	UpdateTimeEnd   string
}

type UpdateMenu struct {
	ParentId interface{}
	Name     interface{}
	Icon     interface{}
	Sort     interface{}
	Status   interface{}
}
