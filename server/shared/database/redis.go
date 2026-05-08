package database

import (
	"context"
	"strconv"

	"foodRatingSystem/shared/config"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func ConnectRedis() {
	cfg := config.AppConfig
	redisDB, _ := strconv.Atoi("0")
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: cfg.RedisPassword,
		DB:       redisDB,
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		panic(err)
	}
}
