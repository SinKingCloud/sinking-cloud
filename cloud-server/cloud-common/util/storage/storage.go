package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"
)

// Provider 表示底层对象存储厂商或协议族。
type Provider string

const (
	ProviderAliyunOSS  Provider = "aliyun_oss"
	ProviderTencentCOS Provider = "tencent_cos"
	ProviderAmazonS3   Provider = "amazon_s3"
	ProviderMinIO      Provider = "minio"
	ProviderQiniuKodo  Provider = "qiniu_kodo"
	ProviderHuaweiOBS  Provider = "huawei_obs"
	ProviderGoogleGCS  Provider = "google_cloud_storage"
	ProviderAzureBlob  Provider = "azure_blob"
)

// UploadProtocol 表示大文件上传底层采用的协议模型。
type UploadProtocol string

const (
	UploadProtocolMultipart UploadProtocol = "multipart"
	UploadProtocolBlock     UploadProtocol = "block"
	UploadProtocolResumable UploadProtocol = "resumable"
)

// ChunkPolicy 描述当前驱动推荐的分块策略。
type ChunkPolicy struct {
	DefaultSize  int64 `json:"defaultSize"`
	MinSize      int64 `json:"minSize"`
	SizeMultiple int64 `json:"sizeMultiple"`
	MaxChunks    int   `json:"maxChunks"`
}

// Capabilities 描述当前驱动支持的能力集合。
type Capabilities struct {
	SmallUpload       bool             `json:"smallUpload"`
	LargeUpload       bool             `json:"largeUpload"`
	ResumeUpload      bool             `json:"resumeUpload"`
	UploadStatusQuery bool             `json:"uploadStatusQuery"`
	SignedURL         bool             `json:"signedUrl"`
	CopyObject        bool             `json:"copyObject"`
	BatchDelete       bool             `json:"batchDelete"`
	Protocols         []UploadProtocol `json:"protocols,omitempty"`
	ChunkPolicy       ChunkPolicy      `json:"chunkPolicy,omitempty"`
}

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

// UploadRequest 描述一次简单上传请求。
type UploadRequest struct {
	Key                string            `json:"key"`
	Body               io.Reader         `json:"-"`
	Size               int64             `json:"size"`
	ContentType        string            `json:"contentType"`
	CacheControl       string            `json:"cacheControl"`
	ContentDisposition string            `json:"contentDisposition"`
	Metadata           map[string]string `json:"metadata,omitempty"`
}

// LargeUploadRequest 描述一次由 SDK 托管的大文件上传请求。
// 该接口适合“服务端已经拿到整个文件流”的场景，底层驱动会选择最合适的分段协议。
type LargeUploadRequest struct {
	UploadRequest
	ChunkSize int64 `json:"chunkSize"`
}

// StartUploadSessionRequest 描述初始化上传会话的请求。
type StartUploadSessionRequest struct {
	Key                string            `json:"key"`
	Size               int64             `json:"size"`
	ChunkSize          int64             `json:"chunkSize"`
	ContentType        string            `json:"contentType"`
	CacheControl       string            `json:"cacheControl"`
	ContentDisposition string            `json:"contentDisposition"`
	Metadata           map[string]string `json:"metadata,omitempty"`
}

// UploadSession 是可序列化的上传会话信息。
// 它既适用于 S3/OSS/COS/OBS 这类 multipart 模型，也适用于 GCS resumable、Azure block blob 这类不同协议。
type UploadSession struct {
	Provider           Provider          `json:"provider"`
	Protocol           UploadProtocol    `json:"protocol"`
	Bucket             string            `json:"bucket"`
	Key                string            `json:"key"`
	ProviderKey        string            `json:"providerKey"`
	SessionID          string            `json:"sessionId"`
	Size               int64             `json:"size"`
	ChunkSize          int64             `json:"chunkSize"`
	ContentType        string            `json:"contentType"`
	CacheControl       string            `json:"cacheControl"`
	ContentDisposition string            `json:"contentDisposition"`
	Metadata           map[string]string `json:"metadata,omitempty"`
	Opaque             map[string]string `json:"opaque,omitempty"`
	CreatedAt          time.Time         `json:"createdAt"`
}

