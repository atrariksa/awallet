package drivers

import (
	"github.com/atrariksa/awallet/configs"
	"github.com/go-redis/redis"
)

func GetRedisClient(cfg *configs.Config) *redis.Client {
	c := cfg.Cache
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Host + ":" + c.Port,
		Password:     c.Password,
		DB:           c.Namespace,
		DialTimeout:  c.DialTimeout,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		IdleTimeout:  c.IdleTimeout,
		MaxConnAge:   c.MaxConnAge,
		MinIdleConns: c.MinIdleConns,
	})

	return rdb
}
