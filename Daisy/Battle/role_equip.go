package main

import (
	"Daisy/Const"
	"Daisy/Data"
	"Daisy/ErrorCode"
	"Daisy/Fight/attraffix"
	"Daisy/ItemProto"
	"Daisy/Proto"
	"fmt"
)

type IAttribute interface {
	Calc(pos int32, item *Proto.Item) bool
}

type AttributeCalc struct {
	AttrAffixList []attraffix.AttrAffix // 属性词条列表
	BornBuffs     []uint32              // 出生buff列表
	suit          map[uint32]uint32     //套装
}

func (tc *AttributeCalc) Init() {
	tc.AttrAffixList = make([]attraffix.AttrAffix, 0)
	tc.BornBuffs = make([]uint32, 0)
	tc.suit = make(map[uint32]uint32)
}

func (tc *AttributeCalc) Calc(pos int32, item *Proto.Item) bool {
	//计算词缀属性
	for _, v := range item.EquipmentData.Affixes {
		if v.AffixEffectType == Const.AffixEffectType_Attr {
			data := attraffix.AttrAffix{}
			data.Field = attraffix.Field(v.PropertyID)
			if v.AffixParam == Const.AffixValueType_B {
				data.ParaB = float64(v.Value)
			} else {
				data.ParaA = v.Value
			}
			tc.AttrAffixList = append(tc.AttrAffixList, data)
		} else {
			tc.BornBuffs = append(tc.BornBuffs, v.PropertyID)
		}
	}

	//计算套装
	return true
}

func (user *_User) RPC_BuildEquipItem(buildID string, pos int32, itemID string) int32 {
	user.Debug(fmt.Sprintf("RPC_BuildEquipItem %s %d %s", buildID, pos, itemID))

	if user.role == nil {
		return ErrorCode.RoleIsNil
	}
	return user.role.BuildEquipItem(buildID, pos, itemID)
}

func (r *_Role) BuildEquipItem(buildID string, pos int32, itemID string) int32 {
	build, ok := r.prop.Data.BuildMap[buildID]
	if !ok {
		return ErrorCode.NotFindBuild
	}

	//特工是否拥有
	if build.SpecialAgentID == 0 {
		return ErrorCode.NoSpecialAgent
	}
	SpecialAgent, found := r.prop.Data.SpecialAgentList[build.SpecialAgentID]
	if found == false {
		return ErrorCode.NoSpecialAgent
	}

	var item ItemProto.IItem
	if item = r.GetItemFromPack(itemID); item == nil {
		return ErrorCode.NotFindBuildEquip
	}

	equipConfig, found := Data.GetEquipConfig().EquipMent_ConfigItems[item.GetConfigID()]
	if found == false {
		return ErrorCode.NotFindBuildEquip
	}

	//等级
	if equipConfig.ItemLevel > SpecialAgent.Base.Level {
		return ErrorCode.BuildEquipNotEnoughLevel
	}

	//位置
	if r.CheckBuildEquipmentPosition(equipConfig.Position, pos) == false {
		return ErrorCode.BuildEquipPosWrong
	}
	//同一个戒指不能放两个槽
	if Proto.ContainerEnum_Ring1 == Proto.ContainerEnum_EquipPos(equipConfig.Position) {
		_id, has := build.EquipmentMap[int32(Proto.ContainerEnum_Ring2)]
		if has == true && _id == itemID {
			return ErrorCode.BuildEquipPosWrong
		}
	} else if Proto.ContainerEnum_Ring2 == Proto.ContainerEnum_EquipPos(equipConfig.Position) {
		_id, has := build.EquipmentMap[int32(Proto.ContainerEnum_Ring1)]
		if has == true && _id == itemID {
			return ErrorCode.BuildEquipPosWrong
		}
	}

	//装备 引用处理
	//build 先移除 后增加
	//老build 内道具 引用处理
	oldEquipID, has := build.EquipmentMap[pos]
	if has == true {
		//checkBagSpace build里面如果涉及到装备回放外部背包检测
		if !r.HasEnoughSpace(1, int32(Proto.ContainerEnum_EquipBag)) {
			return ErrorCode.EquipBagNotEnoughSpace
		}

		oldItem := r.GetItemFromPack(oldEquipID)
		if oldItem == nil {
			r.Error("被build装备的装备道具 但是在build装备背包中未找到。数据错乱!!!")
		} else {
			if oldItem.GetData().ExpandData.InUse == 0 {
				r.Error("build装备背包中 装备道具 引用计数应该不为 0 。数据错乱!!!")
			} else {
				r.itemDelInUse(oldItem, build.BuildID)
			}
		}
	}

	//新build 内道具 引用处理
	r.itemAddInUse(item, build.BuildID)

	r.prop.SyncBuildEquipItem(buildID, pos, itemID)
	r.refreshBuildFightAttr(buildID)
	return ErrorCode.Success
}

