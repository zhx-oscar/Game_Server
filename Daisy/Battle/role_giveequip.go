package main

import (
	"Daisy/Const"
	"Daisy/ErrorCode"
	"Daisy/Proto"
	"fmt"
	"math"
)

const characterCountLimit uint32 = 80

const giveEquipCountLimit uint32 = 6

//RPC_GiveEquip RPC赠送装备
func (user *_User) RPC_GiveEquip(itemId, targetId, message string) (int32, uint32) {
	if user.role == nil {
		return ErrorCode.RoleIsNil, 0
	}

	return user.role.GiveEquip(itemId, targetId, message)
}

//GiveEquip 赠送装备
func (r *_Role) GiveEquip(itemId, targetId, message string) (int32, uint32) {
	giveEquipCount := r.prop.Data.ShareSpoils.GiveEquipCount
	leftGiveCount := uint32(math.Max(float64(giveEquipCountLimit-giveEquipCount), 0))
	container := r.GetItemContainerByItemType(uint32(Proto.ItemEnum_Equipment))
	if container == nil {
		return ErrorCode.EquipContainerNotFound, 0
	}

	item := r.GetItemFromPack(itemId)
	if item == nil {
		return ErrorCode.EquipNotExist, leftGiveCount
	}

	if r.GetID() == targetId {
		return ErrorCode.CantGiveToSelf, leftGiveCount
	}

	targetFound := false
	for _, info := range item.GetData().EquipmentData.OwnerTeamMemberList {
		if info.ID == targetId {
			targetFound = true
			break
		}
	}

	if !targetFound {
		return ErrorCode.TargetNotInList, leftGiveCount
	}

	if Const.UTF8Width(message) > int(characterCountLimit) {
		return ErrorCode.CharacterCountExceed, leftGiveCount
	}

	for buildId := range r.prop.Data.BuildMap {
		if _, ok := item.GetData().ExpandData.RelationBuildList[buildId]; ok {
			return ErrorCode.EquipUsedInBuild, leftGiveCount
		}
	}

	if giveEquipCount >= giveEquipCountLimit {
		return ErrorCode.GiveCountOverLimit, leftGiveCount
	}

	// 移除装备
	_item := container.RemoveItemByID(itemId)
	if _item == nil {
		return ErrorCode.RemoveEquipFailed, leftGiveCount
	}

	if len(message) == 0 {
		message = fmt.Sprintf("%s送了你一份礼物，快去看看吧", r.prop.Data.Base.Name)
	}

	var items []*Proto.MailAttachment
	items = append(items, &Proto.MailAttachment{Data: _item.GetData()})

	// 发送邮件
	r.SendMail(fmt.Sprintf("来自%s的礼物", r.prop.Data.Base.Name), message, targetId, "", false, items)

	giveEquipCount++
	r.prop.SyncSetGiveEquipCount(giveEquipCount)
	leftGiveCount--

	return ErrorCode.Success, leftGiveCount
}
