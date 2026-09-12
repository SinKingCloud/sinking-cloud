package permission

type UpdatePermission struct {
	Name *string
	Path *string
}

type SelectPermission struct {
	Name            *string
	Code            *string
	Path            *string
	CreateTimeStart *string
	CreateTimeEnd   *string
	UpdateTimeStart *string
	UpdateTimeEnd   *string
}
