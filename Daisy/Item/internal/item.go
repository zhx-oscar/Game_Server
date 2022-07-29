package internal

import (
	"Daisy/ItemProto"
	"Daisy/Proto"
	"fmt"
)

func NewItemContainer() IItemContainer {
	return &Container{}
}

type Container struct {
	owner   IItemContainerOwner
	control IItemContainerControl
	ItemProto.IContainer
}

func (c *Container) Init(owner IItemContainerOwner, ContainerType int32, maxNum uint32, control IItemContainerControl, data ItemProto.IContainer) {
	c.owner = owner
	c.control = control
	c.IContainer = data
}

// UpdataItem 修改物品数据
func (c *Container) UpdataItem(item ItemProto.IItem) {
	c.control.OnItemChange([]ItemProto.IItem{item}, c)
}

// RemoveItem 移除指定位置物品
func (c *Container) RemoveItem(pos int32) ItemProto.IItem {
	item := c.GetIItem(pos)
	if item == nil {
		return nil
	}

	c.control.OnItemRemoveFromContainer([]ItemProto.IItem{item}, c)
	return item
}

func (c *Container) pileItem(item ItemProto.IItem) bool {
	//整摞不处理堆叠
	if item.GetNum() >= item.GetMaxPileNum() {
		return false
	}

	for _, itemdata := range c.GetItemMap() {
		if item.GetNum() == 0 {
			break
		}
		_item := ItemProto.CreateIItemByData(itemdata)
		if item.CanPile(_item) == false {
			continue
		}
		if _item.GetNum() >= _item.GetMaxPileNum() {
			continue
		}

		if _item.GetNum()+item.GetNum() > _item.GetMaxPileNum() {
			item.SetNum(item.GetNum() - (_item.GetMaxPileNum() - _item.GetNum()))
			_item.SetNum(_item.GetMaxPileNum())
			c.control.OnItemChange([]ItemProto.IItem{_item}, c)
			return true

		} else {
			_item.SetNum(_item.GetNum() + item.GetNum())
			item.SetNum(0)
			c.control.OnItemChange([]ItemProto.IItem{_item}, c)
			return true
		}
	}

	return false
}

//PileItemNum  堆叠，返回成功堆叠数量,外面数据同步
func (c *Container) pileItemNum(configID, typ, num uint32) (uint32, ItemProto.IItem) {
	var add uint32

	item := ItemProto.CreateIItemByID(configID, typ)
	itemMaxPileNum := item.GetMaxPileNum()

	if num >= itemMaxPileNum {
		//只处理不足一堆数量的堆叠，省点效率
		num = num % itemMaxPileNum
	}
	if num == 0 {
		return 0, nil
	}
	//处理堆叠
	for _, itemdata := range c.GetItemMap() {
		if num == 0 {
			break
		}
		_item := ItemProto.CreateIItemByData(itemdata)
		if _item.CanPile(item) == false {
			continue
		}
		if _item.GetNum() >= itemMaxPileNum {
			continue
		}
		if _item.GetNum()+num > itemMaxPileNum {
			num -= itemMaxPileNum - _item.GetNum()
			add += itemMaxPileNum - _item.GetNum()
			_item.SetNum(itemMaxPileNum)
			return add, _item
		} else {
			_item.SetNum(_item.GetNum() + num)
			add += num
			num = 0
			return add, _item
		}
	}
	return add, nil
}

//addItem 找坐标和同步数据
func (c *Container) addItem(item ItemProto.IItem) bool {
	max := c.GetMaxSpace()
	var found bool
	//包裹格子从1开始
	for i := int32(1); i <= int32(max); i++ {
		_, ok := c.GetItemMap()[i]
		if ok == false {
			item.SetPos(i)
			item.GetData().ExpandData.ContainerType = Proto.ContainerEnum_Type(c.GetType())
			found = true
			break
		}
	}
	if found == false {
		//满了
		return false
	}

	fmt.Println("AddItem ID:", item.GetData().Base.ID, "Type:", item.GetData().Base.Type)
	c.control.OnItemAddToContainer([]ItemProto.IItem{item}, c)

	return true
}

// AddItem 增加物品
func (c *Container) AddItem(item ItemProto.IItem) (ItemProto.IItem, bool) {
	if c.control.CanAddItemNum(item.GetConfigID(), item.GetType(), item.GetNum()) == false {
		return nil, false
	}

	max := c.GetMaxSpace()

	//处理堆叠
	c.pileItem(item)

	if item.GetNum() == 0 {
		return nil, true //全部堆叠
	}

	if len(c.GetItemMap()) >= int(max) {
		return nil, false
	}

	if c.addItem(item) == false{
		return nil, false
	}
	return item,true
}

