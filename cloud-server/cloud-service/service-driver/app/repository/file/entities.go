package file

type SelectFile struct {
	ChannelId       *int64
	SceneId         *int64
	Ext             *string
	Hash            *string
	Status          *int
	Name            *string
	Path            *string
	CreateTimeStart *string
	CreateTimeEnd   *string
	UpdateTimeStart *string
	UpdateTimeEnd   *string
}

type UpdateFile struct {
	Status *int
}
