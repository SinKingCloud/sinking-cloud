package dashboard

import (
	"encoding/json"
	"service-admin/app/constant"

	commonConstant "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
)

// GetServiceCount 获取服务数量
func (s *service) GetServiceCount() map[string]int {
	data := map[string]int{}
	if cacheData, exists := s.cache.Get(constant.CacheNameWithDashboardServiceCount); exists && cacheData != "" {
		if err := json.Unmarshal([]byte(cacheData), &data); err == nil {
			return data
		}
	}
	services := []string{
		commonConstant.ServiceGateway,
		commonConstant.ServiceAuth,
		commonConstant.ServiceSocket,
		commonConstant.ServiceCron,
		commonConstant.ServicePay,
		commonConstant.ServiceDriver,
		commonConstant.ServiceUser,
		commonConstant.ServiceAdmin,
	}
	for _, serviceName := range services {
		nodes, err := s.consul.GetServiceNodes(serviceName)
		if err != nil {
			data[serviceName] = 0
			continue
		}
		data[serviceName] = len(nodes)
	}
	if cacheData, err := json.Marshal(data); err == nil {
		s.cache.SetWithExpire(constant.CacheNameWithDashboardServiceCount, string(cacheData), constant.CacheTimeWithDashboardServiceCount)
	}
	return data
}
