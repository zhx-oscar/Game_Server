package main

import (
	"Daisy/Battle/drop"
	"Daisy/Const"
	"Daisy/Data"
	"Daisy/ErrorCode"
	"Daisy/ItemProto"
	"Daisy/Proto"
)

func (user *_User) RPC_PurchaseBox(id uint32) int32 {
	if user.role == nil {
		return ErrorCode.RoleIsNil
	}
	return user.role.PurchaseLogic(id)
}

func (user *_User) RPC_OpenBox(id uint32) (*Proto.SupplyAwardData, bool, int32){
	if user.role == nil {
		return nil,false,ErrorCode.RoleIsNil
	}
	return user.role.OpenBoxLogic(id)
}

func (r *_Role) OpenBoxLogic(id uint32)(*Proto.SupplyAwardData, bool, int32){
	tmp, ok := Data.GetSupplyConfig().SupplyBox_ConfigItems[id]
	if !ok{
		return nil, false, ErrorCode.GiftBagUnknown
	}
	str, k := r.prop.Data.SupplyInfo.BoxMold[id]
	if !k || str.HoldNum <= 0{
		return nil, false, ErrorCode.GiftBagNotEnough
	}
	if str.OpenNum < tmp.OpenLimit {
		dropID := Data.GetSupplyConfig().SupplyBox_ConfigItems[id].DropID
		tmpDrop := &drop.Drop{}
		ok, items := tmpDrop.Drop(dropID, 0, 0)
		if !ok {
			return nil, false, ErrorCode.Failure
		}
		award, flag := r.AddItemListOrMail(items, "补给", "奖励")
		r.prop.SyncChangeSupplyNum(id, 1)
		return award, flag, ErrorCode.Success
	}else{
		return nil, false, ErrorCode.GiftBagOpenUpperLimit
	}
}

func (r *_Role) PurchaseLogic(id uint32) int32{
	tmp, ok := Data.GetSupplyConfig().SupplyDiscount_ConfigItems[id]
	if !ok{
		return ErrorCode.GiftBagUnknown
	}
	if tmp.UpperLimit == 0{
		return r.Purchase(id)
	}else{
		if r.prop.Data.SupplyInfo.DiscountNum[id] < tmp.UpperLimit{
			return r.DiscountPurchase(id)
		}else{
			return r.Purchase(id)
		}
	}
}

func (r *_Role) DiscountPurchase(id uint32) int32{
	tmp := Data.GetSupplyConfig().SupplyDiscount_ConfigItems[id]
	switch tmp.CurrencyType {
	case Const.Gold:
		if r.RemoveGold(tmp.DiscountPrice, Const.SupplyBoxCost) == true{
			r.prop.SyncAddSupplyNum(tmp.SupplyTypeID,tmp.SupplyCount)
			r.prop.SyncAddDiscountNum(id,1)
			return ErrorCode.Success
		}else{
			return ErrorCode.GoldNotEnough
		}
	case Const.Diamond:
		if r.RemoveDiamond(tmp.DiscountPrice, Const.SupplyBoxCost) == true{
			r.prop.SyncAddSupplyNum(tmp.SupplyTypeID, tmp.SupplyCount)
			r.prop.SyncAddDiscountNum(id,1)
			return ErrorCode.Success
		}else{
			return ErrorCode.DiamondNotEnough
		}
	default:
		return ErrorCode.Failure
	}
}

func (r *_Role) Purchase(id uint32) int32{
	tmp := Data.GetSupplyConfig().SupplyDiscount_ConfigItems[id]
	switch tmp.CurrencyType {
	case Const.Gold:
		if r.RemoveGold(tmp.OriginalCost,Const.SupplyBoxCost) == true{
			r.prop.SyncAddSupplyNum(tmp.SupplyTypeID, tmp.SupplyCount)
			return ErrorCode.Success
		}else{
			return ErrorCode.GoldNotEnough
		}
	case Const.Diamond:
		if r.RemoveDiamond(tmp.OriginalCost,Const.SupplyBoxCost) == true{
			r.prop.SyncAddSupplyNum(tmp.SupplyTypeID, tmp.SupplyCount)
			return ErrorCode.Success
		}else{
			return ErrorCode.DiamondNotEnough
		}
	default:
		return ErrorCode.Failure
	}
}

