package archive

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 定义支持的压缩格式。
const (
	FormatZip      = "zip"
	FormatTar      = "tar"
	FormatGzip     = "gz"
	FormatTgz      = "tgz"
	copyBufferSize = 32 * 1024
)

// Service 提供压缩与解压缩能力。
type Service struct{}

// New 创建一个压缩服务实例。
func New() *Service {
	return &Service{}
}

// GetFormatByExt 根据文件名后缀推断压缩格式。
func (s *Service) GetFormatByExt(filename string) string {
	lowerName := strings.ToLower(filename)
	ext := strings.ToLower(filepath.Ext(lowerName))
	if ext == "" {
		return ""
	}

	ext = ext[1:] // 移除点号
	switch ext {
	case "zip":
		return FormatZip
	case "tar":
		return FormatTar
	case "gz", "gzip":
		if strings.HasSuffix(lowerName, ".tar.gz") || strings.HasSuffix(lowerName, ".tar.gzip") {
			return FormatTgz
		}
		return FormatGzip
	case "tgz":
		return FormatTgz
	default:
		return ""
	}
}

// IsSupportedFormat 判断给定格式是否受支持。
func (s *Service) IsSupportedFormat(format string) bool {
	switch format {
	case FormatZip, FormatTar, FormatGzip, FormatTgz:
		return true
	default:
		return false
	}
}

// Compress 压缩文件或目录到目标归档文件。
// ctx: 上下文，用于取消操作
// srcPaths: 要压缩的源文件或目录路径列表
// destPath: 目标压缩文件路径
// format: 压缩格式（zip、tar、gz、tgz）
// callback: 进度回调函数，参数依次为已处理字节数、总字节数、当前文件、总文件数、当前索引
// 返回 false 表示取消操作
func (s *Service) Compress(ctx context.Context, srcPaths []string, destPath, format string, callback func(current, total int64, currentFile string, totalFiles, currentIndex int64) bool) (err error) {
	if !s.IsSupportedFormat(format) {
		return errors.New("不支持的压缩格式")
	}

	if err := s.checkContext(ctx); err != nil {
		return err
	}

	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	if err := s.checkContext(ctx); err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = os.Remove(destPath)
		}
	}()

	switch format {
	case FormatZip:
		err = s.compressZip(ctx, srcPaths, destPath, callback)
	case FormatTar:
		err = s.compressTar(ctx, srcPaths, destPath, false, callback)
	case FormatGzip:
		if len(srcPaths) != 1 {
			return errors.New("gzip格式只支持压缩单个文件")
		}
		err = s.compressGzip(ctx, srcPaths[0], destPath)
	case FormatTgz:
		err = s.compressTar(ctx, srcPaths, destPath, true, callback)
	default:
		err = errors.New("不支持的压缩格式")
	}

	return err
}

// Extract 解压指定归档文件到目标目录。
// ctx: 上下文，用于取消操作
// srcPath: 要解压的源文件路径
// destPath: 解压目标目录
// callback: 进度回调函数，参数依次为已处理字节数、总字节数、当前文件、总文件数、当前索引
// 返回 false 表示取消操作
func (s *Service) Extract(ctx context.Context, srcPath, destPath string, callback func(current, total int64, currentFile string, totalFiles, currentIndex int64) bool) error {
	format := s.GetFormatByExt(srcPath)
	if !s.IsSupportedFormat(format) {
		return errors.New("不支持的压缩格式")
	}

	if err := s.checkContext(ctx); err != nil {
		return err
	}
	ctxErr := s.checkContext(ctx)

	destExists := true
	if _, err := os.Lstat(destPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			destExists = false
		} else {
			return err
		}
	}

	if err := os.MkdirAll(destPath, 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}
	if err := s.ensureNoSymlinkPath(destPath, destPath); err != nil {
		return err
	}

	if err := s.checkContext(ctx); err != nil {
		return err
	}

	var err error
	switch format {
	case FormatZip:
		err = s.extractZip(ctx, srcPath, destPath, callback)
	case FormatTar:
		err = s.extractTar(ctx, srcPath, destPath, false, callback)
	case FormatGzip:
		err = s.extractGzip(ctx, srcPath, destPath)
	case FormatTgz:
		err = s.extractTar(ctx, srcPath, destPath, true, callback)
	default:
		err = errors.New("不支持的压缩格式")
	}

	if err != nil && !errors.Is(err, ctxErr) && !destExists {
		_ = os.RemoveAll(destPath)
	}

	return err
}

