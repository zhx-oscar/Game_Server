package internal

import (
	"Cinder/Base/linemath"
	"Daisy/Fight/internal/conf"
	"Daisy/Fight/internal/log"
	"Daisy/Proto"
	"fmt"
	"math"
)

// _SkillFlow 技能流程
type _SkillFlow struct {
	scene *Scene
}

// init 初始化
func (flow *_SkillFlow) init(scene *Scene) {
	flow.scene = scene
}

// update 帧更新
func (flow *_SkillFlow) update() {
	for _, pawn := range flow.scene.pawnList {
		if !pawn.IsAlive() {
			continue
		}

		flow.updateOne(pawn.curSkill)
	}
}

// isSkillFlowRunning 是否有技能流程正在执行中
func (flow *_SkillFlow) isSkillFlowRunning() bool {
	for _, pawn := range flow.scene.pawnList {
		if flow.isSkillRunning(pawn) {
			return true
		}
	}

	return false
}

// updateOne 更新一个技能
func (flow *_SkillFlow) updateOne(skill *Skill) {
	if skill == nil {
		return
	}

	switch skill.Stat {
	case Proto.SkillState_Dashing:
		flow.inDashing(skill)
	case Proto.SkillState_ShowTime:
		flow.inShowTime(skill)
	case Proto.SkillState_Before:
		flow.inBefore(skill)
	case Proto.SkillState_Attack:
		flow.inAttack(skill)
	case Proto.SkillState_Later:
		flow.inLater(skill)
	}

	//技能流程中 转向处理
	flow.turn(skill)
}

//turn 技能流程中转向修正
func (flow *_SkillFlow) turn(skill *Skill) {
	if skill == nil {
		return
	}

	//只有单体锁定技能才能使用turn
	if skill._SkillItem.GetAttackType() != conf.AttackType_Single {
		return
	}

	if skill.Config.TimeLineConfig.Turn == nil {
		return
	}

	//turn 未开始
	if flow.scene.NowTime-skill.beginTime < skill.turnBeginTime+(skill.endDashingTime-skill.beginDashingTime) {
		return
	}

	//turn已结束
	if flow.scene.NowTime-skill.beginTime > skill.turnBeginTime+skill.Config.TimeLineConfig.Turn.Duration+(skill.endDashingTime-skill.beginDashingTime) && skill.isAlreadyTurned {
		return
	}

	var target *Pawn
	if len(skill.TargetList) > 0 {
		target = skill.TargetList[0]
	}

	//目标缺失
	if target == nil {
		log.Errorf("turn target is Illegal skillID: %v", skill.Config.SkillMain_Config.ID)
		return
	}

	//if !skill.isAlreadyTurned {
	//	fmt.Println("+++++++++ begin   ")
	//}
	//
	//fmt.Println("+++++++++ 00   ", skill.Caster.Scene.NowTime, skill.Caster.GetAngle()*180/math.Pi, CalcAngle(skill.Caster.GetPos(), target.GetPos())*180/math.Pi)
	//fmt.Println("+++++++++ 11   ", skill.Caster.Scene.NowTime, skill.Caster.GetPos(), target.GetPos())
	//角度一致 或者 目标点和自己重合不需要转向处理
	isNeedTurn := !FloatEqual(float64(skill.Caster.GetAngle()), CalcAngle(skill.Caster.GetPos(), target.GetPos()))
	if Vector2Equal(skill.Caster.GetPos(), target.GetPos()) {
		isNeedTurn = false
	}
	if !isNeedTurn {
		//角度不需要转向调整
		return
	}

	angle, err := flow.getNextTrunedAngle(skill, target)
	if err != nil {
		log.Errorf("getNextTrunedAngle skillID: %v, err: %v", skill.Config.SkillMain_Config.ID, err)
		return
	}
	//fmt.Println("+++++++++ 22  ", skill.Caster.Scene.NowTime, angle*180/math.Pi, skill.Caster.angle*180/math.Pi, CalcAngle(skill.Caster.GetPos(), target.GetPos())*180/math.Pi)
	//if FloatEqual(float64(angle), CalcAngle(skill.Caster.GetPos(), target.GetPos())*180/math.Pi) {
	//	fmt.Println("+++++++++ 33     success")
	//}
	skill.Caster.SetAngle(angle, true, true, false)
	skill.isAlreadyTurned = true
}

