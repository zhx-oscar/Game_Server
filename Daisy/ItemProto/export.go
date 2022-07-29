package ItemProto

import (
	"Daisy/ItemProto/internal"
	"Daisy/Proto"
)

type IItem = internal.IItem

type IContainer = internal.IContainer

var CreateIItemByID = internal.CreateIItemByID
var CreateIItemByData = internal.CreateIItemByData
var CreateIContainer = internal.CreateIContainer
var GetItemMaxNum = internal.GetItemMaxNum
var GetItemMaxPileNum = internal.GetItemMaxPileNum

func GetItemContainerType(itemTyp uint32) int32 {
	var ctyp Proto.ContainerEnum_Type
	switch Proto.ItemEnum_Type(itemTyp) {
	case Proto.ItemEnum_Equipment:
		ctyp = Proto.ContainerEnum_EquipBag
	case Proto.ItemEnum_SkillItem:
		ctyp = Proto.ContainerEnum_SkillBag
	}
	return int32(ctyp)
}

func GetItemBuildContainerType(itemTyp uint32) int32 {
	var ctyp Proto.ContainerEnum_Type
	switch Proto.ItemEnum_Type(itemTyp) {
	case Proto.ItemEnum_Equipment:
		ctyp = Proto.ContainerEnum_BuildEquipBag
	}
	return int32(ctyp)
}

func IsItemBuildContainerType(typ int32) bool {
	return Proto.ContainerEnum_Type(typ) == Proto.ContainerEnum_BuildEquipBag
}
