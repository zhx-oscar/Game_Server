package Data

import (
	"Daisy/DataTables"
	"strconv"
)

// 导出获取各种表数据的接口
// 添加一张表时，需要手动添加导出接口

//GetSkillConfig 获取技能配置
func GetSkillConfig() *DataTables.Skill_Config_Data {
	return inst.Get("Skill_Config").(*DataTables.Skill_Config_Data)
}

// GetBuffConfig 查询buff配置
func GetBuffConfig() *DataTables.Buff_Config_Data {
	return inst.Get("Buff_Config").(*DataTables.Buff_Config_Data)
}

//GetAttackConfig 获取伤害体配置
func GetAttackConfig() *DataTables.Attack_Config_Data {
	return inst.Get("Attack_Config").(*DataTables.Attack_Config_Data)
}

//GetMonsterConfig 获取NPC配置
func GetMonsterConfig() *DataTables.Monster_Config_Data {
	return inst.Get("Monster_Config").(*DataTables.Monster_Config_Data)
}

//GetMassConfig 获取质量配置
func GetMassConfig() *DataTables.Mass_Config_Data {
	return inst.Get("Mass_Config").(*DataTables.Mass_Config_Data)
}

//GetAIData_Config 获取AI黑板配置
func GetAIData_Config() *DataTables.AIData_Config_Data {
	return inst.Get("AIData_Config").(*DataTables.AIData_Config_Data)
}

//GetSceneConfig 获取战场配置
func GetSceneConfig() *DataTables.Scene_Config_Data {
	return inst.Get("Scene_Config").(*DataTables.Scene_Config_Data)
}

//GetSpecialAgentConfig 获取特工配置
func GetSpecialAgentConfig() *DataTables.SpecialAgent_Config_Data {
	return inst.Get("SpecialAgent_Config").(*DataTables.SpecialAgent_Config_Data)
}

//GetPropConfig 获取数值配置
func GetPropConfig() *DataTables.Prop_Config_Data {
	return inst.Get("Prop_Config").(*DataTables.Prop_Config_Data)
}

//GetBeggingConfig 获取乞求配置
func GetBeggingConfig() *DataTables.Begging_Config_Data {
	return inst.Get("Begging_Config").(*DataTables.Begging_Config_Data)
}

//GetSkillIDByBuildID 根据技能原型和BD获取新技能
func GetSkillIDByBuildID(skillID, buildID uint32) uint32 {
	for key, value := range GetSkillConfig().SkillMain_ConfigItems {
		if value.ProtoID == skillID && value.BuildID == buildID {
			return key
		}
	}
	return skillID
}

//GetFightConstConfig 获取战斗数值配置
func GetFightConstConfig() *DataTables.FightConst_Config_Data {
	return inst.Get("FightConst_Config").(*DataTables.FightConst_Config_Data)
}

//GetFightConstConfigIntValue 获取战斗常量配置表 int
func GetFightConstConfigIntValue(id uint32) (int, bool) {
	value, ok := inst.Get("FightConst_Config").(*DataTables.FightConst_Config_Data).FightConst_ConfigItems[id]
	if ok {
		val, err := strconv.Atoi(value.Value)
		return val, err == nil
	}

	return 0, false
}

//GetFightConstConfigStringValue 获取战斗常量配置表 string
func GetFightConstConfigStringValue(id uint32) (string, bool) {
	value, ok := inst.Get("FightConst_Config").(*DataTables.FightConst_Config_Data).FightConst_ConfigItems[id]
	if ok {
		return value.Value, ok
	}

	return "", false
}

//GetGuideConfig 获取新手引导配置
func GetGuideConfig() *DataTables.Guide_Config_Data {
	return inst.Get("Guide_Config").(*DataTables.Guide_Config_Data)
}

//GetDropConfig 获取掉落配置表
func GetDropConfig() *DataTables.Drop_Config_Data {
	return inst.Get("Drop_Config").(*DataTables.Drop_Config_Data)
}

//GetEquipConfig 获取装备配置
func GetEquipConfig() *DataTables.Equip_Config_Data {
	return inst.Get("Equip_Config").(*DataTables.Equip_Config_Data)
}

//GetItemTypeConfig 获取ItemType_Config
func GetItemTypeConfig() *DataTables.ItemType_Config_Data {
	return inst.Get("ItemType_Config").(*DataTables.ItemType_Config_Data)
}

//GetTalentConfig 获取天赋配置
func GetTalentConfig() *DataTables.Talent_Config_Data {
	return inst.Get("Talent_Config").(*DataTables.Talent_Config_Data)
}

//GetFastBattleConfig 获取FastBattle_Config
func GetFastBattleConfig() *DataTables.FastBattle_Config_Data {
	return inst.Get("FastBattle_Config").(*DataTables.FastBattle_Config_Data)
}

//GetTargetStrategyConfig 获取 TargetStrategy_Config
func GetTargetStrategyConfig() *DataTables.TargetStrategy_Config_Data {
	return inst.Get("TargetStrategy_Config").(*DataTables.TargetStrategy_Config_Data)
}

//GetSeasonConfig 获取 Session_Config
func GetSeasonConfig() *DataTables.Season_Config_Data {
	return inst.Get("Season_Config").(*DataTables.Season_Config_Data)
}

const (
	MailConfigMaxNum            uint32 = 1 // 邮件数量
	MailConfigExpireTime        uint32 = 2 // 普通邮件过期时间，秒数（7天）
	MailConfigAttachExpireTime  uint32 = 3 // 附件过期时间，秒数(30天)
	MailConfigTitleLen          uint32 = 4 // 邮件标题长度，字符数量。汉字占2个字符，英文数字占1个字符
	MailConfigContentLen        uint32 = 5 // 邮件内容长度
	MainConfigAttachLen         uint32 = 6 // 邮件附件长度
	MailConfigHasReadExpireTime uint32 = 7 // 已读邮件过期时间，秒数（7天）
)

func GetMailGeneralConfig(id uint32) uint32 {
	config, ok := inst.Get("Mail_Config").(*DataTables.Mail_Config_Data).MailGeneral_ConfigItems[uint32(id)]
	if !ok {
		return 0
	}
	return config.Value
}

func GetSupplyConfig() *DataTables.Supply_Config_Data {
	return inst.Get("Supply_Config").(*DataTables.Supply_Config_Data)
}

//GetTitleConfig 获取Title_Config
func GetTitleConfig() *DataTables.Title_Config_Data {
	return inst.Get("Title_Config").(*DataTables.Title_Config_Data)
}

//GetPlayerUpgradeConfig 获取PlayerUpgrade
func GetPlayerUpgradeConfig() *DataTables.PlayerUpgrade_Config_Data {
	return inst.Get("PlayerUpgrade_Config").(*DataTables.PlayerUpgrade_Config_Data)
}
