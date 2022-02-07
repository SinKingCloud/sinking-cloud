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
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-spider/app/util/origin"
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
	ReferProxyTimer(util.Conf.GetString(constant.SpiderHttpsProxy))
	num := util.Conf.GetInt(constant.SpiderTianYanChaThreadNum)
	log.Println("爬虫服务协程启动,数量：", num, ",初始最大连续错误次数", util.Conf.GetInt(constant.SpiderMaxFailNum))
	for i := 1; i <= num; i++ {
		go spiderRun()
	}
}

func spiderRun() {
	host := GetProxy()
	for true {
		jobs := GetJob(1)
		for _, v := range jobs {
			failNum := 0
			tempJob := v
			continueNum := 0
			//爬取内容
			for true {
				tempJob.LastIndex++
				code := commerce.GetBusinessCode(
					strconv.Itoa(v.Code)+
						strconv.Itoa(v.Subject)+
						fmt.Sprintf("%07d", tempJob.LastIndex), true)
				tempJob.LastRunTime = str.DateTime(time.Now())
				if len(code) != 15 {
					continueNum++
					continue
				}
				data, err := tianyancha.Search(code, host)
				if err == origin.ProxyExpire {
					host = GetProxy()
					data, err = tianyancha.Search(code, host)
				}
				if err == origin.Success && len(data.CompanyList) > 0 {
					failNum = 0
					continueNum = 0
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
					maxNum := util.Conf.GetInt(constant.SpiderMaxFailNum)
					if failNum == maxNum {
						//更新任务信息
						tempJob.LastIndex = tempJob.LastIndex - maxNum - continueNum - 1
						if tempJob.LastIndex <= 0 {
							tempJob.LastIndex = 0
						}
						tempJob.Status = 1
					}
				}
				if failNum == 0 || tempJob.Status == 1 {
					util.Database.Db.Model(&tempJob).Updates(map[string]interface{}{
						"last_index":    tempJob.LastIndex,
						"update_time":   tempJob.LastRunTime,
						"last_run_time": tempJob.LastRunTime,
						"status":        tempJob.Status,
					})
				}
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
