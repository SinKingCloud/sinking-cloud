package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

type Redis struct {
	Client *redis.Client
}

// NewRedis 实例化redis
func NewRedis(host string, port string, pwd string, db int64) (*Redis, error) {
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
	})

	// 测试连接
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("Redis连接失败: %v", err)
	}

	// 清空db,防止污染系统
	client.FlushDB(ctx)

	RedisClient = client
	return &Redis{Client: client}, nil
}

// Get 获取缓存
func (client *Redis) Get(key string) (string, bool) {
	ctx := context.Background()
	reply, err := client.Client.Get(ctx, key).Result()
	if err != nil || reply == "" {
		return "", false
	}
	return reply, true
}

// SetWithExpire 设置缓存(带过期时间)
func (client *Redis) SetWithExpire(key string, value string, ttl int64) bool {
	ctx := context.Background()
	expiration := time.Duration(ttl) * time.Millisecond
	err := client.Client.Set(ctx, key, value, expiration).Err()
	return err == nil
}

// Set 设置缓存(不过期)
func (client *Redis) Set(key string, value string) bool {
	ctx := context.Background()
	err := client.Client.Set(ctx, key, value, 0).Err()
	return err == nil
}

// SetExpire 设置过期时间
func (client *Redis) SetExpire(key string, ttl int64) bool {
	ctx := context.Background()
	expiration := time.Duration(ttl) * time.Millisecond
	result := client.Client.PExpire(ctx, key, expiration).Val()
	return result
}

// Lock 上锁
func (client *Redis) Lock(key string, ttl int64) bool {
	ctx := context.Background()
	expiration := time.Duration(ttl) * time.Millisecond
	result := client.Client.SetNX(ctx, key, 1, expiration).Val()
	return result
}

// UnLock 解锁
func (client *Redis) UnLock(key string) bool {
	return client.Delete(key)
}

// IsLock 是否lock状态
func (client *Redis) IsLock(key string) bool {
	ctx := context.Background()
	reply, err := client.Client.Get(ctx, key).Int()
	if err == nil && reply == 1 {
		return true
	}
	return false
}

// Remember 获取缓存,不存在则设置
func (client *Redis) Remember(key string, fun func() string, ttl int64) string {
	value, exists := client.Get(key)
	if !exists {
		value = fun()
		client.SetWithExpire(key, value, ttl)
	}
	return value
}

// Delete 删除缓存
func (client *Redis) Delete(key string) bool {
	ctx := context.Background()
	err := client.Client.Del(ctx, key).Err()
	return err == nil
}

// Close 关闭连接
func (client *Redis) Close() error {
	return client.Client.Close()
}

// Inc 自增
func (client *Redis) Inc(key string) int {
	ctx := context.Background()
	reply, err := client.Client.Incr(ctx, key).Result()
	if err != nil {
		return 0
	}
	return int(reply)
}

// Dec 自减
func (client *Redis) Dec(key string) int {
	ctx := context.Background()
	reply, err := client.Client.Decr(ctx, key).Result()
	if err != nil {
		return 0
	}
	return int(reply)
}
