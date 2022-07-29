package internal

import (
	"Daisy/Fight/internal/conf"
	"Daisy/Fight/internal/log"
	"Daisy/Proto"
	b3 "github.com/magicsea/behavior3go"
	b3config "github.com/magicsea/behavior3go/config"
	b3core "github.com/magicsea/behavior3go/core"
	"math/rand"
)

//GetCurSkill 获取当前技能
type GetCurSkill struct {
	b3core.Action
	isrand int
}

func (bev *GetCurSkill) Initialize(setting *b3config.BTNodeCfg) {
	bev.Action.Initialize(setting)
	bev.isrand = setting.GetPropertyAsInt("isrand")
}

// OnTick 循环
func (bev *GetCurSkill) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[GetCurSkill] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	var skillIndex int

	//必杀技
	skillList := pawn.GetUsableUltimateSkills()
	if len(skillList) > 0 {
		if bev.isrand == 1 {
			skillIndex = rand.Intn(len(skillList))
		}

		tick.Blackboard.Set(curSkill, skillList[skillIndex], "", "")
		return b3.SUCCESS
	}

	skillList = pawn.GetUsableSuperSkills()
	if len(skillList) > 0 {
		useableskill := bev.getFirstUsableSkill(pawn, skillList)
		if useableskill == nil {
			return b3.FAILURE
		}

		tick.Blackboard.Set(curSkill, useableskill, "", "")
		return b3.SUCCESS
	}

	//skillList = pawn.GetUsableNormalAttacks()
	//if len(skillList) > 0 {
	//	tick.Blackboard.Set(curSkill, skillList[0], "", "")
	//	return b3.SUCCESS
	//}

	return b3.FAILURE
}

//getFirstUsableSkill 获取第一个可用技能 参照技能栏顺序
func (bev *GetCurSkill) getFirstUsableSkill(pawn *Pawn, useableSkillList []*_SkillItem) *_SkillItem {
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
			if val.Config.SkillValue_Config.ID == useableSkill.Config.SkillValue_Config.ID {
				return useableSkill
			}
		}
	}

	return nil
}

//CastBloadSkill 释放黑板中的技能
type CastBloadSkill struct {
	b3core.Action
	runningstate int
}

func (bev *CastBloadSkill) Initialize(setting *b3config.BTNodeCfg) {
	bev.Action.Initialize(setting)
	bev.runningstate = setting.GetPropertyAsInt("runningstate")
}

// OnTick 循环
func (bev *CastBloadSkill) OnTick(tick *b3core.Tick) b3.Status {
	defer func() {
		tick.Blackboard.Remove(enemyList)
		tick.Blackboard.Remove(attackPos)
		//tick.Blackboard.Remove(attackTargetPos)
	}()

	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[CastBloadSkill] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	skill, ok := getSkillBoard(tick, curSkill)
	if !ok {
		log.Error("*************[CastBloadSkill] getSkillBoard fail")
		return b3.FAILURE
	}
	target, ok := getPawnBoard(tick, attackTarget)
	if !ok {
		log.Error("*************[CastBloadSkill] getPawnBoard fail")
		return b3.FAILURE
	}

	isAlreadyCast := tick.Blackboard.GetBool(isAlreadyCastSkill, "", "")
	if !isAlreadyCast {

		//目标死亡了 技能释放失败
		if !target.IsAlive() {
			return b3.FAILURE
		}

		//自己处于被击状态不能释放技能
		if pawn.State.CantMove {
			return b3.FAILURE
		}

		//攻击目标不处于攻击范围以内
		if !skill.TargetInAttackRange(target) {
			log.Debug("目标处于技能攻击范围以外 当前技能释放失败 ", pawn.UID, skill.Config.SkillMain_Config.ID, skill.Config.Name)
			{
				if skill.Config.CanDash {
					log.Debug("技能可以冲刺， 自己与目标距离：", Distance(skill.Caster.GetPos(), target.GetPos()), " 技能包含冲刺的射程：", skill.Caster.Attr.CollisionRadius+target.Attr.CollisionRadius+skill.GetMinCastDistance()+skill.Caster.Info.FastDist, skill.Caster.Attr.CollisionRadius+target.Attr.CollisionRadius+skill.GetMaxCastDistance()+skill.Caster.Info.FastDist)
				} else {
					log.Debug("技能不能冲刺， 自己与目标距离：", Distance(skill.Caster.GetPos(), target.GetPos()), " 技能不包含冲刺的射程：", skill.Caster.Attr.CollisionRadius+target.Attr.CollisionRadius+skill.GetMinCastDistance(), skill.Caster.Attr.CollisionRadius+target.Attr.CollisionRadius+skill.GetMaxCastDistance())
				}
			}
			return b3.FAILURE
		}

		//如果当前释放技能是 超能技或者必杀技
		if (skill.Config.SkillKind == conf.SkillKind_Super || skill.Config.SkillKind == conf.SkillKind_Ultimate) && pawn.IsNormalAttackRunning() {

			//个人必杀不能打断
			isInterruptNormalSkill := pawn.getBlackboardValueByKey(eB_UltimateSkillInterruptNormalSkill)
			if skill.Config.SkillKind == conf.SkillKind_Ultimate && !(isInterruptNormalSkill != nil && isInterruptNormalSkill.(bool)) {
				return b3.FAILURE
			}

			//个人超能不能打断
			isInterruptNormalSkill = pawn.getBlackboardValueByKey(eB_SuperSkillInterruptNormalSkill)
			if skill.Config.SkillKind == conf.SkillKind_Super && !(isInterruptNormalSkill != nil && isInterruptNormalSkill.(bool)) {
				return b3.FAILURE
			}

			pawn.BreakCurSkill(pawn, Proto.SkillBreakReason_Normal)
		}

		//fmt.Println("+++++++++++  ", pawn.IsMoving(), pawn.UID)
		//if pawn.IsMoving() {
		//	pawn.Stop()
		//}

		//使用技能之前 冲刺所需参数处理
		bev.doDashingParam(tick, skill, target)

		result := pawn.UseSkill(skill, target.GetPos(), []*Pawn{target})
		if !result {
			log.Debug("useSkill fail ", pawn.UID, skill.Config.Name)
		}
	}

	tick.Blackboard.Remove(isAlreadyCastSkill)
	//fmt.Println("+++++++ing  ", pawn.UID, bev.runningstate, skill.Config.Name, pawn.CanCast(), pawn.IsUsingSkill())

	if bev.runningstate > 0 {
		if pawn.IsSkillRunning() {
			if bev.isBreakNormalAttackRunning(tick) {
				pawn.BreakNormalAttackCombo(pawn)
			}

			tick.Blackboard.Set(isAlreadyCastSkill, true, "", "")
			return b3.RUNNING
		}
	}

	//fmt.Println("+++++++end  ", pawn.UID, bev.runningstate, skill.Config.Name, pawn.CanCast(), pawn.IsSkillRunning())
	//fmt.Println("+++++++++++ UseSkill   ", pawn.UID, pawn.GetPos())
	return b3.SUCCESS
}

