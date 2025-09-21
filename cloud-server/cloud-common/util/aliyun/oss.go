package aliyun

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// 重试配置
const (
	maxRetries = 3
	retryDelay = time.Second
)

// OSSConfig OSS配置信息
type OSSConfig struct {
	Endpoint        string // OSS地域节点
	AccessKeyID     string // AccessKey ID
	AccessKeySecret string // AccessKey Secret
	BucketName      string // 存储桶名称
	BasePath        string // 基础路径前缀
}

// OSSClient OSS客户端
type OSSClient struct {
	client *oss.Client
	bucket *oss.Bucket
	config *OSSConfig
}

// FileInfo 文件信息
type FileInfo struct {
	Key          string    `json:"key"`          // 文件键名
	Size         int64     `json:"size"`         // 文件大小
	LastModified time.Time `json:"lastModified"` // 最后修改时间
	ETag         string    `json:"etag"`         // 文件ETag
	ContentType  string    `json:"contentType"`  // 文件类型
	URL          string    `json:"url"`          // 访问URL
}

// UploadResult 上传结果
type UploadResult struct {
	Key         string `json:"key"`         // 文件键名
	URL         string `json:"url"`         // 访问URL
	Size        int64  `json:"size"`        // 文件大小
	ContentType string `json:"contentType"` // 文件类型
	ETag        string `json:"etag"`        // 文件ETag
}

// MultipartUploadResult 分片上传结果
type MultipartUploadResult struct {
	UploadResult
	UploadID string `json:"uploadId"` // 上传ID
	Parts    int    `json:"parts"`    // 分片数量
}

// NewOSSClient 创建OSS客户端
func NewOSSClient(config *OSSConfig) (*OSSClient, error) {
	if config == nil {
		return nil, fmt.Errorf("OSS配置不能为空")
	}

	if config.Endpoint == "" || config.AccessKeyID == "" ||
		config.AccessKeySecret == "" || config.BucketName == "" {
		return nil, fmt.Errorf("OSS配置信息不完整")
	}

	// 创建OSS客户端
	client, err := oss.New(config.Endpoint, config.AccessKeyID, config.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("创建OSS客户端失败: %v", err)
	}

	// 获取存储桶
	bucket, err := client.Bucket(config.BucketName)
	if err != nil {
		return nil, fmt.Errorf("获取存储桶失败: %v", err)
	}

	return &OSSClient{
		client: client,
		bucket: bucket,
		config: config,
	}, nil
}

// buildKey 构建文件键名
func (c *OSSClient) buildKey(filename string) string {
	key := filename
	if c.config.BasePath != "" {
		basePath := strings.TrimRight(c.config.BasePath, "/")
		key = basePath + "/" + strings.TrimLeft(filename, "/")
	}
	return key
}

// getFileURL 获取文件访问URL
func (c *OSSClient) getFileURL(key string) string {
	// 如果是自定义域名，直接拼接
	if !strings.Contains(c.config.Endpoint, "aliyuncs.com") {
		return fmt.Sprintf("https://%s/%s", c.config.Endpoint, key)
	}

	// 标准OSS域名格式
	endpoint := strings.TrimPrefix(c.config.Endpoint, "http://")
	endpoint = strings.TrimPrefix(endpoint, "https://")
	return fmt.Sprintf("https://%s.%s/%s", c.config.BucketName, endpoint, key)
}

// retryOperation 重试操作
func (c *OSSClient) retryOperation(operation func() error) error {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if err := operation(); err != nil {
			lastErr = err
			if i < maxRetries-1 {
				time.Sleep(retryDelay * time.Duration(i+1)) // 指数退避
			}
			continue
		}
		return nil
	}
	return lastErr
}

