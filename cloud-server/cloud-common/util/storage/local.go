package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/file"
)

// LocalConfig 本地存储配置
type LocalConfig struct {
	RootPath string `json:"rootPath"` // 本地存储根路径
	BaseURL  string `json:"baseUrl"`  // 访问基础URL
	BasePath string `json:"basePath"` // 基础路径前缀
}

// LocalClient 本地存储客户端
type LocalClient struct {
	disk   *file.Disk
	config *LocalConfig
}

// NewLocalClient 创建本地存储客户端
func NewLocalClient(config *LocalConfig) (Storage, error) {
	if config == nil {
		return nil, fmt.Errorf("本地存储配置不能为空")
	}

	if config.RootPath == "" {
		return nil, fmt.Errorf("本地存储根路径不能为空")
	}

	// 验证路径安全性
	absPath, err := filepath.Abs(config.RootPath)
	if err != nil {
		return nil, fmt.Errorf("无效的根路径: %v", err)
	}
	config.RootPath = absPath

	// 创建根目录（如果不存在）
	if err := os.MkdirAll(config.RootPath, 0755); err != nil {
		return nil, fmt.Errorf("创建根目录失败: %v", err)
	}

	// 测试目录是否可写
	testFile := filepath.Join(config.RootPath, ".write_test")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return nil, fmt.Errorf("目录不可写: %v", err)
	}
	os.Remove(testFile) // 清理测试文件

	// 创建文件操作实例
	disk := file.NewDisk(config.RootPath)

	return &LocalClient{
		disk:   disk,
		config: config,
	}, nil
}

// buildKey 构建文件键名
func (c *LocalClient) buildKey(filename string) string {
	// 防止路径遍历攻击
	filename = filepath.Clean(filename)
	if strings.Contains(filename, "..") {
		// 记录安全警告日志
		filename = strings.ReplaceAll(filename, "..", "")
	}

	key := filename
	if c.config.BasePath != "" {
		basePath := strings.TrimRight(c.config.BasePath, "/")
		key = basePath + "/" + strings.TrimLeft(filename, "/")
	}
	return key
}

// getFileURL 获取文件访问URL
func (c *LocalClient) getFileURL(key string) string {
	if c.config.BaseURL == "" {
		return key
	}

	baseURL := strings.TrimRight(c.config.BaseURL, "/")
	return baseURL + "/" + strings.TrimLeft(key, "/")
}

// UploadFile 简单上传文件
func (c *LocalClient) UploadFile(filename string, file io.Reader, contentType string) (*UploadResult, error) {
	key := c.buildKey(filename)

	// 确保目录存在
	dir := filepath.Dir(key)
	if dir != "." && dir != "/" {
		if err := c.disk.CreateDir(dir); err != nil {
			return nil, fmt.Errorf("创建目录失败: %v", err)
		}
	}

	// 创建文件
	localFile, err := c.disk.OpenFile(key, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %v", err)
	}
	defer localFile.Close()

	// 复制数据
	written, err := io.Copy(localFile, file)
	if err != nil {
		return nil, fmt.Errorf("写入文件失败: %v", err)
	}

	return &UploadResult{
		Key:         key,
		URL:         c.getFileURL(key),
		Size:        written,
		ContentType: contentType,
		ETag:        fmt.Sprintf("%x", time.Now().UnixNano()), // 简单的ETag
	}, nil
}

// UploadFileFromPath 从本地路径上传文件
func (c *LocalClient) UploadFileFromPath(filename string, localPath string) (*UploadResult, error) {
	open, err := os.Open(localPath)
	if err != nil {
		return nil, fmt.Errorf("打开本地文件失败: %v", err)
	}
	defer open.Close()

	// 获取文件信息
	fileInfo, err := open.Stat()
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %v", err)
	}

	// 根据文件扩展名获取ContentType
	contentType := getContentType(filepath.Ext(localPath))

	result, err := c.UploadFile(filename, open, contentType)
	if err != nil {
		return nil, err
	}

	result.Size = fileInfo.Size()
	return result, nil
}