// checkContext 检查上下文是否已经取消。
func (s *Service) checkContext(ctx context.Context) error {
	if ctx == nil {
		return nil
	}
	return ctx.Err()
}

// copyWithContext 在复制流时周期性检查上下文取消状态。
func (s *Service) copyWithContext(ctx context.Context, dst io.Writer, src io.Reader) (int64, error) {
	buf := make([]byte, copyBufferSize)
	var written int64

	for {
		if err := s.checkContext(ctx); err != nil {
			return written, err
		}

		nr, readErr := src.Read(buf)
		if nr > 0 {
			if err := s.checkContext(ctx); err != nil {
				return written, err
			}

			nw, writeErr := dst.Write(buf[:nr])
			written += int64(nw)
			if writeErr != nil {
				return written, writeErr
			}
			if nw != nr {
				return written, io.ErrShortWrite
			}
		}

		if readErr == io.EOF {
			return written, nil
		}
		if readErr != nil {
			return written, readErr
		}
	}
}

// copyFileToWriter 将本地文件内容写入目标 writer。
func (s *Service) copyFileToWriter(ctx context.Context, srcPath string, dst io.Writer) (int64, error) {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return 0, err
	}

	written, copyErr := s.copyWithContext(ctx, dst, srcFile)
	closeErr := srcFile.Close()
	if copyErr != nil {
		return written, copyErr
	}
	if closeErr != nil {
		return written, closeErr
	}

	return written, nil
}

// writeExtractedFile 将归档中的文件内容写入解压目标文件。
func (s *Service) writeExtractedFile(ctx context.Context, destPath string, perm os.FileMode, src io.Reader) (int64, error) {
	destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return 0, err
	}

	written, copyErr := s.copyWithContext(ctx, destFile, src)
	closeErr := destFile.Close()
	if copyErr != nil {
		return written, copyErr
	}
	if closeErr != nil {
		return written, closeErr
	}

	return written, nil
}

// collectSourceStats 统计源路径中的文件总数和总大小。
func (s *Service) collectSourceStats(ctx context.Context, paths []string) (int64, int64, error) {
	var totalFiles int64
	var totalSize int64

	for _, path := range paths {
		err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if err := s.checkContext(ctx); err != nil {
				return err
			}

			totalFiles++
			if !info.IsDir() {
				totalSize += info.Size()
			}
			return nil
		})
		if err != nil {
			return 0, 0, err
		}
	}

	return totalFiles, totalSize, nil
}

// buildArchiveEntryPath 生成文件在归档内部使用的相对路径。
func (s *Service) buildArchiveEntryPath(srcPath, currentPath string) (string, error) {
	relPath, err := filepath.Rel(filepath.Dir(srcPath), currentPath)
	if err != nil {
		return "", err
	}
	return filepath.ToSlash(relPath), nil
}

// resolveExtractPath 计算归档条目在目标目录中的实际落盘路径。
func (s *Service) resolveExtractPath(destRoot, entryName string) (string, error) {
	cleanDestRoot := filepath.Clean(destRoot)
	targetPath := filepath.Join(cleanDestRoot, filepath.Clean(entryName))

	relPath, err := filepath.Rel(cleanDestRoot, targetPath)
	if err != nil {
		return "", err
	}
	if relPath == ".." || strings.HasPrefix(relPath, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("非法的文件路径: %s", entryName)
	}

	return targetPath, nil
}

// ensureNoSymlinkPath 检查目标路径链路中是否包含符号链接。
func (s *Service) ensureNoSymlinkPath(destRoot, targetPath string) error {
	cleanDestRoot := filepath.Clean(destRoot)

	rootInfo, err := os.Lstat(cleanDestRoot)
	if err == nil && rootInfo.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("目标目录不能是符号链接: %s", cleanDestRoot)
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	relPath, err := filepath.Rel(cleanDestRoot, targetPath)
	if err != nil {
		return err
	}
	if relPath == "." {
		return nil
	}

	currentPath := cleanDestRoot
	for _, part := range strings.Split(relPath, string(os.PathSeparator)) {
		currentPath = filepath.Join(currentPath, part)

		info, statErr := os.Lstat(currentPath)
		if statErr != nil {
			if errors.Is(statErr, os.ErrNotExist) {
				return nil
			}
			return statErr
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("路径包含符号链接: %s", currentPath)
		}
	}

	return nil
}

