package conf

import (
	"Daisy/DataTables"
	"Daisy/Proto"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
)

// beHitPath 受击配置路径
var beHitPath = "../res/AnimatorController"

// SetBeHitPath 设置受击配置路径
func SetBeHitPath(path string) {
	beHitPath = path
}

// BeHitTimeLine 受击时间轴
type BeHitTimeLine struct {
	BeforeTime   uint32
	MiddleTime   uint32
	AfterTime    uint32
	Time         uint32
	CollapseTime uint32
	Distance     float32
	EndTime      uint32
}

// BeHitConfig 受击配置
type BeHitConfig struct {
	Hit    BeHitTimeLine
	Broken BeHitTimeLine
	Down   BeHitTimeLine
	Float  BeHitTimeLine
	Stun   BeHitTimeLine
	Block  BeHitTimeLine
}

// LoadModelConfig 加载Model受击配置
func (conf *BeHitConfig) LoadModelConfig(npcExcelConf *DataTables.Monster_Config_Data, modelID uint32) error {
	if modelConfig, ok := npcExcelConf.Model_ConfigItems[modelID]; !ok {
		return fmt.Errorf("在工作簿【Model_Config】中找不到model【%d】", modelID)
	} else {
		if modelConfig.AnimatorController == "" {
			return nil
		}

		if beHitConfig, err := loadBeHitTimeLine(modelConfig.AnimatorController); err != nil {
			return fmt.Errorf("加载受击配置【%s】错误，%s", modelConfig.AnimatorController, err.Error())
		} else {
			*conf = *beHitConfig
		}
	}

	return nil
}

// LoadRoleConfig 加载role受击配置
func (conf *BeHitConfig) LoadRoleConfig(jobExcelConf *DataTables.SpecialAgent_Config_Data, jobID uint32) error {
	if jobBaseConfig, ok := jobExcelConf.SpecialAgent_ConfigItems[jobID]; !ok {
		return fmt.Errorf("在工作簿【JobBase_Config】中找不到职业【%d】", jobID)
	} else {
		if jobBaseConfig.AnimatorController == "" {
			return nil
		}

		if beHitConfig, err := loadBeHitTimeLine(jobBaseConfig.AnimatorController); err != nil {
			return fmt.Errorf("加载受击配置【%s】错误，%s", jobBaseConfig.AnimatorController, err.Error())
		} else {
			*conf = *beHitConfig
		}
	}

	return nil
}

// GetBeHitTime 获取受击时间
func (conf *BeHitConfig) GetBeHitTime(hitType Proto.HitType_Enum) uint32 {
	switch hitType {
	case Proto.HitType_Hit:
		return conf.Hit.BeforeTime + conf.Hit.MiddleTime + conf.Hit.AfterTime
	case Proto.HitType_Broken:
		return conf.Broken.BeforeTime + conf.Broken.MiddleTime + conf.Broken.AfterTime
	case Proto.HitType_Down:
		return conf.Down.BeforeTime + conf.Down.MiddleTime + conf.Down.AfterTime
	case Proto.HitType_Float:
		if riseTime, ok := GetConfigs().GetConstConUint32Value(ConstExcel_HitFloatRiseTime); ok {
			return riseTime*2 + conf.Float.AfterTime
		}

		return 0
	case Proto.HitType_Stun:
		return conf.Stun.BeforeTime + conf.Stun.MiddleTime + conf.Stun.AfterTime
	case Proto.HitType_Block:
		return conf.Block.Time
	case Proto.HitType_BlockBreak:
		return conf.Block.CollapseTime
	}

	return 0
}

func loadBeHitTimeLine(fileName string) (*BeHitConfig, error) {
	var beHitConfig BeHitConfig
	filePath := path.Join(beHitPath, fileName+".json")
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &beHitConfig)
	if err != nil {
		return nil, err
	}

	return &beHitConfig, nil
}

// loadRoleBeHitConfigs 加载role受击配置
func loadRoleBeHitConfigs(jobExcelConf *DataTables.SpecialAgent_Config_Data) (map[uint32]*BeHitConfig, error) {
	beHitConfs := map[uint32]*BeHitConfig{}

	for jobBaseID := range jobExcelConf.SpecialAgent_ConfigItems {
		beHitConf := &BeHitConfig{}
		if err := beHitConf.LoadRoleConfig(jobExcelConf, jobBaseID); err != nil {
			return nil, fmt.Errorf("刷新Job受击配置失败，错误信息：%s", err.Error())
		}

		beHitConfs[jobBaseID] = beHitConf
	}

	return beHitConfs, nil
}

// loadMonsterBeHitConfigs 加载npc受击配置
func loadMonsterBeHitConfigs(npcExcelConf *DataTables.Monster_Config_Data) (map[uint32]*BeHitConfig, error) {
	beHitConfs := map[uint32]*BeHitConfig{}

	for npcID, logicConf := range npcExcelConf.Logic_ConfigItems {
		beHitConf := &BeHitConfig{}
		if err := beHitConf.LoadModelConfig(npcExcelConf, logicConf.ModelID); err != nil {
			return nil, fmt.Errorf("刷新Npc受击配置失败，错误信息：%s", err.Error())
		}

		beHitConfs[npcID] = beHitConf
	}

	return beHitConfs, nil
}
