package cache

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type Redis struct {
	Pool *redis.Pool
}

// NewRedis 实例化redis
func NewRedis(host string, port string, pwd string, db int64) Cache {
	pool := &redis.Pool{
		//连接方法
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host+":"+port)
			if err != nil {
				return nil, err
			}
			if pwd != "" {
				if _, err := c.Do("AUTH", pwd); err != nil {
					err := c.Close()
					if err != nil {
						return nil, err
					}
					return nil, err
				}
			}
			_, err = c.Do("SELECT", db)
			if err != nil {
				panic(err)
				return nil, err
			}
			return c, nil
		},
		DialContext:  nil,
		TestOnBorrow: nil,
		//最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态。
		MaxIdle: 5000,
		//最大的激活连接数，表示同时最多有N个连接
		MaxActive: 10000,
		//最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭
		IdleTimeout:     0,
		Wait:            false,
		MaxConnLifetime: 0,
	}
	return &Redis{Pool: pool}
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