// openTarStream 打开 tar 或 tar.gz 文件并返回对应的读取流。
func (s *Service) openTarStream(srcPath string, isGzipped bool) (*os.File, *tar.Reader, io.Closer, error) {
	file, err := os.Open(srcPath)
	if err != nil {
		return nil, nil, nil, err
	}

	if !isGzipped {
		return file, tar.NewReader(file), nil, nil
	}

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		_ = file.Close()
		return nil, nil, nil, err
	}

	return file, tar.NewReader(gzipReader), gzipReader, nil
}

// closeTarStream 关闭 tar 读取过程中创建的文件和压缩流。
func (s *Service) closeTarStream(file *os.File, streamCloser io.Closer) {
	if streamCloser != nil {
		_ = streamCloser.Close()
	}
	if file != nil {
		_ = file.Close()
	}
}

// scanTarArchive 预扫描 tar 包，统计条目数和展开后的总大小。
func (s *Service) scanTarArchive(ctx context.Context, srcPath string, isGzipped bool) (int64, int64, error) {
	file, tarReader, streamCloser, err := s.openTarStream(srcPath, isGzipped)
	if err != nil {
		return 0, 0, err
	}
	defer s.closeTarStream(file, streamCloser)

	var totalFiles int64
	var totalSize int64
	for {
		if err := s.checkContext(ctx); err != nil {
			return 0, 0, err
		}

		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, 0, err
		}

		totalFiles++
		if header.FileInfo().Mode().IsRegular() {
			totalSize += header.Size
		}
	}

	return totalFiles, totalSize, nil
}

