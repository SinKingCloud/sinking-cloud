package cache

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type Redis struct {
	Pool *redis.Pool
}

// Get 获取缓存
func (redis *Redis) Get(key string) (string, bool) {
	c1 := redis.Pool.Get()
	reply, err := c1.Do("get", key)
	defer redis.closeRedisClient(c1)
	if err != nil || reply == nil {
		return "", false
	}
	return string(reply.([]byte)), true
}

// SetWithExpire 设置缓存(带过期时间)
func (redis *Redis) SetWithExpire(key string, value string, ttl int64) bool {
	c1 := redis.Pool.Get()
	_, err := c1.Do("set", key, value)
	if ttl > 0 {
		_, _ = c1.Do("pexpire", key, ttl)
	}
	defer redis.closeRedisClient(c1)
	if err == nil {
		return true
	}
	return false
}

// Set 设置缓存(不过期)
func (redis *Redis) Set(key string, value string) bool {
	c1 := redis.Pool.Get()
	_, err := c1.Do("set", key, value)
	defer redis.closeRedisClient(c1)
	if err == nil {
		return true
	}
	return false
}

// Remember 获取缓存,不存在则设置
func (redis *Redis) Remember(key string, fun func() string, ttl int64) string {
	value, exists := redis.Get(key)
	if exists {
		return value
	}
	value = fun()
	if ttl > 0 {
		redis.SetWithExpire(key, value, ttl)
	}
	redis.Set(key, value)
	return value
}

// Delete 删除缓存
func (redis *Redis) Delete(key string) bool {
	c1 := redis.Pool.Get()
	res, err := c1.Do("del", key)
	fmt.Println(res, err)
	defer redis.closeRedisClient(c1)
	if err == nil {
		return true
	}
	return false
}

// closeRedisClient 关闭链接
func (redis *Redis) closeRedisClient(c1 redis.Conn) {
	err := c1.Close()
	if err != nil {
		return
	}
}
