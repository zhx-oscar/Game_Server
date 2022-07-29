package ServerConfig

import (
	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("Log.Dir", "../log/")
	viper.SetDefault("Log.Level", "debug")
	viper.SetDefault("Redis.Addr", "127.0.0.1:6379")
	viper.SetDefault("Redis.Password", "")
	viper.SetDefault("MongoDB.Addr", "127.0.0.1:27017")
	viper.SetDefault("MongoDB.DataBase", "Game")
	viper.SetDefault("MongoDB.User", "")
	viper.SetDefault("MongoDB.Password", "")
	viper.SetDefault("Login.ListenAddr", "0.0.0.0:8080")
	viper.SetDefault("Login.AutoCreate", true)
	viper.SetDefault("ETCD.Addr", "127.0.0.1:2379")
	viper.SetDefault("NSQ.Addr", "127.0.0.1:4150")
	viper.SetDefault("NSQ.Lookup", "127.0.0.1:4161")
	viper.SetDefault("NSQ.Admin", "127.0.0.1:4171")
	viper.SetDefault("NATS.Addr", "127.0.0.1:4222")
	viper.SetDefault("LoadBlance.OnOff", false)
	viper.SetDefault("LoadBlance.LimitValue", 1.0)
	viper.SetDefault("TrafficControl.MaxTrafficCount", 50000)
	viper.SetDefault("TrafficControl.On_Off", false)

	viper.SetConfigFile("server.json")
	if err := viper.MergeInConfig(); err != nil {
		log.Error("Read server.json failed, use default value")
	}
}
