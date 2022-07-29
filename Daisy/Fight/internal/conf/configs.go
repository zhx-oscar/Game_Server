package conf

import (
	"Daisy/DataTables"
	"Daisy/Proto"
	"strconv"
)

// config 全局配置
var configs *Configs

// GetConfigs 获取全局配置
func GetConfigs() *Configs {
	return configs
}

// Configs 全局配置
type Configs struct {
	// excel配置表
	skillExcelConf        *DataTables.Skill_Config_Data
	buffExcelConf         *DataTables.Buff_Config_Data
	attackExcelConf       *DataTables.Attack_Config_Data
	specialAgentExcelConf *DataTables.SpecialAgent_Config_Data
	monsterExcelConf      *DataTables.Monster_Config_Data
	aiDataExcelConf       *DataTables.AIData_Config_Data
	propValueExcelConf    *DataTables.Prop_Config_Data
	constExcelConf        *DataTables.FightConst_Config_Data
	massExcelConf         *DataTables.Mass_Config_Data
	targetStrategy        *DataTables.TargetStrategy_Config_Data

	// 技能与buff配置
	skillConfs map[uint32]*SkillConfig
	buffConfs  map[uint32]*BuffConfig

	// 受击类配置
	playerBeHitConfs  map[uint32]*BeHitConfig
	monsterBeHitConfs map[uint32]*BeHitConfig
	npcActConfs       map[uint32]*ActConfig

	// AIData
	aiDataConfs map[uint32]*AIInfo

	// pawnConfig
	pawnConfs map[Proto.PawnType_Enum]map[uint32]*PawnConfig

	// innerAttackTemplate 内置伤害体模板
	innerAttackTemplate map[InnerAttackID]*AttackTemplate
}

//GetPawnConfig 获取pawnConfig
func (confs *Configs) GetPawnConfig(pawnType Proto.PawnType_Enum, id uint32) (*PawnConfig, bool) {
	conf, ok := confs.pawnConfs[pawnType]
	if !ok {
		return nil, false
	}

	pawnConfig, ok := conf[id]
	if !ok {
		return nil, false
	}

	return pawnConfig, true
}

//GetSkillExcelConfig 获取技能excel原生数据
func (confs *Configs) GetSkillExcelConfig() *DataTables.Skill_Config_Data {
	return confs.skillExcelConf
}

//GetBuffExcelConfig 获取buff excel原生数据
func (confs *Configs) GetBuffExcelConfig() *DataTables.Buff_Config_Data {
	return confs.buffExcelConf
}

//GetAttackExcelConfig 获取伤害体excel原生数据
func (confs *Configs) GetAttackExcelConfig() *DataTables.Attack_Config_Data {
	return confs.attackExcelConf
}

//GetMonsterExcelConfig 获取怪物excel原生数据
func (confs *Configs) GetMonsterExcelConfig() *DataTables.Monster_Config_Data {
	return confs.monsterExcelConf
}

//GetSpecialAgentExcelConfig 获取特工excel原生数据
func (confs *Configs) GetSpecialAgentExcelConfig() *DataTables.SpecialAgent_Config_Data {
	return confs.specialAgentExcelConf
}

//GetPropValueExcelConfig 获取数值配置excel原生数据
func (confs *Configs) GetPropValueExcelConfig() *DataTables.Prop_Config_Data {
	return confs.propValueExcelConf
}

//GetMassExcelConfig 获取质量配置excel原生数据
func (confs *Configs) GetMassExcelConfig() *DataTables.Mass_Config_Data {
	return confs.massExcelConf
}

//GetTargetStrategyConfig 获取目标决策配置表信息
func (confs *Configs) GetTargetStrategyConfig() *DataTables.TargetStrategy_Config_Data {
	return confs.targetStrategy
}

// GetSkillConfig 查询技能配置
func (confs *Configs) GetSkillConfig(valueID uint32) (*SkillConfig, bool) {
	conf, ok := confs.skillConfs[valueID]
	return conf, ok
}

// GetBuffConfig 查询buff配置
func (confs *Configs) GetBuffConfig(mainID uint32) (*BuffConfig, bool) {
	conf, ok := confs.buffConfs[mainID]
	return conf, ok
}

// GetPlayerBeHitConfig 查询玩家受击配置
func (confs *Configs) GetPlayerBeHitConfig(jobID uint32) (*BeHitConfig, bool) {
	conf, ok := confs.playerBeHitConfs[jobID]
	return conf, ok
}

// GetNpcBeHitConfig 查询Npc受击配置
func (confs *Configs) GetNpcBeHitConfig(npcID uint32) (*BeHitConfig, bool) {
	conf, ok := confs.monsterBeHitConfs[npcID]
	return conf, ok
}

// GetNpcActConfig 查询Npc表演配置
func (confs *Configs) GetNpcActConfig(npcID uint32) (*ActConfig, bool) {
	conf, ok := confs.npcActConfs[npcID]
	return conf, ok
}

// GetAIDataConf 查询行为树外部AI配置
func (confs *Configs) GetAIDataConf(key uint32) (*AIInfo, bool) {
	val, ok := confs.aiDataConfs[key]
	return val, ok
}

//GetConstConf 获取常量配置
func (confs *Configs) GetConstConf(key uint32) (string, bool) {
	if confs.constExcelConf == nil {
		return "", false
	}

	val, ok := confs.constExcelConf.FightConst_ConfigItems[key]
	if ok {
		return val.Value, ok
	}

	return "", false
}

//GetConstConIntValue 获取常量配置表 int
func (confs *Configs) GetConstConIntValue(key uint32) (int, bool) {
	if confs.constExcelConf == nil {
		return 0, false
	}

	value, ok := confs.constExcelConf.FightConst_ConfigItems[key]
	if ok {
		val, err := strconv.Atoi(value.Value)
		return val, err == nil
	}

	return 0, false
}

//GetConstConUint32Value 获取常量配置表 uint32
func (confs *Configs) GetConstConUint32Value(key uint32) (uint32, bool) {
	if confs.constExcelConf == nil {
		return 0, false
	}

	value, ok := confs.constExcelConf.FightConst_ConfigItems[key]
	if ok {
		val, err := strconv.Atoi(value.Value)
		return uint32(val), err == nil
	}

	return 0, false
}

//GetConstConFloat64Value 获取常量配置表 float64
func (confs *Configs) GetConstConFloat64Value(key uint32) (float64, bool) {
	if confs.constExcelConf == nil {
		return 0, false
	}

	value, ok := confs.constExcelConf.FightConst_ConfigItems[key]
	if ok {
		val, err := strconv.ParseFloat(value.Value, 64)
		return val, err == nil
	}

	return 0, false
}

// GetInnerAttackTemplate 获取内置伤害体模板
func (confs *Configs) GetInnerAttackTemplate(id InnerAttackID) (*AttackTemplate, bool) {
	tmpl, ok := confs.innerAttackTemplate[id]
	return tmpl, ok
}
