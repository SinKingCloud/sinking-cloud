package commerce

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-spider/app/constant"
	"time"
)

func LockWaitRun(fun func()) {
	for true {
		res := util.Cache.Lock(constant.SpiderLock, constant.SpiderLockMaxTime)
		if res {
			break
		} else {
			time.Sleep(1000 * time.Millisecond)
		}
	}
	fun()
	util.Cache.UnLock(constant.SpiderLock)
}
