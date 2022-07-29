package internal

import (
	"Daisy/Data"
	"Daisy/Proto"
)

func CreateIContainer(data *Proto.ItemContainer) IContainer {
	return &ContainerData{
		Data: data,
	}
}

type ContainerData struct {
	Data *Proto.ItemContainer
}

func (c *ContainerData) SetData(data *Proto.ItemContainer) {
	c.Data = data
}

func (c *ContainerData) GetData() *Proto.ItemContainer {
	return c.Data
}

func (c *ContainerData) GetType() int32 {
	return int32(c.Data.Type)
}

func (c *ContainerData) GetMaxSpace() uint32 {
	config,ok := Data.GetItemTypeConfig().ItemPackMaxSpace_ConfigItems[uint32(c.GetType())]
	if ok == true {
		return config.MaxSpace
	}
	return 10000
}

func (c *ContainerData) GetItemMap() map[int32]*Proto.Item {
	return c.Data.ItemMap
}

func (c *ContainerData) GetItemByPos(pos int32) *Proto.Item {
	item, ok := c.Data.ItemMap[pos]
	if ok == false {
		return nil
	}
	return item
}

func (c *ContainerData) GetItemByID(id string) *Proto.Item {
	pos, ok := c.Data.Id2Pos[id]
	if ok == false {
		return nil
	}

	return c.GetItemByPos(pos)
}

func (c *ContainerData) Traversal(cb func(item IItem) (bool, interface{})) interface{} {
	for _, itemdata := range c.GetData().ItemMap {
		item := CreateIItemByData(itemdata)
		if r, s := cb(item); r == false {
			return s
		}
	}
	return nil
}
