package main

import (
	"Daisy/Const"
	"Daisy/Data"
	"Daisy/DataTables"
	"Daisy/ErrorCode"
	"Daisy/Fight"
	"Daisy/Proto"
	"strconv"
	"time"
)

//calcSkillValueID 计算skillValueID
func (r *_Role) calcSkillValueID(skillMainID, skillLv uint32) uint32 {
	data, ok := Data.SkillData[skillMainID]
	if !ok {
		r.Error("unknown skillMainID ", skillMainID)
		return 0
	}

	skillValueCfg, ok := data[skillLv]
	if !ok {
		r.Error("unknown SkillLv !! skillMainID + skillLv is : ", skillMainID, skillLv)
		return 0
	}

	return skillValueCfg.ID
}

//GetMoneyByMoneyID 获取对应货币数量
func (r *_Role) GetMoneyByMoneyID(id uint32) uint64 {
	switch id {
	case Const.Gold:
		return r.prop.Data.Gold
	case Const.Diamond:
		return r.prop.Data.Diamond
	default:
		r.Error("unknown money id:", id)
		return 0
	}
}

//UseMoneyByMoneyID 使用对应货币数量
func (r *_Role) UseMoneyByMoneyID(id uint32, num uint64) {
	switch id {
	case Const.Gold:
		r.prop.SyncRemoveGold(uint32(num))
	case Const.Diamond:
		r.prop.SyncRemoveDiamond(uint32(num))
	default:
		r.Error("unknown money id:", id)
	}
}

//RPC_ChangeBuildSkill 修改build内技能
func (user *_User) RPC_ChangeBuildSkill(buildID string, skillID, pos uint32) int32 {
	if user.role == nil {
		return ErrorCode.RoleIsNil
	}

	return user.role.changeBuildSkill(buildID, skillID, pos)
}

//changeBuildSkill 修改build内技能
func (r *_Role) changeBuildSkill(buildID string, skillID, pos uint32) int32 {
	build, ok := r.prop.Data.BuildMap[buildID]
	if !ok {
		return ErrorCode.NotFindBuild
	}

	skillEntity, ok := r.prop.Data.SkillsLearned[skillID]
	if !ok {
		return ErrorCode.BuildSkillNotLearned
	}

	//技能道具
	skillEntityConfig, ok := Data.GetSkillConfig().SkillEntity_ConfigItems[skillEntity.SkillMainID]
	if !ok {
		return ErrorCode.SkillItemNotFind
	}

	//技能主体信息
	skill, ok := Data.GetSkillConfig().SkillMain_ConfigItems[skillEntity.SkillMainID]
	if !ok {
		return ErrorCode.SkillItemNotFind
	}

	//所属特工类型检测
	if !r.checkSkillAssociatedSpecialAgent(build, skillEntityConfig) {
		return ErrorCode.SkillNotMatchSpecialAgent
	}

	var resultCode int32
	switch skill.SkillKind {
	case Fight.SkillKind_Ultimate:
		//必杀技处理
		resultCode = r.changeBuildUltimateSkill(buildID, skillID)

	//case Fight.SkillKind_Gifted:
	//	//天赋技处理
	//	resultCode = r.changeBuildGiftedSkill(buildID, skillID)

	case Fight.SkillKind_Super:
		//超能技处理
		resultCode = r.changeBuildSuperSkill(build, buildID, skillID, pos)

	default:
		r.Error("装备的技能类型错误 ", skill.SkillKind)
		return ErrorCode.BuilldSkillKindNotMatch
	}

	r.refreshBuildFightAttr(build.BuildID)
	return resultCode
}

