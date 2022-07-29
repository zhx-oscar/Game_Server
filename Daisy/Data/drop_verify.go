package Data

import (
	"Daisy/DataTables"
	"fmt"
	"strconv"
	"strings"
)

// DropConfigData 掉落表二次处理
var DropConfigData map[uint32]*DropEntry

func init() {
	DropConfigData = make(map[uint32]*DropEntry)
}

const (
	DropItemTypeProp = 1 //掉落物品类型是道具
	DropItemTypeLib  = 2 //掉落物品类型是掉落库
)

const (
	DropEntryRuleTypeWeight = 1 //掉落规则是权重掉落
	DropEntryRuleTypeSolo   = 2 //掉落规则是独立掉落
)

const (
	DropItemConditionTypeJob   = 1 //掉落条件满足职业
	DropItemConditionTypeLevel = 2 //掉落条件满足等级
)
const DropRuleWeight = 10000 // 权重占比分母

// DropEntry 表格里掉落条目,公告id暂时不处理
type DropEntry struct {
	EntryId    uint32      // 掉落id
	EntryRule  uint32      // 掉落规则
	EntryTimes uint32      // 掉落次数
	EntryItems []_DropItem // 掉落物品
}

// DropItem 掉落物品
type _DropItem struct {
	ItemProbability uint32        //掉落概率
	ItemId          uint32        //掉落物品id
	ItemDropType    uint32        //掉落物品类型，可能是道具可能是道具库
	ItemProp        _DropProp     //掉落物品如果不是道具，则该值为空
	ItemCondition   DropCondition //针对不同概率掉落的条件
}

// DropProp 掉落道具
type _DropProp struct {
	PropType   uint32 //道具类型，如果是道具库，则道具类型为0
	PropMinNum uint32 //道具最小掉落数量
	PropMaxNum uint32 //道具最大掉落数量
}

// DropCondition 掉落条件
type DropCondition struct {
	ConditionContext map[uint32]interface{} // 或者用 []string 表示多个条件 1-2；2-1
	ConditionIsOr    bool                   // 多个条件是否是或
}

// 将表格里的数据组织成自己的数据结构
// 组织的时候可能会panic
func HandleDropConfig() bool {
	excedata := GetDropConfig()
	dropConfigData := make(map[uint32]*DropEntry)

	for i := 1; i <= len(excedata.Drop_ConfigItems); i++ {
		val := excedata.Drop_ConfigItems[uint32(i)]
		entry, ok := dropConfigData[val.DropBoxID]
		if ok {
			if len(entry.EntryItems) < 1 {
				fmt.Printf("掉落表 掉落库条目加入失败 %d\n", val.DropBoxID)
				return false
			}
			probability := entry.EntryItems[len(entry.EntryItems)-1].ItemProbability + val.Probability
			dropItem, ok := readItemInfo(val, probability)
			if !ok {
				return false
			}

			if entry.EntryRule == DropEntryRuleTypeWeight {
				if dropItem.ItemProbability > DropRuleWeight {
					//计算规则为权重掉落的时候概率不能超过10000
					fmt.Printf("掉落表 计算规则为权重掉落时概率不能超过10000，现已超过 掉落库id %d- 概率总和 %d\n", val.DropBoxID, dropItem.ItemProbability)
					return false
				}
			}
			entry.EntryItems = append(entry.EntryItems, dropItem)

		} else {
			temp := &DropEntry{
				EntryItems: make([]_DropItem, 0),
				EntryId:    val.DropBoxID,
				EntryRule:  val.Rule,
				EntryTimes: val.Times,
			}
			dropItem, ok := readItemInfo(val, val.Probability)
			if !ok {
				return false
			}
			temp.EntryItems = append(temp.EntryItems, dropItem)
			dropConfigData[val.DropBoxID] = temp
		}
	}
	if !checkDeadLock(dropConfigData) {
		return false
	}
	DropConfigData = dropConfigData
	return true
}

// 读取道具信息
func readPropInfo(config *DataTables.Drop_Config) _DropProp {
	// 读取道具
	dropProp := _DropProp{
		PropType:   config.ItemType,
		PropMinNum: config.MinNum,
		PropMaxNum: config.MaxNum,
	}
	return dropProp
}