//getNextTrunedAngle 获取下一次转向角度
func (flow *_SkillFlow) getNextTrunedAngle(skill *Skill, target *Pawn) (float32, error) {
	if skill == nil {
		return 0, fmt.Errorf("skill is Illegal")
	}

	if target == nil || !target.IsAlive() {
		return 0, fmt.Errorf("target is Illegal")
	}

	//目标角度
	targetAngle := float32(CalcAngle(skill._SkillItem.Caster.GetPos(), target.GetPos()))
	t := 1 / float32(flow.scene.secFrames)
	if flow.scene.NowTime-skill.inStatTime > skill.turnBeginTime+skill.Config.TimeLineConfig.Turn.Duration+(skill.endDashingTime-skill.beginDashingTime) {
		t = float32(skill.Config.TimeLineConfig.Turn.Duration%uint32(1000/float32(flow.scene.secFrames))) / 1000
	}

	var resultAngle float32
	trunAngle := t * skill.Config.TimeLineConfig.Turn.Speed

	//目标角度距离自己最近的方向是否是逆时针方向
	isCCW := BetweenAngle(targetAngle, skill.Caster.angle, AddAngle(skill.Caster.angle, math.Pi))

	if isCCW {
		if FloatGreaterEqual(float64(trunAngle), 2*math.Pi) {
			resultAngle = targetAngle
		} else {
			resultAngle = AddAngle(skill.Caster.angle, trunAngle)
			//检测最终角度 和 目标角度 是否越界
			if !BetweenAngle(resultAngle, skill.Caster.angle, targetAngle) {
				resultAngle = targetAngle
			}
		}
	} else {
		if FloatGreaterEqual(float64(trunAngle), 2*math.Pi) {
			resultAngle = targetAngle
		} else {
			resultAngle = AddAngle(skill.Caster.angle, -trunAngle)
			//检测最终角度 和 目标角度 是否越界
			if !BetweenAngle(resultAngle, targetAngle, skill.Caster.angle) {
				resultAngle = targetAngle
			}
		}
	}

	return resultAngle, nil
}

// inReady 准备阶段
func (flow *_SkillFlow) inReady(skill *Skill) {
	caster := skill.Caster

	// 设置技能cd结束时间
	if skill.Config.CoolDown > 0 {
		skill.cdEndTime = caster.Scene.NowTime + skill.Config.CoolDown
	}

	// 记录回放
	actionUseSkill := &Proto.UseSkill{
		CasterId:       caster.UID,
		SkillId:        skill.Config.ValueID(),
		SkipBefore:     skill.skipBefore,
		CombineCasters: skill.combineSkillReadyMembers,
	}

	if len(skill.TargetList) > 0 {
		actionUseSkill.TargetId = skill.TargetList[0].UID
	}

	actionUseSkill.TargetPos = &Proto.Position{
		X: skill.TargetPos.X,
		Y: skill.TargetPos.Y,
	}

	flow.scene.PushAction(actionUseSkill)

	// 记录action
	skill.actionUseSKill = actionUseSkill

	// 发送事件
	skill.Caster.Events.EmitSkillReadyFinish(skill)
	if skill.Stat == Proto.SkillState_Wait {
		return
	}

	// 技能有特写
	if skill.Config.TimeLineConfig.ShowTime > 0 {
		// 进入特写阶段
		flow.changeFlowStat(skill, Proto.SkillState_ShowTime, Proto.SkillEndReason_Normal, nil)
		return
	}

	// 开始冲刺
	flow.beginDashing(skill)
}

// inDashing 冲刺阶段
func (flow *_SkillFlow) inDashing(skill *Skill) {
	pawn := skill.Caster

	if !pawn.isDashingEnd() {
		return
	}

	//冲刺结束时间
	skill.endDashingTime = pawn.Scene.NowTime

	// 更新目标位置
	if skill.GetAttackType() == conf.AttackType_Single {
		skill.TargetPos = skill.TargetList[0].GetPos()
	}

	if !Vector2Equal(skill.TargetPos, pawn.GetPos()) {
		pawn.SetAngle(float32(CalcAngle(pawn.GetPos(), skill.TargetPos)), true, true, true)
	}

	// 发送事件
	pawn.Events.EmitSkillDashingFinish(skill)
	if skill.Stat == Proto.SkillState_Wait {
		return
	}

	// 发送事件
	pawn.Events.EmitSkillInStart(skill)
	if skill.Stat == Proto.SkillState_Wait {
		return
	}

	// 跳过前摇
	if skill.skipBefore {
		// 连招进入命中阶段
		flow.changeFlowStat(skill, Proto.SkillState_Attack, Proto.SkillEndReason_Normal, nil)
		if skill.Stat == Proto.SkillState_Wait {
			return
		}
	} else {
		// 进入前摇阶段
		flow.changeFlowStat(skill, Proto.SkillState_Before, Proto.SkillEndReason_Normal, nil)
		if skill.Stat == Proto.SkillState_Wait {
			return
		}
	}
}