//checkSkillAssociatedSpecialAgent 检测技能特工关联
func (r *_Role) checkSkillAssociatedSpecialAgent(build *Proto.BuildData, skillEntityCfg *DataTables.SkillEntity_Config) bool {
	agentCfg, ok := Data.GetSpecialAgentConfig().SpecialAgent_ConfigItems[build.SpecialAgentID]
	if !ok {
		return false
	}

	propCfg, ok := Data.GetPropConfig().PropValue_ConfigItems[agentCfg.PropValueID]
	if !ok {
		return false
	}

	//数值系统类型列表 特工ID列表  0 非0  特定特工才能使用
	if len(skillEntityCfg.SpecialAgentType) == 0 && len(skillEntityCfg.SpecialAgentID) != 0 {
		//特工ID 匹配
		for _, specialAgentID := range skillEntityCfg.SpecialAgentID {
			if specialAgentID == build.SpecialAgentID {
				return true
			}
		}
	}

	//数值系统类型列表 特工ID列表  0 0  全部通用
	if len(skillEntityCfg.SpecialAgentType) == 0 && len(skillEntityCfg.SpecialAgentID) == 0 {
		return true
	}

	//数值系统类型列表 特工ID列表  非0 0  数值系统类型检测
	if len(skillEntityCfg.SpecialAgentType) > 0 && len(skillEntityCfg.SpecialAgentID) == 0 {
		//数值系统类型列表 匹配
		for _, specialAgentType := range skillEntityCfg.SpecialAgentType {
			if specialAgentType == propCfg.Type {
				return true
			}
		}
	}

	//数值系统类型列表 特工ID列表  非0 非0  ||关系
	if len(skillEntityCfg.SpecialAgentType) > 0 && len(skillEntityCfg.SpecialAgentID) > 0 {
		//数值系统类型列表 匹配
		for _, specialAgentType := range skillEntityCfg.SpecialAgentType {
			if specialAgentType == propCfg.Type {
				return true
			}
		}

		//特工ID 匹配
		for _, specialAgentID := range skillEntityCfg.SpecialAgentID {
			if specialAgentID == build.SpecialAgentID {
				return true
			}
		}
	}

	return false
}

//changeBuildUltimateSkill 修改build内技能--必杀技
func (r *_Role) changeBuildUltimateSkill(buildID string, skillID uint32) int32 {
	r.prop.SyncUpdateBuildUltimateSkillID(buildID, skillID)
	return ErrorCode.Success
}

//changeBuildSuperSkill 修改build内技能--超能技
func (r *_Role) changeBuildSuperSkill(build *Proto.BuildData, buildID string, skillID, pos uint32) int32 {

	//槽位数量最大上限配置获取
	maxPos, ok := Data.GetSpecialAgentConfig().SpecialAgentConst_ConfigItems[Const.SpecialAgent_buildSuperSlillCount]
	if !ok {
		return ErrorCode.BuildSuperSkillCountMaxNotFind
	}

	//槽位数量最大上限配置检测
	if pos > maxPos.Value {
		return ErrorCode.OverLimitBuildSuperSkillMaxNotFind
	}

	//槽位 响应客户端lua配合 从1开始
	if pos == 0 {
		return ErrorCode.BuilldSkillSuperPosError
	}

	//查找当前位置之前是否装备的有技能
	oldSkillID, oldExist := build.Skill.SuperSkill[pos]
	//当前位置装备的技能 和 新技能是同一个技能 不做任何处理
	if oldExist && oldSkillID == skillID {
		return ErrorCode.Success
	}

	var newItemOldPos uint32
	for index, val := range build.Skill.SuperSkill {
		if val == skillID {
			newItemOldPos = index
			break
		}
	}
	//检测新超能技是否已经被当前build使用
	r.prop.SyncSwapBuildSuper(buildID, skillID, newItemOldPos, oldSkillID, pos)

	return ErrorCode.Success

}

//getBuildBattleSkills 获取build内参与战斗的技能数据 目前返回 必杀技、天赋技、超能技
func (r *_Role) getBuildBattleSkills() (ultimateSkills, giftedSkill, superSkills []uint32) {
	fightingBuild, ok := r.prop.Data.BuildMap[r.prop.Data.FightingBuildID]
	if !ok {
		return
	}

	//必杀技处理
	if fightingBuild.Skill.UltimateSkillID != 0 {
		skillData, ok := r.prop.Data.SkillsLearned[fightingBuild.Skill.UltimateSkillID]
		if !ok {
			r.Errorf("出战build:%v 必杀技已经装备,但是已学技能列表中并没有查找到对应技能:%v", r.prop.Data.FightingBuildID, fightingBuild.Skill.UltimateSkillID)
			return
		}

		//等级生效
		skillValueID := r.calcSkillValueID(skillData.SkillMainID, skillData.Lv)

		ultimateSkills = append(ultimateSkills, skillValueID)
	}

	//天赋技处理

	//超能技处理
	maxPos, ok := Data.GetSpecialAgentConfig().SpecialAgentConst_ConfigItems[Const.SpecialAgent_buildSuperSlillCount]
	if ok {
		for pos := uint32(1); pos <= maxPos.Value; pos++ {
			skillIID := fightingBuild.Skill.SuperSkill[pos]
			if skillIID == 0 {
				continue
			}

			skillData, ok := r.prop.Data.SkillsLearned[skillIID]
			if !ok {
				r.Errorf("出战build:%v 对应超能技槽位:%v 已经装备,但是背包中并没有查找到对应的超能技道具:%v", r.prop.Data.FightingBuildID, pos, skillIID)
				return
			}

			//等级生效
			skillValueID := r.calcSkillValueID(skillData.SkillMainID, skillData.Lv)

			superSkills = append(superSkills, skillValueID)
		}
	}

	return
}

