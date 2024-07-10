package initializers

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		fmt.Printf("Could not connect to Redis: %v", err)
	}
}
