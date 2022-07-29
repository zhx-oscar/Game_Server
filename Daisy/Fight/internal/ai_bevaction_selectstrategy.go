package internal

import (
	"Daisy/Fight/internal/log"
	b3 "github.com/magicsea/behavior3go"
	b3config "github.com/magicsea/behavior3go/config"
	b3core "github.com/magicsea/behavior3go/core"
	"math/rand"
)

//RandomSelectEnemy 随机选择一个敌人
type RandomSelectEnemy struct {
	b3core.Action
}

func (bev *RandomSelectEnemy) Initialize(setting *b3config.BTNodeCfg) {

	bev.Action.Initialize(setting)
}

// OnTick 循环
func (bev *RandomSelectEnemy) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[RandomSelectEnemy] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	//之前攻击目标且还活着，不用换目标
	oldEnemy, ok := getPawnBoard(tick, attackTarget)
	if ok && oldEnemy.IsAlive() {
		return b3.SUCCESS
	}

	enemys := pawn.GetEnemyList()
	var aliveEnemys []*Pawn
	for _, enemy := range enemys {
		if enemy.IsAlive() && !pawn.cantBeSelect(enemy) {
			aliveEnemys = append(aliveEnemys, enemy)
		}
	}

	if len(aliveEnemys) > 0 {
		enemy := aliveEnemys[rand.Intn(len(aliveEnemys))]
		pawn.setAttackTarget(tick, enemy)
		return b3.SUCCESS
	}

	return b3.FAILURE
}

//selectHpLowestEnemyFromBoard 从黑板敌人列表中选择一个血量最低的敌人
type selectHpLowestEnemyFromBoard struct {
	b3core.Action
}

func (bev *selectHpLowestEnemyFromBoard) Initialize(setting *b3config.BTNodeCfg) {

	bev.Action.Initialize(setting)
}

// OnTick 循环
func (bev *selectHpLowestEnemyFromBoard) OnTick(tick *b3core.Tick) b3.Status {

	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[selectHpLowestEnemyFromBoard] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	var target *Pawn
	enemys, ok := getEnemyListBoard(tick, enemyList)
	if !ok {
		log.Error("*************[selectHpLowestEnemyFromBoard] getEnemyListBoard fail", pawn.UID)
		return b3.FAILURE
	}

	for _, enemy := range enemys {
		if enemy.IsAlive() && !pawn.cantBeSelect(enemy) {
			if target == nil {
				target = enemy
				continue
			}

			if enemy.Attr.CurHP > target.Attr.CurHP {
				target = enemy
			}
		}
	}

	//未找到目标
	if target == nil {
		log.Error("*************[selectHpLowestEnemyFromBoard] target is nil")
		return b3.FAILURE
	}

	pawn.setAttackTarget(tick, target)
	return b3.SUCCESS
}

//SelectEnemyWithMinHP 全场选择一个血量最低的敌人
type SelectEnemyWithMinHP struct {
	b3core.Action
}

func (bev *SelectEnemyWithMinHP) Initialize(setting *b3config.BTNodeCfg) {

	bev.Action.Initialize(setting)
}

// OnTick 循环
func (bev *SelectEnemyWithMinHP) OnTick(tick *b3core.Tick) b3.Status {

	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[SelectEnemyWithMinHP] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	var target *Pawn
	enemys := pawn.GetEnemyList()

	for _, enemy := range enemys {
		if enemy.IsAlive() && !pawn.cantBeSelect(enemy) {
			if target == nil {
				target = enemy
				continue
			}

			if enemy.Attr.CurHP < target.Attr.CurHP {
				target = enemy
			}
		}
	}

	//未找到目标
	if target == nil {
		log.Error("*************[SelectEnemyWithMinHP] target is nil")
		return b3.FAILURE
	}

	pawn.setAttackTarget(tick, target)
	return b3.SUCCESS
}

//SelectSelf 选择自己作为技能施放目标
type SelectSelf struct {
	b3core.Action
}

func (bev *SelectSelf) Initialize(setting *b3config.BTNodeCfg) {
	bev.Action.Initialize(setting)
}

