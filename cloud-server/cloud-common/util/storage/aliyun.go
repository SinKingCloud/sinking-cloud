package storage

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type aliyunOSSDriver struct {
	client *oss.Client
	bucket *oss.Bucket
	config Config
}

func (d *aliyunOSSDriver) provider() Provider {
	return ProviderAliyunOSS
}

func (d *aliyunOSSDriver) capabilities() capabilities {
	return capabilities{
		smallUpload:       true,
		largeUpload:       true,
		resumeUpload:      true,
		uploadStatusQuery: true,
		signedURL:         true,
		copyObject:        true,
		batchDelete:       true,
		protocols:         []uploadProtocol{uploadProtocolMultipart},
		chunkPolicy: chunkPolicy{
			defaultSize: defaultChunkSize,
			minSize:     minCompatibleChunkSize,
			maxChunks:   10000,
		},
	}
}

func (d *aliyunOSSDriver) upload(ctx context.Context, request uploadRequest) (*ObjectInfo, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if request.body == nil {
		return nil, errUploadBodyRequired
	}

	providerKey, err := buildProviderKey(d.config.Prefix, request.key)
	if err != nil {
		return nil, err
	}

	options := d.buildObjectOptions(
		detectContentType(request.contentType, request.key),
		request.cacheControl,
		request.contentDisposition,
		request.metadata,
	)

	if err := d.bucket.PutObject(providerKey, request.body, options...); err != nil {
		return nil, fmt.Errorf("upload object failed: %w", err)
	}

	return d.statByProviderKey(ctx, providerKey)
}

func (d *aliyunOSSDriver) startUploadSession(ctx context.Context, request startUploadSessionRequest) (*uploadSession, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}

	providerKey, err := buildProviderKey(d.config.Prefix, request.key)
	if err != nil {
		return nil, err
	}

	chunkSize := normalizeChunkSize(request.chunkSize, d.capabilities().chunkPolicy, request.size)
	contentType := detectContentType(request.contentType, request.key)
	options := d.buildObjectOptions(contentType, request.cacheControl, request.contentDisposition, request.metadata)

	result, err := d.bucket.InitiateMultipartUpload(providerKey, options...)
	if err != nil {
		return nil, fmt.Errorf("init upload session failed: %w", err)
	}

	return &uploadSession{
		provider:           ProviderAliyunOSS,
		protocol:           uploadProtocolMultipart,
		bucket:             d.config.Bucket,
		key:                request.key,
		providerKey:        providerKey,
		sessionID:          result.UploadID,
		size:               request.size,
		chunkSize:          chunkSize,
		contentType:        contentType,
		cacheControl:       request.cacheControl,
		contentDisposition: request.contentDisposition,
		metadata:           cloneMetadata(request.metadata),
		createdAt:          time.Now(),
	}, nil
}

func (d *aliyunOSSDriver) uploadChunk(ctx context.Context, request uploadChunkRequest) (*uploadedChunk, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if err := d.validateUploadChunkRequest(request); err != nil {
		return nil, err
	}

	part, err := d.bucket.UploadPart(d.uploadSessionHandle(request.session), request.body, request.size, request.index)
	if err != nil {
		return nil, fmt.Errorf("upload chunk failed: %w", err)
	}

	return &uploadedChunk{
		index:        request.index,
		offset:       request.offset,
		size:         request.size,
		etag:         strings.Trim(part.ETag, "\""),
		lastModified: time.Now(),
	}, nil
}

func (d *aliyunOSSDriver) getUploadStatus(ctx context.Context, session uploadSession) (*uploadStatus, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if err := d.validateUploadSession(session); err != nil {
		return nil, err
	}

	result, err := d.bucket.ListUploadedParts(d.uploadSessionHandle(session))
	if err != nil {
		return nil, fmt.Errorf("get upload status failed: %w", err)
	}

	chunks := make([]uploadedChunk, 0, len(result.UploadedParts))
	var committedBytes int64
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

	return &uploadStatus{
		session:        session,
		uploadedChunks: chunks,
		committedBytes: committedBytes,
	}, nil
}

