package database

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client
var Ctx = context.Background()

func InitRedis() {
	addr := os.Getenv("REDIS_ADDR")
	pass := os.Getenv("REDIS_PASS")
	db := os.Getenv("REDIS_DB")

	if addr == "" {
		addr = "127.0.0.1:6379"
	}
	if pass == "" {
		pass = ""
	}
	if db == "" {
		db = "0"
	}

	// 初始化 Redis 客户端
	Rdb = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     pass,
		DB:           0,
		PoolSize:     20, // 最大连接数
		MinIdleConns: 5,  // 最小空闲连接
	})

	// 检测连接是否成功
	if _, err := Rdb.Ping(Ctx).Result(); err != nil {
		log.Fatalf("Redis 连接失败 ❌: %v", err)
	} else {
		log.Println("Redis 连接成功 ✅")
	}
}
