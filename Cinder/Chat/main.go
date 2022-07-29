package main

import (
	_ "Cinder/Base/Log"
	"Cinder/Chat/rpcproc/logic"
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
	"os/signal"
	"strings"
	"syscall"

	log "github.com/cihub/seelog"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	pflag.Int("Chat.MaxFriendCount", -1, "最大好友数，负数表示无限制")
	pflag.Int("Chat.MaxAddFriendReq", -1, "最大可接收加好友请求数，负数表示无限制")
}

func main() {
	defer log.Flush()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	areaID := "1"
	serverID := "1"

	var opt struct {
		Call      func(string)      `long:"call" description:"callback"`
		ParameMap map[string]string `long:"ParameMap" description:"A map from string to string"`
	}
	opt.Call = func(value string) {
		log.Debug("Chat ParameMap in callback: ", value)
	}
	_, err := flags.ParseArgs(&opt, os.Args[1:])
	if err != nil {
		panic(fmt.Sprintf("Parse error: %v", err))
	}

	for key, val := range opt.ParameMap {
		viper.Set(key, val)
	}

	if v := viper.GetString("Chat.AreaID"); v != "" {
		areaID = v
	}
	if v := viper.GetString("Chat.ServerID"); v != "" {
		serverID = v
	}

	log.Debugf("chat server is starting (area=%s, id=%s)...", areaID, serverID)
	if err := serverInit(areaID, serverID); err != nil {
		log.Error("failed to start server ", err)
		return
	}

	// 启动时创建索引
	if err := logic.EnsureMongoDBIndexes(); err != nil {
		log.Errorf("failed to ensure mongodb indexes: %v", err)
		return
	}
	defer logic.SaveAllToDBOnExit() // 退出时存DB

	log.Debug("chat server is running...")
	<-c

	log.Debug("chat server is closing...")
	serverDestroy()
}
