package internal

import (
	"Daisy/Data"
	"Daisy/Proto"
)

func CreateIItemByData(data *Proto.Item) IItem {
	if data == nil {
		return nil
	}
	return &ItemData{
		Data: data,
	}
}

func CreateIItemByID(configID, typ uint32) IItem {
	item := &Proto.Item{
		Base: &Proto.ItemBase{
			ConfigID: configID,
			Type:     Proto.ItemEnum_Type(typ),
		},
	}
	return &ItemData{
		Data: item,
	}
}

func GetItemMaxNum(configID, typ uint32) uint32 {
	switch Proto.ItemEnum_Type(typ) {
	case Proto.ItemEnum_SkillItem:
		config,ok := Data.GetSkillConfig().SkillItem_ConfigItems[configID]
		if ok == true{
			return config.MaxNum
		}
	}
	return 0
}

func GetItemMaxPileNum(configID, typ uint32) uint32 {
	switch Proto.ItemEnum_Type(typ) {
	case Proto.ItemEnum_SkillItem:
		config,ok := Data.GetSkillConfig().SkillItem_ConfigItems[configID]
		if ok == true{
			return config.MaxPileNum
		}
	}
	return 1
}

type ItemData struct {
	Data *Proto.Item
}

func (i *ItemData) SetData(data *Proto.Item) {
	i.Data = data
}

func (i *ItemData) GetData() *Proto.Item {
	return i.Data
}

func (i *ItemData) GetID() string {
	return i.Data.Base.ID
}

func (i *ItemData) GetType() uint32 {
	return uint32(i.Data.Base.Type)
}

func (i *ItemData) GetConfigID() uint32 {
	return i.Data.Base.ConfigID
}

func (i *ItemData) GetPos() int32 {
	return i.Data.Base.Pos
}

func (i *ItemData) SetPos(pos int32) {
	i.Data.Base.Pos = pos
}

func (i *ItemData) GetNum() uint32 {
	return i.Data.Base.Num
}

func (i *ItemData) SetNum(num uint32) {
	i.Data.Base.Num = num
}

func (i *ItemData) GetMaxNum() uint32 {
	return GetItemMaxNum(i.Data.Base.ConfigID, uint32(i.Data.Base.Type))
}

func (i *ItemData) GetMaxPileNum() uint32 {
	return GetItemMaxPileNum(i.Data.Base.ConfigID, uint32(i.Data.Base.Type))
}

//CanPile 是否可以堆叠
func (i *ItemData) CanPile(item IItem) bool {
	return i.Data.Base.ConfigID == item.GetConfigID() && uint32(i.Data.Base.Type) == item.GetType()
}