// UploadFileFromMultipart 从multipart.File上传文件
func (c *LocalClient) UploadFileFromMultipart(filename string, file multipart.File, header *multipart.FileHeader) (*UploadResult, error) {
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

// MultipartUpload 分片上传文件（本地存储不需要分片，直接合并写入）
func (c *LocalClient) MultipartUpload(filename string, file io.Reader, fileSize int64, partSize int64) (*MultipartUploadResult, error) {
	// 本地存储不需要真正的分片上传，直接写入文件
	result, err := c.UploadFile(filename, file, "application/octet-stream")
	if err != nil {
		return nil, err
	}

	return &MultipartUploadResult{
		UploadResult: *result,
		UploadID:     fmt.Sprintf("local-%x", time.Now().UnixNano()),
		Parts:        1, // 本地存储只有一个"分片"
	}, nil
}

// DeleteFile 删除文件
func (c *LocalClient) DeleteFile(filename string) error {
	key := c.buildKey(filename)

	if !c.disk.Exists(key) {
		return fmt.Errorf("文件不存在: %s", key)
	}

	if err := c.disk.Delete(key); err != nil {
		return fmt.Errorf("删除文件失败: %v", err)
	}

	return nil
}

// DeleteFiles 批量删除文件
func (c *LocalClient) DeleteFiles(filenames []string) error {
	var errors []string

	for _, filename := range filenames {
		if err := c.DeleteFile(filename); err != nil {
			errors = append(errors, err.Error())
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("批量删除失败: %s", strings.Join(errors, "; "))
	}

	return nil
}

// GetFileInfo 获取文件信息
func (c *LocalClient) GetFileInfo(filename string) (*FileInfo, error) {
	key := c.buildKey(filename)

	fileInfo, err := c.disk.FileInfo(key)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %v", err)
	}

	return &FileInfo{
		Key:          key,
		Size:         fileInfo.Size,
		LastModified: time.Unix(fileInfo.UpdateTime, 0),
		ETag:         fmt.Sprintf("%x-%d", fileInfo.UpdateTime, fileInfo.Size),
		ContentType:  getContentType(filepath.Ext(filename)),
		URL:          c.getFileURL(key),
	}, nil
}

// ListFiles 获取文件列表
func (c *LocalClient) ListFiles(prefix string, maxKeys int) ([]*FileInfo, error) {
	result, err := c.ListFilesWithPage(prefix, "", maxKeys)
	if err != nil {
		return nil, err
	}
	return result.Files, nil
}

// ListFilesWithPage 分页获取文件列表
func (c *LocalClient) ListFilesWithPage(prefix string, marker string, maxKeys int) (*ListFilesResult, error) {
	if maxKeys <= 0 {
		maxKeys = 100
	}

	// 构建完整前缀路径
	fullPrefix := prefix
	if c.config.BasePath != "" {
		basePath := strings.TrimRight(c.config.BasePath, "/")
		if prefix != "" {
			fullPrefix = basePath + "/" + strings.TrimLeft(prefix, "/")
		} else {
			fullPrefix = basePath
		}
	}

	// 如果前缀为空，列出根目录
	if fullPrefix == "" {
		fullPrefix = "."
	}

	// 获取文件列表
	files, err := c.disk.FileList(fullPrefix)
	if err != nil {
		return nil, fmt.Errorf("获取文件列表失败: %v", err)
	}

	// 根据文件名排序，以确保分页的一致性
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name < files[j].Name
	})

	var result []*FileInfo
	skipMode := marker != ""
	count := 0
	nextMarker := ""
	isTruncated := false

	for _, fileInfo := range files {
		// 跳过目录
		if fileInfo.IsDir {
			continue
		}

		// 构建完整路径
		var fullPath string
		if fullPrefix == "." {
			fullPath = fileInfo.Name
		} else {
			fullPath = fullPrefix + "/" + fileInfo.Name
		}

		// 如果在跳过模式，跳过直到找到marker
		if skipMode {
			if fullPath <= marker {
				continue
			}
			skipMode = false
		}

		// 如果已经达到最大数量
		if count >= maxKeys {
			nextMarker = fullPath
			isTruncated = true
			break
		}

		result = append(result, &FileInfo{
			Key:          fullPath,
			Size:         fileInfo.Size,
			LastModified: time.Unix(fileInfo.UpdateTime, 0),
			ETag:         fmt.Sprintf("%x-%d", fileInfo.UpdateTime, fileInfo.Size),
			ContentType:  getContentType(filepath.Ext(fileInfo.Name)),
			URL:          c.getFileURL(fullPath),
		})

		count++
	}

	return &ListFilesResult{
		Files:       result,
		NextMarker:  nextMarker,
		IsTruncated: isTruncated,
		TotalCount:  0, // 本地存储不实时计算总数
	}, nil
}

