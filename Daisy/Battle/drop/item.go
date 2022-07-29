package drop

import (
	"Cinder/Base/Util"
	"Daisy/Data"
	"Daisy/Proto"
)

//CreateItem 创建道具数据
func (drop *Drop) CreateItem(configID uint32, typ uint32, num uint32, args ...uint32) *Proto.Item {
	item := Proto.Item{}
	item.Base = &Proto.ItemBase{}
	item.ExpandData = &Proto.ItemExpand{
		InUse:             0,
		RelationBuildList: map[string]bool{},
	}
	_typ := Proto.ItemEnum_Type(typ)
	item.Base.ID = Util.GetGUID()
	item.Base.Type = _typ
	item.Base.Num = num
	item.Base.ConfigID = configID

	switch _typ {
	case Proto.ItemEnum_Equipment:
		_,ok := Data.GetEquipConfig().EquipMent_ConfigItems[configID]
		if ok == false {
			return nil
		}
		drop.MakeEquipment(&item, args...)
	case Proto.ItemEnum_SkillItem:
		_, ok := Data.GetSkillConfig().SkillItem_ConfigItems[configID]
		if ok == false {
			return nil
		}
		item.SkillItemData = &Proto.SkillItem{}
	}
	return &item
}
