package main

import (
	"Daisy/Data"
	"Daisy/Fight"
	"fmt"
	"strings"
)

func fightInit(confPath, serverPath string) *FightConf {
	// 战斗双方信息
	conf := &FightConf{
		Enemy:             map[int]*Stand{},
		Mine:              map[int]*Stand{},
		_MineStandPosMap:  map[int]*StandPos{},
		_EnemyStandPosMap: map[int]*StandPos{},
	}

	// 加载战斗双方信息
	conf.LoadConfig(confPath)

	// 加载战斗配置
	fightLoadConf(conf, serverPath)

	return conf
}

func fightLoadConf(conf *FightConf, serverPath string) {
	// Windows \ 替换成 /
	serverPath = strings.ReplaceAll(serverPath, "\\", "/")

	fmt.Println("fightSimulator serverPath is :", serverPath)

	// 加载Excel数据
	if err := Data.LoadDataTables(serverPath + "res/DataTables"); err != nil {
		panic(fmt.Sprintf("ErrorCode:404 读取DataTables配置出错，%s", err.Error()))
	}

	// 加载战场配置
	conf.LoadBattlefield(serverPath + "res/MapData/" + conf.Battlefield)

	// 加载战斗系统配置
	Fight.LoadConfig(serverPath)

	// 加载站位
	conf.LoadStandPos()
}

func fightSetLog(path string, enable bool) {
	if !enable {
		Fight.SetLogLevel(Fight.Level_Off)
		return
	}

	if path != "" {
		Fight.SetLogLevel(Fight.Level_File)
		Fight.SetLogPath(path)
	} else {
		Fight.SetLogLevel(Fight.Level_Console)
	}
}

// fightRun 运行战斗系统
func fightRun(conf FightConf) []byte {
	fightResult, err := Fight.Play(conf.BuildSceneInfo())
	if err != nil {
		panic(fmt.Sprintf("模拟战斗出错，%s", err.Error()))
	}
	fightResult.CommitID = CommitID
	fightResult.Progress = conf.BattleAreaID

	fightData, err := fightResult.Marshal()
	if err != nil {
		panic(fmt.Sprintf("序列化战报出错，%s", err.Error()))
	}

	return fightData
}
