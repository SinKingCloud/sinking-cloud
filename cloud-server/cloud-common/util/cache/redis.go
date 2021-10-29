package cache

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

var redisPool *redis.Pool

func SetRedisPool(pool *redis.Pool) {
	redisPool = pool
}

// RedisGet 获取缓存
func RedisGet(key string) (string, bool) {
	c1 := redisPool.Get()
	reply, err := c1.Do("get", key)
	defer closeRedisClient(c1)
	if err != nil {
		return "", false
	}
	return string(reply.([]byte)), true
}

// RedisSetWithExpire 设置缓存(带过期时间)
func RedisSetWithExpire(key string, value string, ttl int64) bool {
	c1 := redisPool.Get()
	_, err := c1.Do("set", key, value)
	if ttl > 0 {
		_, _ = c1.Do("pexpire", key, ttl)
	}
	defer closeRedisClient(c1)
	if err == nil {
		return true
	}
	return false
}

// RedisSet 设置缓存(不过期)
func RedisSet(key string, value string) bool {
	c1 := redisPool.Get()
	_, err := c1.Do("set", key, value)
	defer closeRedisClient(c1)
	if err == nil {
		return true
	}
	return false
}

// RedisRemember 获取缓存,不存在则设置
func RedisRemember(key string, fun func() string, ttl int64) string {
	value, exists := RedisGet(key)
	if exists {
		return value
	}
	value = fun()
	if ttl > 0 {
		RedisSetWithExpire(key, value, ttl)
	}
	RedisSet(key, value)
	return value
}

// RedisDelete 删除缓存
func RedisDelete(key string) bool {
	c1 := redisPool.Get()
	res, err := c1.Do("del", key)
	fmt.Println(res, err)
	defer closeRedisClient(c1)
	if err == nil {
		return true
	}
	return false
}

// closeRedisClient 关闭链接
func closeRedisClient(c1 redis.Conn) {
	err := c1.Close()
	if err != nil {
		return
	}
}
