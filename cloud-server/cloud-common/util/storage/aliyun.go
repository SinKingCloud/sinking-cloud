package storage

import (
	"fmt"
	"io"
	"mime/multipart"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/aliyun"
)

// AliyunOSSConfig 阿里云OSS配置
type AliyunOSSConfig struct {
	Endpoint        string // OSS地域节点
	AccessKeyID     string // AccessKey ID
	AccessKeySecret string // AccessKey Secret
	BucketName      string // 存储桶名称
	BasePath        string // 基础路径前缀
}

// AliyunOSSClient 阿里云OSS客户端
type AliyunOSSClient struct {
	ossClient *aliyun.OSSClient
}

// NewAliyunOSSClient 创建阿里云OSS存储客户端
func NewAliyunOSSClient(config *AliyunOSSConfig) (Storage, error) {
	if config == nil {
		return nil, fmt.Errorf("OSS配置不能为空")
	}

	// 转换为aliyun包的配置
	ossConfig := &aliyun.OSSConfig{
		Endpoint:        config.Endpoint,
		AccessKeyID:     config.AccessKeyID,
		AccessKeySecret: config.AccessKeySecret,
		BucketName:      config.BucketName,
		BasePath:        config.BasePath,
	}

	// 创建OSS客户端
	ossClient, err := aliyun.NewOSSClient(ossConfig)
	if err != nil {
		return nil, err
	}

	return &AliyunOSSClient{
		ossClient: ossClient,
	}, nil
}

// UploadFile 简单上传文件
func (c *AliyunOSSClient) UploadFile(filename string, file io.Reader, contentType string) (*UploadResult, error) {
	result, err := c.ossClient.UploadFile(filename, file, contentType)
	if err != nil {
		return nil, err
	}

	// 转换为storage包的结果类型
	return &UploadResult{
		Key:         result.Key,
		URL:         result.URL,
		Size:        result.Size,
		ContentType: result.ContentType,
		ETag:        result.ETag,
	}, nil
}

// UploadFileFromPath 从本地路径上传文件
func (c *AliyunOSSClient) UploadFileFromPath(filename string, localPath string) (*UploadResult, error) {
	result, err := c.ossClient.UploadFileFromPath(filename, localPath)
	if err != nil {
		return nil, err
	}

	// 转换为storage包的结果类型
	return &UploadResult{
		Key:         result.Key,
		URL:         result.URL,
		Size:        result.Size,
		ContentType: result.ContentType,
		ETag:        result.ETag,
	}, nil
}

// UploadFileFromMultipart 从multipart.File上传文件
func (c *AliyunOSSClient) UploadFileFromMultipart(filename string, file multipart.File, header *multipart.FileHeader) (*UploadResult, error) {
	result, err := c.ossClient.UploadFileFromMultipart(filename, file, header)
	if err != nil {
		return nil, err
	}

	// 转换为storage包的结果类型
	return &UploadResult{
		Key:         result.Key,
		URL:         result.URL,
		Size:        result.Size,
		ContentType: result.ContentType,
		ETag:        result.ETag,
	}, nil
}

// MultipartUpload 分片上传文件
func (c *AliyunOSSClient) MultipartUpload(filename string, file io.Reader, fileSize int64, partSize int64) (*MultipartUploadResult, error) {
	result, err := c.ossClient.MultipartUpload(filename, file, fileSize, partSize)
	if err != nil {
		return nil, err
	}

	// 转换为storage包的结果类型
	return &MultipartUploadResult{
		UploadResult: UploadResult{
			Key:         result.Key,
			URL:         result.URL,
			Size:        result.Size,
			ContentType: result.ContentType,
			ETag:        result.ETag,
		},
		UploadID: result.UploadID,
		Parts:    result.Parts,
	}, nil
}

// DeleteFile 删除文件
func (c *AliyunOSSClient) DeleteFile(filename string) error {
	return c.ossClient.DeleteFile(filename)
}

// DeleteFiles 批量删除文件
func (c *AliyunOSSClient) DeleteFiles(filenames []string) error {
	return c.ossClient.DeleteFiles(filenames)
}

