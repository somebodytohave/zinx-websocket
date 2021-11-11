package global

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/sun-fight/zinx-websocket/zlog"
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
		zlog.Fatal("redis connect ping failed, err:", err)
	} else {
		zlog.Info("redis connect ping response:", pong)
		Redis = client
	}
}
