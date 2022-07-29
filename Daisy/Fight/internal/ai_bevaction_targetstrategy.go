package internal

import (
	"Daisy/DataTables"
	"Daisy/Fight/internal/log"
	b3 "github.com/magicsea/behavior3go"
	b3config "github.com/magicsea/behavior3go/config"
	b3core "github.com/magicsea/behavior3go/core"
	"math/rand"
)

//GetUsableSkill 获取当前技能
type GetUsableSkill struct {
	b3core.Action
}

//UsableSkill 可用技能数据   技能+目标
type UsableSkill struct {
	skill  *_SkillItem
	target *Pawn
}

func (bev *GetUsableSkill) Initialize(setting *b3config.BTNodeCfg) {
	bev.Action.Initialize(setting)
}

// OnTick 循环
func (bev *GetUsableSkill) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[GetCurSkill] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	useableSkill := pawn.getUseableSkill(true)
	if useableSkill == nil {
		return b3.FAILURE
	}

	pawn.setAttackTarget(tick, useableSkill.target)
	tick.Blackboard.Set(curSkill, useableSkill.skill, "", "")
	return b3.SUCCESS
}

//getUseableSkill 获取可用技能 isJudgeAttackDistance：是否判断技能攻击距离
func (ai *_PawnBehavior) getUseableSkill(isJudgeAttackDistance bool) *UsableSkill {
	var useableskill *UsableSkill
	var baseSkillList []*_SkillItem
	var useableSkillList []*UsableSkill

	//必杀技
	baseSkillList = ai.owner.GetUsableUltimateSkills()
	//是否攻击距离生效
	if isJudgeAttackDistance {
		useableSkillList = ai.getSkillListInSkillRange(baseSkillList)
	}

	if len(useableSkillList) > 0 {
		skillIndex := rand.Intn(len(useableSkillList))
		useableskill = useableSkillList[skillIndex]
	}

	if useableskill == nil {
		//超能技
		baseSkillList = ai.owner.GetUsableSuperSkills()

		//是否攻击距离生效
		if isJudgeAttackDistance {
			useableSkillList = ai.getSkillListInSkillRange(baseSkillList)
		}

		useableskill = ai.getFirstUsableSkill(ai.owner, useableSkillList)
	}

	return useableskill
}

//getSkillListInSkillRange 目标是否在当前技能攻击范围以内
func (ai *_PawnBehavior) getSkillListInSkillRange(skillList []*_SkillItem) []*UsableSkill {
	var result []*UsableSkill

	for _, skill := range skillList {
		target := ai.getSkillAttackTarget(skill)
		if target == nil {
			log.Error("技能 通过目标策略表 找不到攻击目标 ", skill.Config.SkillMain_Config.ID)
			continue
		}

		if skill.TargetInAttackRange(target) {
			temp := &UsableSkill{
				skill:  skill,
				target: target,
			}
			result = append(result, temp)
		}
	}

	return result
}

//getFirstUsableSkill 获取第一个可用技能 参照技能栏顺序
func (ai *_PawnBehavior) getFirstUsableSkill(pawn *Pawn, useableSkillList []*UsableSkill) *UsableSkill {
	if len(useableSkillList) == 0 || pawn == nil {
		return nil
	}

	//基础超能技列表
	baseSkillList := pawn.superSkillList
	// 超载状态
	if pawn.State.OverDrive {
		baseSkillList = pawn.overDriveSuperSkillList
	}

	//参照基础技能列表顺序，第一个可用技能
	for _, val := range baseSkillList {
		for _, useableSkill := range useableSkillList {
			if val.Config.SkillValue_Config.ID == useableSkill.skill.Config.SkillValue_Config.ID {
				return useableSkill
			}
		}
	}

	return nil
}

//getSkillAttackTarget 获取当前技能对应的目标
func (ai *_PawnBehavior) getSkillAttackTarget(skill *_SkillItem) *Pawn {

	strategyCfg, ok := ai.owner.Scene.GetTargetStrategyConfig().TargetStrategy_ConfigItems[skill.Config.SkillMain_Config.ID]
	if !ok {
		log.Error("not find GetTargetStrategyConfig.TargetStrategy_ConfigItems id:", skill.Config.SkillMain_Config.ID)
		return nil
	}

	//原始基础目标列表处理
	baseTargetList := map[uint32]*Pawn{}
	for _, pawn := range ai.owner.Scene.pawnList {
		if !pawn.IsBackground() && pawn.IsAlive() {
			baseTargetList[pawn.UID] = pawn
		}
	}

	for _, baseStrategyID := range strategyCfg.StrategyID {
		ai.getTargetList_BaseStrategy(baseStrategyID, baseTargetList)
	}

	//当最终有多个目标符合条件时，随机取一个
	for _, pawn := range baseTargetList {
		return pawn
	}

	return nil
}