// UploadFile 简单上传文件
func (c *OSSClient) UploadFile(filename string, file io.Reader, contentType string) (*UploadResult, error) {
	key := c.buildKey(filename)

	// 设置上传选项
	var options []oss.Option
	if contentType != "" {
		options = append(options, oss.ContentType(contentType))
	}

	// 带重试的上传文件
	err := c.retryOperation(func() error {
		return c.bucket.PutObject(key, file, options...)
	})
	if err != nil {
		return nil, fmt.Errorf("上传文件失败: %v", err)
	}

	// 获取文件信息
	props, err := c.bucket.GetObjectDetailedMeta(key)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %v", err)
	}

	// 获取文件大小
	var size int64
	if sizeStr := props.Get("Content-Length"); sizeStr != "" {
		_, _ = fmt.Sscanf(sizeStr, "%d", &size)
	}

	return &UploadResult{
		Key:         key,
		URL:         c.getFileURL(key),
		Size:        size,
		ContentType: props.Get("Content-Type"),
		ETag:        strings.Trim(props.Get("ETag"), "\""),
	}, nil
}

// UploadFileFromPath 从本地路径上传文件
func (c *OSSClient) UploadFileFromPath(filename string, localPath string) (*UploadResult, error) {
	file, err := os.Open(localPath)
	if err != nil {
		return nil, fmt.Errorf("打开本地文件失败: %v", err)
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %v", err)
	}

	// 根据文件扩展名获取ContentType
	contentType := getContentType(filepath.Ext(localPath))

	result, err := c.UploadFile(filename, file, contentType)
	if err != nil {
		return nil, err
	}

	result.Size = fileInfo.Size()
	return result, nil
}

// UploadFileFromMultipart 从multipart.File上传文件
func (c *OSSClient) UploadFileFromMultipart(filename string, file multipart.File, header *multipart.FileHeader) (*UploadResult, error) {
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = getContentType(filepath.Ext(header.Filename))
	}

	result, err := c.UploadFile(filename, file, contentType)
	if err != nil {
		return nil, err
	}

	result.Size = header.Size
	return result, nil
}

// MultipartUpload 分片上传文件
func (c *OSSClient) MultipartUpload(filename string, file io.Reader, fileSize int64, partSize int64) (*MultipartUploadResult, error) {
	key := c.buildKey(filename)

	// 限制分片大小以防止内存溢出
	const maxPartSize = 100 * 1024 * 1024 // 100MB
	if partSize > maxPartSize {
		partSize = maxPartSize
	}
	if partSize < 5*1024*1024 { // OSS最小分片大小5MB
		partSize = 5 * 1024 * 1024
	}

	// 初始化分片上传
	imur, err := c.bucket.InitiateMultipartUpload(key)
	if err != nil {
		return nil, fmt.Errorf("初始化分片上传失败: %v", err)
	}

	// 分片上传
	var parts []oss.UploadPart
	var partNumber = 1
	var uploadedSize int64 = 0

	// 复用缓冲区以减少内存分配
	buffer := make([]byte, partSize)

	for uploadedSize < fileSize {
		// 计算当前分片大小
		currentPartSize := partSize
		if uploadedSize+partSize > fileSize {
			currentPartSize = fileSize - uploadedSize
		}

		// 读取分片数据
		n, err := io.ReadFull(file, buffer[:currentPartSize])
		if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
			// 取消分片上传
			_ = c.bucket.AbortMultipartUpload(imur)
			return nil, fmt.Errorf("读取分片数据失败: %v", err)
		}

		// 上传分片
		part, err := c.bucket.UploadPart(imur, bytes.NewReader(buffer[:n]), int64(n), partNumber)
		if err != nil {
			// 取消分片上传
			_ = c.bucket.AbortMultipartUpload(imur)
			return nil, fmt.Errorf("上传分片失败: %v", err)
		}

		parts = append(parts, part)
		partNumber++
		uploadedSize += int64(n)
	}

	// 完成分片上传
	cmur, err := c.bucket.CompleteMultipartUpload(imur, parts)
	if err != nil {
		// 取消分片上传
		_ = c.bucket.AbortMultipartUpload(imur)
		return nil, fmt.Errorf("完成分片上传失败: %v", err)
	}

	return &MultipartUploadResult{
		UploadResult: UploadResult{
			Key:         key,
			URL:         c.getFileURL(key),
			Size:        fileSize,
			ContentType: getContentType(filepath.Ext(filename)),
			ETag:        strings.Trim(cmur.ETag, "\""),
		},
		UploadID: imur.UploadID,
		Parts:    len(parts),
	}, nil
}

