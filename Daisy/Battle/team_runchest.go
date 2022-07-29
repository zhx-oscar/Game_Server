package main

import (
	"Cinder/Base/Const"
	"Cinder/Base/User"
	"Daisy/Battle/sceneconfigdata"
	"Daisy/ErrorCode"
	"Daisy/ItemProto"
	"Daisy/Prop"
	"Daisy/Proto"
	"math/rand"
	"time"
)

type _RunChestModel struct {
	items        map[string]uint32
	chestEndTime time.Time
}

func (team *_Team) EnterRunChest(param *sceneconfigdata.TriggerParamChest) {
	if team.runChest != nil {
		team.Errorf("EnterRunChest 上一个宝箱未结束")
		return
	}

	//根据概率是否刷新宝箱
	if int(param.Probability) <= rand.Intn(10000) {
		return
	}

	team.runChest = &_RunChestModel{
		items: make(map[string]uint32),
	}
	for key := range team.prop.Data.Base.Members {
		team.runChest.items[key] = param.DropID
	}
	team.runChest.chestEndTime = time.Now().Add(10 * time.Second)

	chestLoc := &Proto.PVector3{X: param.ChestLoc.X, Y: param.ChestLoc.Y, Z: param.ChestLoc.Z}
	chestRot := &Proto.PVector3{X: param.ChestRot.X, Y: param.ChestRot.Y, Z: param.ChestRot.Z}
	team.TraversalUser(func(user User.IUser) bool {
		user.Rpc(Const.Agent, "RPC_RunChestAppear", chestLoc, chestRot)
		return true
	})

	team.Info("EnterRunChest")
}

func (team *_Team) LeaveRunChest() {
	if team.runChest == nil {
		team.Errorf("LeaveRunChest 宝箱已经结束")
		return
	}

	for key := range team.prop.Data.Base.Members {
		if dropID, ok := team.runChest.items[key]; ok {
			team.lootRunChest(key, dropID)
		}
	}

	team.runChest = nil

	team.Info("LeaveRunChest")
}

func (team *_Team) LoopRunChest() {
	if team.runChest != nil && time.Now().After(team.runChest.chestEndTime) {
		team.LeaveRunChest()
		team.Error("跑图宝箱超时")
	}
}

func (team *_Team) OnLootRunChest(roleID string) (int32, *Proto.Items) {
	if team.runChest == nil {
		team.Errorf("roleID:%s 非法拾取该宝箱")
		return ErrorCode.InvalidLootRunChest, nil
	}

	dropID, ok := team.runChest.items[roleID]
	if !ok {
		team.Errorf("roleID:%s 非法拾取该宝箱")
		return ErrorCode.InvalidLootRunChest, nil
	}

	items := team.lootRunChest(roleID, dropID)
	delete(team.runChest.items, roleID)

	finish := true
	for key := range team.runChest.items {
		if _, err := team.GetUser(key); err == nil {
			finish = false
		}
	}

	if finish {
		team.LeaveRunChest()
	}

	team.Infof("OnLootRunChest roleID:%s success", roleID)
	return ErrorCode.Success, items
}

func (team *_Team) lootRunChest(roleID string, dropID uint32) *Proto.Items {
	ia := team.GetActorByUserID(roleID)
	if ia == nil {
		team.Error("get actor err:")
		return nil
	}

	role := ia.(*_Role)
	items, err := role.Drop(dropID)
	if err != nil {
		role.Error("生成掉落err:", err)
		return nil
	}

	// 创建队伍成员信息列表
	var memberInfoList []*Proto.OwnerInfo
	for id := range team.prop.Data.Base.Members {
		owner := &Proto.OwnerInfo{
			ID:   id,
			Name: team.GetActorByUserID(id).GetProp().(*Prop.RoleProp).Data.Base.Name,
		}
		memberInfoList = append(memberInfoList, owner)
	}

	itemList := &Proto.Items{}
	itemData := role.AddItemList(items)
	for _, val := range itemData {
		if val.GetNum == 0 {
			continue
		}
		itemList.Items = append(itemList.Items, val.ItemData)
		if val.ItemData.Base.Type == Proto.ItemEnum_Equipment {
			val.ItemData.EquipmentData.OwnerTeamMemberList = memberInfoList
			role.UpdataItem(ItemProto.CreateIItemByData(val.ItemData))
		}
	}

	return itemList
}
