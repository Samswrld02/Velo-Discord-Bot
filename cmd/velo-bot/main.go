package main

import (
	"fmt"
	"os"

	redis "github.com/Samswrld02/Velo-Discord-Bot/redis-connection"
)



func main() {
	fmt.Printf("hello world! This is the beginning of my discord bot!")   

	redisC := redis.NewRedisClient(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PW"))

	redisC.Test()
}

