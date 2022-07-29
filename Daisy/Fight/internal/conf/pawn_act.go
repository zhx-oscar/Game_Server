package conf

import (
	"Daisy/DataTables"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
)

// actTimelinePath 表演时间轴配置路径
var actTimelinePath = "../res/Timeline"

// SetActTimelinePath 出生或死亡表演时间轴配置路径
func SetActTimelinePath(path string) {
	actTimelinePath = path
}

// ActConfig 表演配置
type ActConfig struct {
	Born      ActTimeLine
	Dead      ActTimeLine
	OverDrive ActTimeLine
	WeakBegin ActTimeLine
	WeakLoop  ActTimeLine
	WeakEnd   ActTimeLine
}

// ActTimeLine 表演配置
type ActTimeLine struct {
	Time uint32
}

// loadActConfig 加载表演配置
func (conf *ActConfig) loadActConfig(npcExcelConf *DataTables.Monster_Config_Data, npcID uint32) error {
	if logicConfig, ok := npcExcelConf.Logic_ConfigItems[npcID]; !ok {
		return fmt.Errorf("在工作簿【Logic_Config】中找不到npcid【%d】", npcID)
	} else {
		var err error
		if logicConfig.BornTimeline != "" {
			if conf.Born, err = loadActTimeLineInfo(npcID, "DeadOrAlive/", logicConfig.BornTimeline); err != nil {
				return fmt.Errorf("加载表演配置【%s】错误，%s", logicConfig.BornTimeline, err.Error())
			}
		}

		if logicConfig.DeadTimeline != "" {
			if conf.Dead, err = loadActTimeLineInfo(npcID, "DeadOrAlive/", logicConfig.DeadTimeline); err != nil {
				return fmt.Errorf("加载表演配置【%s】错误，%s", logicConfig.DeadTimeline, err.Error())
			}
		}

		if logicConfig.OverDriveTimline != "" {
			if conf.OverDrive, err = loadActTimeLineInfo(npcID, "OverDrive/", logicConfig.OverDriveTimline); err != nil {
				return fmt.Errorf("加载表演配置【%s】错误，%s", logicConfig.OverDriveTimline, err.Error())
			}
		}

		if logicConfig.WeakBeginTimline != "" {
			if conf.WeakBegin, err = loadActTimeLineInfo(npcID, "Weak/", logicConfig.WeakBeginTimline); err != nil {
				return fmt.Errorf("加载表演配置【%s】错误，%s", logicConfig.WeakBeginTimline, err.Error())
			}
		}

		if logicConfig.WeakLoopTimline != "" {
			if conf.WeakLoop, err = loadActTimeLineInfo(npcID, "Weak/", logicConfig.WeakLoopTimline); err != nil {
				return fmt.Errorf("加载表演配置【%s】错误，%s", logicConfig.WeakLoopTimline, err.Error())
			}
		}

		if logicConfig.WeakEndTimline != "" {
			if conf.WeakEnd, err = loadActTimeLineInfo(npcID, "Weak/", logicConfig.WeakEndTimline); err != nil {
				return fmt.Errorf("加载表演配置【%s】错误，%s", logicConfig.WeakEndTimline, err.Error())
			}
		}
	}

	return nil
}

func loadActTimeLineInfo(npcID uint32, folderName, fileName string) (ActTimeLine, error) {
	var timeLine ActTimeLine
	filePath := path.Join(actTimelinePath, folderName, fileName+".json")
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return timeLine, fmt.Errorf("在工作簿【Monster_Config】中的怪物【%d】的表演TimeLine配置错误，%s", npcID, err.Error())
	}

	err = json.Unmarshal(file, &timeLine)
	if err != nil {
		return timeLine, fmt.Errorf("在工作簿【Monster_Config】中的怪物【%d】的表演TimeLine配置错误，%s", npcID, err.Error())
	}

	return timeLine, nil
}

// loadNpcActConfigs 加载npc表演配置
func loadNpcActConfigs(npcExcelConf *DataTables.Monster_Config_Data) (map[uint32]*ActConfig, error) {
	actConfs := map[uint32]*ActConfig{}

	for npcID := range npcExcelConf.Logic_ConfigItems {
		actConf := &ActConfig{}
		if err := actConf.loadActConfig(npcExcelConf, npcID); err != nil {
			return nil, fmt.Errorf("刷新Npc表演配置失败，错误信息：%s", err.Error())
		}

		actConfs[npcID] = actConf
	}

	return actConfs, nil
}