// 读取条件信息
// 如果表格修改的话，这个函数需要改变
func readConditionInfo(config *DataTables.Drop_Config) (DropCondition, bool) {
	ok := true
	// 读取条件
	dropCondition := DropCondition{
		ConditionIsOr: config.Relation,
	}
	condition := make(map[uint32]interface{}, 2)

	if !fullUpCondition(config.Condition1, &condition, config.Param1) {
		ok = false
	}
	if !fullUpCondition(config.Condition2, &condition, config.Param2) {
		ok = false
	}

	dropCondition.ConditionContext = condition
	return dropCondition, ok
}

// 转换条件类型
func fullUpCondition(conditionType uint32, dropCondition *map[uint32]interface{}, context string) bool {
	switch conditionType {
	case DropItemConditionTypeJob:
		job, err := strconv.Atoi(context)
		if err != nil {
			fmt.Printf("掉落表 掉落条件是职业参数转换失败, 掉落条件%d - 掉落参数%v \n", conditionType, context)
			return false
		}
		(*dropCondition)[conditionType] = uint32(job)
	case DropItemConditionTypeLevel:
		args := strings.Split(context, ",")
		if len(args) != 2 {
			fmt.Printf("掉落表 掉落条件读取参数失败, 掉落条件%d - 掉落参数%v \n", conditionType, context)
			return false
		}
		levels := make([]uint32, 2)
		minLevel, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("掉落表 掉落条件是最低等级参数转换失败, 掉落条件%d - 掉落参数%v \n", conditionType, context)
			return false
		}
		maxLevel, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("掉落表 掉落条件是最高等级参数转换失败, 掉落条件%d - 掉落参数%v \n", conditionType, context)
			return false
		}
		levels[0] = uint32(minLevel)
		levels[1] = uint32(maxLevel)

		(*dropCondition)[conditionType] = levels
	case 0:
		// 是0就不存，表示表格里这行参数没填
	default:
		fmt.Printf("掉落表 不支持的掉落条件, 掉落条件是 %d\n", conditionType)
		return false
	}
	return true
}

// 读取掉落栏信息
// 读取条件时可能会panic
func readItemInfo(config *DataTables.Drop_Config, probability uint32) (_DropItem, bool) {
	ok := true
	dropItem := _DropItem{
		ItemProbability: probability,
		ItemId:          config.Key,
		ItemDropType:    config.Type,
	}

	// 读取条件
	dropItem.ItemCondition, ok = readConditionInfo(config)

	// 读取道具
	if config.Type == DropItemTypeProp {
		dropItem.ItemProp = readPropInfo(config)
	}

	return dropItem, ok
}

// 表格读完后组织好结构检查一遍结构里的东西
// 检查表格里的死锁
func checkDeadLock(dropConfigData map[uint32]*DropEntry) bool {

	for index, _ := range dropConfigData {
		nestStack := make(map[uint32]uint32)
		ok := checkDropEntry(nestStack, index, dropConfigData)
		if !ok {
			return ok
		}
		nestStack = nil
	}
	return true
}

func checkDropEntry(nestStack map[uint32]uint32, dropId uint32, dropConfigData map[uint32]*DropEntry) bool {

	if !checkRepeatDropId(nestStack, dropId) {
		return false
	}
	dropEntry, ok := dropConfigData[dropId]
	if !ok {
		fmt.Printf("掉落表 该掉落库条目未配置, 掉落库id是 %d\n", dropId)
		return false
	}

	for _, item := range dropEntry.EntryItems {
		if item.ItemDropType == DropItemTypeLib {
			ok = checkDropEntry(nestStack, item.ItemId, dropConfigData)
			if !ok {
				return ok
			}
			delete(nestStack, item.ItemId)
		}
	}
	return true
}

func checkRepeatDropId(nestStack map[uint32]uint32, dropid uint32) bool {
	_, ok := nestStack[dropid]
	if ok {
		fmt.Printf("掉落表 死循环检查未通过，重复掉落库id %d\n", dropid)
		return false
	}
	nestStack[dropid] = 1
	return true
}
