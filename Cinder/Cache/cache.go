package Cache

import (
	_ "Cinder/Base/ServerConfig"

	log "github.com/cihub/seelog"
	"github.com/go-redis/redis/v7"
	"github.com/spf13/viper"
)

func init() {
	if err := initRedis(); err != nil {
		log.Error("Init redis err ", err)
	}
}

///////////////////////////////////////////////////////////////////////////////

// Redis相关

var RedisDB *redis.Client

func initRedis() error {
	redisOptions := &redis.Options{}
	redisOptions.Addr = viper.GetString("Redis.Addr")
	redisOptions.Password = viper.GetString("Redis.Password")

	RedisDB = redis.NewClient(redisOptions)

	_, err := RedisDB.Ping().Result()
	return err
}
