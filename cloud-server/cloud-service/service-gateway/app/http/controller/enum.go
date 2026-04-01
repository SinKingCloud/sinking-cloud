package controller

import (
	"strings"
	"sync"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Enum 枚举数据
func Enum(c *context.Context) {
	var form struct {
		Name string `json:"name" default:"" validate:"required" label:"枚举名称"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	query := make(map[string][]string, 10)
	for _, v := range strings.Split(strings.ToLower(form.Name), ",") {
		parts := strings.SplitN(v, ".", 2)
		if len(parts) == 2 {
			query[parts[0]] = append(query[parts[0]], parts[1])
		}
	}
	if len(query) <= 0 {
		c.Error("枚举数据不存在")
		return
	}
	const maxGoroutine = 10
	queue := make(chan struct{}, maxGoroutine)
	data := make(map[string]interface{}, 10)
	var mu sync.Mutex
	var wg sync.WaitGroup
	for srvName, enumNames := range query {
		queue <- struct{}{}
		wg.Add(1)
		go func(srv string, names []string) {
			defer func() {
				<-queue
				wg.Done()
			}()
			temp := service.Enum.Data(srv, names)
			if temp == nil {
				return
			}
			mu.Lock()
			for k1, v1 := range temp {
				data[srv+"."+k1] = v1
			}
			mu.Unlock()
		}(srvName, enumNames)
	}
	wg.Wait()
	c.SuccessWithData("获取数据成功", data)
}
