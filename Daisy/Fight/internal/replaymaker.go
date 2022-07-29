package internal

import (
	"Daisy/Fight/internal/log"
	"Daisy/Proto"
	"fmt"
	"unsafe"
)

// _ReplayMaker 回放生成器
type _ReplayMaker struct {
	scene               *Scene
	replay, debugReplay Proto.FightReplay
}

// init 初始化
func (maker *_ReplayMaker) init(scene *Scene) {
	maker.scene = scene
}

// PushAction 添加动作
func (maker *_ReplayMaker) PushAction(action interface{}) {
	fightAction := &Proto.FightAction{}

	switch pbAction := action.(type) {
	case *Proto.SummonPawn:
		fightAction.Type = Proto.ActionType_SummonPawn
		fightAction.ActionSummonPawn = pbAction

		log.Debug("战场ID：", uintptr(unsafe.Pointer(maker.scene)),
			"，指令：召唤",
			"，时间(ms)：", maker.scene.NowTime,
			"，SelfId：", fightAction.ActionSummonPawn.SelfId)

	case *Proto.SetTarget:
		fightAction.Type = Proto.ActionType_SetTarget
		fightAction.ActionSetTarget = pbAction

		log.Debug("战场ID：", uintptr(unsafe.Pointer(maker.scene)),
			"，指令：设置目标",
			"，时间(ms)：", maker.scene.NowTime,
			"，SelfId：", fightAction.ActionSetTarget.SelfId,
			"，TargetId：", fightAction.ActionSetTarget.TargetId)

	case *Proto.MoveBegin:
		fightAction.Type = Proto.ActionType_MoveBegin
		fightAction.ActionMoveBegin = pbAction

		var x, y float32
		if fightAction.ActionMoveBegin.LookAtPos != nil {
			x = fightAction.ActionMoveBegin.LookAtPos.X
			y = fightAction.ActionMoveBegin.LookAtPos.Y
		}

		log.Debug("战场ID：", uintptr(unsafe.Pointer(maker.scene)),
			"，指令：移动开始",
			"，时间(ms)：", maker.scene.NowTime,
			"，PawnID：", fightAction.ActionMoveBegin.SelfId,
			"，MoveMode：", fightAction.ActionMoveBegin.MoveMode,
			"，FanAngle：", fightAction.ActionMoveBegin.Angle,
			"，Speed：", fightAction.ActionMoveBegin.Speed,
			"，目标点：{", fightAction.ActionMoveBegin.Pos.X, " ", fightAction.ActionMoveBegin.Pos.Y, "}",
			"，ExpectMoveEndTime：", fightAction.ActionMoveBegin.ExpectMoveEndTime,
			"，ActualMoveEndTime：", fightAction.ActionMoveBegin.ActualMoveEndTime,
			"，迂回中心点：{", x, " ", y, "}")

	case *Proto.FixMoveData:
		fightAction.Type = Proto.ActionType_FixMoveData
		fightAction.ActionFixMoveData = pbAction

		var x, y float32
		if fightAction.ActionFixMoveData.Pos != nil {
			x = fightAction.ActionFixMoveData.Pos.X
			y = fightAction.ActionFixMoveData.Pos.Y
		}

		log.Debug("战场ID：", uintptr(unsafe.Pointer(maker.scene)),
			"，指令：移动修正",
			"，时间(ms)：", maker.scene.NowTime,
			"，PawnID：", fightAction.ActionFixMoveData.SelfId,
			"，FanAngle：", fightAction.ActionFixMoveData.Angle,
			"，位置：{", x, " ", y, "}")

	case *Proto.MoveEnd:
		fightAction.Type = Proto.ActionType_MoveEnd
		fightAction.ActionMoveEnd = pbAction

		log.Debug("战场ID：", uintptr(unsafe.Pointer(maker.scene)),
			"，指令：移动结束",
			"，时间(ms)：", maker.scene.NowTime,
			"，PawnID：", fightAction.ActionMoveEnd.SelfId,
			"，MoveMode：", fightAction.ActionMoveEnd.MoveMode,
			"，位置: {", fightAction.ActionMoveEnd.Pos.X, " ", fightAction.ActionMoveEnd.Pos.Y, "}")

	case *Proto.UseSkill:
		fightAction.Type = Proto.ActionType_UseSkill
		fightAction.ActionUseSkill = pbAction

		log.Debug("战场ID：", uintptr(unsafe.Pointer(maker.scene)),
			"，指令：使用技能",
			"，时间(ms)：", maker.scene.NowTime,
			"，施法者：", fightAction.ActionUseSkill.CasterId,
			"，跳过前摇： ", fightAction.ActionUseSkill.SkipBefore,
			"，跳过后摇：", fightAction.ActionUseSkill.SkipLater,
			"，技能ID：", fightAction.ActionUseSkill.SkillId,
			"，目标ID：", fightAction.ActionUseSkill.TargetId)

	case *Proto.BreakSkill:
		fightAction.Type = Proto.ActionType_BreakSkill
		fightAction.ActionBreakSkill = pbAction

		log.Debug("战场ID：", uintptr(unsafe.Pointer(maker.scene)),
			"，指令：打断技能",
			"，时间(ms)：", maker.scene.NowTime,
			"，施法者：", fightAction.ActionBreakSkill.CasterId,
			"，目标ID：", fightAction.ActionBreakSkill.TargetId,
			"，被打断技能ID：", fightAction.ActionBreakSkill.SkillId)

	case *Proto.NewAttack:
		fightAction.Type = Proto.ActionType_NewAttack
		fightAction.ActionNewAttack = pbAction

		log.Debug("战场ID：", uintptr(unsafe.Pointer(maker.scene)),
			"，指令：创建伤害体",
			"，时间(ms)：", maker.scene.NowTime,
			"，施法者：", fightAction.ActionNewAttack.CasterId,
			"，来源：", fightAction.ActionNewAttack.Src,
			"，技能ID：", fightAction.ActionNewAttack.SkillId,
			"，技能伤害索引：", fightAction.ActionNewAttack.Index,
			"，BuffID：", fightAction.ActionNewAttack.BuffId,
			"，伤害体ID：", fightAction.ActionNewAttack.AttackId,
			"，目标ID：", fightAction.ActionNewAttack.TargetId,
			"，目标位置：", fightAction.ActionNewAttack.TargetPos,
			"，位移结束时间：", fightAction.ActionNewAttack.MoveEndTime,
			"，分组：", fightAction.ActionNewAttack.GroupID)

	case *Proto.DelAttack:
		fightAction.Type = Proto.ActionType_DelAttack
		fightAction.ActionDelAttack = pbAction

		log.Debug("战场ID：", uintptr(unsafe.Pointer(maker.scene)),
			"，指令：删除伤害体",
			"，时间(ms)：", maker.scene.NowTime,
			"，伤害体ID：", fightAction.ActionDelAttack.AttackId)

	case *Proto.AttackHit:
		fightAction.Type = Proto.ActionType_AttackHit
		fightAction.ActionAttackHit = pbAction

		log.Debug("战场ID：", uintptr(unsafe.Pointer(maker.scene)),
			"，指令：伤害体命中",
			"，时间(ms)：", maker.scene.NowTime,
			"，目标ID：", fightAction.ActionAttackHit.TargetId,
			"，伤害体ID：", fightAction.ActionAttackHit.AttackId,
			"，命中次数：", fightAction.ActionAttackHit.HitTimes,
			"，伤害类型(bits)：", fightAction.ActionAttackHit.DamageBit,
			"，扣除血量：", fightAction.ActionAttackHit.DamageHP,
			"，扣除护盾：", fightAction.ActionAttackHit.DamageHPShield)

	case *Proto.BeHit:
		fightAction.Type = Proto.ActionType_BeHit
		fightAction.ActionBeHit = pbAction

		log.Debug("战场ID：", uintptr(unsafe.Pointer(maker.scene)),
			"，指令：受击",
			"，时间(ms)：", maker.scene.NowTime,
			"，目标ID：", fightAction.ActionBeHit.TargetId,
			"，伤害体ID：", fightAction.ActionBeHit.AttackId,
			"，受击类型：", fightAction.ActionBeHit.HitType)

	case *Proto.AddBuff:
		fightAction.Type = Proto.ActionType_AddBuff
		fightAction.ActionAddBuff = pbAction

		log.Debug("战场ID：", uintptr(unsafe.Pointer(maker.scene)),
			"，指令：施加Buff",
			"，时间(ms)：", maker.scene.NowTime,
			"，施法者：", fightAction.ActionAddBuff.CasterId,
			"，目标ID：", fightAction.ActionAddBuff.TargetId,
			"，BuffID：", fightAction.ActionAddBuff.BuffId)

	case *Proto.RemoveBuff:
		fightAction.Type = Proto.ActionType_RemoveBuff
		fightAction.ActionRemoveBuff = pbAction

		log.Debug("战场ID：", uintptr(unsafe.Pointer(maker.scene)),
			"，指令：移除Buff",
			"，时间(ms)：", maker.scene.NowTime,
			"，目标ID：", fightAction.ActionRemoveBuff.SelfId,
			"，BuffID：", fightAction.ActionRemoveBuff.BuffId)

	case *Proto.ChangeAttr:
		fightAction.Type = Proto.ActionType_ChangeAttr
		fightAction.ActionChangeAttr = pbAction

		log.Debug("战场ID：", uintptr(unsafe.Pointer(maker.scene)),
			"，指令：修改属性",
			"，时间(ms)：", maker.scene.NowTime,
			"，PawnID：", fightAction.ActionChangeAttr.SelfId,
			"，属性类型：", fightAction.ActionChangeAttr.AttrType,
			"，旧值：", fightAction.ActionChangeAttr.OldValue,
			"，新值：", fightAction.ActionChangeAttr.NewValue)

	case *Proto.ChangeStat:
		fightAction.Type = Proto.ActionType_ChangeStat
		fightAction.ActionChangeStat = pbAction

		log.Debug("战场ID：", uintptr(unsafe.Pointer(maker.scene)),
			"，指令：修改状态",
			"，时间(ms)：", maker.scene.NowTime,
			"，PawnID：", fightAction.ActionChangeStat.SelfId,
			"，状态类型：", fightAction.ActionChangeStat.StatType,
			"，状态值：", fightAction.ActionChangeStat.StatValue)

	case *Proto.FightBegin:
		fightAction.Type = Proto.ActionType_FightBegin
		fightAction.ActionFightBegin = pbAction

		log.Debug("战场ID：", uintptr(unsafe.Pointer(maker.scene)),
			"，指令：战斗开始",
			"，时间(ms)：", maker.scene.NowTime)

	case *Proto.FightEnd:
		fightAction.Type = Proto.ActionType_FightEnd
		fightAction.ActionFightEnd = pbAction

		log.Debug("战场ID：", uintptr(unsafe.Pointer(maker.scene)),
			"，指令：战斗结束",
			"，时间(ms)：", maker.scene.NowTime)

	case *Proto.ChangeSkillState:
		fightAction.Type = Proto.ActionType_ChangeSkillState
		fightAction.ActionChangeSkillState = pbAction

		log.Debug("战场ID：", uintptr(unsafe.Pointer(maker.scene)),
			"，指令：修改技能状态",
			"，时间(ms)：", maker.scene.NowTime,
			"，PawnID：", fightAction.ActionChangeSkillState.CasterId,
			"，技能ID：", fightAction.ActionChangeSkillState.SkillId,
			"，状态：", fightAction.ActionChangeSkillState.State)

	case *Proto.CombineSkillEndTime:
		fightAction.Type = Proto.ActionType_CombineSkillEndTime
		fightAction.ActionCombineSkillEndTime = pbAction

	case *Proto.CombineSkillPoint:
		fightAction.Type = Proto.ActionType_CombineSkillPoint
		fightAction.ActionCombineSkillPoint = pbAction

	default:
		panic(fmt.Sprintf("invalid action type :%T", action))
	}

	framesLen := len(maker.replay.FrameList)

	if framesLen <= 0 || maker.scene.nowFrames > maker.replay.FrameList[framesLen-1].Sequence {
		if framesLen > 0 && len(maker.replay.FrameList[framesLen-1].ActionList) <= 0 {
			maker.replay.FrameList[framesLen-1].Sequence = maker.scene.nowFrames
		} else {
			frame := &Proto.FightFrame{}
			frame.Sequence = maker.scene.nowFrames

			maker.replay.FrameList = append(maker.replay.FrameList, frame)
		}
	}

	if len(maker.replay.FrameList) > 0 {
		frame := maker.replay.FrameList[len(maker.replay.FrameList)-1]
		frame.ActionList = append(frame.ActionList, fightAction)
	}
}

