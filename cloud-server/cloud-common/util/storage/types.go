package storage

import "time"

// Provider 表示底层对象存储厂商类型。
type Provider string

const (
	ProviderAliyunOSS Provider = "aliyun_oss"
)

// Config 是统一的云对象存储配置。
type Config struct {
	Provider        Provider `json:"provider"`        // 指定底层存储厂商类型。
	Endpoint        string   `json:"endpoint"`        // 是厂商提供的接入节点。
	Region          string   `json:"region"`          // 是存储所在地域，部分厂商签名时会用到。
	Bucket          string   `json:"bucket"`          // 是对象存储空间名称。
	AccessKeyID     string   `json:"accessKeyId"`     // 是访问密钥 ID。
	AccessKeySecret string   `json:"accessKeySecret"` // 是访问密钥 Secret。
	SecurityToken   string   `json:"securityToken"`   // 是临时凭证场景下的安全令牌。
	Prefix          string   `json:"prefix"`          // 是统一拼接到对象 key 前的逻辑前缀。
	Domain          string   `json:"domain"`          // 是自定义访问域名，配置后优先用于生成对象访问地址。
	UseHTTPS        bool     `json:"useHttps"`        // 控制默认访问地址是否使用 HTTPS。
}

// AliyunOSSConfig 是阿里云 OSS 的便捷配置。
type AliyunOSSConfig struct {
	Endpoint        string `json:"endpoint"`        // 是阿里云 OSS 地域节点。
	Region          string `json:"region"`          // 是阿里云 OSS 所在地域。
	Bucket          string `json:"bucket"`          // 是阿里云 OSS Bucket 名称。
	AccessKeyID     string `json:"accessKeyId"`     // 是阿里云访问密钥 ID。
	AccessKeySecret string `json:"accessKeySecret"` // 是阿里云访问密钥 Secret。
	SecurityToken   string `json:"securityToken"`   // 是 STS 临时授权场景下的安全令牌。
	Prefix          string `json:"prefix"`          // 是统一拼接到对象 key 前的逻辑前缀。
	Domain          string `json:"domain"`          // 是对象访问时使用的自定义域名。
	UseHTTPS        bool   `json:"useHttps"`        // 控制生成访问地址时是否使用 HTTPS。
}

// ObjectInfo 描述一个存储对象的公共信息。
type ObjectInfo struct {
	Key          string            `json:"key"`                // 是对业务层暴露的逻辑对象路径。
	ProviderKey  string            `json:"providerKey"`        // 是底层厂商实际使用的对象路径。
	Size         int64             `json:"size"`               // 是对象大小，单位为字节。
	LastModified time.Time         `json:"lastModified"`       // 是对象最后修改时间。
	ETag         string            `json:"etag"`               // 是对象内容标识，通常可用于校验变化。
	ContentType  string            `json:"contentType"`        // 是对象的内容类型。
	URL          string            `json:"url"`                // 是根据当前配置生成的对象访问地址。
	Metadata     map[string]string `json:"metadata,omitempty"` // 是对象自定义元数据。
}

// ListRequest 描述对象列表查询请求。
type ListRequest struct {
	Prefix string `json:"prefix"` // 只返回指定逻辑前缀下的对象。
	Marker string `json:"marker"` // 是上一页返回的游标，用于继续翻页。
	Limit  int    `json:"limit"`  // 是单页最多返回的对象数量。
}

// ListResult 描述对象列表查询结果。
type ListResult struct {
	Objects    []ObjectInfo `json:"objects"`    // 是当前页返回的对象列表。
	NextMarker string       `json:"nextMarker"` // 是下一页查询所需的游标。
	HasMore    bool         `json:"hasMore"`    // 表示后续是否还有更多对象。
}

// CopyResult 描述复制对象操作结果。
type CopyResult struct {
	SourceKey      string    `json:"sourceKey"`      // 是复制源对象的逻辑路径。
	DestinationKey string    `json:"destinationKey"` // 是复制目标对象的逻辑路径。
	ETag           string    `json:"etag"`           // 是复制后目标对象的内容标识。
	CopiedAt       time.Time `json:"copiedAt"`       // 是复制完成时间。
}

// UploadState 描述一个上传任务的简化状态。
type UploadState struct {
	ID             string `json:"id"`                      // 是上传任务 ID，可继续用于续传或完成上传。
	Key            string `json:"key"`                     // 是当前上传任务对应的逻辑对象路径。
	Size           int64  `json:"size"`                    // 是声明的文件总大小，单位为字节。
	ChunkSize      int64  `json:"chunkSize"`               // 是当前任务使用的分片大小。
	UploadedParts  []int  `json:"uploadedParts,omitempty"` // 是已经成功上传的分片序号列表。
	CommittedBytes int64  `json:"committedBytes"`          // 是云端已确认写入的总字节数。
	Done           bool   `json:"done"`                    // 表示上传任务是否已经完成。
}
