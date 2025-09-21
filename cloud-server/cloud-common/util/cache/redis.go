package cache

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

var RedisClient *redis.Pool

type Redis struct {
	Pool *redis.Pool
}

// NewRedis 实例化redis
func NewRedis(host string, port string, pwd string, db int64) *Redis {
	address := fmt.Sprintf("%v:%v", host, port)
	dbOption := redis.DialDatabase(int(db))
	pwdOption := redis.DialPassword(pwd)
	pool := &redis.Pool{
		//连接方法
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address, pwdOption, dbOption)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		//最大的激活连接数，表示同时最多有N个连接
		MaxActive: 0,
		//最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态。
		MaxIdle: 100,
		//最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭
		IdleTimeout:     0,
		Wait:            false,
		MaxConnLifetime: 0,
	}
	RedisClient = pool
	//清空db,防止污染系统
	_, _ = RedisClient.Get().Do("FLUSHDB")
	return &Redis{Pool: pool}
}

// Get 获取缓存
func (client *Redis) Get(key string) (string, bool) {
	c1 := client.Pool.Get()
	defer client.closeRedisClient(c1)
	reply, err := redis.String(c1.Do("get", key))
	if err != nil || reply == "" {
		return "", false
	}
	return reply, true
}

// SetWithExpire 设置缓存(带过期时间)
func (client *Redis) SetWithExpire(key string, value string, ttl int64) bool {
	c1 := client.Pool.Get()
	defer client.closeRedisClient(c1)
	_, err := c1.Do("set", key, value, "px", ttl)
	if err == nil {
		return true
	}
	return false
}

// Set 设置缓存(不过期)
func (client *Redis) Set(key string, value string) bool {
	c1 := client.Pool.Get()
	defer client.closeRedisClient(c1)
	_, err := c1.Do("set", key, value)
	if err == nil {
		return true
	}
	return false
}

// SetExpire 设置过期时间
func (client *Redis) SetExpire(key string, ttl int64) bool {
	c1 := client.Pool.Get()
	defer client.closeRedisClient(c1)
	num, err := redis.Int(c1.Do("pexpire", key, ttl))
	return err == nil && num > 0
}

// Lock 上锁
func (client *Redis) Lock(key string, ttl int64) bool {
	c1 := client.Pool.Get()
	defer client.closeRedisClient(c1)
	_, err := redis.String(c1.Do("set", key, 1, "px", ttl, "nx"))
	if err != nil {
		return false
	}
	return true
}

// UnLock 解锁
func (client *Redis) UnLock(key string) bool {
	return client.Delete(key)
}

// IsLock 是否lock状态
func (client *Redis) IsLock(key string) bool {
	c1 := client.Pool.Get()
	defer client.closeRedisClient(c1)
	reply, err := redis.Int(c1.Do("get", key))
	if err != nil && reply == 1 {
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
	c1 := client.Pool.Get()
	_, err := c1.Do("del", key)
	defer client.closeRedisClient(c1)
	if err == nil {
		return true
	}
	return false
}

// closeRedisClient 关闭链接
func (client *Redis) closeRedisClient(c1 redis.Conn) {
	err := c1.Close()
	if err != nil {
		return
	}
}

// Inc 自增
func (client *Redis) Inc(key string) int {
	c1 := client.Pool.Get()
	defer client.closeRedisClient(c1)
	reply, err := redis.Int(c1.Do("incr", key))
	if err != nil {
		return 0
	}
	return reply
}

// Dec 自增
func (client *Redis) Dec(key string) int {
	c1 := client.Pool.Get()
	defer client.closeRedisClient(c1)
	reply, err := redis.Int(c1.Do("decr", key))
	if err != nil {
		return 0
	}
	return reply
}
