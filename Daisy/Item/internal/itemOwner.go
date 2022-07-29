package internal

import (
	"Daisy/ItemProto"
	"Daisy/Proto"
)

type ContainerOwner struct {
	containerMap map[int32]IItemContainer
}

func (c *ContainerOwner) GetContainer(typ int32) IItemContainer {
	container, ok := c.containerMap[typ]
	if ok == true {
		return container
	}
	return nil
}

func (c *ContainerOwner) AddContainer(typ int32, maxNum uint32, control IItemContainerControl) IItemContainer {
	data := control.AddContainerData(int32(typ), maxNum)
	container := NewItemContainer()
	container.Init(c, int32(typ), maxNum, control, data)

	if c.containerMap == nil {
		c.containerMap = make(map[int32]IItemContainer)
	}
	c.containerMap[int32(typ)] = container
	return container
}

//HasEnoughSpace 是否有足够的空间
//typ 每个包裹都有足够空间才返回true
func (c *ContainerOwner) HasEnoughSpace(num uint32, typ ...int32) bool {
	for _, v := range typ {
		container := c.GetContainer(v)
		if container != nil {
			if container.HasEnoughSpace(num) == false {
				return false
			}
		}
	}
	return true
}

//RemoveItemFromPack 从目标包裹删除道具
func (c *ContainerOwner) RemoveItemFromPack(itemID string) bool {
	removeItem := c.Traversal(func(container IItemContainer) (bool, interface{}) {
		_item := container.RemoveItemByID(itemID)
		if _item == nil {
			return true, nil
		}
		return false, _item
	})
	return removeItem != nil
}

//GetItemFromPack 从目标包裹获取道具
func (c *ContainerOwner) GetItemFromPack(itemID string) ItemProto.IItem {
	item := c.Traversal(func(container IItemContainer) (bool, interface{}) {
		_item := container.GetIItemByID(itemID)
		if _item == nil {
			return true, nil
		}
		return false, _item
	})

	if item == nil {
		return nil
	}
	ii, ok := item.(ItemProto.IItem)
	if ok == false {
		return nil
	}
	return ii
}

//GetItemContainerByItemType 根据道具类型找到可存放的包裹
func (c *ContainerOwner) GetItemContainerByItemType(typ uint32) IItemContainer {
	return c.GetContainer(ItemProto.GetItemContainerType(typ))
}

//CanAddItem 能否放入item
func (c *ContainerOwner) canAddItem(item ItemProto.IItem) bool {
	container := c.GetItemContainerByItemType(item.GetType())
	if container != nil {
		num := container.GetCanAddNum(item.GetConfigID(), item.GetType(), item.GetNum())
		if num != item.GetNum() {
			return false
		}
		return container.HasEnoughSpace(1)
	}
	return false
}

//AddItemToPack 往目标包裹添加道具
func (c *ContainerOwner) AddItemToPack(item ItemProto.IItem) bool {
	container := c.GetItemContainerByItemType(item.GetType())
	if container != nil {
		_, ok := container.AddItem(item)
		return ok
	}
	return false
}

func (c *ContainerOwner) RemoveItem(configID uint32, typ uint32, num uint32) uint32 {
	container := c.GetItemContainerByItemType(typ)
	if container != nil {
		removeNum, _ := container.RemoveItemNum(configID, typ, num)
		return removeNum
	}
	return 0
}

func (c *ContainerOwner) AddItem(configID uint32, typ uint32, num uint32, args ...uint32) uint32 {
	container := c.GetItemContainerByItemType(typ)
	if container != nil {
		num = container.GetCanAddNum(configID, typ, num)
		return container.AddItemNum(configID, typ, num, args...)
	}
	return 0
}

func (c *ContainerOwner) UpdataItem(item ItemProto.IItem) {
	c.Traversal(func(container IItemContainer) (bool, interface{}) {
		_item := container.GetIItemByID(item.GetID())
		if _item == nil {
			return true, nil
		}
		container.UpdataItem(item)
		return false, nil
	})
}

func (c *ContainerOwner) GetItemNum(configID uint32, typ uint32) uint32 {
	container := c.GetItemContainerByItemType(typ)
	if container != nil {
		return container.GetItemNum(configID, typ)
	}
	return 0
}

func (c *ContainerOwner) getItem(typ uint32, pos int32) ItemProto.IItem {
	container := c.GetItemContainerByItemType(typ)
	if container != nil {
		item := container.GetIItem(pos)
		return item
	}
	return nil
}

