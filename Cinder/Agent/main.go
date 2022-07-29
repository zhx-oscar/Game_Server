package main

import (
	"Cinder/Base/Core"
	_ "Cinder/Base/Log"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"

	log "github.com/cihub/seelog"
	"github.com/jessevdk/go-flags"
	"github.com/spf13/viper"
)

func main() {
	defer log.Flush()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	log.Debug("agent server starting ... ")

	areaID := "1"
	serverID := "1"

	var opt struct {
		Call      func(string)      `long:"call" description:"callback"`
		ParameMap map[string]string `long:"ParameMap" description:"A map from string to string"`
	}
	opt.Call = func(value string) {
		log.Debug("Agent ParameMap in callback: ", value)
	}
	_, err := flags.ParseArgs(&opt, os.Args[1:])
	if err != nil {
		panic(fmt.Sprintf("Parse error: %v", err))
	}

	for key, val := range opt.ParameMap {
		viper.Set(key, val)
	}

	if v := viper.GetString("Agent.AreaID"); v != "" {
		areaID = v
	}
	if v := viper.GetString("Agent.ServerID"); v != "" {
		serverID = v
	}

	if err := serverInit(areaID, serverID); err != nil {
		log.Errorf("start server failed ", err)
		return
	}

	log.Debug("agent server is running ....", Core.Inst.GetServiceID())

	go func() {
		_ = http.ListenAndServe("0.0.0.0:30881", nil)
	}()

	<-c

	log.Debug("agent server closing ...")

	serverDestroy()
}