func (r *_Role) AddItemListOrMail(items []*Proto.DropMaterial, title string, content string)(*Proto.SupplyAwardData, bool){
	_items := items
	award := &Proto.SupplyAwardData{}
	item := r.GetItemsTypeSum(_items)
	var flag bool
	var left []*Proto.DropMaterial
	data := r.AddItemList(items)
	if len(data) != 0 {
		_data := data
		for _, val := range data {
			awardData := &Proto.SupplyAwardItem{
				ID:       val.ItemData.Base.ID,
				Type:     val.ItemData.Base.Type,
				Num:      val.GetNum,
			}
			award.SupplyAwardItems = append(award.SupplyAwardItems, awardData)
		}
		database := r.GetDataTypeSum(_data)
		for k, v := range item {
			for _, v1 := range database {
				if v.MaterialType == uint32(v1.ItemData.Base.Type) && v.MaterialId == v1.OldConfigID{
					if item[k].MaterialNum > v1.GetNum {
						material := &Proto.DropMaterial{
							MaterialId:   v.MaterialId,
							MaterialType: v.MaterialType,
							MaterialNum:  item[k].MaterialNum - v1.GetNum,
						}
						item[k].MaterialNum = 0
						left = append(left, material)
					}else{
						item[k].MaterialNum = 0
					}
				}
			}
		}
		for k1, val := range item{
			if item[k1].MaterialNum != 0{
				left = append(left,val)
			}
		}
	}else{
		left = item
	}
	if len(left) != 0{
		leftItems := r.CreatProp(left)
		list := make([]*Proto.MailAttachment, 0)
		for _,v := range leftItems {
			list = append(list, &Proto.MailAttachment{
				Data:v.GetData(),
			})
		}
		r.SendMail(title, content, r.GetID(), r.prop.Data.Base.Name,false, list)
		for _, v := range list {
			award.SupplyAwardMail = append(award.SupplyAwardMail, v.Data)
		}
		flag = true
	}
	return award, flag
}

func (r *_Role) GetItemsTypeSum(_items []*Proto.DropMaterial)[]*Proto.DropMaterial{
	for i := 0; i < len(_items); i++{
		for j := i+1; j < len(_items); j++{
			if _items[i].MaterialType == _items[j].MaterialType && _items[i].MaterialId == _items[j].MaterialId{
				_items[i].MaterialNum += _items[j].MaterialNum
				_items[j].MaterialNum = 0
			}
		}
	}
	var ret []*Proto.DropMaterial
	for _, v := range _items{
		if v.MaterialNum != 0 {
			ret = append(ret, v)
		}
	}
	return ret
}

func(r *_Role) GetDataTypeSum(_data []*Proto.GetItemData)[]*Proto.GetItemData{
	for i := 0; i < len(_data); i++{
		for j := i+1; j < len(_data); j++{
			if _data[i].ItemData.Base.Type == _data[j].ItemData.Base.Type && _data[i].OldConfigID == _data[j].OldConfigID{
				_data[i].GetNum += _data[j].GetNum
				_data[j].GetNum = 0
			}
		}
	}
	var ret []*Proto.GetItemData
	for _, v := range _data{
		if v.GetNum!= 0 {
			ret = append(ret, v)
		}
	}
	return ret
}

func (r *_Role) CreatProp(left []*Proto.DropMaterial)[]ItemProto.IItem{
	var ret []ItemProto.IItem
	for _, v := range left{
		maxNum := ItemProto.GetItemMaxPileNum(v.MaterialId, v.MaterialType)
		var num uint32
		if v.MaterialNum%maxNum != 0 {
			num = v.MaterialNum/maxNum + 1
		}else {
			num = v.MaterialNum/maxNum
		}
		for i := 0; uint32(i) < num; i++{
			if v.MaterialNum > maxNum{
				_item := r.CreateItem(v.MaterialId, v.MaterialType, maxNum)
				if _item != nil {
					ret = append(ret, _item)
				} else {
					r.Error("找不到道具:", v)
				}
				v.MaterialNum -= maxNum
			}else{
				_item := r.CreateItem(v.MaterialId, v.MaterialType, v.MaterialNum)
				if _item != nil {
					ret = append(ret, _item)
				} else {
					r.Error("找不到道具:", v)
				}
			}
		}
	}
	return ret
}
