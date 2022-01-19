package index

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
)

func (ControllerIndex) Index(c *sinking_server.Context) {
	//producer.NewJobLogProducer().Build(&model.JobLog{
	//	CityCode:    time.Now().String(),
	//	SubjectCode: time.Now().String(),
	//	Index:       time.Now().String(),
	//	RegCode:     time.Now().String(),
	//	CompanyName: time.Now().String(),
	//	Status:      0,
	//}).Send()
	c.Success("请求成功")
}
