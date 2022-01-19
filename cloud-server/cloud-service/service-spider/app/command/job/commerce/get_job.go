package commerce

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-spider/app/model"
	"time"
)

// GetJob 获取任务
func GetJob(num int) []model.Job {
	var jobs []model.Job
	//使用分布式锁获取任务列表
	LockWaitRun(func() {
		util.Database.Db.
			Where("status=0").
			Order("last_run_time asc,id asc").
			Limit(num).
			Find(&jobs)
		for _, v := range jobs {
			util.Database.Db.Model(&model.Job{}).Where("id=?", v.Id).Update("last_run_time", time.Now().Format("2006-01-02 15:04:05"))
		}
	})
	return jobs
}
