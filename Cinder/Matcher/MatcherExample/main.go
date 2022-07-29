package main

import (
	_ "Cinder/Base/Log"
	"Cinder/Matcher/matcherlib"
	"flag"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/viper"

	log "github.com/cihub/seelog"
)

func main() {
	defer log.Flush()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	log.Debug("matcher server starting ... ")

	areaID := "1"
	serverID := "1"
	flag.StringVar(&areaID, "area", "1", "大区ID")
	flag.StringVar(&serverID, "id", "1", "服务器ID")
	flag.Parse()
	if v := viper.GetString("Matcher.AreaID"); v != "" {
		areaID = v
	}
	if v := viper.GetString("Matcher.ServerID"); v != "" {
		serverID = v
	}

	// 初始化进程信息
	if err := matcherlib.Init(areaID, serverID, &RPCProc{}, &RoomEventHandler{}); err != nil {
		log.Error("start server failed ", err)
		return
	}
	defer matcherlib.Destroy()

	go func() {
		log.Debug(http.ListenAndServe("localhost:16060", nil))
	}()

	log.Debug("matcher server is running .... ", serverID)
	<-c
	log.Debug("matcher server closing ...")
}
