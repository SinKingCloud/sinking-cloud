package str

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"
)

// DateTime 自定义日期时间类型，适配JSON和数据库操作
type DateTime time.Time

// MarshalJSON 实现JSON序列化接口
func (t *DateTime) MarshalJSON() ([]byte, error) {
	// 处理零值情况，避免序列化出无效时间
	if time.Time(*t).IsZero() {
		return []byte("null"), nil
	}
	stamp := fmt.Sprintf("\"%s\"", time.Time(*t).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

// UnmarshalJSON 实现JSON反序列化接口
func (t *DateTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*t = DateTime(time.Time{}) // 赋值零值
		return nil
	}
	timeStr := strings.Trim(string(data), "\"")
	parsedTime, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return fmt.Errorf("解析时间失败: %w", err)
	}
	*t = DateTime(parsedTime)
	return nil
}

// Value 实现数据库写入接口（关键优化：值接收器）
func (t DateTime) Value() (driver.Value, error) {
	// 处理零值，避免写入无效时间
	if time.Time(t).IsZero() {
		return nil, nil
	}
	return time.Time(t).Format("2006-01-02 15:04:05"), nil
}

// Scan 实现数据库读取接口（保持指针接收器，需要修改接收者）
func (t *DateTime) Scan(v interface{}) error {
	switch vt := v.(type) {
	case nil:
		*t = DateTime(time.Time{})
		return nil
	case time.Time:
		*t = DateTime(vt)
		return nil
	case string:
		parsedTime, err := time.Parse("2006-01-02 15:04:05", vt)
		if err != nil {
			return fmt.Errorf("字符串时间解析失败: %w", err)
		}
		*t = DateTime(parsedTime)
		return nil
	case []byte:
		parsedTime, err := time.Parse("2006-01-02 15:04:05", string(vt))
		if err != nil {
			return fmt.Errorf("字节时间解析失败: %w", err)
		}
		*t = DateTime(parsedTime)
		return nil
	default:
		return errors.New("不支持的时间类型")
	}
}

// String 实现字符串输出
func (t *DateTime) String() string {
	if t == nil {
		return ""
	}
	return time.Time(*t).Format("2006-01-02 15:04:05")
}