// OnTick 循环
func (bev *SelectSelf) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[SelectSelf] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	pawn.setAttackTarget(tick, pawn)
	return b3.SUCCESS
}

//SelectNearestEnemy 选择自己最近的一名敌人写入目标黑板
type SelectNearestEnemy struct {
	b3core.Action
}

func (bev *SelectNearestEnemy) Initialize(setting *b3config.BTNodeCfg) {
	bev.Action.Initialize(setting)
}

// OnTick 循环
func (bev *SelectNearestEnemy) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[SelectNearestEnemy] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	var target *Pawn
	var enemyMinDis float32

	enemys := pawn.GetEnemyList()
	for _, enemy := range enemys {
		if enemy.IsAlive() && !pawn.cantBeSelect(enemy) {
			if target == nil {
				target = enemy
				enemyMinDis = DistancePawn(pawn, enemy)
				continue
			}

			enemyDis := DistancePawn(pawn, enemy)
			if enemyDis < enemyMinDis {
				target = enemy
				enemyMinDis = enemyDis
			}
		}
	}

	if target == nil {
		return b3.FAILURE
	}

	pawn.setAttackTarget(tick, target)
	return b3.SUCCESS
}

//SelectLowestHPPartner 选择血量最低的伙伴写入目标黑板
type SelectLowestHPPartner struct {
	b3core.Action
}

func (bev *SelectLowestHPPartner) Initialize(setting *b3config.BTNodeCfg) {
	bev.Action.Initialize(setting)
}

// OnTick 循环
func (bev *SelectLowestHPPartner) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[SelectLowestHPPartner] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	var target *Pawn
	var lowestHP int64

	pawnList := pawn.GetPartnerList()
	for _, partner := range pawnList {
		if partner.IsAlive() && !pawn.cantBeSelect(partner) {
			if target == nil {
				target = partner
				lowestHP = partner.Attr.CurHP
				continue
			}

			if partner.Attr.CurHP < lowestHP {
				target = partner
				lowestHP = partner.Attr.CurHP
			}
		}
	}

	if target == nil {
		return b3.FAILURE
	}

	pawn.setAttackTarget(tick, target)
	return b3.SUCCESS
}

//RandSelectEnemyReast 纯随机选择一个敌人
type RandSelectEnemyReast struct {
	b3core.Action
}

func (bev *RandSelectEnemyReast) Initialize(setting *b3config.BTNodeCfg) {
	bev.Action.Initialize(setting)
}

// OnTick 循环
func (bev *RandSelectEnemyReast) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[RandSelectEnemyReast] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	enemys := pawn.GetEnemyList()
	var aliveEnemys []*Pawn
	for _, enemy := range enemys {
		if enemy.IsAlive() && !pawn.cantBeSelect(enemy) {
			aliveEnemys = append(aliveEnemys, enemy)
		}
	}

	if len(aliveEnemys) > 0 {
		enemy := aliveEnemys[rand.Intn(len(aliveEnemys))]
		pawn.setAttackTarget(tick, enemy)
		return b3.SUCCESS
	}

	return b3.FAILURE
}

//SelectEnemyByMinAttr 根据属性数值最低选择敌人
type SelectEnemyByMinAttr struct {
	b3core.Action
	AttrTypeKey string
}

func (bev *SelectEnemyByMinAttr) Initialize(setting *b3config.BTNodeCfg) {
	bev.Action.Initialize(setting)
	bev.AttrTypeKey = setting.GetPropertyAsString("AttrTypeKey")
}

// OnTick 循环
func (bev *SelectEnemyByMinAttr) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[SelectEnemyByMinAttr] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	value, ok := pawn.getBlackboardValueByKey(bev.AttrTypeKey).(string)
	if !ok {
		log.Error("*************[SelectEnemyByMinAttr] 外部黑板 AttrTypeKey 没有配置")
		return b3.ERROR
	}

	var target *Pawn
	enemys := pawn.GetEnemyList()

	switch value {
	case eB_Attr_HP:
		for _, enemy := range enemys {
			if enemy.IsAlive() && !pawn.cantBeSelect(enemy) {
				if target == nil {
					target = enemy
					continue
				}

				if enemy.Attr.CurHP < target.Attr.CurHP {
					target = enemy
				}
			}
		}
	}

	//未找到目标
	if target == nil {
		log.Error("*************[SelectEnemyByMinAttr] target is nil")
		return b3.FAILURE
	}

	pawn.setAttackTarget(tick, target)
	return b3.SUCCESS
}