// PushDebugAction 添加debug动作
func (maker *_ReplayMaker) PushDebugAction(action interface{}) {
	fightAction := &Proto.FightAction{}

	switch pbAction := action.(type) {
	case *Proto.ChangeAttr:
		fightAction.Type = Proto.ActionType_ChangeAttr
		fightAction.ActionChangeAttr = pbAction

		log.Debug("战场ID：", uintptr(unsafe.Pointer(maker.scene)),
			"，指令：模拟器DEBUG（修改属性）",
			"，时间(ms)：", maker.scene.NowTime,
			"，PawnID：", fightAction.ActionChangeAttr.SelfId,
			"，属性类型：", fightAction.ActionChangeAttr.AttrType,
			"，旧值：", fightAction.ActionChangeAttr.OldValue,
			"，新值：", fightAction.ActionChangeAttr.NewValue)

	case *Proto.DebugInfo:
		fightAction.Type = Proto.ActionType_DebugInfo
		fightAction.ActionDebugInfo = pbAction

		log.Debug("战场ID：", uintptr(unsafe.Pointer(maker.scene)),
			"，指令：模拟器DEBUG（逻辑信息）",
			"，时间(ms)：", maker.scene.NowTime,
			"，信息：", fightAction.ActionDebugInfo.Message)

	case *Proto.AttackShowAoe:
		fightAction.Type = Proto.ActionType_AttackShowAoe
		fightAction.ActionAttackShowAoe = pbAction

		log.Debug("战场ID：", uintptr(unsafe.Pointer(maker.scene)),
			"，指令：模拟器DEBUG（显示伤害体Aoe）",
			"，时间(ms)：", maker.scene.NowTime,
			"，数据：", fmt.Sprintf("%+v", fightAction.ActionAttackShowAoe))

	case *Proto.AttackMoveAoe:
		fightAction.Type = Proto.ActionType_AttackMoveAoe
		fightAction.ActionAttackMoveAoe = pbAction

		log.Debug("战场ID：", uintptr(unsafe.Pointer(maker.scene)),
			"，指令：模拟器DEBUG（移动伤害体Aoe）",
			"，时间(ms)：", maker.scene.NowTime,
			"，数据：", fmt.Sprintf("%+v", fightAction.ActionAttackMoveAoe))
	}

	if maker.scene.SimulatorMode() {
		framesLen := len(maker.debugReplay.FrameList)

		if framesLen <= 0 || maker.scene.nowFrames > maker.debugReplay.FrameList[framesLen-1].Sequence {
			if framesLen > 0 && len(maker.debugReplay.FrameList[framesLen-1].ActionList) <= 0 {
				maker.debugReplay.FrameList[framesLen-1].Sequence = maker.scene.nowFrames
			} else {
				frame := &Proto.FightFrame{}
				frame.Sequence = maker.scene.nowFrames

				maker.debugReplay.FrameList = append(maker.debugReplay.FrameList, frame)
			}
		}

		if len(maker.debugReplay.FrameList) > 0 {
			frame := maker.debugReplay.FrameList[len(maker.debugReplay.FrameList)-1]
			frame.ActionList = append(frame.ActionList, fightAction)
		}
	}
}

