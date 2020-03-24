package db

import (
	"fmt"

	"github.com/go-redis/redis"
)

type RedisCon struct {
	Connection *redis.Client
}

var connection RedisCon

func Redisdb() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	connection.Connection = client

	return connection.Connection
}

func GetUserId() string {
	value, error := connection.Connection.Get("asdjhsajdh-asdsada-22324-dsf").Result()

	if error != nil {
		fmt.Println(error)
		return ""
	}

	return value
}
