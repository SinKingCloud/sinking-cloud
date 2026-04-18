package storage

// Client 是统一的云对象存储客户端。
type Client struct {
	driver         driver // 是当前客户端绑定的底层厂商实现。
	uploadIDSecret []byte // 用于对上传任务 ID 做签名校验。
}

// New 根据统一配置创建云对象存储客户端。
func New(config Config) (*Client, error) {
	driver, err := newDriver(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		driver:         driver,
		uploadIDSecret: buildUploadIDSecret(config),
	}, nil
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

// capabilities 返回当前驱动声明的能力集合。
func (c *Client) capabilities() capabilities {
	return c.driver.capabilities()
}

// encodeUploadID 将内部上传会话编码成可对外传递的任务 ID。
func (c *Client) encodeUploadID(session uploadSession) (string, error) {
	return encodeUploadID(c.uploadIDSecret, session)
}

// decodeUploadID 将外部传入的任务 ID 还原为内部上传会话。
func (c *Client) decodeUploadID(uploadID string) (uploadSession, error) {
	return decodeUploadID(c.uploadIDSecret, uploadID)
}
