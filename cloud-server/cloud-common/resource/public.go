package resource

import (
	"embed"
)

var (
	//go:embed ip/*
	Ip embed.FS
)