func (c *ContainerOwner) AddItemList(items []*Proto.DropMaterial, args ...uint32) []*Proto.GetItemData{
	var itemMap map[int32][]*Proto.DropMaterial
	itemMap = make(map[int32][]*Proto.DropMaterial)

	//整理下items,区分包裹类型，合并同类道具
	for _, v := range items {
		//包裹
		container := c.GetItemContainerByItemType(v.MaterialType)
		if container == nil {
			return nil
		}
		_c, ok := itemMap[container.GetType()]
		if ok == false {
			list := make([]*Proto.DropMaterial, 0)
			itemMap[container.GetType()] = list
			_c, _ = itemMap[container.GetType()] //应该不会没有把
		}
		found := false
		for k2, v2 := range _c {
			if v2.MaterialId == v.MaterialId && v2.MaterialType == v.MaterialType {
				//同类数量合并
				itemMap[container.GetType()][k2].MaterialNum += v.MaterialNum
				itemMap[container.GetType()][k2].MaterialNum = container.GetCanAddNum(v2.MaterialId, v2.MaterialType, itemMap[container.GetType()][k2].MaterialNum)
				found = true
				break
			}
		}
		if found == false {
			itemMap[container.GetType()] = append(itemMap[container.GetType()], &Proto.DropMaterial{
				MaterialId:v.MaterialId,
				MaterialType:v.MaterialType,
				MaterialNum:container.GetCanAddNum(v.MaterialId, v.MaterialType, v.MaterialNum),
			})
		}
	}

	list := make([]*Proto.GetItemData, 0)
	for k, param := range itemMap {
		container := c.GetContainer(int32(k))
		if container == nil {
			return nil
		}
		_list := container.AddItemList(param)
		if len(_list) > 0 {
			list = append(list, _list...)
		}
	}
	return list
}

func (c *ContainerOwner) RemoveItemList(items []ItemProto.IItem) []ItemProto.IItem{
	var itemMap map[int32][]ItemProto.IItem
	itemMap = make(map[int32][]ItemProto.IItem)
	removeList := make([]ItemProto.IItem, 0)

	//整理下items,区分包裹类型，合并同类道具
	for _, v := range items {
		//包裹
		container := c.GetItemContainerByItemType(v.GetType())
		if container == nil {
			return removeList
		}
		_, ok := itemMap[container.GetType()]
		if ok == false {
			list := make([]ItemProto.IItem, 0)
			itemMap[container.GetType()] = list
		}
		itemMap[container.GetType()] = append(itemMap[container.GetType()], v)
	}

	for k, param := range itemMap {
		container := c.GetContainer(int32(k))
		if container == nil {
			return nil
		}
		_list := container.RemoveItemList(param)
		if len(_list) > 0 {
			removeList = append(removeList, _list...)
		}
	}
	return removeList
}

func (c *ContainerOwner) CanAddItemList(items []*Proto.DropMaterial) ([]int32,bool) {
	var itemMap map[int32][]*Proto.DropMaterial
	itemMap = make(map[int32][]*Proto.DropMaterial)
	fullCList := make([]int32, 0)
	//整理下items,区分包裹类型，合并同类道具
	for _, v := range items {
		//包裹
		container := c.GetItemContainerByItemType(v.MaterialType)
		if container == nil {
			return nil,false
		}
		_c, ok := itemMap[container.GetType()]
		if ok == false {
			list := make([]*Proto.DropMaterial, 0)
			itemMap[container.GetType()] = list
			_c, _ = itemMap[container.GetType()] //应该不会没有把
		}
		found := false
		for k2, v2 := range _c {
			if v2.MaterialId == v.MaterialId && v2.MaterialType == v.MaterialType {
				//同类数量合并
				itemMap[container.GetType()][k2].MaterialNum += v.MaterialNum
				found = true
			}
			break
		}
		if found == false {
			itemMap[container.GetType()] = append(itemMap[container.GetType()], &Proto.DropMaterial{
				MaterialId:v.MaterialId,
				MaterialType:v.MaterialType,
				MaterialNum:v.MaterialNum,
			})
		}
	}

	can := true
	for k, param := range itemMap {
		container := c.GetContainer(int32(k))
		if container == nil {
			return nil,false
		}
		if container.CanAddItemList(param) == false {
			fullCList = append(fullCList, k)
			can = false
			continue
		}
	}

	return fullCList, can
}

func (c *ContainerOwner) CanAddItemIList(items []ItemProto.IItem) ([]int32,bool) {
	list := make([]*Proto.DropMaterial, 0)
	for _, v := range items {
		list = append(list, &Proto.DropMaterial{
			MaterialId:		v.GetConfigID(),
			MaterialType:	v.GetType(),
			MaterialNum:	v.GetNum(),
		})
	}
	return c.CanAddItemList(list)
}

func (c *ContainerOwner) AddItems(items []ItemProto.IItem) []ItemProto.IItem {
	list := make([]ItemProto.IItem, 0)
	for _, v := range items {
		if c.AddItemToPack(v) == true {
			list = append(list, v)
		}
	}
	return list
}

func (c *ContainerOwner) Move(item ItemProto.IItem, dst int32) bool {
	dstContainer := c.GetContainer(dst)
	if dstContainer == nil || dstContainer.CanAddItems([]ItemProto.IItem{item}) == false {
		return false
	}
	c.RemoveItemFromPack(item.GetID())
	dstContainer.AddItem(item)

	//Todo 跟客户端定义 Move 同步数据协议

	return true
}

func (c *ContainerOwner) Traversal(cb func(container IItemContainer) (bool, interface{})) interface{} {
	for _, container := range c.containerMap {
		if r, s := cb(container); r == false {
			return s
		}
	}
	return nil
}
