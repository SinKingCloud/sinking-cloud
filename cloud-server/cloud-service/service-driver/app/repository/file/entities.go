package file

type SelectFile struct {
	ChannelId       string
	SceneId         string
	Ext             string
	Hash            string
	Status          string
	Name            string
	Path            string
	CreateTimeStart string
	CreateTimeEnd   string
	UpdateTimeStart string
	UpdateTimeEnd   string
}

type UpdateFile struct {
	Status interface{}
}
