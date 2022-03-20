package logs

import "log"

type SystemLog struct {
}

func (s SystemLog) Println(data ...interface{}) {
	log.Println(data)
}
