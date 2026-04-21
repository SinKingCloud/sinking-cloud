package str

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
)

// ByteTool 字节流工具
type ByteTool struct {
	reader io.Reader
}

// NewByteTool 实例化工具类
func NewByteTool(reader io.Reader) *ByteTool {
	return &ByteTool{reader: reader}
}

// Md5 获取md5
func (s *ByteTool) Md5() (string, error) {
	if s == nil || s.reader == nil {
		return "", errors.New("reader不能为空")
	}
	hash := md5.New()
	buffer := make([]byte, 32*1024)
	if _, err := io.CopyBuffer(hash, s.reader, buffer); err != nil {
		return "", err
	}
	md5Str := hex.EncodeToString(hash.Sum(nil))
	return md5Str, nil
}