// GetFileInfo 获取文件信息
func (c *AliyunOSSClient) GetFileInfo(filename string) (*FileInfo, error) {
	result, err := c.ossClient.GetFileInfo(filename)
	if err != nil {
		return nil, err
	}

	// 转换为storage包的结果类型
	return &FileInfo{
		Key:          result.Key,
		Size:         result.Size,
		LastModified: result.LastModified,
		ETag:         result.ETag,
		ContentType:  result.ContentType,
		URL:          result.URL,
	}, nil
}

// ListFiles 获取文件列表
func (c *AliyunOSSClient) ListFiles(prefix string, maxKeys int) ([]*FileInfo, error) {
	results, err := c.ossClient.ListFiles(prefix, maxKeys)
	if err != nil {
		return nil, err
	}

	// 转换为storage包的结果类型
	var files []*FileInfo
	for _, result := range results {
		files = append(files, &FileInfo{
			Key:          result.Key,
			Size:         result.Size,
			LastModified: result.LastModified,
			ETag:         result.ETag,
			ContentType:  result.ContentType,
			URL:          result.URL,
		})
	}

	return files, nil
}

// FileExists 检查文件是否存在
func (c *AliyunOSSClient) FileExists(filename string) (bool, error) {
	return c.ossClient.FileExists(filename)
}

// ListFilesWithPage 分页获取文件列表
func (c *AliyunOSSClient) ListFilesWithPage(prefix string, marker string, maxKeys int) (*ListFilesResult, error) {
	result, err := c.ossClient.ListFilesWithPage(prefix, marker, maxKeys)
	if err != nil {
		return nil, err
	}

	// 转换为storage包的结果类型
	var files []*FileInfo
	for _, file := range result.Files {
		files = append(files, &FileInfo{
			Key:          file.Key,
			Size:         file.Size,
			LastModified: file.LastModified,
			ETag:         file.ETag,
			ContentType:  file.ContentType,
			URL:          file.URL,
		})
	}

	return &ListFilesResult{
		Files:       files,
		NextMarker:  result.NextMarker,
		IsTruncated: result.IsTruncated,
		TotalCount:  result.TotalCount,
	}, nil
}

// ListAllFiles 获取所有文件（自动处理分页）
func (c *AliyunOSSClient) ListAllFiles(prefix string) ([]*FileInfo, error) {
	results, err := c.ossClient.ListAllFiles(prefix)
	if err != nil {
		return nil, err
	}

	// 转换为storage包的结果类型
	var files []*FileInfo
	for _, result := range results {
		files = append(files, &FileInfo{
			Key:          result.Key,
			Size:         result.Size,
			LastModified: result.LastModified,
			ETag:         result.ETag,
			ContentType:  result.ContentType,
			URL:          result.URL,
		})
	}

	return files, nil
}

// CopyFile 复制文件
func (c *AliyunOSSClient) CopyFile(sourceKey string, destinationKey string) (*CopyResult, error) {
	result, err := c.ossClient.CopyFile(sourceKey, destinationKey)
	if err != nil {
		return nil, err
	}

	// 转换为storage包的结果类型
	return &CopyResult{
		SourceKey:      result.SourceKey,
		DestinationKey: result.DestinationKey,
		CopyTime:       result.CopyTime,
		ETag:           result.ETag,
	}, nil
}

// MoveFile 移动文件（复制后删除源文件）
func (c *AliyunOSSClient) MoveFile(sourceKey string, destinationKey string) (*CopyResult, error) {
	result, err := c.ossClient.MoveFile(sourceKey, destinationKey)
	if err != nil {
		return nil, err
	}

	// 转换为storage包的结果类型
	return &CopyResult{
		SourceKey:      result.SourceKey,
		DestinationKey: result.DestinationKey,
		CopyTime:       result.CopyTime,
		ETag:           result.ETag,
	}, nil
}

// GetSignedURL 获取文件的签名URL（用于私有读取）
func (c *AliyunOSSClient) GetSignedURL(filename string, expireSeconds int64) (string, error) {
	return c.ossClient.GetSignedURL(filename, expireSeconds)
}