// AddItemNum 增加物品
func (c *Container) AddItemNum(configID, typ, num uint32, args ...uint32) uint32 {
	if c.control.CanAddItemNum(configID, uint32(typ), num) == false {
		return 0
	}

	var add uint32

	max := c.GetMaxSpace()
	maxPileNum := ItemProto.GetItemMaxPileNum(configID, typ)

	//处理堆叠
	add, _item := c.pileItemNum(configID, typ, num)
	if _item != nil {
		c.control.OnItemChange([]ItemProto.IItem{_item}, c)
	}

	//还有剩余
	if num > add {
		num -= add
		for {
			if len(c.GetItemMap()) >= int(max) {
				break
			}

			//新创建道具
			var itemNum uint32
			if num > maxPileNum {
				itemNum = maxPileNum
			} else {
				itemNum = num
			}
			ii := c.control.CreateItem(configID, uint32(typ), itemNum, args...)
			if ii == nil {
				break
			}
			if ok := c.addItem(ii); ok == true {
				num -= itemNum
				add += itemNum
			}
			if num == 0 {
				break
			}
		}
	}

	return add
}

func (c *Container) AddItemList(items []*Proto.DropMaterial, args ...uint32) []*Proto.GetItemData {
	list := make([]*Proto.GetItemData, 0)
	maxpos := int32(0)
	additems := make([]ItemProto.IItem, 0)
	changeitems := make([]ItemProto.IItem, 0)

	for _, v := range items {
		if c.control.CanAddItemNum(v.MaterialId, v.MaterialType, v.MaterialNum) == false {
			continue
		}
		left := v.MaterialNum

		pileNum, item := c.pileItemNum(v.MaterialId, v.MaterialType, v.MaterialNum)
		if item != nil {
			changeitems = append(changeitems, item)
			list = append(list, &Proto.GetItemData{
				GetNum:   pileNum,
				OldConfigID:v.MaterialId,
				ItemData: item.GetData(),
			})

		}
		if pileNum < v.MaterialNum {
			left -= pileNum
			maxPileNum := ItemProto.GetItemMaxPileNum(v.MaterialId, v.MaterialType)
			num := left / maxPileNum
			if left%maxPileNum != 0 {
				num += 1
			}
			for i := uint32(0); i < num; i++ {
				if left == 0 {
					break
				}

				itemNum := maxPileNum
				if left < maxPileNum {
					itemNum = left
				}
				left -= itemNum
				item2 := c.control.CreateItem(v.MaterialId, v.MaterialType, itemNum, args...)
				if item2 != nil {
					//获取坐标
					pos := c.GetNewPos(maxpos)
					if pos == 0 {
						//满了
						break
					}
					maxpos = pos

					item2.SetPos(pos)
					item2.GetData().ExpandData.ContainerType = Proto.ContainerEnum_Type(ItemProto.GetItemContainerType(item2.GetType()))

					additems = append(additems, item2)
					list = append(list, &Proto.GetItemData{
						GetNum:   itemNum,
						OldConfigID:v.MaterialId,
						ItemData: item2.GetData(),
					})

				}
			}
		}
	}
	c.control.OnItemChange(changeitems, c)
	c.control.OnItemAddToContainer(additems, c)
	return list
}

func (c *Container) RemoveItemList(items []ItemProto.IItem) []ItemProto.IItem {
	c.control.OnItemRemoveFromContainer(items, c)
	return items
}

//CanAddItemNum 包裹是否有足够空间
func (c *Container) CanAddItemNum(configID, typ, num uint32) bool {

	var canAdd uint32
	max := c.GetMaxSpace()
	_item := ItemProto.CreateIItemByID(configID, typ)
	itemMaxPileNum := _item.GetMaxPileNum()

	//判断最大可拥有数量
	itemMaxNum := _item.GetMaxNum()
	if itemMaxNum != 0 && itemMaxNum < num+c.GetItemNum(configID, typ) {
		return false
	}

	//空间足够大
	if itemMaxPileNum*(max-uint32(len(c.GetItemMap()))) >= num {
		return true
	}

	//处于临界值，需要看看能否堆上去
	for _, itemdata := range c.GetItemMap() {
		item := ItemProto.CreateIItemByData(itemdata)
		if item.CanPile(_item) {
			continue
		}
		if item.GetNum() >= item.GetMaxPileNum() {
			continue
		}
		canAdd += item.GetMaxPileNum() - item.GetNum()
		if canAdd >= num {
			return true
		}
	}

	if int(max) <= len(c.GetItemMap()) {
		return false
	}

	return (canAdd + itemMaxPileNum*(max-uint32(len(c.GetItemMap())))) >= num
}

