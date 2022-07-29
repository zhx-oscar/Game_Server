package main

import (
	_ "Cinder/Base/Log"
	"Cinder/DBAgent"
	"Daisy/Prop"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/spf13/viper"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	log "github.com/cihub/seelog"
)

func main() {
	defer log.Flush()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	log.Debug("db agent server starting ... ")

	areaID := "1"
	serverID := "1"

	var opt struct {
		Call      func(string)      `long:"call" description:"callback"`
		ParameMap map[string]string `long:"ParameMap" description:"A map from string to string"`
	}
	opt.Call = func(value string) {
		log.Debug("DBAgent ParameMap in callback: ", value)
	}
	_, err := flags.ParseArgs(&opt, os.Args[1:])
	if err != nil {
		panic(fmt.Sprintf("Parse error: %v", err))
	}

	for key, val := range opt.ParameMap {
		viper.Set(key, val)
	}

	if v := viper.GetString("DBAgent.AreaID"); v != "" {
		areaID = v
	}
	if v := viper.GetString("DBAgent.ServerID"); v != "" {
		serverID = v
	}

	// 初始化进程信息
	if err = DBAgent.Init(areaID, serverID, &_RPCProc{}); err != nil {
		log.Error("start server failed ", err)
		return
	}

	Prop.RegisterPropType()

	log.Debug("db agent server is running .... ", serverID)

	go func() {
		_ = http.ListenAndServe("0.0.0.0:6142", nil)
	}()

	<-c

	log.Debug("db agent server closing ...")

	DBAgent.Destroy()
}
