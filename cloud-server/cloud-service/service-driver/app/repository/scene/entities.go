package scene

type SelectScene struct {
	Name            *string
	ChannelId       *int64
	Status          *int
	Login           *int
	Remark          *string
	CreateTimeStart *string
	CreateTimeEnd   *string
	UpdateTimeStart *string
	UpdateTimeEnd   *string
}

type UpdateScene struct {
	Name      *string
	MaxSize   *int64
	AllowExt  *string
	Login     *int
	ChannelId *int64
	Remark    *string
	Status    *int
}
