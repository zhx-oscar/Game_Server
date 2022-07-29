package DB

import (
	"github.com/go-redis/redis/v7"
	"github.com/spf13/viper"
)

var RedisDB *redis.Client

func init() {
	option := &redis.Options{}

	if v := viper.GetString("LogicRedis.Addr"); v != "" {
		option.Addr = v
	}
	if v := viper.GetString("LogicRedis.Password"); v != "" {
		option.Password = v
	}

	RedisDB = redis.NewClient(option)
	_, err := RedisDB.Ping().Result()
	if err != nil {
		panic(err)
	}
}
