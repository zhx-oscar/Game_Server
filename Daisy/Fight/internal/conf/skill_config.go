package conf

import (
	"Daisy/DataTables"
	"fmt"
)

// SkillConfig 技能配置
type SkillConfig struct {
	*DataTables.SkillMain_Config
	*DataTables.SkillValue_Config
	TimeLineConfig *SkillTimeLine
	TemplateConfig *AttackTemplate
}

// MainID 技能主ID
func (conf *SkillConfig) MainID() uint32 {
	return conf.SkillMain_Config.ID
}

// ValueID 技能数值ID
func (conf *SkillConfig) ValueID() uint32 {
	return conf.SkillValue_Config.ID
}

// loadConfig 加载技能配置
func (conf *SkillConfig) loadConfig(skillExcelConf *DataTables.Skill_Config_Data, attackExcelConf *DataTables.Attack_Config_Data, valueID uint32) error {
	var ok bool
	var err error

	// 加载技能数值配置
	if conf.SkillValue_Config, ok = skillExcelConf.SkillValue_ConfigItems[valueID]; !ok {
		return fmt.Errorf("在工作簿【SkillValue_Config】中找不到技能【%d】", valueID)
	}

	// 技能逻辑ID
	mainID := conf.SkillID

	// 加载技能逻辑配置
	if conf.SkillMain_Config, ok = skillExcelConf.SkillMain_ConfigItems[mainID]; !ok {
		return fmt.Errorf("在工作簿【SkillMain_Config】中找不到技能【%d】的技能逻辑【%d】", valueID, mainID)
	}

	// 加载时间轴配置
	if conf.SkillMain_Config.Timeline == "" {
		return fmt.Errorf("在工作簿【SkillMain_Config】中技能逻辑【%d】未配置Timeline", mainID)
	}

	if conf.TimeLineConfig, err = loadSkillTimelineConfig(conf.SkillMain_Config.Timeline); err != nil {
		return fmt.Errorf("在工作簿【SkillMain_Config】中加载技能逻辑【%d】的Timeline失败，%s", mainID, err.Error())
	}

	// 加载伤害体模板配置
	if conf.TemplateID <= 0 {
		return fmt.Errorf("在工作簿【SkillMain_Config】中技能逻辑【%d】未配置TemplateID", mainID)
	}

	if conf.TemplateConfig, err = loadAttackTemplateConfig(attackExcelConf, conf.TemplateID, conf.TemplateArgs, conf.TimeLineConfig.Attacks); err != nil {
		return fmt.Errorf("在工作簿【SkillMain_Config】中加载技能逻辑【%d】的伤害体模板失败，%s", mainID, err.Error())
	}

	return nil
}

// loadSkillConfigs 加载所有技能配置
func loadSkillConfigs(skillExcelConf *DataTables.Skill_Config_Data, attackExcelConf *DataTables.Attack_Config_Data) (map[uint32]*SkillConfig, error) {
	skillConf := map[uint32]*SkillConfig{}

	for _, valueConf := range skillExcelConf.SkillValue_ConfigItems {
		skillConfs := &SkillConfig{}
		if err := skillConfs.loadConfig(skillExcelConf, attackExcelConf, valueConf.ID); err != nil {
			return nil, fmt.Errorf("刷新技能配置失败，错误信息：%s", err.Error())
		}

		skillConf[valueConf.ID] = skillConfs
	}

	return skillConf, nil
}
