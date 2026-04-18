package storage

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

const (
	defaultChunkSize       int64 = 8 << 20
	minCompatibleChunkSize int64 = 5 << 20
	uploadIDVersion              = "v1"
)

var (
	errUnsupportedProvider      = errors.New("storage provider is not implemented")
	errInvalidObjectKey         = errors.New("invalid storage object key")
	errLargeUploadNotSupported  = errors.New("large upload is not supported")
	errUploadStatusNotSupported = errors.New("upload status query is not supported")
	errCopyObjectNotSupported   = errors.New("copy object is not supported")
	errSignedURLNotSupported    = errors.New("signed URL is not supported")
	errUploadBodyRequired       = errors.New("upload body is required")
	errInvalidUploadSession     = errors.New("invalid upload session")
	errInvalidUploadChunk       = errors.New("invalid upload chunk request")
	errInvalidCompleteUpload    = errors.New("invalid complete upload request")
	errInvalidUploadID          = errors.New("invalid upload ID")
)

func checkContext(ctx context.Context) error {
	if ctx == nil {
		return nil
	}
	return ctx.Err()
}

func normalizeObjectKey(key string) (string, error) {
	key = strings.TrimSpace(strings.ReplaceAll(key, "\\", "/"))
	key = strings.TrimPrefix(key, "/")
	if key == "" {
		return "", errInvalidObjectKey
	}

	cleanKey := path.Clean(key)
	if cleanKey == "." || cleanKey == "/" || cleanKey == "" || cleanKey == ".." || strings.HasPrefix(cleanKey, "../") {
		return "", errInvalidObjectKey
	}

	return strings.TrimPrefix(cleanKey, "/"), nil
}

func normalizePrefix(prefix string) string {
	prefix = strings.TrimSpace(strings.ReplaceAll(prefix, "\\", "/"))
	prefix = strings.Trim(prefix, "/")
	if prefix == "" {
		return ""
	}

	cleanPrefix := path.Clean(prefix)
	if cleanPrefix == "." || cleanPrefix == "/" {
		return ""
	}

	return strings.Trim(cleanPrefix, "/")
}

func normalizeListPrefix(prefix string) string {
	prefix = strings.TrimSpace(strings.ReplaceAll(prefix, "\\", "/"))
	prefix = strings.Trim(prefix, "/")
	if prefix == "" {
		return ""
	}

	cleanPrefix := path.Clean(prefix)
	if cleanPrefix == "." || cleanPrefix == "/" || cleanPrefix == ".." || strings.HasPrefix(cleanPrefix, "../") {
		return ""
	}

	return strings.Trim(cleanPrefix, "/")
}

func buildProviderKey(prefix, key string) (string, error) {
	logicalKey, err := normalizeObjectKey(key)
	if err != nil {
		return "", err
	}

	prefix = normalizePrefix(prefix)
	if prefix == "" {
		return logicalKey, nil
	}

	return path.Join(prefix, logicalKey), nil
}

func buildProviderPrefix(prefix, listPrefix string) string {
	basePrefix := normalizePrefix(prefix)
	logicalPrefix := normalizeListPrefix(listPrefix)

	switch {
	case basePrefix == "" && logicalPrefix == "":
		return ""
	case basePrefix == "":
		return logicalPrefix
	case logicalPrefix == "":
		return basePrefix + "/"
	default:
		return path.Join(basePrefix, logicalPrefix)
	}
}

func stripProviderPrefix(prefix, providerKey string) string {
	prefix = normalizePrefix(prefix)
	providerKey = strings.Trim(strings.ReplaceAll(providerKey, "\\", "/"), "/")

	if prefix == "" {
		return providerKey
	}
	if providerKey == prefix {
		return ""
	}

	prefixWithSlash := prefix + "/"
	if strings.HasPrefix(providerKey, prefixWithSlash) {
		return strings.TrimPrefix(providerKey, prefixWithSlash)
	}

	return providerKey
}

