package cache

type Cache interface {
	Get(key string) (string, bool)
	SetWithExpire(key string, value string, ttl int64) bool
	SetExpire(key string, ttl int64) bool
	Set(key string, value string) bool
	Remember(key string, value string, ttl int64) string
	Delete(key string) bool
	Lock(key string, ttl int64) bool
	UnLock(key string) bool
}