//RPC_UpgradeSkill 升级技能
func (user *_User) RPC_UpgradeSkill(skillID uint32) int32 {
	if user.role == nil {
		return ErrorCode.RoleIsNil
	}

	return user.role.UpgradeSkill(skillID)
}

//UpgradeSkill 升级技能
func (r *_Role) UpgradeSkill(skillID uint32) int32 {

	skillData, ok := r.prop.Data.SkillsLearned[skillID]
	if !ok {
		return ErrorCode.BuildSkillNotLearned
	}

	//技能道具
	skillEntityCfg, ok := Data.GetSkillConfig().SkillEntity_ConfigItems[skillID]
	if !ok {
		return ErrorCode.SkillItemNotFind
	}

	//是否可以升级 判断
	errorCode := r.canUpgradeSkill(skillID)
	if errorCode != ErrorCode.Success {
		return errorCode
	}

	//扣除货币
	if skillEntityCfg.UpCoinType > 0 {
		r.UseMoneyByMoneyID(skillEntityCfg.UpCoinType, uint64(skillEntityCfg.UpCoinNum[skillData.Lv]))
	}

	//扣除对应技能材料
	if skillEntityCfg.UpItemID > 0 {
		r.RemoveItem(skillEntityCfg.UpItemID, uint32(Proto.ItemEnum_SkillItem), skillEntityCfg.UpItemNum[skillData.Lv])
	}

	//升级++
	skillData.Lv++
	r.prop.SyncUpdateAlreadyLearnedSkill(skillData)
	r.upgradeSkillRedPoint(skillID)
	//r.Debug("+++++++++++ UpgradeSkill success ", r.getItemNumBySkillBag(item), r.prop.Data.Gold, r.prop.Data.Diamond)
	return ErrorCode.Success
}

//canUpgradeSkill 目标技能是否可以升级
func (r *_Role) canUpgradeSkill(skillID uint32) int32 {
	skillData, ok := r.prop.Data.SkillsLearned[skillID]
	if !ok {
		return ErrorCode.BuildSkillNotLearned
	}

	//技能道具
	skillEntityCfg, ok := Data.GetSkillConfig().SkillEntity_ConfigItems[skillID]
	if !ok {
		return ErrorCode.SkillItemNotFind
	}

	//是否达到最大等级
	if skillData.Lv >= skillEntityCfg.TopSkillLevel {
		return ErrorCode.BuildUpgradeSkillMaxLv
	}

	//验证 等级数量 和 不同等级材料消耗 配置是否对等
	if (skillEntityCfg.UpCoinType > 0 && skillEntityCfg.TopSkillLevel != uint32(len(skillEntityCfg.UpCoinNum))) || (skillEntityCfg.UpItemID > 0 && skillEntityCfg.TopSkillLevel != uint32(len(skillEntityCfg.UpItemNum))) {
		return ErrorCode.BuildUpgradeSkillmaterialsNumNotMatchMaxLv
	}

	//有货币消耗配置 + 货币检测是否满足
	if skillEntityCfg.UpCoinType > 0 && r.GetMoneyByMoneyID(skillEntityCfg.UpCoinType) < uint64(skillEntityCfg.UpCoinNum[skillData.Lv]) {
		return ErrorCode.BuildUpgradeSkillCoinNotEnough
	}

	//有材料消耗配置 升级需要的技能道具数量是否足够  未被使用的技能才能用作升级材料
	if skillEntityCfg.UpItemID > 0 && r.GetItemNum(skillEntityCfg.UpItemID, uint32(Proto.ItemEnum_SkillItem)) < skillEntityCfg.UpItemNum[skillData.Lv] {
		return ErrorCode.BuildUpgradeSkillmaterialsNumNotEnough
	}

	return ErrorCode.Success
}

