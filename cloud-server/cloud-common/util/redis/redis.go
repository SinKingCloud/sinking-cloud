package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

type Client struct {
	client *redis.Client
}

// NewClient 实例化redis
func NewClient(host string, port string, pwd string, db int64) (*Client, error) {
	address := fmt.Sprintf("%v:%v", host, port)
	client := redis.NewClient(&redis.Options{
		Addr:            address,
		Password:        pwd,
		DB:              int(db),
		PoolSize:        100,             // 连接池大小
		MinIdleConns:    10,              // 最小空闲连接数
		MaxIdleConns:    50,              // 最大空闲连接数
		ConnMaxLifetime: 5 * time.Minute, // 连接最大生存时间
		ConnMaxIdleTime: 1 * time.Minute, // 连接最大空闲时间
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled, // 显式禁用维护通知
		},
	})
	// 测试连接
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("Redis连接失败: %v", err)
	}
	// 清空db,防止污染系统
	client.FlushDB(ctx)
	return &Client{client: client}, nil
}

// Get 获取缓存
func (c *Client) Get(key string) (string, bool) {
	ctx := context.Background()
	reply, err := c.client.Get(ctx, key).Result()
	if err != nil || reply == "" {
		return "", false
	}
	return reply, true
}

// SetWithExpire 设置缓存(带过期时间)
func (c *Client) SetWithExpire(key string, value string, ttl int64) bool {
	ctx := context.Background()
	expiration := time.Duration(ttl) * time.Millisecond
	err := c.client.Set(ctx, key, value, expiration).Err()
	return err == nil
}

// Set 设置缓存(不过期)
func (c *Client) Set(key string, value string) bool {
	ctx := context.Background()
	err := c.client.Set(ctx, key, value, 0).Err()
	return err == nil
}

// SetExpire 设置过期时间
func (c *Client) SetExpire(key string, ttl int64) bool {
	ctx := context.Background()
	expiration := time.Duration(ttl) * time.Millisecond
	result := c.client.PExpire(ctx, key, expiration).Val()
	return result
}

// Lock 上锁
func (c *Client) Lock(key string, ttl int64) bool {
	ctx := context.Background()
	expiration := time.Duration(ttl) * time.Millisecond
	result := c.client.SetNX(ctx, key, 1, expiration).Val()
	return result
}

// UnLock 解锁
func (c *Client) UnLock(key string) bool {
	return c.Delete(key)
}

// IsLock 是否lock状态
func (c *Client) IsLock(key string) bool {
	ctx := context.Background()
	reply, err := c.client.Get(ctx, key).Int()
	if err == nil && reply == 1 {
		return true
	}
	return false
}

// Remember 获取缓存,不存在则设置
func (c *Client) Remember(key string, fun func() string, ttl int64) string {
	value, exists := c.Get(key)
	if !exists {
		value = fun()
		c.SetWithExpire(key, value, ttl)
	}
	return value
}

// Delete 删除缓存
func (c *Client) Delete(key string) bool {
	ctx := context.Background()
	err := c.client.Del(ctx, key).Err()
	return err == nil
}

// Close 关闭连接
func (c *Client) Close() error {
	return c.client.Close()
}

// Inc 自增
func (c *Client) Inc(key string) int {
	ctx := context.Background()
	reply, err := c.client.Incr(ctx, key).Result()
	if err != nil {
		return 0
	}
	return int(reply)
}

// Dec 自减
func (c *Client) Dec(key string) int {
	ctx := context.Background()
	reply, err := c.client.Decr(ctx, key).Result()
	if err != nil {
		return 0
	}
	return int(reply)
}

// GetClient 获取底层Redis客户端
func (c *Client) GetClient() *redis.Client {
	return c.client
}