func normalizeChunkSize(chunkSize int64, policy chunkPolicy, objectSize int64) int64 {
	if policy.defaultSize <= 0 {
		policy.defaultSize = defaultChunkSize
	}
	if policy.minSize <= 0 {
		policy.minSize = minCompatibleChunkSize
	}

	if chunkSize <= 0 {
		chunkSize = policy.defaultSize
	}
	if chunkSize < policy.minSize {
		chunkSize = policy.minSize
	}

	if policy.maxChunks > 0 && objectSize > 0 {
		minSizeByCount := (objectSize + int64(policy.maxChunks) - 1) / int64(policy.maxChunks)
		if minSizeByCount > chunkSize {
			chunkSize = minSizeByCount
		}
	}

	if policy.sizeMultiple > 0 {
		chunkSize = roundUpMultiple(chunkSize, policy.sizeMultiple)
	}

	if chunkSize < policy.minSize {
		chunkSize = policy.minSize
	}
	return chunkSize
}

func detectContentType(contentType, key string) string {
	if strings.TrimSpace(contentType) != "" {
		return contentType
	}

	if guessed := mime.TypeByExtension(strings.ToLower(path.Ext(key))); guessed != "" {
		return guessed
	}

	return "application/octet-stream"
}

func cloneMetadata(metadata map[string]string) map[string]string {
	if len(metadata) == 0 {
		return nil
	}

	cloned := make(map[string]string, len(metadata))
	for key, value := range metadata {
		cloned[key] = value
	}

	return cloned
}

func sortUploadedChunks(chunks []uploadedChunk) []uploadedChunk {
	sorted := make([]uploadedChunk, len(chunks))
	copy(sorted, chunks)
	sort.Slice(sorted, func(i, j int) bool {
		left := sorted[i]
		right := sorted[j]

		switch {
		case left.index > 0 && right.index > 0 && left.index != right.index:
			return left.index < right.index
		case left.offset != right.offset:
			return left.offset < right.offset
		default:
			return left.index < right.index
		}
	})
	return sorted
}

func ensurePositiveExpire(expire time.Duration) time.Duration {
	if expire <= 0 {
		return 15 * time.Minute
	}
	return expire
}

func roundUpMultiple(value, multiple int64) int64 {
	if multiple <= 0 {
		return value
	}
	remainder := value % multiple
	if remainder == 0 {
		return value
	}
	return value + multiple - remainder
}

func newDriver(config Config) (driver, error) {
	switch config.Provider {
	case ProviderAliyunOSS:
		return newAliyunOSSDriver(config)
	default:
		return nil, fmt.Errorf("%w: %s", errUnsupportedProvider, config.Provider)
	}
}

func newAliyunOSSDriver(config Config) (*aliyunOSSDriver, error) {
	if strings.TrimSpace(config.Endpoint) == "" ||
		strings.TrimSpace(config.Bucket) == "" ||
		strings.TrimSpace(config.AccessKeyID) == "" ||
		strings.TrimSpace(config.AccessKeySecret) == "" {
		return nil, fmt.Errorf("aliyun oss config is incomplete")
	}

	config.UseHTTPS = resolveHTTPSPreference(config.UseHTTPS, config.Endpoint, config.Domain)
	config.Prefix = normalizePrefix(config.Prefix)

	clientOptions := make([]oss.ClientOption, 0, 4)
	if strings.TrimSpace(config.SecurityToken) != "" {
		clientOptions = append(clientOptions, oss.SecurityToken(config.SecurityToken))
	}
	if strings.TrimSpace(config.Region) != "" {
		clientOptions = append(clientOptions, oss.Region(config.Region), oss.AuthVersion(oss.AuthV4))
	}

	client, err := oss.New(config.Endpoint, config.AccessKeyID, config.AccessKeySecret, clientOptions...)
	if err != nil {
		return nil, fmt.Errorf("create aliyun oss client failed: %w", err)
	}

	bucket, err := client.Bucket(config.Bucket)
	if err != nil {
		return nil, fmt.Errorf("open aliyun oss bucket failed: %w", err)
	}

	return &aliyunOSSDriver{
		client: client,
		bucket: bucket,
		config: config,
	}, nil
}