//CanAddItemList 包裹是否有足够空间
func (c *Container) CanAddItemList(items []*Proto.DropMaterial) bool {
	//默认items 不存在重复的道具
	max := c.GetMaxSpace()
	list := make([]Proto.DropMaterial, 0)
	var need uint32
	for _, param := range items {
		iitem := ItemProto.CreateIItemByID(param.MaterialId, param.MaterialType)
		itemMaxNum := iitem.GetMaxNum()
		if itemMaxNum != 0 && param.MaterialNum+c.GetItemNum(param.MaterialId, param.MaterialType) > itemMaxNum {
			//最大拥有数量判断
			return false
		}
		itemMaxPileNum := iitem.GetMaxPileNum()
		need += param.MaterialNum / itemMaxPileNum
		if param.MaterialNum%itemMaxPileNum != 0 {
			list = append(list, Proto.DropMaterial{MaterialId: param.MaterialId, MaterialType: param.MaterialType, MaterialNum: param.MaterialNum % itemMaxPileNum})
		}
	}

	//空间足够
	if max >= uint32(len(c.GetItemMap()))+need+uint32(len(list)) {
		return true
	}

	//空间有限，需要判断能否堆上去
	for _, v := range list {
		_item := ItemProto.CreateIItemByID(v.MaterialId, v.MaterialType)

		//处理堆叠
		for _, itemdata := range c.GetItemMap() {
			item := ItemProto.CreateIItemByData(itemdata)
			if item.CanPile(_item) == false {
				continue
			}
			if item.GetNum() >= item.GetMaxPileNum() {
				continue
			}
			//放不下
			if item.GetNum()+v.MaterialNum > item.GetMaxPileNum() {
				need += 1
			}
			break
		}
	}

	return max >= uint32(len(c.GetItemMap()))+need
}

func (c *Container) CanAddItems(items []ItemProto.IItem) bool {
	//默认items 不存在重复的道具
	max := c.GetMaxSpace()

	if max >= uint32(len(c.GetItemMap())+len(items)) {
		return true
	}

	var need uint32

	for _, v := range items {
		itemMaxPileNum := v.GetMaxPileNum()
		if itemMaxPileNum == v.GetNum() {
			need += 1
			continue
		}

		//处理堆叠
		for _, itemdata := range c.GetItemMap() {
			item := ItemProto.CreateIItemByData(itemdata)
			if item.CanPile(v) {
				continue
			}
			if item.GetNum() >= item.GetMaxPileNum() {
				continue
			}
			if item.GetMaxPileNum() < item.GetNum()+v.GetNum() {
				need += 1
			}
			break
		}
	}

	return max >= uint32(len(c.GetItemMap()))+need
}

func (c *Container) GetIItem(pos int32) ItemProto.IItem {
	itemdata := c.GetItemByPos(pos)
	if itemdata == nil {
		return nil
	}

	return ItemProto.CreateIItemByData(itemdata)
}

// RemoveItemNum 删除指定数量的指定物品
func (c *Container) RemoveItemNum(configID, typ, num uint32) (uint32, bool) {
	if num <= 0 {
		return 0, false
	}

	// 先找出指定数量的物品所在格子，再删除
	var removePos []ItemProto.IItem
	value := num
	list := make([]ItemProto.IItem, 0)

	for _, itemdata := range c.GetItemMap() {
		item := ItemProto.CreateIItemByData(itemdata)
		if item.GetConfigID() == configID && item.GetType() == typ {
			if value < item.GetNum() {
				item.SetNum(item.GetNum() - value)
				value = 0
				list = append(list, item)
				break
			} else {
				removePos = append(removePos, item)
				value -= item.GetNum()
			}
			if value == 0 {
				break
			}
		}
	}

	c.control.OnItemRemoveFromContainer(removePos, c)
	c.control.OnItemChange(list, c)

	if value > 0 {
		// 数量不足
		return num - value, false
	}

	return num, true
}

func (c *Container) GetItemNum(configID, typ uint32) uint32 {
	num := uint32(0)
	for _, itemdata := range c.GetItemMap() {
		item := ItemProto.CreateIItemByData(itemdata)
		if item != nil && item.GetConfigID() == configID && item.GetType() == typ {
			num += item.GetNum()
		}
	}

	return num
}

func (c *Container) HasEnoughSpace(num uint32) bool {
	return uint32(len(c.GetItemMap()))+num <= c.GetMaxSpace()
}

func (c *Container) GetIItemByID(id string) ItemProto.IItem {
	itemdata := c.GetItemByID(id)
	if itemdata == nil {
		return nil
	}

	return ItemProto.CreateIItemByData(itemdata)
}

func (c *Container) RemoveItemByID(id string) ItemProto.IItem {
	item := c.GetIItemByID(id)
	if item != nil {
		return c.RemoveItem(item.GetPos())
	}
	return nil
}

func (c *Container) GetNewPos(from int32) int32 {
	//包裹格子从1开始
	for i := from + 1; i <= int32(c.GetMaxSpace()); i++ {
		_, ok := c.GetItemMap()[i]
		if ok == false {
			return i
		}
	}
	return 0
}

//GetCanAddNum 获还得能添加的数量
//传入 计划添加的 num 数量
//返回本次能添加的数量
func (c *Container) GetCanAddNum(configID, typ, num uint32) uint32 {
	maxNum := ItemProto.GetItemMaxNum(configID, typ)
	if maxNum != 0 {
		_num := c.GetItemNum(configID, typ)
		if _num+num <= maxNum {
			return num
		} else if maxNum > _num {
			return maxNum - _num
		} else {
			return 0
		}
	}
	return num
}

func (c *Container) Traversal(cb func(item ItemProto.IItem) bool) {
	for _, itemdata := range c.GetItemMap() {
		if !cb(ItemProto.CreateIItemByData(itemdata)) {
			return
		}
	}
}