func (d *aliyunOSSDriver) completeUpload(ctx context.Context, request completeUploadRequest) (*ObjectInfo, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if err := d.validateUploadSession(request.session); err != nil {
		return nil, err
	}
	if len(request.chunks) == 0 {
		return nil, errInvalidCompleteUpload
	}

	parts := make([]oss.UploadPart, 0, len(request.chunks))
	for _, chunk := range request.chunks {
		if chunk.index <= 0 || strings.TrimSpace(chunk.etag) == "" {
			return nil, errInvalidCompleteUpload
		}
		parts = append(parts, oss.UploadPart{
			PartNumber: chunk.index,
			ETag:       chunk.etag,
		})
	}

	sort.Sort(oss.UploadParts(parts))
	if _, err := d.bucket.CompleteMultipartUpload(d.uploadSessionHandle(request.session), parts); err != nil {
		return nil, fmt.Errorf("complete upload failed: %w", err)
	}

	return d.statByProviderKey(ctx, request.session.providerKey)
}

func (d *aliyunOSSDriver) abortUpload(ctx context.Context, session uploadSession) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	if err := d.validateUploadSession(session); err != nil {
		return err
	}
	if err := d.bucket.AbortMultipartUpload(d.uploadSessionHandle(session)); err != nil {
		return fmt.Errorf("abort upload failed: %w", err)
	}
	return nil
}

func (d *aliyunOSSDriver) delete(ctx context.Context, key string) error {
	if err := checkContext(ctx); err != nil {
		return err
	}

	providerKey, err := buildProviderKey(d.config.Prefix, key)
	if err != nil {
		return err
	}
	if err := d.bucket.DeleteObject(providerKey); err != nil {
		return fmt.Errorf("delete object failed: %w", err)
	}
	return nil
}

func (d *aliyunOSSDriver) deleteMany(ctx context.Context, keys []string) error {
	if err := checkContext(ctx); err != nil {
		return err
	}

	providerKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		providerKey, err := buildProviderKey(d.config.Prefix, key)
		if err != nil {
			return err
		}
		providerKeys = append(providerKeys, providerKey)
	}

	if _, err := d.bucket.DeleteObjects(providerKeys); err != nil {
		return fmt.Errorf("batch delete objects failed: %w", err)
	}
	return nil
}

func (d *aliyunOSSDriver) stat(ctx context.Context, key string) (*ObjectInfo, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}

	providerKey, err := buildProviderKey(d.config.Prefix, key)
	if err != nil {
		return nil, err
	}

	return d.statByProviderKey(ctx, providerKey)
}

func (d *aliyunOSSDriver) exists(ctx context.Context, key string) (bool, error) {
	if err := checkContext(ctx); err != nil {
		return false, err
	}

	providerKey, err := buildProviderKey(d.config.Prefix, key)
	if err != nil {
		return false, err
	}

	exists, err := d.bucket.IsObjectExist(providerKey)
	if err != nil {
		return false, fmt.Errorf("check object exists failed: %w", err)
	}
	return exists, nil
}

func (d *aliyunOSSDriver) list(ctx context.Context, request ListRequest) (*ListResult, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}

	limit := request.Limit
	if limit <= 0 {
		limit = 100
	}

	options := []oss.Option{
		oss.Prefix(buildProviderPrefix(d.config.Prefix, request.Prefix)),
		oss.MaxKeys(limit),
	}

	if marker := normalizeListPrefix(request.Marker); marker != "" {
		options = append(options, oss.Marker(buildProviderPrefix(d.config.Prefix, marker)))
	}

	result, err := d.bucket.ListObjects(options...)
	if err != nil {
		return nil, fmt.Errorf("list objects failed: %w", err)
	}

	objects := make([]ObjectInfo, 0, len(result.Objects))
	for _, object := range result.Objects {
		objects = append(objects, d.objectFromListObject(object))
	}

	return &ListResult{
		Objects:    objects,
		NextMarker: stripProviderPrefix(d.config.Prefix, result.NextMarker),
		HasMore:    result.IsTruncated,
	}, nil
}

func (d *aliyunOSSDriver) copy(ctx context.Context, sourceKey, destinationKey string) (*CopyResult, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}

	srcProviderKey, err := buildProviderKey(d.config.Prefix, sourceKey)
	if err != nil {
		return nil, err
	}
	dstProviderKey, err := buildProviderKey(d.config.Prefix, destinationKey)
	if err != nil {
		return nil, err
	}

	if _, err := d.bucket.CopyObject(srcProviderKey, dstProviderKey); err != nil {
		return nil, fmt.Errorf("copy object failed: %w", err)
	}

	info, err := d.statByProviderKey(ctx, dstProviderKey)
	if err != nil {
		return nil, err
	}

	return &CopyResult{
		SourceKey:      sourceKey,
		DestinationKey: destinationKey,
		ETag:           info.ETag,
		CopiedAt:       time.Now(),
	}, nil
}

