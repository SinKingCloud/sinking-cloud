package job

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-spider/app/model"
	"strconv"
	"sync"
)

var lock bool

// Reset 重置消息下标
func (ControllerJob) Reset(c *sinking_server.Context) {
	if lock {
		c.Error("任务正在进行")
		return
	}
	go func() {
		lock = true
		wg := &sync.WaitGroup{}
		var jobs []model.Job
		util.Database.Db.Find(&jobs)
		tasks := make(chan model.Job)
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(group *sync.WaitGroup) {
				for task := range tasks {
					//查询系统最大记录
					var jobLog model.JobLog
					util.Database.Db.Where("`city_code`=? and `subject_code`=?", task.Code, task.Subject).Order("`index` desc").First(&jobLog)
					if jobLog.Id <= 0 {
						task.LastIndex = 0
					} else {
						task.LastIndex, _ = strconv.Atoi(jobLog.Index)
					}
					util.Database.Db.Updates(&task)
				}
				group.Done()
			}(wg)
		}
		for _, v := range jobs {
			tasks <- v
		}
		close(tasks)
		wg.Wait()
		lock = false
	}()
	c.Success("已下发重置任务")
}