//getTargetList_BaseStrategy 通过策略获取目标
func (ai *_PawnBehavior) getTargetList_BaseStrategy(baseStrategyID uint32, baseTargetList map[uint32]*Pawn) {
	if len(baseTargetList) == 0 {
		return
	}

	baseStrategyCfg, ok := ai.owner.Scene.GetTargetStrategyConfig().StrategyBase_ConfigItems[baseStrategyID]
	if !ok {
		log.Error("not find GetTargetStrategyConfig.StrategyBase_ConfigItems id:", baseStrategyID)
		baseTargetList = map[uint32]*Pawn{}
		return
	}

	BaseStrategyTargetList := ai.getTargetListByCamp(baseStrategyCfg.Camp)
	ai.getTargetListByCalcType(baseStrategyCfg, BaseStrategyTargetList)

	//外部传入的基础目标列表 和 当前策略得到的目标列表做 交集处理
	for key := range baseTargetList {
		if _, ok = BaseStrategyTargetList[key]; !ok {
			delete(baseTargetList, key)
		}
	}
}

const (
	strategyCamp_Enemy           = 1 //敌方
	strategyCamp_Friend          = 2 //友方包含自己
	strategyCamp_FriendExcludeMe = 3 //友方不包含自己
	strategyCamp_Owner           = 4 //自己
)

//getTargetListByCamp 通过阵营获取基础目标列表
func (ai *_PawnBehavior) getTargetListByCamp(campID uint32) map[uint32]*Pawn {
	result := map[uint32]*Pawn{}

	switch campID {
	case strategyCamp_Enemy:
		list := ai.owner.GetEnemyList()
		for _, pawn := range list {
			if pawn.IsAlive() && !ai.owner.cantBeSelect(pawn) {
				result[pawn.UID] = pawn
			}
		}
	case strategyCamp_Friend:
		list := ai.owner.GetPartnerList()
		for _, pawn := range list {
			if pawn.IsAlive() && !ai.owner.cantBeSelect(pawn) {
				result[pawn.UID] = pawn
			}
		}
	case strategyCamp_FriendExcludeMe:
		list := ai.owner.GetPartnerList()
		for _, pawn := range list {
			if pawn.UID != ai.owner.UID && pawn.IsAlive() && !ai.owner.cantBeSelect(pawn) {
				result[pawn.UID] = pawn
			}
		}
	case strategyCamp_Owner:
		result[ai.owner.UID] = ai.owner
	default:
		log.Error("unknown StrategyConfig.campID: ", campID)
	}

	return result
}

const (
	strategyType_Distance = 1
	strategyType_HP       = 2
)

//getCalcValueByStrategyType 获取不同策略类型对应的计算类型数值
func (ai *_PawnBehavior) getCalcValueByStrategyType(pawnA, pawnB *Pawn, cfg *DataTables.StrategyBase_Config) uint32 {
	if cfg == nil || pawnA == nil || pawnB == nil {
		return 0
	}

	switch cfg.StrategyType {
	case strategyType_Distance:
		//距离放大一百倍
		return uint32(DistancePawn(pawnA, pawnB) * 100)
	case strategyType_HP:
		return uint32(pawnA.Attr.CurHP)
	default:
		return 0
	}
}

const (
	strategyCalcType_Max      = 1 //最大
	strategyCalcType_Min      = 2 //最小
	strategyCalcType_Interval = 3 //区间
)

//getTargetListByCalcType 获取基础目标列表--计算类型
func (ai *_PawnBehavior) getTargetListByCalcType(cfg *DataTables.StrategyBase_Config, baseTargetList map[uint32]*Pawn) {
	if baseTargetList == nil || cfg == nil {
		baseTargetList = map[uint32]*Pawn{}
		return
	}

	switch cfg.CalcType {
	case strategyCalcType_Max:
		var target *Pawn
		var maxValue uint32
		for _, pawn := range baseTargetList {
			curValue := ai.getCalcValueByStrategyType(pawn, ai.owner, cfg)
			if target == nil {
				target = pawn
				maxValue = curValue
			} else if curValue > maxValue {
				target = pawn
				maxValue = curValue
			}
		}

		for key, val := range baseTargetList {
			if target == nil || val.UID != target.UID {
				delete(baseTargetList, key)
			}
		}
	case strategyCalcType_Min:
		var target *Pawn
		var minValue uint32
		for _, pawn := range baseTargetList {
			curValue := ai.getCalcValueByStrategyType(pawn, ai.owner, cfg)
			if target == nil {
				target = pawn
				minValue = curValue
			} else if curValue < minValue {
				target = pawn
				minValue = curValue
			}
		}

		for key, val := range baseTargetList {
			if target == nil || val.UID != target.UID {
				delete(baseTargetList, key)
			}
		}
	case strategyCalcType_Interval:
		for key, pawn := range baseTargetList {
			curValue := ai.getCalcValueByStrategyType(pawn, ai.owner, cfg)
			if curValue < cfg.Param1 || curValue > cfg.Param2 {
				delete(baseTargetList, key)
			}
		}
	default:
		return
	}
}