//CalcEquipAttr 遍历计算装备属性
func (r *_Role) CalcEquip(buildID string) *AttributeCalc {
	attribute := &AttributeCalc{}

	build, has := r.prop.Data.BuildMap[buildID]
	if has == false {
		return attribute
	}

	for pos, id := range build.EquipmentMap {
		item := r.GetItemFromPack(id)
		if item == nil {
			continue
		}
		if ok := attribute.Calc(pos, item.GetData()); ok == false {
			break
		}
	}
	return attribute
}

func (r *_Role) CheckBuildEquipmentPosition(configPos uint32, buildPos int32) bool {
	if configPos == 7 {
		return Proto.ContainerEnum_Ring1 == Proto.ContainerEnum_EquipPos(buildPos) || Proto.ContainerEnum_Ring2 == Proto.ContainerEnum_EquipPos(buildPos)
	} else {
		return configPos == uint32(buildPos)
	}
}

func (r *_Role) CanEquipItem(buildPos int32, item *Proto.Item) bool {
	//Todo 判断等级位置等等
	return true
}

type BuffAffix struct {
	BuffID    uint32
	Score     uint32
	ScoreType uint32
}

type EquipAffix struct {
	AttrList []attraffix.AttrAffix
	BuffList []BuffAffix
}

func (r *_Role) GetEquipAttrAffix(itemID string) EquipAffix {
	var data EquipAffix
	data.AttrList = make([]attraffix.AttrAffix, 0)
	data.BuffList = make([]BuffAffix, 0)

	item := r.GetItemFromPack(itemID)
	for _, v := range item.GetData().EquipmentData.Affixes {
		if v.AffixEffectType == Const.AffixEffectType_Attr {
			//属性
			attr := attraffix.AttrAffix{
				Field: attraffix.Field(v.PropertyID),
			}
			if v.AffixParam == Const.AffixValueType_B {
				attr.ParaB = float64(v.Value)
			} else {
				attr.ParaA = v.Value
			}
			data.AttrList = append(data.AttrList, attr)
		} else if v.AffixEffectType == Const.AffixEffectType_Buff {
			buff := BuffAffix{
				BuffID: v.PropertyID,
			}
			config, ok := Data.GetEquipConfig().EquipAffix_ConfigItems[v.AffixID]
			if ok == true {
				buff.Score = config.BuffAffixScore
				buff.ScoreType = config.BuffAffixType
			}
			data.BuffList = append(data.BuffList, buff)
		}
	}
	return data
}

func (r *_Role) GetLucky() uint32 {
	build, ok := r.prop.Data.BuildMap[r.prop.Data.FightingBuildID]
	if ok == false {
		return 0
	}
	return uint32(build.FightAttr.Lucky)
}

//RPC_SellEquip 售卖装备，传入装备id，返回得到的金币和错误码
func (user *_User) RPC_SellEquip(itemId string) (uint32, int32) {
	user.Debug(fmt.Sprintf("RPC_SellEquip %s", itemId))

	if user.role == nil {
		return 0, ErrorCode.RoleIsNil
	}

	return user.role.SellEquip([]string{itemId})
}

func (r *_Role) SellEquip(list []string) (uint32, int32) {
	var money uint32
	container := r.GetItemContainerByItemType(uint32(Proto.ItemEnum_Equipment))
	if container == nil {
		return 0, -1
	}
	currencyType := uint32(0)
	removeList := make([]ItemProto.IItem, 0)
	for _, v := range list {
		item := container.GetIItemByID(v)
		if item == nil {
			continue
		}
		config, ok := Data.GetEquipConfig().EquipMent_ConfigItems[item.GetConfigID()]
		if ok == false {
			continue
		}
		currencyType = config.CurrencyType
		money += config.Price
		removeList = append(removeList, item)
	}

	container.RemoveItemList(removeList)
	r.AddMoney(money, currencyType, Const.SELL_EQUIPMENT)
	return money, ErrorCode.Success
}

//GetFightingBDEquips 获取上阵特工的装备列表，用于战报
func (r *_Role) GetFightingBDEquips() map[int32]uint32 {
	equips := make(map[int32]uint32)
	if bd, ok := r.prop.Data.BuildMap[r.prop.Data.FightingBuildID]; ok {
		for key, value := range bd.EquipmentMap {
			if item := r.GetItemFromPack(value); item != nil {
				equips[key] = item.GetConfigID()
			}
		}
	}

	return equips
}
