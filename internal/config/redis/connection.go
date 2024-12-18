package redis

import (
	"fmt"

	"github.com/DanyAdhi/learn-golang/internal/config"
	"github.com/redis/go-redis/v9"
)

func Connect() *redis.Client {
	rdConnStr := fmt.Sprintf(
		"redis://:%s@%s:%s",
		config.AppConfig.REDIS_PASSWORD,
		config.AppConfig.REDIS_HOST,
		config.AppConfig.REDIS_PORT,
	)
	opt, err := redis.ParseURL(rdConnStr)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(opt)
	return rdb
}
