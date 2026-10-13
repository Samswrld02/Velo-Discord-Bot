package redisconnection

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)


type RedisC struct {
	rdb *redis.Client
	ctx *context.Context
}

func NewRedisClient(addr string, password string, db int) (*RedisC) {
	//create context 
	ctx := context.Background()

	//create connection
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		Password: password,
		DB: db,
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

func (r *RedisC) MakeTestData(queue string)  string{
	command := Command{Name: "Sam", Lastname: "test"}
	payload, err := json.Marshal(command)

	err = r.rdb.RPush(*r.ctx, queue, payload ).Err()
	if err != nil {
		panic(err)
	}

	return string(payload)
}

func (r *RedisC) ReadFromQueue(queue string) []byte{

	//will block till event has been shot into queue
	val, err := r.rdb.BLPop(*r.ctx, 0, queue).Result()

	if err != nil {
		panic(err)
	}

	return []byte(val[1])
}
