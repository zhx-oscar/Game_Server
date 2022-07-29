package main

import (
	"Cinder/Base/Prop"
	"Daisy/Battle/drop"
	"Daisy/Const"
	"Daisy/Item"
	"Daisy/ItemProto"
	"Daisy/Proto"
)

//role对item的操作方法，查看 Item/internal/interface.go/IItemContainerOwner

func (r *_Role) InitItemContainer() {
	//Todo 读配置，获得各包裹初始大小
	//目前先随意写个数字
	r.AddContainer(int32(Proto.ContainerEnum_Common), 0, r)

	r.AddContainer(int32(Proto.ContainerEnum_EquipBag), 0, r)

	r.AddContainer(int32(Proto.ContainerEnum_SkillBag), 0, r)

	r.AddContainer(int32(Proto.ContainerEnum_BuildEquipBag), 0, r)
}

//OnItemAddToContainer 获得新道具时被调用
func (r *_Role) OnItemAddToContainer(items []ItemProto.IItem, container Item.IItemContainer) {
	itemsData := &Proto.Items{Items: make([]*Proto.Item, 0)}
	for _, ii := range items {

		//新获得技能书 自动学习技能并且技能书不会加入背包
		if r.learnSkill(ii.GetData()) {
			//学习扣除1之后，如果数量达到0需要扣除
			if ii.GetData().Base.Num == 0 {
				continue
			}
		}

		itemsData.Items = append(itemsData.Items, ii.GetData())
	}
	r.prop.SyncAddItemToContainer(container.GetType(), itemsData, Prop.Target_Client)

	//真正添加之后，再抛出事件
	r.FireLocalEvent(Const.Event_addItem, items)
}

//OnItemRemoveFromContainer 删除道具时被调用
func (r *_Role) OnItemRemoveFromContainer(items []ItemProto.IItem, container Item.IItemContainer) {
	posData := &Proto.Int32Array{Data: make([]int32, 0)}
	for _, ii := range items {
		posData.Data = append(posData.Data, ii.GetPos())
	}
	r.prop.SyncDelItemFromContainer(container.GetType(), posData, Prop.Target_Client)

	//真正移除之后，再抛出事件
	r.FireLocalEvent(Const.Event_removeItem, items)
}

//OnItemChange 道具数量改变时被调用
func (r *_Role) OnItemChange(items []ItemProto.IItem, container Item.IItemContainer) {
	itemsData := &Proto.Items{Items: make([]*Proto.Item, 0)}
	buildItemsData := &Proto.Items{Items: make([]*Proto.Item, 0)}
	for _, ii := range items {
		//装备道具特殊处理build背包
		if ii.GetType() == uint32(Proto.ItemEnum_Equipment) {
			if ii.GetData().ExpandData.InUse == 0 && ItemProto.IsItemBuildContainerType(container.GetType()) == true {
				//需要挪走
				r.Move(ii, ItemProto.GetItemContainerType(ii.GetType()))
				return
			} else if ii.GetData().ExpandData.InUse == 1 && ItemProto.IsItemBuildContainerType(container.GetType()) == false {
				//需要挪走
				r.Move(ii, ItemProto.GetItemBuildContainerType(ii.GetType()))
				return
			}
		}

		if ii.GetData().ExpandData.InUse == 0 {
			itemsData.Items = append(itemsData.Items, ii.GetData())
		} else {
			buildItemsData.Items = append(buildItemsData.Items, ii.GetData())
		}
	}
	if len(itemsData.Items) > 0 {
		r.prop.SyncUpdateItemToContainer(container.GetType(), itemsData, Prop.Target_Client)
	}
	if len(buildItemsData.Items) > 0 {
		r.prop.SyncUpdateItemToContainer(container.GetType(), buildItemsData, Prop.Target_Client)
	}

	//真正更新之后，再抛出事件
	r.FireLocalEvent(Const.Event_updateItem, items)
}

//CanAddItemNum 检查是否能放指定数量道具
func (r *_Role) CanAddItemNum(configID uint32, typ uint32, num uint32) bool {

	//ToDo 道具重量等特性判断

	container := r.GetItemContainerByItemType(typ)
	return container.CanAddItemNum(configID, typ, num)
}

func (r *_Role) CanAddItem(item ItemProto.IItem) bool {
	//ToDo 道具重量等特性判断

	container := r.GetItemContainerByItemType(item.GetType())
	return container.CanAddItems([]ItemProto.IItem{item})
}

//CreateItem 创建道具数据
func (r *_Role) CreateItem(configID uint32, typ uint32, num uint32, args ...uint32) ItemProto.IItem {
	drop := drop.Drop{}

	args = append(args, r.GetLucky()) //生成装备需要幸运值

	return ItemProto.CreateIItemByData(drop.CreateItem(configID, uint32(typ), num, args...))
}

func (r *_Role) AddContainerData(typ int32, maxNum uint32) ItemProto.IContainer {
	return ItemProto.CreateIContainer(r.prop.AddContainerData(typ, maxNum))
}

func (r *_Role) GetContainerData(typ int32) ItemProto.IContainer {
	data, err := r.prop.GetContainerData(typ)
	if err == nil {
		return ItemProto.CreateIContainer(data)
	}
	return nil
}