// UploadChunkRequest 描述上传单个数据块的请求。
// Index 建议从 1 开始；Offset 用于需要按字节偏移续传的协议；ChunkID 用于 Azure Block Blob 等需要显式块 ID 的场景。
type UploadChunkRequest struct {
	Session  UploadSession `json:"session"`
	Index    int           `json:"index"`
	Offset   int64         `json:"offset"`
	Size     int64         `json:"size"`
	ChunkID  string        `json:"chunkId"`
	IsFinal  bool          `json:"isFinal"`
	Body     io.Reader     `json:"-"`
	Checksum string        `json:"checksum,omitempty"`
}

// UploadedChunk 描述一个已经上传成功的数据块。
type UploadedChunk struct {
	Index        int               `json:"index"`
	Offset       int64             `json:"offset"`
	Size         int64             `json:"size"`
	ChunkID      string            `json:"chunkId,omitempty"`
	ETag         string            `json:"etag,omitempty"`
	Checksum     string            `json:"checksum,omitempty"`
	LastModified time.Time         `json:"lastModified"`
	Opaque       map[string]string `json:"opaque,omitempty"`
}

// UploadStatus 描述上传会话当前状态。
// 对于 multipart / block 协议，UploadedChunks 是已上传分块列表。
// 对于 resumable 协议，UploadedChunks 可以为空，CommittedBytes 用于表达服务端已持久化的字节范围。
type UploadStatus struct {
	Session        UploadSession   `json:"session"`
	UploadedChunks []UploadedChunk `json:"uploadedChunks,omitempty"`
	CommittedBytes int64           `json:"committedBytes"`
	Finalized      bool            `json:"finalized"`
}

