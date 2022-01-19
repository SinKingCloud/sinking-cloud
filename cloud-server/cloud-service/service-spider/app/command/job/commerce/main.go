package commerce

import (
	"encoding/json"
	"fmt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-spider/app/command/queues/producer"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-spider/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-spider/app/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-spider/app/util/commerce"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-spider/app/util/origin/tianyancha"
	"log"
	"strconv"
	"time"
)

//Run 执行爬虫监测
func Run() {
	go run()
}

func run() {
	num := util.Conf.GetInt(constant.ApplicationSpiderThreadNum)
	log.Println("爬虫服务协程启动,数量：", num, ",初始最大连续错误次数", util.Conf.GetInt(constant.ApplicationSpiderMaxFailNum))
	for i := 1; i <= num; i++ {
		go spiderRun()
	}
}

func spiderRun() {
	for true {
		jobs := GetJob(1)
		for _, v := range jobs {
			failNum := 0
			tempJob := v
			//爬取内容
			for true {
				tempJob.LastIndex++
				code := commerce.GetBusinessCode(
					strconv.Itoa(v.Code)+
						strconv.Itoa(v.Subject)+
						fmt.Sprintf("%07d", tempJob.LastIndex), true)
				tempJob.LastRunTime = str.DateTime(time.Now())
				if len(code) != 15 {
					continue
				}
				data, err := tianyancha.Search(code)
				if err == nil && len(data.CompanyList) > 0 {
					failNum = 0
					//保存数据
					for _, v1 := range data.CompanyList {
						text, _ := json.Marshal(v1)
						info := model.JobLog{
							CityCode:    strconv.Itoa(v.Code),
							SubjectCode: strconv.Itoa(v.Subject),
							Index:       strconv.Itoa(tempJob.LastIndex),
							RegCode:     code,
							TycId:       strconv.Itoa(v1.Id),
							CreditCode:  v1.CreditCode,
							CompanyName: v1.Name,
							Content:     string(text),
						}
						producer.NewJobLogProducer().Build(&info).Send()
					}
				} else {
					failNum++
					maxNum := util.Conf.GetInt(constant.ApplicationSpiderMaxFailNum)
					if failNum == maxNum {
						//更新任务信息
						tempJob.LastIndex = tempJob.LastIndex - maxNum - 1
						if tempJob.LastIndex <= 0 {
							tempJob.LastIndex = 0
						}
						tempJob.Status = 1
					}
				}
				util.Database.Db.Updates(&tempJob)
				if tempJob.Status == 1 {
					break
				}
			}
		}
		if len(jobs) <= 0 {
			time.Sleep(1 * time.Second)
		}
	}
}