// DeleteFile 删除文件
func (c *OSSClient) DeleteFile(filename string) error {
	key := c.buildKey(filename)

	err := c.bucket.DeleteObject(key)
	if err != nil {
		return fmt.Errorf("删除文件失败: %v", err)
	}

	return nil
}

// DeleteFiles 批量删除文件
func (c *OSSClient) DeleteFiles(filenames []string) error {
	if len(filenames) == 0 {
		return nil
	}

	keys := make([]string, len(filenames))
	for i, filename := range filenames {
		keys[i] = c.buildKey(filename)
	}

	_, err := c.bucket.DeleteObjects(keys)
	if err != nil {
		return fmt.Errorf("批量删除文件失败: %v", err)
	}

	return nil
}

// GetFileInfo 获取文件信息
func (c *OSSClient) GetFileInfo(filename string) (*FileInfo, error) {
	key := c.buildKey(filename)

	// 获取文件元数据
	props, err := c.bucket.GetObjectDetailedMeta(key)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %v", err)
	}

	// 解析文件大小
	var size int64
	if sizeStr := props.Get("Content-Length"); sizeStr != "" {
		_, _ = fmt.Sscanf(sizeStr, "%d", &size)
	}

	// 解析最后修改时间
	var lastModified time.Time
	if modifiedStr := props.Get("Last-Modified"); modifiedStr != "" {
		lastModified, _ = time.Parse(time.RFC1123, modifiedStr)
	}

	return &FileInfo{
		Key:          key,
		Size:         size,
		LastModified: lastModified,
		ETag:         strings.Trim(props.Get("ETag"), "\""),
		ContentType:  props.Get("Content-Type"),
		URL:          c.getFileURL(key),
	}, nil
}

// ListFiles 获取文件列表
func (c *OSSClient) ListFiles(prefix string, maxKeys int) ([]*FileInfo, error) {
	result, err := c.ListFilesWithPage(prefix, "", maxKeys)
	if err != nil {
		return nil, err
	}
	return result.Files, nil
}

// ListFilesResult 文件列表结果
type ListFilesResult struct {
	Files       []*FileInfo // 文件列表
	NextMarker  string      // 下一页标记
	IsTruncated bool        // 是否还有更多数据
	TotalCount  int         // 总数量（OSS不支持，始终为0）
}

// ListFilesWithPage 分页获取文件列表
func (c *OSSClient) ListFilesWithPage(prefix string, marker string, maxKeys int) (*ListFilesResult, error) {
	if maxKeys <= 0 {
		maxKeys = 100
	}

	// 构建前缀
	fullPrefix := prefix
	if c.config.BasePath != "" {
		basePath := strings.TrimRight(c.config.BasePath, "/")
		if prefix != "" {
			fullPrefix = basePath + "/" + strings.TrimLeft(prefix, "/")
		} else {
			fullPrefix = basePath + "/"
		}
	}

	// 构建查询选项
	options := []oss.Option{
		oss.Prefix(fullPrefix),
		oss.MaxKeys(maxKeys),
	}

	// 添加分页标记
	if marker != "" {
		options = append(options, oss.Marker(marker))
	}

	// 列举文件
	lsRes, err := c.bucket.ListObjects(options...)
	if err != nil {
		return nil, fmt.Errorf("获取文件列表失败: %v", err)
	}

	var files []*FileInfo
	for _, object := range lsRes.Objects {
		files = append(files, &FileInfo{
			Key:          object.Key,
			Size:         object.Size,
			LastModified: object.LastModified,
			ETag:         strings.Trim(object.ETag, "\""),
			ContentType:  object.Type,
			URL:          c.getFileURL(object.Key),
		})
	}

	return &ListFilesResult{
		Files:       files,
		NextMarker:  lsRes.NextMarker,
		IsTruncated: lsRes.IsTruncated,
		TotalCount:  0, // OSS不支持获取总数
	}, nil
}

