package redisconnection

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)


type RedisC struct {
	rdb *redis.Client
	ctx *context.Context
}

func NewRedisClient(addr string, password string) (*RedisC) {
	//create context 
	ctx := context.Background()

	//create connection
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		Password: password,
		DB: 0,
	})
	return &RedisC{
		rdb: rdb,
		ctx: &ctx,
	}
}

type Command struct {
	Name string `json:"name"`
	Lastname string `json:"lastname"`
}

func (r *RedisC) Test() {

	command := Command{Name: "Sam", Lastname: "jaboury"}
	payload, err := json.Marshal(command)

	err = r.rdb.RPush(*r.ctx, "jsonlist", payload ).Err()
	if err != nil {
		panic(err)
	}

	val, err := r.rdb.BLPop(*r.ctx, 0, "jsonlist").Result()

	if err != nil {
		panic(err)
	}

	var got Command

	if err := json.Unmarshal([]byte(val[1]), &got); err != nil {
		panic(err)
	}

	out, _ := json.Marshal(got)

	fmt.Println("test", string(out) )
}
