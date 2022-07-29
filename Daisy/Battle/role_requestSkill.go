package main

import (
	"Cinder/Base/Const"
	"Cinder/Base/Util"
	"Daisy/Data"
	"Daisy/DataTables"
	"Daisy/ErrorCode"
	"Daisy/Fight"
	"Daisy/Prop"
	"Daisy/Proto"
	"fmt"
	"math/rand"
	"time"
)

const superSkillReceiveLimit uint32 = 6

const giveSuperSkillToTargetLimit uint32 = 4

const ultimateSkillReceiveLimit uint32 = 1

const giveUltimateSkillToTargetLimit uint32 = 1

//RPC_RequestSkill RPC乞求技能
func (user *_User) RPC_RequestSkill(skillItemId uint32) int32 {
	if user.role == nil {
		return ErrorCode.RoleIsNil
	}

	return user.role.RequestSkill(skillItemId)
}

// RequestSkill 乞求技能
func (r *_Role) RequestSkill(skillItemId uint32) int32 {
	team := r.GetSpace().(*_Team)
	requestSkillCount := r.prop.Data.ShareSpoils.RequestSkillCount
	if requestSkillCount >= requestSkillCountLimit {
		team.Infof("RequestSkill role:%s 今日乞求次数已经超过上限", r.GetID())
		return ErrorCode.RequestCountOverLimit
	}

	if len(team.prop.Data.Base.Members) == 1 {
		team.Infof("RequestSkill role:%s 不在队伍中", r.GetID())
		return ErrorCode.FirstJoinATeam
	}

	var skillItemCfg *DataTables.SkillItem_Config
	var skillEntityCfg *DataTables.SkillEntity_Config
	var skillMainCfg *DataTables.SkillMain_Config
	var ok bool
	skillItemCfg, ok = Data.GetSkillConfig().SkillItem_ConfigItems[skillItemId]
	if !ok {
		team.Errorf("RequestSkill role:%s 技能%d SkillItem配置找不到", r.GetID(), skillItemCfg.SkillID)
		return ErrorCode.SkillItemNotFind
	}

	skillEntityCfg, ok = Data.GetSkillConfig().SkillEntity_ConfigItems[skillItemCfg.SkillID]
	if !ok {
		team.Errorf("RequestSkill role:%s 技能%d SkillEntity配置找不到", r.GetID(), skillItemCfg.SkillID)
		return ErrorCode.SkillItemNotFind
	}

	skillMainCfg, ok = Data.GetSkillConfig().SkillMain_ConfigItems[skillItemCfg.SkillID]
	if !ok {
		team.Errorf("RequestSkill role:%s 技能%d SkillMain配置找不到", r.GetID(), skillItemCfg.SkillID)
		return ErrorCode.SkillItemNotFind
	}

	skillLevel := uint32(0)
	if skillData, ok := r.prop.Data.SkillsLearned[skillItemCfg.SkillID]; ok {
		//是否达到最大等级
		if skillData.Lv >= skillEntityCfg.TopSkillLevel {
			team.Infof("RequestSkill role:%s 乞求技能%d已经达到最大等级 ", r.GetID(), skillItemCfg.SkillID)
			return ErrorCode.BuildUpgradeSkillMaxLv
		}

		skillLevel = skillData.Lv
	}

	rollNum := uint32(rand.Intn(10000))
	sum := uint32(0)
	messageId := uint32(1)
	for id, conf := range Data.GetBeggingConfig().Begging_ConfigItems {
		if rollNum >= sum && rollNum < sum+conf.Probability {
			messageId = id
			break
		}

		sum += conf.Probability
	}

	skillReceivedLimit := superSkillReceiveLimit
	if skillMainCfg.SkillKind == Fight.SkillKind_Super {
		skillReceivedLimit = superSkillReceiveLimit
	}

	if skillMainCfg.SkillKind == Fight.SkillKind_Ultimate {
		skillReceivedLimit = ultimateSkillReceiveLimit
	}

	requestSkillCount++
	uid := Util.GetGUID()
	requestSkill := &Proto.RequestSkill{
		RequestMessageId:        messageId,
		RequestCreateTimeStamp:  time.Now().Unix(),
		RequestSkillItemId:      skillItemId,
		SkillCurLevel:           skillLevel,
		ReceivedSkillCount:      0,
		SkillReceivedCountLimit: skillReceivedLimit,
		RequestTimeOutTimeStamp: r.prop.Data.ShareSpoils.RequestSkillResetTimeStamp,
	}

	r.prop.SyncAddRequestSkill(uid, requestSkill, requestSkillCount)

	user := r.GetOwnerUser()
	if user == nil {
		r.Error("[findNewTitle] role's user is nil")
	} else {
		user.Rpc(Const.Game, "RPC_RequestSkillMsg", r.GetID(), uid)
		// todo 拿到返回值并且做判断，返回Fail或者Success
	}

	team.Infof("RequestSkill role:%s 发送技能请求消息成功", r.GetID())
	return ErrorCode.Success
}

