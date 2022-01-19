package job

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-spider/app/command/job/commerce"

// Run 执行异步任务
func Run() {
	//启动全量爬取任务
	commerce.Run()
}
