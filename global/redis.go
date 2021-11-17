package global

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func InitRedis() {
	redisCfg := GlobalObject.RedisConfig
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password, // no password set
		DB:       redisCfg.DB,       // use default DB
	})
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		Glog.Fatal("redis connect ping failed, err:", zap.Error(err))
	} else {
		Glog.Info("redis connect ping response:" + pong)
		Redis = client
	}
}