// ListAllFiles 获取所有文件（自动处理分页）
func (c *LocalClient) ListAllFiles(prefix string) ([]*FileInfo, error) {
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
func (c *LocalClient) FileExists(filename string) (bool, error) {
	key := c.buildKey(filename)
	exists := c.disk.Exists(key) && c.disk.IsFile(key)
	return exists, nil
}

// GetSignedURL 获取文件的签名URL（本地存储直接返回URL）
func (c *LocalClient) GetSignedURL(filename string, expireSeconds int64) (string, error) {
	key := c.buildKey(filename)

	if !c.disk.Exists(key) {
		return "", fmt.Errorf("文件不存在: %s", key)
	}

	// 本地存储不需要签名，直接返回URL
	// 可以在URL中添加过期时间参数（需要业务层验证）
	baseURL := c.getFileURL(key)

	// 添加过期时间参数
	expireTime := time.Now().Add(time.Duration(expireSeconds) * time.Second).Unix()

	u, err := url.Parse(baseURL)
	if err != nil {
		return baseURL, nil // 如果解析失败，返回原始URL
	}

	query := u.Query()
	query.Set("expire", fmt.Sprintf("%d", expireTime))
	u.RawQuery = query.Encode()

	return u.String(), nil
}

// CopyFile 复制文件
func (c *LocalClient) CopyFile(sourceKey string, destinationKey string) (*CopyResult, error) {
	srcKey := c.buildKey(sourceKey)
	dstKey := c.buildKey(destinationKey)

	// 检查源文件是否存在
	if !c.disk.Exists(srcKey) || !c.disk.IsFile(srcKey) {
		return nil, fmt.Errorf("源文件不存在: %s", srcKey)
	}

	// 使用原生文件操作进行复制
	srcPath := c.disk.Path(srcKey, true, true) // 获取完整路径
	dstPath := c.disk.Path(dstKey, true, true)

	// 确保目标目录存在
	dstDir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return nil, fmt.Errorf("创建目标目录失败: %v", err)
	}

	// 打开源文件
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return nil, fmt.Errorf("打开源文件失败: %v", err)
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return nil, fmt.Errorf("创建目标文件失败: %v", err)
	}
	defer dstFile.Close()

	// 复制文件内容
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return nil, fmt.Errorf("复制文件内容失败: %v", err)
	}

	// 获取文件信息
	fileInfo, err := c.disk.FileInfo(dstKey)
	if err != nil {
		return nil, fmt.Errorf("获取目标文件信息失败: %v", err)
	}

	return &CopyResult{
		SourceKey:      srcKey,
		DestinationKey: dstKey,
		CopyTime:       time.Now(),
		ETag:           fmt.Sprintf("%x-%d", fileInfo.UpdateTime, fileInfo.Size),
	}, nil
}

// MoveFile 移动文件（复制后删除源文件）
func (c *LocalClient) MoveFile(sourceKey string, destinationKey string) (*CopyResult, error) {
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
