package storage

import (
	"context"
	"fmt"
	"time"
)

// Delete 删除单个对象。
func (c *Client) Delete(ctx context.Context, key string) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	return c.driver.delete(ctx, key)
}

// DeleteMany 批量删除对象。
func (c *Client) DeleteMany(ctx context.Context, keys []string) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}

	if c.capabilities().batchDelete {
		return c.driver.deleteMany(ctx, keys)
	}

	for _, key := range keys {
		if err := c.driver.delete(ctx, key); err != nil {
			return err
		}
	}
	return nil
}

// Stat 获取对象元信息。
func (c *Client) Stat(ctx context.Context, key string) (*ObjectInfo, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	return c.driver.stat(ctx, key)
}

// Exists 判断对象是否存在。
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	if err := checkContext(ctx); err != nil {
		return false, err
	}
	return c.driver.exists(ctx, key)
}

// List 按前缀分页列出对象。
func (c *Client) List(ctx context.Context, request ListRequest) (*ListResult, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	return c.driver.list(ctx, request)
}

// Copy 复制对象。
func (c *Client) Copy(ctx context.Context, sourceKey, destinationKey string) (*CopyResult, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if !c.capabilities().copyObject {
		return nil, errCopyObjectNotSupported
	}

	return c.driver.copy(ctx, sourceKey, destinationKey)
}

// Move 移动对象，默认通过复制后删除实现。
func (c *Client) Move(ctx context.Context, sourceKey, destinationKey string) (*CopyResult, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if !c.capabilities().copyObject {
		return nil, errCopyObjectNotSupported
	}

	result, err := c.driver.copy(ctx, sourceKey, destinationKey)
	if err != nil {
		return nil, err
	}

	if err := c.driver.delete(ctx, sourceKey); err != nil {
		_ = c.driver.delete(ctx, destinationKey)
		return nil, fmt.Errorf("delete source object failed after copy: %w", err)
	}

	return result, nil
}

// SignedURL 生成对象的临时访问地址。
func (c *Client) SignedURL(ctx context.Context, key string, expire time.Duration) (string, error) {
	if err := checkContext(ctx); err != nil {
		return "", err
	}
	if !c.capabilities().signedURL {
		return "", errSignedURLNotSupported
	}

	return c.driver.signedURL(ctx, key, ensurePositiveExpire(expire))
}
