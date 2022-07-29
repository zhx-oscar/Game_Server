package main

import (
	"Daisy/Battle/drop"
	"bufio"
	"fmt"
	log "github.com/cihub/seelog"
	"os"
	"strconv"
	"strings"
)

func main() {
	my_drop := drop.Drop{}

	fmt.Println("按照掉落盒子id 职业 等级 （用空格区分）进行输入，eg： 1001 1 20 ,")
	fmt.Println("结束请输入 over")

	for {
		reader := bufio.NewReader(os.Stdin)
		str, _ := reader.ReadString('\n')

		if str == "over\r\n" {
			return
		}

		args := strings.Split(str, " ")
		dropid, err := strconv.Atoi(args[0])
		job, err := strconv.Atoi(args[1])
		level, err := strconv.Atoi(args[2])

		//可选参数
		mLevel, err1 := strconv.Atoi(args[3])
		mType, err2 := strconv.Atoi(args[4])
		locky, err3 := strconv.Atoi(args[5])

		if err1 != nil {
			mLevel = 1
		}
		if err2 != nil {
			mType = 1
		}
		if err3 != nil {
			locky = 50
		}

		if err != nil {
			log.Error("读取输入失败，返回")
		}
		log.Infof("%d %d %d,%d,%d,%d", uint32(dropid), uint32(job), uint32(level), uint32(mLevel), uint32(mType), uint32(locky))
		log.Info("输出是 ---------------------------------------------------")

		ok, h := my_drop.Drop(uint32(dropid), uint32(job), uint32(level))
		if ok {
			for i, val := range h {
				log.Infof("次数 %d 道具id %d 道具类型 %d 道具数量 %d", i, val.MaterialId, val.MaterialType, val.MaterialNum)
				item := my_drop.CreateItem(val.MaterialId, val.MaterialType, val.MaterialNum, uint32(mLevel), uint32(mType), uint32(locky))
				if item != nil {
					log.Debugf("掉落生成道具,%s, %d, %d\n", item.Base.ID, item.Base.ConfigID, item.Base.Type)
				} else {
					log.Debugf("掉落生成道具失败\n")
				}
			}
		} else {
			log.Errorf("dropone failed")
		}
	}
}