//doDashingParam 处理当前技能冲刺参数
func (bev *CastBloadSkill) doDashingParam(tick *b3core.Tick, skill *_SkillItem, target *Pawn) {
	if skill.Config.CanDash {
		tick.Blackboard.Set(dash_skill, skill, "", "")
		tick.Blackboard.Set(dash_target, target, "", "")
	}
}

//isBreakNormalAttackRunning 是否中断当前普攻连击
func (bev *CastBloadSkill) isBreakNormalAttackRunning(tick *b3core.Tick) bool {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		return false
	}

	if !pawn.IsNormalAttackRunning() {
		return false
	}

	isInterruptNormalSkill := pawn.getBlackboardValueByKey(eB_SuperSkillInterruptNormalSkill)
	//配置超能技可以打断普攻连击
	if isInterruptNormalSkill != nil && isInterruptNormalSkill.(bool) {
		if skillList := pawn.GetUsableSuperSkills(); len(skillList) > 0 {
			return true
		}
	}

	isInterruptNormalSkill = pawn.getBlackboardValueByKey(eB_UltimateSkillInterruptNormalSkill)
	//配置必杀技可以打断普攻连击
	if isInterruptNormalSkill != nil && isInterruptNormalSkill.(bool) {
		if skillList := pawn.GetUsableUltimateSkills(); len(skillList) > 0 {
			return true
		}
	}

	//如果处在地方AOE预警范围内 立马终端普攻连击

	return false
}

//GetSkillByIndex 通过填写的技能索引找到对应技能并且写入黑板
type GetSkillByIndex struct {
	b3core.Action
	skillIndex int
	skillType  int
}

func (bev *GetSkillByIndex) Initialize(setting *b3config.BTNodeCfg) {
	bev.Action.Initialize(setting)
	bev.skillIndex = setting.GetPropertyAsInt("skillIndex")
	bev.skillType = setting.GetPropertyAsInt("skillType")
}

// OnTick 循环
func (bev *GetSkillByIndex) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[GetSkillByIndex] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	var skill *_SkillItem
	switch bev.skillType {
	case skillType_super:
		skill = bev.GetSkillBySkillType(pawn, pawn.GetSuperSkillList())
	case skillType_ultimate:
		skill = bev.GetSkillBySkillType(pawn, pawn.GetUltimateSkillList())
	default:
	}

	if skill != nil {
		tick.Blackboard.Set(curSkill, skill, "", "")
		return b3.SUCCESS
	}

	return b3.FAILURE
}

func (bev *GetSkillByIndex) GetSkillBySkillType(pawn *Pawn, skillList []*_SkillItem) *_SkillItem {
	if pawn == nil {
		log.Error("*************[GetSkillByIndex] trans Pawn fail")
		return nil
	}

	if len(skillList) > 0 && bev.skillIndex < len(skillList) {
		skill := skillList[bev.skillIndex]
		if pawn.CanUseSkill(skill) {
			return skill
		}

		return nil
	}

	return nil
}

//GetRandNormalAttack 随机获取可用普攻，并把普攻写入黑板.
type GetRandNormalAttack struct {
	b3core.Action
}

func (bev *GetRandNormalAttack) Initialize(setting *b3config.BTNodeCfg) {
	bev.Action.Initialize(setting)
}

// OnTick 循环
func (bev *GetRandNormalAttack) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[GetRandNormalAttack] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	skillList := pawn.GetUsableNormalAttacks()

	if len(skillList) == 0 {
		return b3.FAILURE
	}

	index := rand.Intn(len(skillList))

	tick.Blackboard.Set(curSkill, skillList[index], "", "")
	return b3.SUCCESS
}
