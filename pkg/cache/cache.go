package cache

type Cache interface {
	Set(key string, value interface{}, timeout int) error
	Get(key string, to interface{}) error
	Del(key string) error
	Exist(key string) (bool, error)
	SetExpireTime(key string, timeout int) error
	Close()
	SetNx(key string, value interface{}, timeout int) (bool, error)
}

var gCache Cache

func init() {
	var err error
	gCache, err = InitRedisFromViper()
	if err != nil {
		panic(err)
	}
}

func GetCache() Cache {
	return gCache
}
