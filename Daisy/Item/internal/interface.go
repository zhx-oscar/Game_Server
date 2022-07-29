package internal

import (
	"Daisy/ItemProto"
	"Daisy/Proto"
)

type IItemContainer interface {
	Init(owner IItemContainerOwner, ContainerType int32, maxNum uint32, control IItemContainerControl, data ItemProto.IContainer)

	GetType() int32

	UpdataItem(item ItemProto.IItem)

	RemoveItem(pos int32) ItemProto.IItem
	AddItem(item ItemProto.IItem) (ItemProto.IItem, bool)

	GetIItem(pos int32) ItemProto.IItem
	GetIItemByID(id string) ItemProto.IItem
	RemoveItemByID(id string) ItemProto.IItem

	AddItemList(items []*Proto.DropMaterial, args ...uint32) []*Proto.GetItemData
	AddItemNum(configID, typ, num uint32, args ...uint32) uint32
	CanAddItemNum(configID, typ, num uint32) bool
	RemoveItemNum(configID, typ, num uint32) (uint32, bool)
	RemoveItemList(items []ItemProto.IItem) []ItemProto.IItem

	CanAddItemList(items []*Proto.DropMaterial) bool
	CanAddItems(items []ItemProto.IItem) bool

	GetItemNum(configID, typ uint32) uint32
	HasEnoughSpace(num uint32) bool
	GetNewPos(from int32) int32

	GetCanAddNum(configID, typ, num uint32) uint32

	Traversal(cb func(item ItemProto.IItem) bool)
}

type IItemContainerOwner interface {
	GetContainer(typ int32) IItemContainer
	AddContainer(typ int32, maxNum uint32, control IItemContainerControl) IItemContainer

	GetItemFromPack(itemID string) ItemProto.IItem
	RemoveItemFromPack(itemID string) bool
	HasEnoughSpace(num uint32, typ ...int32) bool
	AddItemToPack(item ItemProto.IItem) bool

	RemoveItem(configID uint32, typ uint32, num uint32) uint32
	AddItem(configID uint32, typ uint32, num uint32, args ...uint32) uint32
	AddItems(items []ItemProto.IItem) []ItemProto.IItem

	UpdataItem(item ItemProto.IItem)
	GetItemNum(configID uint32, typ uint32) uint32
	AddItemList(items []*Proto.DropMaterial, args ...uint32) []*Proto.GetItemData
	RemoveItemList(items []ItemProto.IItem) []ItemProto.IItem
	CanAddItemList(items []*Proto.DropMaterial) ([]int32,bool)
	CanAddItemIList(items []ItemProto.IItem) ([]int32,bool)
	Move(item ItemProto.IItem, dst int32) bool

	Traversal(cb func(container IItemContainer) (bool, interface{})) interface{}
}

type IItemContainerControl interface {
	OnItemAddToContainer(items []ItemProto.IItem, container IItemContainer)
	OnItemChange(items []ItemProto.IItem, container IItemContainer)
	OnItemRemoveFromContainer(items []ItemProto.IItem, container IItemContainer)

	CanAddItemNum(configID uint32, typ uint32, num uint32) bool
	CanAddItem(item ItemProto.IItem) bool

	CreateItem(configID uint32, typ uint32, num uint32, args ...uint32) ItemProto.IItem

	GetContainerData(typ int32) ItemProto.IContainer
	AddContainerData(typ int32, maxNum uint32) ItemProto.IContainer
}