// CompleteUploadRequest 描述完成上传会话的请求。
// 对于需要显式完成的协议，应传入已经上传成功的数据块列表。
// 对于 GCS 这类最终块上传即可完成的协议，Chunks 可以为空，驱动自行处理。
type CompleteUploadRequest struct {
	Session UploadSession   `json:"session"`
	Chunks  []UploadedChunk `json:"chunks,omitempty"`
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

// Client 是统一的云对象存储客户端。
type Client struct {
	driver driver
}

// New 根据统一配置创建云对象存储客户端。
func New(config Config) (*Client, error) {
	driver, err := newDriver(config)
	if err != nil {
		return nil, err
	}

	return &Client{driver: driver}, nil
}

// NewAliyunOSS 使用阿里云 OSS 配置创建客户端。
func NewAliyunOSS(config AliyunOSSConfig) (*Client, error) {
	return New(Config{
		Provider:        ProviderAliyunOSS,
		Endpoint:        config.Endpoint,
		Region:          config.Region,
		Bucket:          config.Bucket,
		AccessKeyID:     config.AccessKeyID,
		AccessKeySecret: config.AccessKeySecret,
		SecurityToken:   config.SecurityToken,
		Prefix:          config.Prefix,
		Domain:          config.Domain,
		UseHTTPS:        config.UseHTTPS,
	})
}

// Provider 返回当前客户端绑定的云厂商类型。
func (c *Client) Provider() Provider {
	return c.driver.provider()
}

// Capabilities 返回当前驱动支持的能力集合。
func (c *Client) Capabilities() Capabilities {
	return c.driver.capabilities()
}

// Upload 执行小文件或普通对象上传。
func (c *Client) Upload(ctx context.Context, request UploadRequest) (*ObjectInfo, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if request.Body == nil {
		return nil, ErrUploadBodyRequired
	}

	return c.driver.upload(ctx, request)
}

// UploadLarge 执行由 SDK 托管的大文件上传。
// 底层会走 StartUploadSession -> UploadChunk -> CompleteUpload 这条通用链路。
func (c *Client) UploadLarge(ctx context.Context, request LargeUploadRequest) (_ *ObjectInfo, err error) {
	if err = checkContext(ctx); err != nil {
		return nil, err
	}
	if request.Body == nil {
		return nil, ErrUploadBodyRequired
	}

	capabilities := c.Capabilities()
	if !capabilities.LargeUpload {
		return nil, ErrLargeUploadNotSupported
	}

	chunkSize := normalizeChunkSize(request.ChunkSize, capabilities.ChunkPolicy, request.Size)
	if request.Size > 0 && request.Size <= chunkSize {
		return c.Upload(ctx, request.UploadRequest)
	}

	session, err := c.StartUploadSession(ctx, StartUploadSessionRequest{
		Key:                request.Key,
		Size:               request.Size,
		ChunkSize:          chunkSize,
		ContentType:        request.ContentType,
		CacheControl:       request.CacheControl,
		ContentDisposition: request.ContentDisposition,
		Metadata:           request.Metadata,
	})
	if err != nil {
		return nil, err
	}

	completed := false
	defer func() {
		if !completed {
			_ = c.AbortUpload(context.Background(), *session)
		}
	}()

	buffer := make([]byte, int(chunkSize))
	uploadedChunks := make([]UploadedChunk, 0, int(maxInt64(1, capacityBySize(request.Size, chunkSize))))
	index := 1
	var offset int64

	for {
		if err = checkContext(ctx); err != nil {
			return nil, err
		}

		n, readErr := io.ReadFull(request.Body, buffer)
		if readErr == io.EOF && n == 0 {
			break
		}
		if readErr != nil && readErr != io.ErrUnexpectedEOF {
			return nil, fmt.Errorf("read upload chunk failed: %w", readErr)
		}
		if n == 0 {
			break
		}

		isFinal := false
		if request.Size > 0 {
			isFinal = offset+int64(n) >= request.Size
		} else if readErr == io.ErrUnexpectedEOF {
			isFinal = true
		}

		chunk, uploadErr := c.UploadChunk(ctx, UploadChunkRequest{
			Session: *session,
			Index:   index,
			Offset:  offset,
			Size:    int64(n),
			IsFinal: isFinal,
			Body:    bytes.NewReader(buffer[:n]),
		})
		if uploadErr != nil {
			return nil, uploadErr
		}

		uploadedChunks = append(uploadedChunks, *chunk)
		index++
		offset += int64(n)

		if readErr == io.ErrUnexpectedEOF {
			break
		}
	}

	if len(uploadedChunks) == 0 {
		return c.Upload(ctx, request.UploadRequest)
	}

	result, err := c.CompleteUpload(ctx, CompleteUploadRequest{
		Session: *session,
		Chunks:  uploadedChunks,
	})
	if err != nil {
		return nil, err
	}

	completed = true
	return result, nil
}

// StartUploadSession 初始化一个可断点续传的上传会话。
func (c *Client) StartUploadSession(ctx context.Context, request StartUploadSessionRequest) (*UploadSession, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}

	capabilities := c.Capabilities()
	if !capabilities.LargeUpload {
		return nil, ErrLargeUploadNotSupported
	}

	request.ChunkSize = normalizeChunkSize(request.ChunkSize, capabilities.ChunkPolicy, request.Size)
	return c.driver.startUploadSession(ctx, request)
}

// UploadChunk 上传单个数据块。
func (c *Client) UploadChunk(ctx context.Context, request UploadChunkRequest) (*UploadedChunk, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if !c.Capabilities().LargeUpload {
		return nil, ErrLargeUploadNotSupported
	}
	if request.Body == nil {
		return nil, ErrUploadBodyRequired
	}

	return c.driver.uploadChunk(ctx, request)
}

// GetUploadStatus 查询上传会话状态。
func (c *Client) GetUploadStatus(ctx context.Context, session UploadSession) (*UploadStatus, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if !c.Capabilities().UploadStatusQuery {
		return nil, ErrUploadStatusNotSupported
	}

	status, err := c.driver.getUploadStatus(ctx, session)
	if err != nil {
		return nil, err
	}

	status.UploadedChunks = sortUploadedChunks(status.UploadedChunks)
	return status, nil
}

// CompleteUpload 完成上传会话并合并对象。
func (c *Client) CompleteUpload(ctx context.Context, request CompleteUploadRequest) (*ObjectInfo, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if !c.Capabilities().LargeUpload {
		return nil, ErrLargeUploadNotSupported
	}

	request.Chunks = sortUploadedChunks(request.Chunks)
	return c.driver.completeUpload(ctx, request)
}

// AbortUpload 终止一个未完成的上传会话。
func (c *Client) AbortUpload(ctx context.Context, session UploadSession) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	if !c.Capabilities().LargeUpload {
		return ErrLargeUploadNotSupported
	}

	return c.driver.abortUpload(ctx, session)
}