// compressZip 将文件或目录压缩为 ZIP 文件。
func (s *Service) compressZip(ctx context.Context, srcPaths []string, destPath string, callback func(current, total int64, currentFile string, totalFiles, currentIndex int64) bool) (err error) {
	if err = s.checkContext(ctx); err != nil {
		return err
	}

	totalFiles, totalSize, err := s.collectSourceStats(ctx, srcPaths)
	if err != nil {
		return err
	}
	if err = s.checkContext(ctx); err != nil {
		return err
	}

	zipFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := zipFile.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	zipWriter := zip.NewWriter(zipFile)
	defer func() {
		if closeErr := zipWriter.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	var processedSize int64
	fileIndex := int64(0)
	for _, srcPath := range srcPaths {
		if err = s.checkContext(ctx); err != nil {
			return err
		}

		err = filepath.Walk(srcPath, func(path string, info os.FileInfo, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if err = s.checkContext(ctx); err != nil {
				return err
			}

			if callback != nil && !callback(processedSize, totalSize, path, totalFiles, fileIndex) {
				return errors.New("操作被取消")
			}

			relPath, err2 := s.buildArchiveEntryPath(srcPath, path)
			if err2 != nil {
				return err2
			}

			if info.IsDir() {
				_, err = zipWriter.Create(relPath + "/")
				fileIndex++
				return err
			}

			fileHeader, err2 := zip.FileInfoHeader(info)
			if err2 != nil {
				return err2
			}
			fileHeader.Name = relPath
			fileHeader.Method = zip.Deflate

			writer, err2 := zipWriter.CreateHeader(fileHeader)
			if err2 != nil {
				return err2
			}

			written, err2 := s.copyFileToWriter(ctx, path, writer)
			if err2 != nil {
				return err2
			}

			processedSize += written
			fileIndex++

			if callback != nil && !callback(processedSize, totalSize, path, totalFiles, fileIndex-1) {
				return errors.New("操作被取消")
			}

			return s.checkContext(ctx)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// extractZip 将 ZIP 文件解压到目标目录。
func (s *Service) extractZip(ctx context.Context, srcPath, destPath string, callback func(current, total int64, currentFile string, totalFiles, currentIndex int64) bool) error {
	if err := s.checkContext(ctx); err != nil {
		return err
	}

	reader, err := zip.OpenReader(srcPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = reader.Close()
	}()

	totalFiles := int64(len(reader.File))
	var totalSize int64
	for _, file := range reader.File {
		if !file.FileInfo().IsDir() {
			totalSize += int64(file.UncompressedSize64)
		}
	}
	if err = s.checkContext(ctx); err != nil {
		return err
	}

	var processedSize int64
	for i, file := range reader.File {
		if err = s.checkContext(ctx); err != nil {
			return err
		}

		if callback != nil && !callback(processedSize, totalSize, file.Name, totalFiles, int64(i)) {
			return errors.New("操作被取消")
		}

		if file.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("不支持解压符号链接: %s", file.Name)
		}

		destFilePath, err2 := s.resolveExtractPath(destPath, file.Name)
		if err2 != nil {
			return err2
		}
		if err = s.ensureNoSymlinkPath(destPath, destFilePath); err != nil {
			return err
		}

		if file.FileInfo().IsDir() {
			if err = os.MkdirAll(destFilePath, file.Mode()); err != nil {
				return err
			}
			continue
		}

		if err = os.MkdirAll(filepath.Dir(destFilePath), 0755); err != nil {
			return err
		}

		destFile, err2 := os.OpenFile(destFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err2 != nil {
			return err2
		}

		srcFile, err2 := file.Open()
		if err2 != nil {
			_ = destFile.Close()
			return err2
		}

		written, err2 := s.copyWithContext(ctx, destFile, srcFile)
		srcCloseErr := srcFile.Close()
		destCloseErr := destFile.Close()
		if err2 != nil {
			return err2
		}
		if srcCloseErr != nil {
			return srcCloseErr
		}
		if destCloseErr != nil {
			return destCloseErr
		}

		processedSize += written

		if callback != nil && !callback(processedSize, totalSize, file.Name, totalFiles, int64(i)) {
			return errors.New("操作被取消")
		}
		if err = s.checkContext(ctx); err != nil {
			return err
		}
	}

	return nil
}

// compressTar 将文件或目录压缩为 TAR 或 TGZ 文件。
func (s *Service) compressTar(ctx context.Context, srcPaths []string, destPath string, useGzip bool, callback func(current, total int64, currentFile string, totalFiles, currentIndex int64) bool) (err error) {
	if err = s.checkContext(ctx); err != nil {
		return err
	}

	totalFiles, totalSize, err := s.collectSourceStats(ctx, srcPaths)
	if err != nil {
		return err
	}
	if err = s.checkContext(ctx); err != nil {
		return err
	}

	tarFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := tarFile.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	var tarWriter *tar.Writer
	if useGzip {
		gzipWriter := gzip.NewWriter(tarFile)
		defer func() {
			if closeErr := gzipWriter.Close(); err == nil && closeErr != nil {
				err = closeErr
			}
		}()
		tarWriter = tar.NewWriter(gzipWriter)
	} else {
		tarWriter = tar.NewWriter(tarFile)
	}
	defer func() {
		if closeErr := tarWriter.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	var processedSize int64
	fileIndex := int64(0)
	for _, srcPath := range srcPaths {
		if err = s.checkContext(ctx); err != nil {
			return err
		}

		err = filepath.Walk(srcPath, func(path string, info os.FileInfo, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if err = s.checkContext(ctx); err != nil {
				return err
			}

			if callback != nil && !callback(processedSize, totalSize, path, totalFiles, fileIndex) {
				return errors.New("操作被取消")
			}

			relPath, err2 := s.buildArchiveEntryPath(srcPath, path)
			if err2 != nil {
				return err2
			}

			header, err2 := tar.FileInfoHeader(info, "")
			if err2 != nil {
				return err2
			}
			header.Name = relPath

			if err = tarWriter.WriteHeader(header); err != nil {
				return err
			}

			if info.Mode().IsRegular() {
				written, err3 := s.copyFileToWriter(ctx, path, tarWriter)
				if err3 != nil {
					return err3
				}
				processedSize += written
			}

			fileIndex++

			if callback != nil && info.Mode().IsRegular() && !callback(processedSize, totalSize, path, totalFiles, fileIndex-1) {
				return errors.New("操作被取消")
			}

			return s.checkContext(ctx)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// extractTar 将 TAR 或 TGZ 文件解压到目标目录。
func (s *Service) extractTar(ctx context.Context, srcPath, destPath string, isGzipped bool, callback func(current, total int64, currentFile string, totalFiles, currentIndex int64) bool) error {
	if err := s.checkContext(ctx); err != nil {
		return err
	}

	totalFiles, totalSize, err := s.scanTarArchive(ctx, srcPath, isGzipped)
	if err != nil {
		return err
	}

	file, tarReader, streamCloser, err := s.openTarStream(srcPath, isGzipped)
	if err != nil {
		return err
	}
	defer s.closeTarStream(file, streamCloser)

	var processedSize int64
	for i := int64(0); i < totalFiles; i++ {
		if err = s.checkContext(ctx); err != nil {
			return err
		}

		header, err2 := tarReader.Next()
		if err2 == io.EOF {
			break
		}
		if err2 != nil {
			return err2
		}

		if callback != nil && !callback(processedSize, totalSize, header.Name, totalFiles, i) {
			return errors.New("操作被取消")
		}

		destFilePath, err2 := s.resolveExtractPath(destPath, header.Name)
		if err2 != nil {
			return err2
		}
		if err = s.ensureNoSymlinkPath(destPath, destFilePath); err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			dirMode := os.FileMode(header.Mode)
			if dirMode == 0 {
				dirMode = 0755
			}
			if err = os.MkdirAll(destFilePath, dirMode); err != nil {
				return err
			}
		case tar.TypeSymlink, tar.TypeLink:
			return fmt.Errorf("不支持解压链接文件: %s", header.Name)
		case tar.TypeChar, tar.TypeBlock, tar.TypeFifo:
			return fmt.Errorf("不支持解压特殊文件: %s", header.Name)
		case tar.TypeReg:
			if err = os.MkdirAll(filepath.Dir(destFilePath), 0755); err != nil {
				return err
			}

			written, err3 := s.writeExtractedFile(ctx, destFilePath, os.FileMode(header.Mode), tarReader)
			if err3 != nil {
				return err3
			}
			processedSize += written
		default:
			if header.FileInfo().Mode().IsRegular() {
				if err = os.MkdirAll(filepath.Dir(destFilePath), 0755); err != nil {
					return err
				}

				written, err3 := s.writeExtractedFile(ctx, destFilePath, os.FileMode(header.Mode), tarReader)
				if err3 != nil {
					return err3
				}
				processedSize += written
			}
		}

		if callback != nil && !callback(processedSize, totalSize, header.Name, totalFiles, i) {
			return errors.New("操作被取消")
		}
	}

	return nil
}

// compressGzip 将单个文件压缩为 GZIP 文件。
func (s *Service) compressGzip(ctx context.Context, srcPath string, destPath string) (err error) {
	if err = s.checkContext(ctx); err != nil {
		return err
	}

	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return err
	}
	if srcInfo.IsDir() {
		return errors.New("不能直接用gzip压缩目录，请使用tgz格式")
	}

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := srcFile.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := destFile.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
		if err != nil {
			_ = os.Remove(destPath)
		}
	}()

	gzipWriter := gzip.NewWriter(destFile)
	defer func() {
		if closeErr := gzipWriter.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	gzipWriter.Name = filepath.Base(srcPath)
	_, err = s.copyWithContext(ctx, gzipWriter, srcFile)
	return err
}

// extractGzip 将 GZIP 文件解压到目标目录。
func (s *Service) extractGzip(ctx context.Context, srcPath string, destPath string) (err error) {
	if err = s.checkContext(ctx); err != nil {
		return err
	}

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := srcFile.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	gzipReader, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := gzipReader.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	entryName := gzipReader.Name
	if entryName == "" {
		entryName = filepath.Base(srcPath)
		if strings.HasSuffix(strings.ToLower(entryName), ".gz") {
			entryName = entryName[:len(entryName)-3]
		} else if strings.HasSuffix(strings.ToLower(entryName), ".gzip") {
			entryName = entryName[:len(entryName)-5]
		}
	}

	destFilePath, err := s.resolveExtractPath(destPath, entryName)
	if err != nil {
		return err
	}
	if err = s.ensureNoSymlinkPath(destPath, destFilePath); err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(destFilePath), 0755); err != nil {
		return err
	}

	destFile, err := os.Create(destFilePath)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := destFile.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
		if err != nil {
			_ = os.Remove(destFilePath)
		}
	}()

	_, err = s.copyWithContext(ctx, destFile, gzipReader)
	return err
}
