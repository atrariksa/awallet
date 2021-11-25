package repos

import (
	"sync"

	"github.com/atrariksa/awallet/configs"
	"github.com/go-redis/redis"
)

type Cache struct {
	cfg *configs.Config
	rc  *redis.Client
	sync.Mutex
}

func NewCache(cfg *configs.Config, rc *redis.Client) *Cache {
	return &Cache{
		cfg: cfg,
		rc:  rc,
	}
}

type ICache interface {
	Get(key string) (val []byte, err error)
	Set(key string, val []byte) (err error)
	Del(key string) (err error)
}

func (c *Cache) Get(key string) (val []byte, err error) {
	c.Lock()
	val, err = c.rc.Get(key).Bytes()
	c.Unlock()
	return
}

func (c *Cache) Set(key string, val []byte) (err error) {
	c.Lock()
	_, err = c.rc.Set(key, val, c.cfg.Cache.CacheDuration).Result()
	c.Unlock()
	return
}

func (c *Cache) Del(key string) (err error) {
	c.Lock()
	_, err = c.rc.Del(key).Result()
	c.Unlock()
	return
}