// Delete 删除单个对象。
func (c *Client) Delete(ctx context.Context, key string) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	return c.driver.delete(ctx, key)
}

// DeleteMany 批量删除对象。
func (c *Client) DeleteMany(ctx context.Context, keys []string) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}

	if c.Capabilities().BatchDelete {
		return c.driver.deleteMany(ctx, keys)
	}

	for _, key := range keys {
		if err := c.driver.delete(ctx, key); err != nil {
			return err
		}
	}
	return nil
}

// Stat 获取对象元信息。
func (c *Client) Stat(ctx context.Context, key string) (*ObjectInfo, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	return c.driver.stat(ctx, key)
}

// Exists 判断对象是否存在。
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	if err := checkContext(ctx); err != nil {
		return false, err
	}
	return c.driver.exists(ctx, key)
}

// List 按前缀分页列出对象。
func (c *Client) List(ctx context.Context, request ListRequest) (*ListResult, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	return c.driver.list(ctx, request)
}

// Copy 复制对象。
func (c *Client) Copy(ctx context.Context, sourceKey, destinationKey string) (*CopyResult, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if !c.Capabilities().CopyObject {
		return nil, ErrCopyObjectNotSupported
	}

	return c.driver.copy(ctx, sourceKey, destinationKey)
}

// Move 移动对象，默认通过复制后删除实现。
func (c *Client) Move(ctx context.Context, sourceKey, destinationKey string) (*CopyResult, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if !c.Capabilities().CopyObject {
		return nil, ErrCopyObjectNotSupported
	}

	result, err := c.driver.copy(ctx, sourceKey, destinationKey)
	if err != nil {
		return nil, err
	}

	if err := c.driver.delete(ctx, sourceKey); err != nil {
		_ = c.driver.delete(ctx, destinationKey)
		return nil, fmt.Errorf("delete source object failed after copy: %w", err)
	}

	return result, nil
}

// SignedURL 生成对象的临时访问地址。
func (c *Client) SignedURL(ctx context.Context, key string, expire time.Duration) (string, error) {
	if err := checkContext(ctx); err != nil {
		return "", err
	}
	if !c.Capabilities().SignedURL {
		return "", ErrSignedURLNotSupported
	}

	return c.driver.signedURL(ctx, key, ensurePositiveExpire(expire))
}

// driver 定义底层云厂商驱动必须实现的统一能力。
type driver interface {
	provider() Provider
	capabilities() Capabilities
	upload(ctx context.Context, request UploadRequest) (*ObjectInfo, error)
	startUploadSession(ctx context.Context, request StartUploadSessionRequest) (*UploadSession, error)
	uploadChunk(ctx context.Context, request UploadChunkRequest) (*UploadedChunk, error)
	getUploadStatus(ctx context.Context, session UploadSession) (*UploadStatus, error)
	completeUpload(ctx context.Context, request CompleteUploadRequest) (*ObjectInfo, error)
	abortUpload(ctx context.Context, session UploadSession) error
	delete(ctx context.Context, key string) error
	deleteMany(ctx context.Context, keys []string) error
	stat(ctx context.Context, key string) (*ObjectInfo, error)
	exists(ctx context.Context, key string) (bool, error)
	list(ctx context.Context, request ListRequest) (*ListResult, error)
	copy(ctx context.Context, sourceKey, destinationKey string) (*CopyResult, error)
	signedURL(ctx context.Context, key string, expire time.Duration) (string, error)
}

// newDriver 根据 provider 选择实际驱动。
func newDriver(config Config) (driver, error) {
	switch config.Provider {
	case ProviderAliyunOSS:
		return newAliyunOSSDriver(config)
	case ProviderTencentCOS, ProviderAmazonS3, ProviderMinIO, ProviderQiniuKodo, ProviderHuaweiOBS, ProviderGoogleGCS, ProviderAzureBlob:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedProvider, config.Provider)
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedProvider, config.Provider)
	}
}

// maxInt64 返回两个 int64 中的较大值。
func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// capacityBySize 根据文件大小和 chunk 大小估算切片容量。
func capacityBySize(size, chunkSize int64) int64 {
	if size <= 0 || chunkSize <= 0 {
		return 1
	}
	return (size + chunkSize - 1) / chunkSize
}