// ListAllFiles 获取所有文件列表（自动分页）
func (c *OSSClient) ListAllFiles(prefix string) ([]*FileInfo, error) {
	var allFiles []*FileInfo
	marker := ""

	for {
		result, err := c.ListFilesWithPage(prefix, marker, 1000)
		if err != nil {
			return nil, err
		}

		allFiles = append(allFiles, result.Files...)

		// 如果没有更多数据，退出循环
		if !result.IsTruncated {
			break
		}

		marker = result.NextMarker
	}

	return allFiles, nil
}

// FileExists 检查文件是否存在
func (c *OSSClient) FileExists(filename string) (bool, error) {
	key := c.buildKey(filename)

	exists, err := c.bucket.IsObjectExist(key)
	if err != nil {
		return false, fmt.Errorf("检查文件是否存在失败: %v", err)
	}

	return exists, nil
}

// CopyResult 复制结果
type CopyResult struct {
	SourceKey      string    // 源文件键名
	DestinationKey string    // 目标文件键名
	CopyTime       time.Time // 复制时间
	ETag           string    // 新文件的ETag
}

// CopyFile 复制文件
func (c *OSSClient) CopyFile(sourceKey string, destinationKey string) (*CopyResult, error) {
	// 构建完整的键名
	srcKey := c.buildKey(sourceKey)
	dstKey := c.buildKey(destinationKey)

	// 执行复制
	_, err := c.bucket.CopyObject(srcKey, dstKey)
	if err != nil {
		return nil, fmt.Errorf("复制文件失败: %v", err)
	}

	// 获取新文件信息
	props, err := c.bucket.GetObjectDetailedMeta(dstKey)
	if err != nil {
		return nil, fmt.Errorf("获取复制后文件信息失败: %v", err)
	}

	return &CopyResult{
		SourceKey:      srcKey,
		DestinationKey: dstKey,
		CopyTime:       time.Now(),
		ETag:           strings.Trim(props.Get("ETag"), "\""),
	}, nil
}

// MoveFile 移动文件（复制后删除源文件）
func (c *OSSClient) MoveFile(sourceKey string, destinationKey string) (*CopyResult, error) {
	// 先复制文件
	copyResult, err := c.CopyFile(sourceKey, destinationKey)
	if err != nil {
		return nil, err
	}

	// 删除源文件
	if err := c.DeleteFile(sourceKey); err != nil {
		// 如果删除失败，尝试删除已复制的文件
		_ = c.DeleteFile(destinationKey)
		return nil, fmt.Errorf("移动文件失败，删除源文件时出错: %v", err)
	}

	return copyResult, nil
}

// GetSignedURL 获取文件的签名URL（用于私有读取）
func (c *OSSClient) GetSignedURL(filename string, expireSeconds int64) (string, error) {
	key := c.buildKey(filename)

	signedURL, err := c.bucket.SignURL(key, oss.HTTPGet, expireSeconds)
	if err != nil {
		return "", fmt.Errorf("生成签名URL失败: %v", err)
	}

	return signedURL, nil
}

// getContentType 根据文件扩展名获取ContentType
func getContentType(ext string) string {
	ext = strings.ToLower(ext)
	contentTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".bmp":  "image/bmp",
		".ico":  "image/x-icon",
		".svg":  "image/svg+xml",
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".ppt":  "application/vnd.ms-powerpoint",
		".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		".txt":  "text/plain",
		".html": "text/html",
		".css":  "text/css",
		".js":   "application/javascript",
		".json": "application/json",
		".xml":  "application/xml",
		".zip":  "application/zip",
		".rar":  "application/x-rar-compressed",
		".7z":   "application/x-7z-compressed",
		".tar":  "application/x-tar",
		".gz":   "application/gzip",
		".mp3":  "audio/mpeg",
		".wav":  "audio/wav",
		".mp4":  "video/mp4",
		".avi":  "video/x-msvideo",
		".mov":  "video/quicktime",
		".wmv":  "video/x-ms-wmv",
		".flv":  "video/x-flv",
	}

	if contentType, exists := contentTypes[ext]; exists {
		return contentType
	}

	return "application/octet-stream"
}