func buildUploadIDSecret(config Config) []byte {
	material := strings.Join([]string{
		"storage-upload-id-secret",
		string(config.Provider),
		strings.TrimSpace(config.Endpoint),
		strings.TrimSpace(config.Region),
		strings.TrimSpace(config.Bucket),
		strings.TrimSpace(config.AccessKeyID),
		strings.TrimSpace(config.AccessKeySecret),
		normalizePrefix(config.Prefix),
	}, "\n")

	sum := sha256.Sum256([]byte(material))
	secret := make([]byte, len(sum))
	copy(secret, sum[:])
	return secret
}

func encodeUploadID(secret []byte, session uploadSession) (string, error) {
	payload, err := json.Marshal(uploadSessionPayload{
		Provider:           session.provider,
		Protocol:           string(session.protocol),
		Bucket:             session.bucket,
		Key:                session.key,
		ProviderKey:        session.providerKey,
		SessionID:          session.sessionID,
		Size:               session.size,
		ChunkSize:          session.chunkSize,
		ContentType:        session.contentType,
		CacheControl:       session.cacheControl,
		ContentDisposition: session.contentDisposition,
		Metadata:           cloneMetadata(session.metadata),
		Opaque:             cloneMetadata(session.opaque),
		CreatedAt:          session.createdAt,
	})
	if err != nil {
		return "", err
	}

	signature := signUploadIDPayload(secret, payload)
	return strings.Join([]string{
		uploadIDVersion,
		base64.RawURLEncoding.EncodeToString(payload),
		base64.RawURLEncoding.EncodeToString(signature),
	}, "."), nil
}

func decodeUploadID(secret []byte, uploadID string) (uploadSession, error) {
	if strings.TrimSpace(uploadID) == "" {
		return uploadSession{}, errInvalidUploadID
	}

	parts := strings.Split(uploadID, ".")
	if len(parts) != 3 || parts[0] != uploadIDVersion {
		return uploadSession{}, errInvalidUploadID
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return uploadSession{}, errInvalidUploadID
	}
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return uploadSession{}, errInvalidUploadID
	}
	if !hmac.Equal(signature, signUploadIDPayload(secret, payload)) {
		return uploadSession{}, errInvalidUploadID
	}

	var session uploadSessionPayload
	if err := json.Unmarshal(payload, &session); err != nil {
		return uploadSession{}, errInvalidUploadID
	}
	if strings.TrimSpace(session.SessionID) == "" || strings.TrimSpace(session.ProviderKey) == "" {
		return uploadSession{}, errInvalidUploadID
	}

	return uploadSession{
		provider:           session.Provider,
		protocol:           uploadProtocol(session.Protocol),
		bucket:             session.Bucket,
		key:                session.Key,
		providerKey:        session.ProviderKey,
		sessionID:          session.SessionID,
		size:               session.Size,
		chunkSize:          session.ChunkSize,
		contentType:        session.ContentType,
		cacheControl:       session.CacheControl,
		contentDisposition: session.ContentDisposition,
		metadata:           cloneMetadata(session.Metadata),
		opaque:             cloneMetadata(session.Opaque),
		createdAt:          session.CreatedAt,
	}, nil
}

func signUploadIDPayload(secret, payload []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write(payload)
	return mac.Sum(nil)
}

func buildChunkID(partNumber int) string {
	return base64.RawURLEncoding.EncodeToString([]byte(strconv.Itoa(partNumber)))
}