//RPC_GiveSkill RPC捐赠技能
func (user *_User) RPC_GiveSkill(targetId, uid string) int32 {
	if user.role == nil {
		return ErrorCode.RoleIsNil
	}

	return user.role.GiveSkill(targetId, uid)
}

// GiveSkill 捐赠技能
func (r *_Role) GiveSkill(targetId, uid string) int32 {
	team := r.GetSpace().(*_Team)
	if r.GetID() == targetId {
		team.Infof("RequestSkill role:%s 无法向自己捐赠技能 ", r.GetID())
		return ErrorCode.CantGiveSkillToSelf
	}

	var ok bool
	target, err := team.GetActor(targetId)
	if err != nil {
		team.Infof("RequestSkill role:%s 赠送对象%s已不在队伍中 ", r.GetID(), targetId)
		return ErrorCode.TargetNotInTeam
	}

	targetProp := target.GetProp().(*Prop.RoleProp)
	var requestSkill *Proto.RequestSkill
	requestSkill, ok = targetProp.Data.ShareSpoils.RequestSkillData[uid]
	if !ok {
		team.Infof("RequestSkill 玩家%s的技能请求%s已失效 ", targetId, uid)
		return ErrorCode.SkillRequestNotValid
	}

	if r.GetItemNum(requestSkill.RequestSkillItemId, uint32(Proto.ItemEnum_SkillItem)) <= 0 {
		team.Infof("RequestSkill 玩家%s没有技能%d", r.GetID(), requestSkill.RequestSkillItemId)
		return ErrorCode.DontHaveSkill
	}

	var skillItemCfg *DataTables.SkillItem_Config
	var skillMainCfg *DataTables.SkillMain_Config

	skillItemCfg, ok = Data.GetSkillConfig().SkillItem_ConfigItems[requestSkill.RequestSkillItemId]
	if !ok {
		team.Errorf("RequestSkill SkillItem表中不存在%d", requestSkill.RequestSkillItemId)
		return ErrorCode.SkillItemNotFind
	}

	skillMainCfg, ok = Data.GetSkillConfig().SkillMain_ConfigItems[skillItemCfg.SkillID]
	if !ok {
		team.Errorf("RequestSkill SkillMain表中不存在%d", requestSkill.RequestSkillItemId)
		return ErrorCode.SkillItemNotFind
	}

	giveTargetSkillCount := uint32(0)
	var giveSkillCount *Proto.GiveSkillCount
	giveSkillCount, ok = r.prop.Data.ShareSpoils.GiveSkillToTargetCount[targetId]
	if ok {
		if count, ok := giveSkillCount.GiveSkillToTargetCount[requestSkill.RequestSkillItemId]; ok {
			if skillMainCfg.SkillKind == Fight.SkillKind_Super {
				if count >= giveSuperSkillToTargetLimit {
					team.Infof("RequestSkill role%s赠送给目标%s的技能%d数量超过上限", r.GetID(), targetId, requestSkill.RequestSkillItemId)
					return ErrorCode.GiveSKillToTargetCountReachLimit
				}
			}

			if skillMainCfg.SkillKind == Fight.SkillKind_Ultimate {
				if count >= giveUltimateSkillToTargetLimit {
					team.Infof("RequestSkill role%s赠送给目标%s的技能%d数量超过上限", r.GetID(), targetId, requestSkill.RequestSkillItemId)
					return ErrorCode.GiveSKillToTargetCountReachLimit
				}
			}

			giveTargetSkillCount = count
		}

	}

	if skillMainCfg.SkillKind == Fight.SkillKind_Super {
		if requestSkill.ReceivedSkillCount >= superSkillReceiveLimit {
			team.Infof("RequestSkill 目标%s收到的技能道具%d数量已达到上限", targetId, requestSkill.RequestSkillItemId)
			return ErrorCode.ReceiveSkillCountReachLimit
		}
	}

	if skillMainCfg.SkillKind == Fight.SkillKind_Ultimate {
		if requestSkill.ReceivedSkillCount >= ultimateSkillReceiveLimit {
			team.Infof("RequestSkill 目标%s收到的技能道具%d数量已达到上限", targetId, requestSkill.RequestSkillItemId)
			return ErrorCode.ReceiveSkillCountReachLimit
		}
	}

	if r.prop.Data.ShareSpoils.GiveSkillCount >= giveSkillCountLimit {
		team.Infof("RequestSkill role%s赠送技能数量超过上限", r.GetID())
		return ErrorCode.GiveSkillCountReachLimit
	}

	container := r.GetItemContainerByItemType(uint32(Proto.ItemEnum_SkillItem))
	if container == nil {
		return ErrorCode.SkillContainerNotFound
	}

	// 移除技能
	_, ok = container.RemoveItemNum(requestSkill.RequestSkillItemId, uint32(Proto.ItemEnum_SkillItem), 1)
	if !ok {
		return ErrorCode.RemoveSkillFailed
	}

	teamMateSkillCount := uint32(0)
	if requestSkill, ok := targetProp.Data.ShareSpoils.RequestSkillData[uid]; ok {
		if _, ok := requestSkill.TeamMateGiveSkillCount[r.GetID()]; ok {
			teamMateSkillCount = requestSkill.TeamMateGiveSkillCount[r.GetID()]
		}
	}
	targetProp.SyncSetReceivedSkillCount(uid, requestSkill.ReceivedSkillCount+1, r.GetID(), teamMateSkillCount+1)
	r.prop.SyncSetGiveSkillCount(r.prop.Data.ShareSpoils.GiveSkillCount+1, targetId, requestSkill.RequestSkillItemId, giveTargetSkillCount+1)

	var items []*Proto.MailAttachment
	items = append(items, &Proto.MailAttachment{ItemID: requestSkill.RequestSkillItemId, ItemType: uint32(Proto.ItemEnum_SkillItem), ItemNum: 1})

	// 发送邮件
	r.SendMail(fmt.Sprintf("来自%s的礼物", r.prop.Data.Base.Name), fmt.Sprintf("%s送了你一份礼物，快去看看吧", r.prop.Data.Base.Name), targetId, "", false, items)

	r.addCommanderExp(skillItemCfg.GiftGotCommanderExp)                       // 增加指挥官经验
	r.addFightingSpecialAgentExp(uint64(skillItemCfg.GiftGotSpecialAgentExp)) // 增加当前上阵特工经验
	r.prop.SyncAddGold(skillItemCfg.GiftGotGold)                              // 增加金币

	team.Infof("RequestSkill role:%s RPC捐赠技能消息发送成功", r.GetID())
	return ErrorCode.Success
}
