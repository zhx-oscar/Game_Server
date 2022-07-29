package main

//import "C"
import (
	"encoding/json"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type result struct {
	ErrorNo   int32
	ErrorInfo string
}

func (res result) Show() {
	content, _ := json.Marshal(res)
	fmt.Printf("Result:%s\n", string(content))
}

var CommitID string

func main() {
	//defer func() {
	//	if err := recover(); err != nil {
	//		result{
	//			ErrorNo:   -1,
	//			ErrorInfo: fmt.Sprintf("%v", err),
	//		}.Show()
	//	}
	//}()

	if len(os.Args) < 2 {
		panic("参数数量错误")
	}

	// 初始化随机数种子
	rand.Seed(time.Now().UnixNano())

	initCinder()

	// 服务器配置路径
	serverPath := "../"
	if len(os.Args) >= 4 {
		serverPath = os.Args[3]
	}

	// 初始化战斗系统
	conf := fightInit(os.Args[1], serverPath)

	// 是否打印日志
	enableLog := true
	if len(os.Args) >= 5 {
		t, _ := strconv.Atoi(os.Args[4])
		enableLog = t != 0
	}

	// 设置战斗日志
	fightSetLog("", enableLog)

	// 运行战斗
    fightData := fightRun(*conf)
    
	if len(fightData) <= 0 {
		panic("战斗数据为空")
	}

	if len(os.Args) >= 3 {
		// 战斗结果写入文件
		if err := ioutil.WriteFile(os.Args[2], fightData, os.ModePerm); err != nil {
			panic(fmt.Sprintf("写入战报文件出错，%s", err.Error()))
		}
	}

	result{
		ErrorNo: 0,
	}.Show()
}

func initCinder() {
	viper.SetConfigFile("config.json")
	if err := viper.MergeInConfig(); err != nil {
		log.Error("Read config.json failed, use default value")
		setCinderConfig()
		return
	}
}

func setCinderConfig() {
	viper.SetDefault("LogicRedis.Addr", "127.0.0.1:6379")
	viper.SetDefault("LogicRedis.Password", "")
	viper.SetDefault("FightLog.LogDir", "../log/")
	viper.SetDefault("FightLog.LogLevel", "off")
}