func collectUploadedChunks(fetchPage func(partNumberMarker int) (oss.ListUploadedPartsResult, error)) ([]uploadedChunk, int64, error) {
	chunks := make([]uploadedChunk, 0, 16)
	var committedBytes int64
	partNumberMarker := 0

	for {
		result, err := fetchPage(partNumberMarker)
		if err != nil {
			return nil, 0, err
		}

		for _, part := range result.UploadedParts {
			size := int64(part.Size)
			chunks = append(chunks, uploadedChunk{
				index:        part.PartNumber,
				size:         size,
				etag:         strings.Trim(part.ETag, "\""),
				lastModified: part.LastModified,
			})
			committedBytes += size
		}

		if !result.IsTruncated {
			return chunks, committedBytes, nil
		}

		nextMarker, err := strconv.Atoi(strings.TrimSpace(result.NextPartNumberMarker))
		if err != nil || nextMarker <= partNumberMarker {
			return nil, 0, fmt.Errorf("invalid uploaded parts pagination marker")
		}
		partNumberMarker = nextMarker
	}
}

// uploadSessionPayload 是用于编码上传任务 ID 的可序列化会话结构。
type uploadSessionPayload struct {
	Provider           Provider          `json:"provider"`           // 标识当前会话所属的存储厂商。
	Protocol           string            `json:"protocol"`           // 标识当前会话采用的上传协议。
	Bucket             string            `json:"bucket"`             // 是当前会话所属的 Bucket 名称。
	Key                string            `json:"key"`                // 是业务层看到的逻辑对象路径。
	ProviderKey        string            `json:"providerKey"`        // 是底层厂商实际使用的对象路径。
	SessionID          string            `json:"sessionId"`          // 是厂商返回的原始上传会话 ID。
	Size               int64             `json:"size"`               // 是整个文件总大小，单位为字节。
	ChunkSize          int64             `json:"chunkSize"`          // 是当前会话采用的分片大小。
	ContentType        string            `json:"contentType"`        // 是对象的内容类型。
	CacheControl       string            `json:"cacheControl"`       // 是对象的缓存控制头。
	ContentDisposition string            `json:"contentDisposition"` // 是对象的下载展示方式。
	Metadata           map[string]string `json:"metadata,omitempty"` // 是对象自定义元数据。
	Opaque             map[string]string `json:"opaque,omitempty"`   // 用于承载不同厂商的私有扩展字段。
	CreatedAt          time.Time         `json:"createdAt"`          // 是会话初始化时间。
}

func metadataFromHeaders(headers http.Header, prefix string) map[string]string {
	if len(headers) == 0 {
		return nil
	}

	prefix = strings.ToLower(strings.TrimSpace(prefix))
	if prefix == "" {
		return nil
	}

	metadata := make(map[string]string)
	for key, values := range headers {
		if len(values) == 0 {
			continue
		}

		lowerKey := strings.ToLower(key)
		if !strings.HasPrefix(lowerKey, prefix) {
			continue
		}

		metadata[strings.TrimPrefix(lowerKey, prefix)] = values[0]
	}

	if len(metadata) == 0 {
		return nil
	}
	return metadata
}

func resolveHTTPSPreference(useHTTPS bool, values ...string) bool {
	for _, value := range values {
		value = strings.ToLower(strings.TrimSpace(value))
		switch {
		case strings.HasPrefix(value, "http://"):
			return false
		case strings.HasPrefix(value, "https://"):
			return true
		}
	}

	if useHTTPS {
		return true
	}
	return true
}

func ensureScheme(raw string, useHTTPS bool) string {
	raw = strings.TrimSpace(raw)
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		return raw
	}

	scheme := "https://"
	if !useHTTPS {
		scheme = "http://"
	}

	if parsed, err := url.Parse(raw); err == nil && parsed.Scheme != "" {
		return raw
	}

	return scheme + raw
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func capacityBySize(size, chunkSize int64) int64 {
	if size <= 0 || chunkSize <= 0 {
		return 1
	}
	return (size + chunkSize - 1) / chunkSize
}
