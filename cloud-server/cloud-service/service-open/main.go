package main

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/bootstrap"
)

func main() {
	//servers := make(chan interface{}) //服务器列表
	//g := &sync.WaitGroup{}
	//g.Add(1)
	//go func(g *sync.WaitGroup) {
	//	(&job.Task{Producer: func(channel chan interface{}) {
	//		servers = channel
	//		g.Add(-1)
	//	}, Consumer: func(param interface{}) {
	//		fmt.Println(param)
	//	}, Thread: 50}).Run()
	//	fmt.Println("通道关闭")
	//}(g)
	//g.Wait()
	////go func() {
	////	for {
	////		sinking_web.Try(func() {
	////			servers <- time.Now().Unix()
	////			time.Sleep(1 * time.Second)
	////		}, func(i interface{}) {
	////			util.Log.Println("通道关闭--", i)
	////		})
	////
	////	}
	////}()
	//time.Sleep(10 * time.Second)
	//close(servers)
	//close(servers)
	//time.Sleep(999999 * time.Second)
	//m := md5.Sum([]byte("Qq@1178710004..."))
	//cli := oicq.GetOicqInstance(1178710004, m, nil, nil)
	//fmt.Println(cli.ActiveProject([]oicq.Job{oicq.Pc}, nil))
	//fmt.Println("success")
	//time.Sleep(time.Second * 9999999)
	//os.Exit(0)
	bootstrap.Load()
	app.Run()
}
