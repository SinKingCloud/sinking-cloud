package storage

// Client 是统一的云对象存储客户端。
type Client struct {
	driver driver
}

// New 根据统一配置创建云对象存储客户端。
func New(config Config) (*Client, error) {
	driver, err := newDriver(config)
	if err != nil {
		return nil, err
	}

	return &Client{driver: driver}, nil
}

// NewAliyunOSS 使用阿里云 OSS 配置创建客户端。
func NewAliyunOSS(config AliyunOSSConfig) (*Client, error) {
	return New(Config{
		Provider:        ProviderAliyunOSS,
		Endpoint:        config.Endpoint,
		Region:          config.Region,
		Bucket:          config.Bucket,
		AccessKeyID:     config.AccessKeyID,
		AccessKeySecret: config.AccessKeySecret,
		SecurityToken:   config.SecurityToken,
		Prefix:          config.Prefix,
		Domain:          config.Domain,
		UseHTTPS:        config.UseHTTPS,
	})
}

func (c *Client) capabilities() capabilities {
	return c.driver.capabilities()
}
