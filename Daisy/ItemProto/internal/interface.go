package internal

import "Daisy/Proto"

type IItem interface {
	SetData(data *Proto.Item)
	GetData() *Proto.Item
	GetID() string
	GetType() uint32
	GetConfigID() uint32
	GetPos() int32
	SetPos(pos int32)
	GetNum() uint32
	SetNum(num uint32)
	GetMaxNum() uint32
	GetMaxPileNum() uint32
	CanPile(item IItem) bool
}

type IContainer interface {
	SetData(data *Proto.ItemContainer)
	GetData() *Proto.ItemContainer
	GetType() int32
	GetMaxSpace() uint32
	GetItemMap() map[int32]*Proto.Item
	GetItemByPos(pos int32) *Proto.Item
	GetItemByID(id string) *Proto.Item
	Traversal(cb func(item IItem) (bool, interface{})) interface{}
}