func (d *aliyunOSSDriver) signedURL(ctx context.Context, key string, expire time.Duration) (string, error) {
	if err := checkContext(ctx); err != nil {
		return "", err
	}

	providerKey, err := buildProviderKey(d.config.Prefix, key)
	if err != nil {
		return "", err
	}

	signedURL, err := d.bucket.SignURL(providerKey, oss.HTTPGet, int64(expire.Seconds()))
	if err != nil {
		return "", fmt.Errorf("generate signed url failed: %w", err)
	}
	return signedURL, nil
}

func (d *aliyunOSSDriver) buildObjectOptions(contentType, cacheControl, contentDisposition string, metadata map[string]string) []oss.Option {
	options := make([]oss.Option, 0, len(metadata)+3)

	if strings.TrimSpace(contentType) != "" {
		options = append(options, oss.ContentType(contentType))
	}
	if strings.TrimSpace(cacheControl) != "" {
		options = append(options, oss.CacheControl(cacheControl))
	}
	if strings.TrimSpace(contentDisposition) != "" {
		options = append(options, oss.ContentDisposition(contentDisposition))
	}
	for key, value := range metadata {
		if strings.TrimSpace(key) == "" {
			continue
		}
		options = append(options, oss.Meta(key, value))
	}

	return options
}

func (d *aliyunOSSDriver) uploadSessionHandle(session uploadSession) oss.InitiateMultipartUploadResult {
	return oss.InitiateMultipartUploadResult{
		Bucket:   d.config.Bucket,
		Key:      session.providerKey,
		UploadID: session.sessionID,
	}
}

func (d *aliyunOSSDriver) validateUploadSession(session uploadSession) error {
	if strings.TrimSpace(session.providerKey) == "" || strings.TrimSpace(session.sessionID) == "" {
		return errInvalidUploadSession
	}
	if session.protocol != "" && session.protocol != uploadProtocolMultipart {
		return errInvalidUploadSession
	}
	return nil
}

func (d *aliyunOSSDriver) validateUploadChunkRequest(request uploadChunkRequest) error {
	if err := d.validateUploadSession(request.session); err != nil {
		return err
	}
	if request.body == nil || request.size <= 0 || request.index <= 0 {
		return errInvalidUploadChunk
	}
	return nil
}

func (d *aliyunOSSDriver) statByProviderKey(ctx context.Context, providerKey string) (*ObjectInfo, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}

	props, err := d.bucket.GetObjectDetailedMeta(providerKey)
	if err != nil {
		return nil, fmt.Errorf("get object meta failed: %w", err)
	}

	var size int64
	if sizeText := props.Get("Content-Length"); sizeText != "" {
		_, _ = fmt.Sscanf(sizeText, "%d", &size)
	}

	var lastModified time.Time
	if modifiedText := props.Get("Last-Modified"); modifiedText != "" {
		lastModified, _ = time.Parse(time.RFC1123, modifiedText)
	}

	return &ObjectInfo{
		Key:          stripProviderPrefix(d.config.Prefix, providerKey),
		ProviderKey:  providerKey,
		Size:         size,
		LastModified: lastModified,
		ETag:         strings.Trim(props.Get("ETag"), "\""),
		ContentType:  props.Get("Content-Type"),
		URL:          d.objectURL(providerKey),
		Metadata:     metadataFromHeaders(props, "x-oss-meta-"),
	}, nil
}

func (d *aliyunOSSDriver) objectFromListObject(object oss.ObjectProperties) ObjectInfo {
	return ObjectInfo{
		Key:          stripProviderPrefix(d.config.Prefix, object.Key),
		ProviderKey:  object.Key,
		Size:         object.Size,
		LastModified: object.LastModified,
		ETag:         strings.Trim(object.ETag, "\""),
		URL:          d.objectURL(object.Key),
	}
}

func (d *aliyunOSSDriver) objectURL(providerKey string) string {
	if strings.TrimSpace(d.config.Domain) != "" {
		return strings.TrimRight(ensureScheme(d.config.Domain, d.config.UseHTTPS), "/") + "/" + providerKey
	}

	endpoint := strings.TrimRight(d.config.Endpoint, "/")
	endpoint = strings.TrimPrefix(endpoint, "http://")
	endpoint = strings.TrimPrefix(endpoint, "https://")
	scheme := "https"
	if !d.config.UseHTTPS {
		scheme = "http"
	}

	return fmt.Sprintf("%s://%s.%s/%s", scheme, d.config.Bucket, endpoint, providerKey)
}
