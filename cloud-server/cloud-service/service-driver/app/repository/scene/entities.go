package scene

type SelectScene struct {
	Name            string
	ChannelId       string
	Status          string
	Remark          string
	CreateTimeStart string
	CreateTimeEnd   string
	UpdateTimeStart string
	UpdateTimeEnd   string
}

type UpdateScene struct {
	Name      interface{}
	MaxSize   interface{}
	AllowExt  interface{}
	ChannelId interface{}
	Remark    interface{}
	Status    interface{}
}
