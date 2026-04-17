package storage

import (
	"context"
	"errors"
	"mime"
	"path"
	"sort"
	"strings"
	"time"
)

const (
	// DefaultChunkSize 是默认分块大小，兼顾主流云厂商兼容性和服务端内存占用。
	DefaultChunkSize int64 = 8 << 20
	// MinCompatibleChunkSize 是跨云存储场景下较稳妥的最小分块大小。
	MinCompatibleChunkSize int64 = 5 << 20

	// Deprecated: 请改用 DefaultChunkSize。
	DefaultMultipartPartSize int64 = DefaultChunkSize
	// Deprecated: 请改用 MinCompatibleChunkSize。
	MinCompatiblePartSize int64 = MinCompatibleChunkSize
)

var (
	// ErrUnsupportedProvider 表示当前 provider 还没有实现驱动。
	ErrUnsupportedProvider = errors.New("storage provider is not implemented")
	// ErrInvalidObjectKey 表示对象 key 非法。
	ErrInvalidObjectKey = errors.New("invalid storage object key")
	// ErrLargeUploadNotSupported 表示当前驱动不支持大文件分块上传。
	ErrLargeUploadNotSupported = errors.New("large upload is not supported")
	// ErrUploadStatusNotSupported 表示当前驱动不支持查询上传进度。
	ErrUploadStatusNotSupported = errors.New("upload status query is not supported")
	// ErrCopyObjectNotSupported 表示当前驱动不支持对象复制。
	ErrCopyObjectNotSupported = errors.New("copy object is not supported")
	// ErrSignedURLNotSupported 表示当前驱动不支持签名访问地址。
	ErrSignedURLNotSupported = errors.New("signed URL is not supported")
	// ErrUploadBodyRequired 表示上传必须提供文件流。
	ErrUploadBodyRequired = errors.New("upload body is required")
	// ErrInvalidUploadSession 表示上传会话非法。
	ErrInvalidUploadSession = errors.New("invalid upload session")
	// ErrInvalidUploadChunk 表示上传分块请求非法。
	ErrInvalidUploadChunk = errors.New("invalid upload chunk request")
	// ErrInvalidCompleteUpload 表示完成上传请求非法。
	ErrInvalidCompleteUpload = errors.New("invalid complete upload request")

	// Deprecated: 请改用 ErrLargeUploadNotSupported。
	ErrMultipartNotSupported = ErrLargeUploadNotSupported
	// Deprecated: 请改用 ErrInvalidCompleteUpload 或 ErrInvalidUploadChunk。
	ErrInvalidMultipartUpload = ErrInvalidCompleteUpload
)

// checkContext 检查上下文是否已经取消。
func checkContext(ctx context.Context) error {
	if ctx == nil {
		return nil
	}
	return ctx.Err()
}

// normalizeObjectKey 规范化对象 key，并拦截路径穿越。
func normalizeObjectKey(key string) (string, error) {
	key = strings.TrimSpace(strings.ReplaceAll(key, "\\", "/"))
	key = strings.TrimPrefix(key, "/")
	if key == "" {
		return "", ErrInvalidObjectKey
	}

	cleanKey := path.Clean(key)
	if cleanKey == "." || cleanKey == "/" || cleanKey == "" || cleanKey == ".." || strings.HasPrefix(cleanKey, "../") {
		return "", ErrInvalidObjectKey
	}

	return strings.TrimPrefix(cleanKey, "/"), nil
}

// normalizePrefix 规范化配置中的基础前缀。
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

// normalizeListPrefix 规范化列表查询前缀。
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

// buildProviderKey 将逻辑 key 映射为实际存储 provider 使用的 key。
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

// buildProviderPrefix 将逻辑列表前缀映射为实际存储 provider 使用的前缀。
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

// stripProviderPrefix 将真实 provider key 转回逻辑 key。
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

// normalizeChunkSize 根据驱动分块策略归一化分块大小。
func normalizeChunkSize(chunkSize int64, policy ChunkPolicy, objectSize int64) int64 {
	if policy.DefaultSize <= 0 {
		policy.DefaultSize = DefaultChunkSize
	}
	if policy.MinSize <= 0 {
		policy.MinSize = MinCompatibleChunkSize
	}

	if chunkSize <= 0 {
		chunkSize = policy.DefaultSize
	}
	if chunkSize < policy.MinSize {
		chunkSize = policy.MinSize
	}

	if policy.MaxChunks > 0 && objectSize > 0 {
		minSizeByCount := (objectSize + int64(policy.MaxChunks) - 1) / int64(policy.MaxChunks)
		if minSizeByCount > chunkSize {
			chunkSize = minSizeByCount
		}
	}

	if policy.SizeMultiple > 0 {
		chunkSize = roundUpMultiple(chunkSize, policy.SizeMultiple)
	}

	if chunkSize < policy.MinSize {
		chunkSize = policy.MinSize
	}
	return chunkSize
}

// detectContentType 推断内容类型。
func detectContentType(contentType, key string) string {
	if strings.TrimSpace(contentType) != "" {
		return contentType
	}

	if guessed := mime.TypeByExtension(strings.ToLower(path.Ext(key))); guessed != "" {
		return guessed
	}

	return "application/octet-stream"
}

// cloneMetadata 复制 metadata，避免外部 map 被内部修改。
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

// sortUploadedChunks 对已上传分块结果按序排序。
func sortUploadedChunks(chunks []UploadedChunk) []UploadedChunk {
	sorted := make([]UploadedChunk, len(chunks))
	copy(sorted, chunks)
	sort.Slice(sorted, func(i, j int) bool {
		left := sorted[i]
		right := sorted[j]

		switch {
		case left.Index > 0 && right.Index > 0 && left.Index != right.Index:
			return left.Index < right.Index
		case left.Offset != right.Offset:
			return left.Offset < right.Offset
		default:
			return left.Index < right.Index
		}
	})
	return sorted
}

// ensurePositiveExpire 对签名 URL 有效期做兜底。
func ensurePositiveExpire(expire time.Duration) time.Duration {
	if expire <= 0 {
		return 15 * time.Minute
	}
	return expire
}

// roundUpMultiple 将分块大小向上对齐到指定倍数。
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
