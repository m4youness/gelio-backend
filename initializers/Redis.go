package initializers

import (
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func InitRedis() {
	if os.Getenv("REDIS_URL") == "" {
		log.Fatal("Redis url is empty")
	}

	opt, _ := redis.ParseURL(os.Getenv("REDIS_URL"))

	RedisClient = redis.NewClient(opt)

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
}
