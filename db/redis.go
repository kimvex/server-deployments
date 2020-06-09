package db

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

type RedisCon struct {
	Connection redis.Conn
}

var connection RedisCon

func Redisdb() redis.Conn {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		fmt.Println(err, "Error conection redis")
	}

	connection.Connection = c

	return connection.Connection
}

//GetUserId for get userid
func GetUserId() string {
	value, error := redis.String(connection.Connection.Do("get", "asdjhsajdh-asdsada-22324-dsf"))
	fmt.Println("value", value)

	if error != nil {
		fmt.Println(error)
		return ""
	}

	return value
}
