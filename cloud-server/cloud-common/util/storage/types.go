package storage

import "time"

// Provider 表示底层对象存储厂商类型。
type Provider string

const (
	ProviderAliyunOSS Provider = "aliyun_oss"
)

// Config 是统一的云对象存储配置。
type Config struct {
	Provider        Provider `json:"provider"`
	Endpoint        string   `json:"endpoint"`
	Region          string   `json:"region"`
	Bucket          string   `json:"bucket"`
	AccessKeyID     string   `json:"accessKeyId"`
	AccessKeySecret string   `json:"accessKeySecret"`
	SecurityToken   string   `json:"securityToken"`
	Prefix          string   `json:"prefix"`
	Domain          string   `json:"domain"`
	UseHTTPS        bool     `json:"useHttps"`
}

// AliyunOSSConfig 是阿里云 OSS 的便捷配置。
type AliyunOSSConfig struct {
	Endpoint        string `json:"endpoint"`
	Region          string `json:"region"`
	Bucket          string `json:"bucket"`
	AccessKeyID     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	SecurityToken   string `json:"securityToken"`
	Prefix          string `json:"prefix"`
	Domain          string `json:"domain"`
	UseHTTPS        bool   `json:"useHttps"`
}

// ObjectInfo 描述一个存储对象的公共信息。
type ObjectInfo struct {
	Key          string            `json:"key"`
	ProviderKey  string            `json:"providerKey"`
	Size         int64             `json:"size"`
	LastModified time.Time         `json:"lastModified"`
	ETag         string            `json:"etag"`
	ContentType  string            `json:"contentType"`
	URL          string            `json:"url"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// ListRequest 描述对象列表查询请求。
type ListRequest struct {
	Prefix string `json:"prefix"`
	Marker string `json:"marker"`
	Limit  int    `json:"limit"`
}

// ListResult 描述对象列表查询结果。
type ListResult struct {
	Objects    []ObjectInfo `json:"objects"`
	NextMarker string       `json:"nextMarker"`
	HasMore    bool         `json:"hasMore"`
}

// CopyResult 描述复制对象操作结果。
type CopyResult struct {
	SourceKey      string    `json:"sourceKey"`
	DestinationKey string    `json:"destinationKey"`
	ETag           string    `json:"etag"`
	CopiedAt       time.Time `json:"copiedAt"`
}

// UploadState 描述一个上传任务的简化状态。
type UploadState struct {
	ID             string `json:"id"`
	Key            string `json:"key"`
	Size           int64  `json:"size"`
	ChunkSize      int64  `json:"chunkSize"`
	UploadedParts  []int  `json:"uploadedParts,omitempty"`
	CommittedBytes int64  `json:"committedBytes"`
	Done           bool   `json:"done"`
}
