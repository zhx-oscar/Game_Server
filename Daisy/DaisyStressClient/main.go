package main

import (
	_ "Cinder/Base/Log"
	"Cinder/stats"
	"flag"
	"fmt"
	"time"

	log "github.com/cihub/seelog"
)

var (
	StressID   = 1
	LoginAddr  = "http://127.0.0.1:8080/Login"
	TargetQPS  = 1
	TotalRobot = 1
	Runtime    = 60
	Action     = "auth"
)

func main() {
	flag.IntVar(&StressID, "id", 1, "制定压测机器人ID")
	flag.StringVar(&LoginAddr, "addr", "http://127.0.0.1:8080/Login", "Login地址")
	flag.IntVar(&TargetQPS, "qps", 200, "每秒启动的机器人数量")
	flag.IntVar(&TotalRobot, "total", 1000, "机器人总数")
	flag.IntVar(&Runtime, "runtime", 0, "机器人运行时长, <=0时为长期运行")
	flag.StringVar(&Action, "action", "login", "制定压测场景")
	flag.Parse()

	log.Debug("压测机器人启动 ", StressID)
	log.Debug("Login: ", LoginAddr)
	log.Debug("QPS: ", TargetQPS, " Total: ", TotalRobot)
	log.Debug("场景: ", Action, " 运行时长: ", time.Duration(Runtime)*time.Second)

	stats.Enable()
	stats.Start()

	startInterval := time.Duration(1000/TargetQPS) * time.Millisecond
	for i := 0; i < TotalRobot; i++ {
		client := NewClient(fmt.Sprintf("StressTest%d-%d", StressID, i), "123456")
		go client.RunAction(Action)

		time.Sleep(startInterval)
	}

	log.Debug("机器人启动完成")

	if Runtime > 0 {
		time.Sleep(time.Duration(Runtime) * time.Second)
	} else {
		select {}
	}
}
