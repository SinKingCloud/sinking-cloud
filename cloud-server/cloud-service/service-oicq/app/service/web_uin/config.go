package web_uin

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/service/app_config"
	"strconv"
	"sync"
	"time"
)

var (
	conf     = make(map[string]string)
	confLock = &sync.RWMutex{}
	confOnce = &sync.Once{}
)

func RefreshAppConfig(clearCache bool) {
	confLock.Lock()
	defer confLock.Unlock()
	conf = app_config.GetIns().GetAppConfigs(constant.AppUinId, clearCache)
}

func GetAppConfigString(key string) string {
	confOnce.Do(func() {
		RefreshAppConfig(true)
		go func() {
			for {
				RefreshAppConfig(false)
				time.Sleep(60 * time.Second)
			}
		}()
	})
	confLock.RLock()
	defer confLock.RUnlock()
	return conf[key]
}

func GetAppConfigInt(key string) int {
	n, e := strconv.Atoi(GetAppConfigString(key))
	if e != nil {
		n = 0
	}
	return n
}

func GetAppConfigFloat64(key string) float64 {
	n, e := strconv.ParseFloat(GetAppConfigString(key), 64)
	if e != nil {
		n = 0
	}
	return n
}