// inShowTime 特写阶段
func (flow *_SkillFlow) inShowTime(skill *Skill) {
	// 特写未结束，特写阶段无需缩放
	if flow.scene.NowTime-skill.inStatTime < skill.Config.TimeLineConfig.ShowTime {
		return
	}

	// 开始冲刺
	flow.beginDashing(skill)
}

// inBefore 前摇阶段
func (flow *_SkillFlow) inBefore(skill *Skill) {
	// 前摇未结束
	if flow.scene.NowTime-skill.inStatTime < skill.ZoomAttackTime(skill.Config.TimeLineConfig.BeforeTime) {
		return
	}

	// 进入命中阶段
	flow.changeFlowStat(skill, Proto.SkillState_Attack, Proto.SkillEndReason_Normal, nil)
}

// inAttack 命中阶段
func (flow *_SkillFlow) inAttack(skill *Skill) {
	// 检测命中持续时间
	if flow.scene.NowTime-skill.inStatTime >= skill.ZoomAttackTime(skill.Config.TimeLineConfig.AttackTime+skill.attackExtendTime) {
		// 进入后摇阶段
		flow.changeFlowStat(skill, Proto.SkillState_Later, Proto.SkillEndReason_Normal, nil)
		return
	}

	// 检测创建伤害体
	for i, atkTimeLine := range skill.Config.TimeLineConfig.Attacks {
		if i < 0 || i >= len(skill.Config.TemplateConfig.AttackConfs) {
			continue
		}

		// 不随技能自动创建，已创建或未到创建时间
		if skill.Config.TemplateConfig.AttackConfs[i].NoAuto || skill.attackBegin.Test(int32(i)) ||
			flow.scene.NowTime-skill.inStatTime < skill.ZoomAttackTime(atkTimeLine.Begin) {
			continue
		}

		// 创建伤害体
		_, err := flow.scene.createSkillAttack(skill, uint32(i))
		if err != nil {
			log.Error(err.Error())
		}

		// 标记已创建
		skill.attackBegin.TurnOn(int32(i))
	}
}

// inLater 后摇阶段
func (flow *_SkillFlow) inLater(skill *Skill) {
	// 后摇未结束
	if flow.scene.NowTime-skill.inStatTime < skill.ZoomAttackTime(skill.Config.TimeLineConfig.LaterTime) {
		return
	}

	// 先置空当前技能
	skill.Caster.curSkill = nil

	// 再重置为等待阶段
	flow.changeFlowStat(skill, Proto.SkillState_Wait, Proto.SkillEndReason_Normal, nil)
}

// changeFlowStat 调整技能阶段
func (flow *_SkillFlow) changeFlowStat(skill *Skill, stat Proto.SkillState_Enum, skillEndReason Proto.SkillEndReason_Enum, breakCaster *Pawn) {
	if skill.Stat == stat {
		return
	}

	// 调整阶段
	lastStat := skill.Stat
	skill.Stat = stat
	skill.inStatTime = flow.scene.NowTime

	// 发送事件
	switch stat {
	case Proto.SkillState_Wait:
		// 打断不记录帧
		if skillEndReason != Proto.SkillEndReason_Break {
			flow.scene.PushAction(&Proto.ChangeSkillState{
				CasterId: skill.Caster.UID,
				SkillId:  skill.Config.ValueID(),
				State:    Proto.SkillState_Wait,
			})
		}

		skill.Caster.Events.EmitSkillInEnd(skill, lastStat, skillEndReason, breakCaster)

	case Proto.SkillState_Ready:
		skill.Caster.Events.EmitSkillInReady(skill)

	case Proto.SkillState_Dashing:
		skill.Caster.Events.EmitSkillInDashing(skill)

	case Proto.SkillState_ShowTime:
		flow.scene.PushAction(&Proto.ChangeSkillState{
			CasterId: skill.Caster.UID,
			SkillId:  skill.Config.ValueID(),
			State:    Proto.SkillState_ShowTime,
		})

		skill.Caster.Events.EmitSkillInShowTime(skill)

	case Proto.SkillState_Before:
		// 跳过前摇不记录帧
		if !skill.skipBefore {
			flow.scene.PushAction(&Proto.ChangeSkillState{
				CasterId: skill.Caster.UID,
				SkillId:  skill.Config.ValueID(),
				State:    Proto.SkillState_Before,
			})
		}

		skill.Caster.Events.EmitSkillInBefore(skill)

	case Proto.SkillState_Attack:
		flow.scene.PushAction(&Proto.ChangeSkillState{
			CasterId: skill.Caster.UID,
			SkillId:  skill.Config.ValueID(),
			State:    Proto.SkillState_Attack,
		})

		skill.Caster.Events.EmitSkillInAttack(skill)

	case Proto.SkillState_Later:
		skill.Caster.Events.EmitSkillInLater(skill)
	}
}

