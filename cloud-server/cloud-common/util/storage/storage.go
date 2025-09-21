package storage

import (
	"io"
	"mime/multipart"
	"time"
)

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

// ListFilesResult 文件列表结果
type ListFilesResult struct {
	Files       []*FileInfo `json:"files"`       // 文件列表
	NextMarker  string      `json:"nextMarker"`  // 下一页标记
	IsTruncated bool        `json:"isTruncated"` // 是否还有更多数据
	TotalCount  int         `json:"totalCount"`  // 总数量（如果支持）
}

// CopyResult 复制结果
type CopyResult struct {
	SourceKey      string    `json:"sourceKey"`      // 源文件键名
	DestinationKey string    `json:"destinationKey"` // 目标文件键名
	CopyTime       time.Time `json:"copyTime"`       // 复制时间
	ETag           string    `json:"etag"`           // 新文件的ETag
}

// Storage 云存储接口，定义了统一的云存储操作
type Storage interface {
	// 文件上传
	UploadFile(filename string, file io.Reader, contentType string) (*UploadResult, error)
	UploadFileFromPath(filename string, localPath string) (*UploadResult, error)
	UploadFileFromMultipart(filename string, file multipart.File, header *multipart.FileHeader) (*UploadResult, error)

	// 分片上传
	MultipartUpload(filename string, file io.Reader, fileSize int64, partSize int64) (*MultipartUploadResult, error)

	// 文件删除
	DeleteFile(filename string) error
	DeleteFiles(filenames []string) error

	// 文件信息
	GetFileInfo(filename string) (*FileInfo, error)
	FileExists(filename string) (bool, error)

	// 文件列表
	ListFiles(prefix string, maxKeys int) ([]*FileInfo, error)
	// 分页获取文件列表
	ListFilesWithPage(prefix string, marker string, maxKeys int) (*ListFilesResult, error)
	// 获取所有文件（自动处理分页）
	ListAllFiles(prefix string) ([]*FileInfo, error)

	// 文件操作
	// 复制文件
	CopyFile(sourceKey string, destinationKey string) (*CopyResult, error)
	// 移动文件（复制后删除源文件）
	MoveFile(sourceKey string, destinationKey string) (*CopyResult, error)

	// URL生成
	GetSignedURL(filename string, expireSeconds int64) (string, error)
}
