package types

import "os"

type Config struct {
	RedisAddr string
	RedisPw string
	RedisQueue string
	RedisDb string
	DiscordBotToken string
	//guild to register slash commands on while developing, empty = global
	DiscordDevGuildID string
	//base url of the velo api, e.g. http://localhost/api/v1 (sail)
	VeloAPIURL string
	//shared secret, same value as DISCORD_BOT_API_TOKEN in the velo backend
	VeloBotToken string
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
		DiscordBotToken: returnEnv("DISCORD_BOT_TOKEN", ""),
		DiscordDevGuildID: returnEnv("DISCORD_DEV_GUILD_ID", ""),
		VeloAPIURL: returnEnv("VELO_API_URL", "http://localhost/api/v1"),
		VeloBotToken: returnEnv("VELO_BOT_TOKEN", ""),
	}
}