// beginDashing 开始冲刺
func (flow *_SkillFlow) beginDashing(skill *Skill) {
	caster := skill.Caster

	var targetPos linemath.Vector2

	switch skill.GetAttackType() {
	case conf.AttackType_Single:
		target := skill.TargetList[0]

		// 目标不在施法范围内，则冲刺过去
		if !caster.Equal(target) && !skill.TargetInCastDistance(target) {
			if !caster.startDashing() {
				flow.breakCurSkill(caster, caster, Proto.SkillBreakReason_Normal)
				return
			}

			skill.beginDashingTime = caster.Scene.NowTime
			flow.changeFlowStat(skill, Proto.SkillState_Dashing, Proto.SkillEndReason_Normal, nil)
			return
		}

		targetPos = target.GetPos()

	case conf.AttackType_Aoe:
		// 目标点不在施法范围内，则冲刺过去
		if !skill.PosInCastDistance(skill.TargetPos) {
			if !caster.startDashing() {
				flow.breakCurSkill(caster, caster, Proto.SkillBreakReason_Normal)
				return
			}

			skill.beginDashingTime = caster.Scene.NowTime
			flow.changeFlowStat(skill, Proto.SkillState_Dashing, Proto.SkillEndReason_Normal, nil)
			return
		}

		targetPos = skill.TargetPos
	}

	// 设置技能释放朝向
	if Vector2Equal(caster.GetPos(), targetPos) {
		caster.SetAngle(caster.GetAngle(), true, true, true)
	} else {
		angle := float32(CalcAngle(caster.GetPos(), targetPos))
		caster.SetAngle(angle, true, true, true)
	}

	// 发送事件
	caster.Events.EmitSkillInStart(skill)
	if skill.Stat == Proto.SkillState_Wait {
		return
	}

	// 跳过前摇
	if skill.skipBefore {
		// 连招进入命中阶段
		flow.changeFlowStat(skill, Proto.SkillState_Attack, Proto.SkillEndReason_Normal, nil)
		if skill.Stat == Proto.SkillState_Wait {
			return
		}
	} else {
		// 进入前摇阶段
		flow.changeFlowStat(skill, Proto.SkillState_Before, Proto.SkillEndReason_Normal, nil)
		if skill.Stat == Proto.SkillState_Wait {
			return
		}
	}
}

// stopDashing 停止冲刺
func (flow *_SkillFlow) stopDashing(skill *Skill) {
	skill.Caster.stopDashing()
}

