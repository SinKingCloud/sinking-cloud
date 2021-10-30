package cache

import "github.com/gomodule/redigo/redis"

type Cache interface {
	Get(key string) (string, bool)
	SetWithExpire(key string, value string, ttl int64) bool
	Set(key string, value string) bool
	Remember(key string, fun func() string, ttl int64) string
	Delete(key string) bool
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
