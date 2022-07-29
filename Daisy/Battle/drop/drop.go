package drop

import (
	"Daisy/Data"
	"Daisy/Proto"
	log "github.com/cihub/seelog"
	"math/rand"
	"time"
)

type Drop struct {
}

func init() {
	if !Data.HandleDropConfig() {
		panic("handle drop config error ")
	}
	rand.Seed(time.Now().UnixNano())
}

// Drop 掉落，默认批量掉落
func (drop *Drop) Drop(dropID uint32, job uint32, level uint32) (bool, []*Proto.DropMaterial) {
	result := make([]*Proto.DropMaterial, 0)
	ok := drop.dropBatch(dropID, job, level, &result)
	return ok, result
}

// dropBatch 批量掉落，根据配置的掉落次数
func (drop *Drop) dropBatch(dropID uint32, job uint32, level uint32, result *[]*Proto.DropMaterial) bool {

	config, ok := Data.DropConfigData[dropID]
	if !ok {
		log.Errorf("[drop] dropBatch配置表为空")
		return false
	}
	times := config.EntryTimes

	var i uint32
	for i = 0; i < times; i++ {
		// 从一个范围里随机一个值
		ok2, p := drop.calulateProbability(config)
		if !ok2 {
			return false
		}

		// 二次转换的时候一定要保证概率是对的
		has_drop := false
		condition_ok := true
	loop:
		for _, val := range config.EntryItems {
			// 找到对应的item
			if p < val.ItemProbability {
				// 因为概率排列是按从小到大来的，所以这样写没问题，最好还是做下检测
				condition_ok = drop.checkCondition(val.ItemCondition, job, level)
				if !condition_ok {
					break loop
				}
				switch val.ItemDropType {
				case Data.DropItemTypeProp:
					material := &Proto.DropMaterial{
						MaterialId:   val.ItemId,
						MaterialType: val.ItemProp.PropType,
						MaterialNum:  uint32(myRandInt(int(val.ItemProp.PropMinNum), int(val.ItemProp.PropMaxNum))),
					}
					*result = append(*result, material)
					has_drop = true
					break loop
				case Data.DropItemTypeLib:
					ok = drop.dropBatch(val.ItemId, job, level, result)
					if ok {
						has_drop = true
						break loop
					} else {
						// 这里没有对返回值做处理，如果有一个嵌套掉落失败，立即返回
						log.Errorf("[drop] dropBatch 嵌套中批量掉落失败, 嵌套掉落库id是 %d", val.ItemId)
						return false
					}
				default:
					log.Errorf("[drop]dropBatch 未定义的掉落类型 %d", val.ItemDropType)
					return false
				}
			}
		}
		//发送提示消息，有概率没有掉落
		if !has_drop && condition_ok {
			//log.Debug("[drop] dropBatch 有概率没有掉落")
		}
	}
	return true
}

// 单个掉落，无视配置的掉落次数
// func (drop *_Drop) DropOne(dropID uint32) (bool, *DropMaterial) {}

// calulateProbability 计算概率
func (drop *Drop) calulateProbability(config *Data.DropEntry) (bool, uint32) {
	// 从一个范围里随机一个值
	var p uint32

	switch config.EntryRule {
	case Data.DropEntryRuleTypeWeight:
		p = uint32(rand.Intn(Data.DropRuleWeight))
	case Data.DropEntryRuleTypeSolo:
		if len(config.EntryItems) < 1 {
			log.Errorf("[drop] calulateProbability 该掉落id的掉落条目为空")
			return false, 0
		}
		max := config.EntryItems[len(config.EntryItems)-1].ItemProbability
		p = uint32(rand.Intn(int(max)))
	default:
		log.Errorf("[drop] calulateProbability 未定义的计算规则 %d", config.EntryRule)
		return false, 0
	}
	return true, p
}

// checkCondition 条件检查
func (drop *Drop) checkCondition(need_conditon Data.DropCondition, job uint32, level uint32) bool {

	var match bool
	if need_conditon.ConditionIsOr {
		for key, val := range need_conditon.ConditionContext {
			switch key {
			case Data.DropItemConditionTypeJob:
				if val.(uint32) == job {
					match = true
				} else {
					match = false
				}
			case Data.DropItemConditionTypeLevel:
				arr := val.([]uint32)
				if level >= (arr[0]) && level <= arr[1] {
					match = true
				} else {
					match = false
				}
			default:
				log.Errorf("[drop] checkCondition 未定义的条件类型")
				return false
			}
			if match {
				return true
			}
		}
		log.Infof("[drop] checkCondition 或条件不匹配")
		return false
	} else {
		for key, val := range need_conditon.ConditionContext {
			switch key {
			case Data.DropItemConditionTypeJob:
				if val.(uint32) == job {
					match = true
				} else {
					match = false
				}
			case Data.DropItemConditionTypeLevel:
				arr := val.([]uint32)
				if level >= (arr[0]) && level <= arr[1] {
					match = true
				} else {
					match = false
				}
			default:
				log.Errorf("[drop] checkCondition 未定义的条件类型")
				return false
			}
			if !match {
				log.Infof("[drop] checkCondition 与条件不匹配, 需要的条件是%v-%v", key, val)
				return false
			}
		}
		return true
	}
}

// myRandInt 生成区间随机数 左闭右闭
func myRandInt(min int, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}

	return rand.Intn(max-min+1) + min
}