// breakCurSkill 打断当前技能
func (flow *_SkillFlow) breakCurSkill(caster, pawn *Pawn, breakReason Proto.SkillBreakReason_Enum) bool {
	skill := pawn.curSkill

	// 没有当前技能或当前技能处于等待阶段
	if skill == nil || skill.Stat == Proto.SkillState_Wait {
		return false
	}

	// 冲刺状态结束冲刺
	if skill.Stat == Proto.SkillState_Dashing {
		flow.stopDashing(skill)
	}

	// 记录debug信息
	flow.scene.PushDebugInfo(func() string {
		return fmt.Sprintf("${PawnID:%d}被${PawnID:%d}打断技能${SkillID:%d}，技能流水：%d",
			pawn.UID,
			caster.UID,
			skill.Config.ValueID(),
			skill.UID)
	})

	// 有使用技能帧并且需要保存打断帧
	if skill.actionUseSKill != nil && breakReason != Proto.SkillBreakReason_Combo {
		// 记录回放
		flow.scene.PushAction(&Proto.BreakSkill{
			CasterId: caster.UID,
			TargetId: pawn.UID,
			SkillId:  skill.Config.ValueID(),
		})
	}

	// 先置空当前技能
	pawn.curSkill = nil

	skillEndReason := Proto.SkillEndReason_Break
	switch breakReason {
	case Proto.SkillBreakReason_Normal:
		skillEndReason = Proto.SkillEndReason_Break
	case Proto.SkillBreakReason_Combo:
		skillEndReason = Proto.SkillEndReason_Combo
	case Proto.SkillBreakReason_AIPause:
		skillEndReason = Proto.SkillEndReason_AIPause
	}

	// 再重置为等待阶段
	flow.changeFlowStat(skill, Proto.SkillState_Wait, skillEndReason, caster)

	if breakReason != Proto.SkillBreakReason_Combo {
		// 发送打断成功事件
		caster.Events.EmitBreakTargetSkill(pawn, skill)

		// 发送被打断事件
		pawn.Events.EmitBeBreakSkill(caster, skill)
	}

	// 打断伤害体
	flow.scene.breakSkillAttacks(skill, caster)

	return true
}

// useSkill 使用技能
func (flow *_SkillFlow) useSkill(pawn *Pawn, skillItem *_SkillItem, targetPos linemath.Vector2, targetList []*Pawn, skipBefore bool) bool {
	// 检测技能否使用指定技能
	if !pawn.CanUseSkill(skillItem) {
		return false
	}

	// 创建技能
	skill, err := skillItem.createSkill()
	if err != nil {
		log.Error(err.Error())
		return false
	}

	// 设置释放目标
	switch skill.GetAttackType() {
	case conf.AttackType_Single:
		if len(targetList) <= 0 {
			return false
		}
		skill.TargetList = targetList

	case conf.AttackType_Aoe:
		skill.TargetPos = targetPos
	}

	// 设置跳过前摇
	skill.skipBefore = skipBefore

	// 强制打断当前技能
	flow.breakCurSkill(pawn, pawn, Proto.SkillBreakReason_Normal)

	// 切换当前技能
	pawn.curSkill = skill

	// 记录debug信息
	flow.scene.PushDebugInfo(func() string {
		return fmt.Sprintf("${PawnID:%d}对位置：%+v，目标：%s，释放了技能${SkillID:%d}，技能流水：%d",
			pawn.UID,
			skill.TargetPos,
			func() (targets string) {
				if len(skill.TargetList) <= 0 {
					return "无"
				}

				for _, target := range skill.TargetList {
					if targets != "" {
						targets += "、"
					}
					targets += fmt.Sprintf("${PawnID:%d}", target.UID)
				}
				return
			}(),
			skill.Config.ValueID(),
			skill.UID)
	})

	// 进入准备阶段
	flow.changeFlowStat(skill, Proto.SkillState_Ready, Proto.SkillEndReason_Normal, nil)
	if skill.Stat == Proto.SkillState_Wait {
		return false
	}

	// 执行技能准备
	flow.inReady(skill)
	if skill.Stat == Proto.SkillState_Wait {
		return false
	}

	return true
}

// isSkillRunning 是否正在使用技能
func (flow *_SkillFlow) isSkillRunning(pawn *Pawn) bool {
	return pawn.IsAlive() && pawn.curSkill != nil && pawn.curSkill.Stat != Proto.SkillState_Wait
}

// searchSkillTargets 查询技能目标列表
func (flow *_SkillFlow) searchSkillTargets(skillItem *_SkillItem, targetPos linemath.Vector2) (targets []*Pawn) {
	for _, atkConf := range skillItem.Config.TemplateConfig.AttackConfs {
		if atkConf.Type != conf.AttackType_Aoe || !atkConf.CastRange {
			continue
		}

		// 计算角度
		angel := skillItem.Caster.GetAngle()

		if !Vector2Equal(targetPos, skillItem.Caster.GetPos()) {
			angel = float32(CalcAngle(skillItem.Caster.GetPos(), targetPos))
		}

		targets = append(targets, flow.scene.searchAttackTargets(atkConf.AttackArgs, targetPos, angel, targetPos.Sub(skillItem.Caster.GetPos()).Len(), skillItem.Caster.GetCamp(), skillItem.Caster.Attr.Scale)...)
	}

	return
}
