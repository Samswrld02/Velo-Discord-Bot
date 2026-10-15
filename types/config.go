package types

import "os"

type Config struct {
	RedisAddr string
	RedisPw string
	RedisQueue string
	RedisDb string
	UserChannel string
	DiscordBotToken string
}

func returnEnv(val string, def string) string {
	if val, ok := os.LookupEnv(val); ok {
		return val
	}

	return def
}

func NewConfig() *Config {
	return &Config{
		RedisAddr: returnEnv("REDIS_ADDR", "localhost:6379"),
		RedisPw: returnEnv("REDIS_PW", ""),
		RedisQueue: returnEnv("REDIS_QUEUE", ""),
		RedisDb: returnEnv("REDIS_DB", "0"),
		UserChannel: returnEnv("", ""),
		DiscordBotToken: returnEnv("DISCORD_BOT_TOKEN", ""),
	}
}