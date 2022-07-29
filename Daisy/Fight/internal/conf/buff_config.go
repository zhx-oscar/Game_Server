package conf

import (
	"Daisy/DataTables"
	"fmt"
)

// BuffConfig buff配置
type BuffConfig struct {
	*DataTables.BuffMain_Config
	EffectConfs []*BuffEffect
}

// MainID buff主ID
func (conf *BuffConfig) MainID() uint32 {
	return conf.BuffMain_Config.ID
}

// loadInnerBuffConfig 加载内部buff配置
func (conf *BuffConfig) loadInnerBuffConfig(buffExcelConf *DataTables.Buff_Config_Data, mainID uint32) error {
	var ok bool
	var err error

	// 加载buff逻辑配置
	if conf.BuffMain_Config, ok = buffExcelConf.BuffMain_ConfigItems[mainID]; !ok {
		return fmt.Errorf("在工作簿【BuffMain_Config】中找不到buff【%d】", mainID)
	}

	// 加载buff效果配置
	if conf.EffectConfs, err = loadBuffEffectConfig(conf.BuffEffect); err != nil {
		return fmt.Errorf("在工作簿【BuffMain_Config】中加载buff【%d】的BuffEffect失败，%s", mainID, err.Error())
	}

	return nil
}

// loadConfig 加载buff配置
func (conf *BuffConfig) loadConfig(buffExcelConf *DataTables.Buff_Config_Data, mainID uint32) error {
	var ok bool
	var err error

	// 加载buff逻辑配置
	if conf.BuffMain_Config, ok = buffExcelConf.BuffMain_ConfigItems[mainID]; !ok {
		return fmt.Errorf("在工作簿【BuffMain_Config】中找不到buff【%d】", mainID)
	}

	// 加载buff效果配置
	if conf.EffectConfs, err = loadBuffEffectConfig(conf.BuffEffect); err != nil {
		return fmt.Errorf("在工作簿【BuffMain_Config】中加载buff【%d】的BuffEffect失败，%s", mainID, err.Error())
	}

	return nil
}

// loadBuffConfigs 加载所有buff配置
func loadBuffConfigs(buffExcelConf *DataTables.Buff_Config_Data) (map[uint32]*BuffConfig, error) {
	buffConfs := map[uint32]*BuffConfig{}

	// 加载内置buff
	for buffMainID := InnerBuffID_Begin; buffMainID < InnerBuffID_End; buffMainID++ {
		buffConf := &BuffConfig{}
		if err := buffConf.loadInnerBuffConfig(buffExcelConf, buffMainID); err != nil {
			return nil, fmt.Errorf("刷新buff配置失败，错误信息：%s", err.Error())
		}

		buffConfs[buffMainID] = buffConf
	}

	// 加载其他buff
	for _, mainConf := range buffExcelConf.BuffMain_ConfigItems {
		buffConf := &BuffConfig{}
		if err := buffConf.loadConfig(buffExcelConf, mainConf.ID); err != nil {
			return nil, fmt.Errorf("刷新buff配置失败，错误信息：%s", err.Error())
		}

		buffConfs[mainConf.ID] = buffConf
	}

	return buffConfs, nil
}