//learnSkill 学习新技能
func (r *_Role) learnSkill(skillItem *Proto.Item) bool {
	if skillItem == nil || skillItem.Base.Type != Proto.ItemEnum_SkillItem {
		return false
	}

	skillItemCfg, ok := Data.GetSkillConfig().SkillItem_ConfigItems[skillItem.Base.ConfigID]
	if !ok {
		r.Error("掉落技能道具，但是对应配置不存在 ", skillItem.Base.ConfigID)
		return false
	}

	_, alreadyLearned := r.prop.Data.SkillsLearned[skillItemCfg.SkillID]
	if alreadyLearned {
		return false
	}

	//扣除学习消耗掉的道具
	skillItem.Base.Num--

	r.prop.SyncUpdateAlreadyLearnedSkill(&Proto.SkillData{
		SkillMainID: skillItemCfg.SkillID,
		Lv:          1,
	})
	return true
}

//onSkillItemAdded 获取技能道具事件处理
func (r *_Role) onSkillItemAdded(item *Proto.Item) {
	if item == nil {
		return
	}

	skillItemCfg, ok := Data.GetSkillConfig().SkillItem_ConfigItems[item.Base.ConfigID]
	if !ok {
		r.Error("unknown skillItem id:", item.Base.ConfigID)
		return
	}

	r.upgradeSkillRedPoint(skillItemCfg.SkillID)
}

//onSkillItemRemoved 消耗技能道具事件处理
func (r *_Role) onSkillItemRemoved(item *Proto.Item) {
	if item == nil {
		return
	}

	skillItemCfg, ok := Data.GetSkillConfig().SkillItem_ConfigItems[item.Base.ConfigID]
	if !ok {
		r.Error("unknown skillItem id:", item.Base.ConfigID)
		return
	}

	r.upgradeSkillRedPoint(skillItemCfg.SkillID)
}

//onSkillItemUpdated 更新技能道具事件处理
func (r *_Role) onSkillItemUpdated(item *Proto.Item) {
	if item == nil {
		return
	}

	skillItemCfg, ok := Data.GetSkillConfig().SkillItem_ConfigItems[item.Base.ConfigID]
	if !ok {
		r.Error("unknown skillItem id:", item.Base.ConfigID)
		return
	}

	r.upgradeSkillRedPoint(skillItemCfg.SkillID)
}

//checkAllLearnedSkillRedPoint 检测已学技能可升级红点 用于货币变化的时候
func (r *_Role) checkAllLearnedSkillRedPoint() {
	for skillID := range r.prop.Data.SkillsLearned {
		r.upgradeSkillRedPoint(skillID)
	}
}

//upgradeSkillRedPoint 技能升级红点
func (r *_Role) upgradeSkillRedPoint(skillID uint32) {
	skillMainCfg, ok := Data.GetSkillConfig().SkillMain_ConfigItems[skillID]
	if !ok {
		r.Error("unknown SkillMain id:", skillID)
		return
	}

	//当前技能道具是否可以升级
	canUpgrade := r.canUpgradeSkill(skillID) == ErrorCode.Success

	//upgradeSkillRedPointKey 当前技能对应的红点key
	var upgradeSkillRedPointKey string

	switch skillMainCfg.SkillKind {
	case Fight.SkillKind_Ultimate:
		//必杀技处理
		upgradeSkillRedPointKey = r.getRedPointKey(Const.RedPointType_upgradeUltimateSkill, strconv.Itoa(int(skillID)))
	case Fight.SkillKind_Super:
		//超能技处理
		upgradeSkillRedPointKey = r.getRedPointKey(Const.RedPointType_upgradeSuperSkill, strconv.Itoa(int(skillID)))
	default:
		return
	}

	redPointData, ok := r.prop.Data.RedPointsData[upgradeSkillRedPointKey]

	//没有对应红点 + 当前技能可以升级
	if !ok && canUpgrade {
		redPointData = &Proto.RedPointInfo{
			Value:      1,
			CreateTime: time.Now().Unix(),
		}
		r.prop.SyncAddRedPoint(upgradeSkillRedPointKey, redPointData)
	}

	//已经有对应红点 + 现在当前技能不可升级
	if ok && !canUpgrade {
		r.prop.SyncRemoveRedPoint(upgradeSkillRedPointKey)
	}
}
