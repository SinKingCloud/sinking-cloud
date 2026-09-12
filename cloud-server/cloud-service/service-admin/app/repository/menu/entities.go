package menu

type SelectMenu struct {
	ParentId        *int64
	Code            *string
	Name            *string
	Icon            *string
	Type            *int
	Status          *int
	CreateTimeStart *string
	CreateTimeEnd   *string
	UpdateTimeStart *string
	UpdateTimeEnd   *string
}

type UpdateMenu struct {
	ParentId *int64
	Name     *string
	Icon     *string
	Sort     *int
	Status   *int
}
