package main

import (
	_ "Cinder/Base/Log"
	"Cinder/Mail/dbidx"
	"Cinder/Mail/rpcproc/usersrv"
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/viper"

	log "github.com/cihub/seelog"
)

func main() {
	areaID := "1"
	serverID := "1"

	var opt struct {
		Call      func(string)      `long:"call" description:"callback"`
		ParameMap map[string]string `long:"ParameMap" description:"A map from string to string"`
	}
	opt.Call = func(value string) {
		log.Debug("Mail ParameMap in callback: ", value)
	}
	_, err := flags.ParseArgs(&opt, os.Args[1:])
	if err != nil {
		panic(fmt.Sprintf("Parse error: %v", err))
	}

	for key, val := range opt.ParameMap {
		viper.Set(key, val)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	defer log.Flush()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	log.Debug("mail server is starting...")

	if err := dbidx.EnsureIndexes(); err != nil {
		log.Errorf("failed to ensure index: %v", err)
		return
	}

	// 必须在 serverInit() 之前从DB加载服务器列表
	if err := usersrv.LoadFromDB(); err != nil {
		log.Errorf("failed to load user server IDs from DB %v", err)
		return
	}

	if v := viper.GetString("Mail.AreaID"); v != "" {
		areaID = v
	}
	if v := viper.GetString("Mail.ServerID"); v != "" {
		serverID = v
	}

	if err := serverInit(areaID, serverID); err != nil {
		log.Error("failed to start server ", err)
		return
	}

	log.Debug("mail server is running...")
	<-c

	log.Debug("mail server is closing...")
	serverDestroy()
}