//SelectEnemyByMaxAttr 根据属性数值最高选择敌人
type SelectEnemyByMaxAttr struct {
	b3core.Action
	AttrTypeKey string
}

func (bev *SelectEnemyByMaxAttr) Initialize(setting *b3config.BTNodeCfg) {
	bev.Action.Initialize(setting)
	bev.AttrTypeKey = setting.GetPropertyAsString("AttrTypeKey")
}

// OnTick 循环
func (bev *SelectEnemyByMaxAttr) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[SelectEnemyByMaxAttr] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	value, ok := pawn.getBlackboardValueByKey(bev.AttrTypeKey).(string)
	if !ok {
		log.Error("*************[SelectEnemyByMaxAttr] 外部黑板 AttrTypeKey 没有配置")
		return b3.ERROR
	}

	var target *Pawn
	enemys := pawn.GetEnemyList()

	switch value {
	case eB_Attr_HP:
		for _, enemy := range enemys {
			if enemy.IsAlive() && !pawn.cantBeSelect(enemy) {
				if target == nil {
					target = enemy
					continue
				}

				if enemy.Attr.CurHP > target.Attr.CurHP {
					target = enemy
				}
			}
		}
	}

	//未找到目标
	if target == nil {
		log.Error("*************[SelectEnemyByMaxAttr] target is nil")
		return b3.FAILURE
	}

	pawn.setAttackTarget(tick, target)
	return b3.SUCCESS
}

//SelectFurthestEnemy 选择距离自己最远的一名敌人写入目标黑板
type SelectFurthestEnemy struct {
	b3core.Action
}

func (bev *SelectFurthestEnemy) Initialize(setting *b3config.BTNodeCfg) {
	bev.Action.Initialize(setting)
}

// OnTick 循环
func (bev *SelectFurthestEnemy) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[SelectFurthestEnemy] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	var target *Pawn
	var enemyMaxDis float32

	enemys := pawn.GetEnemyList()
	for _, enemy := range enemys {
		if enemy.IsAlive() && !pawn.cantBeSelect(enemy) {
			if target == nil {
				target = enemy
				enemyMaxDis = DistancePawn(pawn, enemy)
				continue
			}

			enemyDis := DistancePawn(pawn, enemy)
			if enemyDis > enemyMaxDis {
				target = enemy
				enemyMaxDis = enemyDis
			}
		}
	}

	if target == nil {
		return b3.FAILURE
	}

	pawn.setAttackTarget(tick, target)
	return b3.SUCCESS
}

//GetSkillTargetByBlackboard 获取攻击目标 从黑板中获取当前释放的技能
type GetSkillTargetByBlackboard struct {
	b3core.Action
}

func (bev *GetSkillTargetByBlackboard) Initialize(setting *b3config.BTNodeCfg) {
	bev.Action.Initialize(setting)
}

// OnTick 循环
func (bev *GetSkillTargetByBlackboard) OnTick(tick *b3core.Tick) b3.Status {
	pawn, b := tick.GetTarget().(*Pawn)
	if !b {
		log.Error("*************[GetSkillTargetByBlackboard] trans Pawn fail")
		return b3.ERROR
	}

	pawn.pushActionTorsionInfo(bev.GetName())

	skill, ok := getSkillBoard(tick, curSkill)
	if !ok {
		log.Error("*************[GetSkillTargetByBlackboard] getSkillBoard fail")
		return b3.FAILURE
	}

	target := pawn.getSkillAttackTarget(skill)
	if target == nil {
		log.Error("*************[GetSkillTargetByBlackboard] 技能 通过目标策略表 找不到攻击目标 ", skill.Config.SkillMain_Config.ID)
		return b3.FAILURE
	}

	pawn.setAttackTarget(tick, target)
	return b3.SUCCESS
}
