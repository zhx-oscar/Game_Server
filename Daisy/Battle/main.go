package main

import (
	_ "Cinder/Base/Log"
	"Cinder/Space"
	"Daisy/Prop"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	log "github.com/cihub/seelog"
	"github.com/jessevdk/go-flags"
	"github.com/spf13/viper"
)

func main() {
	defer log.Flush()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	log.Debug("space server is start")

	areaID := "1"
	serverID := "1"

	var opt struct {
		Call      func(string)      `long:"call" description:"callback"`
		ParameMap map[string]string `long:"ParameMap" description:"A map from string to string"`
	}
	opt.Call = func(value string) {
		log.Debug("Battle ParameMap in callback: ", value)
	}
	_, err := flags.ParseArgs(&opt, os.Args[1:])
	if err != nil {
		panic(fmt.Sprintf("Parse error: %v", err))
	}

	for key, val := range opt.ParameMap {
		viper.Set(key, val)
	}

	if v := viper.GetString("Battle.AreaID"); v != "" {
		areaID = v
	}
	if v := viper.GetString("Battle.ServerID"); v != "" {
		serverID = v
	}

	if err = Space.Init(areaID, serverID, &_Team{}, &_User{}, &_RPCProc{}); err != nil {
		log.Error("init space service failed ", err)
		return
	}

	Prop.RegisterPropType()

	//var cellSerList []uint32
	//items, _ := cellserv.GetAllList()
	//for _, item := range items {
	//	cellSerList = append(cellSerList, item.Data.ID)
	//}
	//go func() {
	//	// TODO 暂时采用 goroutine 的方式创建世界频道,10s 间隔是为了保证 nsq 连接成功
	//	time.Sleep(10 * time.Second)
	//	channelMgr.InitAllHuntChannels(cellSerList)
	//	channelMgr.PrintHuntChannels()
	//}()
	//
	//stats.Enable()

	log.Debug("space server is running .... ", serverID)

	go func() {
		_ = http.ListenAndServe("0.0.0.0:6143", nil)
	}()

	<-c

	log.Debug("space server closing ...")

	Space.Destroy()
}