// PushDebugInfo 记录debug信息
func (maker *_ReplayMaker) PushDebugInfo(fun func() string) {
	if !maker.scene.SimulatorMode() && !maker.scene.TestMode() {
		return
	}

	debugInfo := &Proto.DebugInfo{
		Message: fun(),
	}

	if debugInfo.Message == "" {
		return
	}

	maker.PushDebugAction(debugInfo)
}

// PushDetailDebugInfo 记录详细debug信息
func (maker *_ReplayMaker) PushDetailDebugInfo(fun func() string) {
	if !maker.scene.SimulatorMode() && !maker.scene.TestMode() {
		return
	}

	debugInfo := &Proto.DebugInfo{
		Message: fun(),
		Type:    Proto.FightLogType_Detail,
	}

	if debugInfo.Message == "" {
		return
	}

	maker.PushDebugAction(debugInfo)
}

// PushProgrammerDebugInfo 记录程序debug信息
func (maker *_ReplayMaker) PushProgrammerDebugInfo(fun func() string) {
	if !maker.scene.SimulatorMode() && !maker.scene.TestMode() {
		return
	}

	debugInfo := &Proto.DebugInfo{
		Message: fun(),
		Type:    Proto.FightLogType_Programmer,
	}

	if debugInfo.Message == "" {
		return
	}

	maker.PushDebugAction(debugInfo)
}
