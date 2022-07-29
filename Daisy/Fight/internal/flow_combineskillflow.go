package internal

import (
	"Daisy/Fight/internal/conf"
	"Daisy/Proto"
)

// CombineSkillCountdownTime 合体必杀技倒计时时长（毫秒）
const CombineSkillCountdownTime = 20000

// _CombineSkillFlow 合体必杀技流程
type _CombineSkillFlow struct {
	scene                        *Scene
	lockUltimateSkillPawn        *Pawn           // 锁定释放必杀技的pawn
	combineSkillCaster           *Pawn           // 合体技施法者
	combineSkillCountdownEndTime uint32          // 合体技倒计时
	combineSkillPointMap         map[uint32]bool // 合体技亮灯表
	combineSkillPointLimit       int             // 合体技亮灯上限
	combineSkillReadyMembers     []uint32        // 准备参与释放合体技的成员列表
}

// init 初始化
func (flow *_CombineSkillFlow) init(scene *Scene) {
	flow.scene = scene
	flow.combineSkillPointMap = make(map[uint32]bool)

	// 统计合体技亮灯上限
	for _, pawn := range flow.scene.formationList[Proto.Camp_Red].PawnList {
		if pawn.IsRole() {
			flow.combineSkillPointLimit++
		}
	}

	// 创建合体技施法者
	if len(flow.scene.formationList[Proto.Camp_Red].Info.CombineSkills) > 0 {
		flow.combineSkillCaster = flow.scene.GetBackgroundPawn(Proto.Camp_Red)
	}
}

// update 帧更新
func (flow *_CombineSkillFlow) update() {
	flow.CombineSkillReadyCast()
	flow.CastCombineSkill()
}

// LockCastUltimateSkill 锁定释放必杀技
func (flow *_CombineSkillFlow) LockCastUltimateSkill(pawn *Pawn) bool {
	if flow.lockUltimateSkillPawn != nil {
		return false
	}

	flow.lockUltimateSkillPawn = pawn

	return true
}

// UnlockCastUltimateSkill 解锁释放必杀技
func (flow *_CombineSkillFlow) UnlockCastUltimateSkill(pawn *Pawn) {
	if flow.lockUltimateSkillPawn != pawn {
		return
	}

	flow.lockUltimateSkillPawn = nil
}

// CastUltimateSkillIsLocked 是否已锁定释放必杀技
func (flow *_CombineSkillFlow) CastUltimateSkillIsLocked() bool {
	return flow.lockUltimateSkillPawn != nil
}

// StartCombineSkillCountdown 开始合体必杀技倒计时
func (flow *_CombineSkillFlow) StartCombineSkillCountdown() {
	// 检测能否开始倒计时
	if flow.combineSkillPointLimit <= 1 || len(flow.combineSkillReadyMembers) > 0 || flow.combineSkillCountdownEndTime > 0 {
		return
	}

	// 设置倒计时结束时间
	flow.combineSkillCountdownEndTime = flow.scene.NowTime + CombineSkillCountdownTime

	// 记录回放
	flow.scene.PushAction(&Proto.CombineSkillEndTime{
		EndTime: flow.combineSkillCountdownEndTime,
		Cancel:  false,
	})
}

// CancelCombineSkillCountdown 取消合体必杀技倒计时
func (flow *_CombineSkillFlow) CancelCombineSkillCountdown() {
	// 检测能否取消倒计时
	if flow.combineSkillCountdownEndTime <= 0 {
		return
	}

	// 重置倒计时时间
	flow.combineSkillCountdownEndTime = 0

	// 记录回放
	flow.scene.PushAction(&Proto.CombineSkillEndTime{
		EndTime: 0,
		Cancel:  true,
	})
}

// CombineSkillIsCountdown 是否正在进行合体必杀技倒计时
func (flow *_CombineSkillFlow) CombineSkillIsCountdown() bool {
	return flow.combineSkillCountdownEndTime > 0
}

// TurnOnCombineSkillPoint 合体必杀技亮灯
func (flow *_CombineSkillFlow) TurnOnCombineSkillPoint(pawn *Pawn) {
	// 检测能否亮灯
	if flow.combineSkillPointLimit <= 1 || len(flow.combineSkillReadyMembers) > 0 {
		return
	}

	// 检测是否已亮灯
	if _, ok := flow.combineSkillPointMap[pawn.UID]; ok {
		return
	}

	// 设置亮灯
	flow.combineSkillPointMap[pawn.UID] = true

	// 记录回放
	flow.scene.PushAction(&Proto.CombineSkillPoint{
		Point: int32(len(flow.combineSkillPointMap)),
	})
}

// TurnOffCombineSkillPoint 合体必杀技灭灯
func (flow *_CombineSkillFlow) TurnOffCombineSkillPoint(pawn *Pawn) {
	// 检测是否已灭灯
	if _, ok := flow.combineSkillPointMap[pawn.UID]; !ok {
		return
	}

	// 设置灭灯
	delete(flow.combineSkillPointMap, pawn.UID)

	// 记录回放
	flow.scene.PushAction(&Proto.CombineSkillPoint{
		Point: int32(len(flow.combineSkillPointMap)),
	})
}

// ClearCombineSkillPoint 所有人合体必杀技灭灯
func (flow *_CombineSkillFlow) ClearCombineSkillPoint() {
	// 检测是否已全部灭灯
	if len(flow.combineSkillPointMap) <= 0 {
		return
	}

	// 重置合体技亮灯表
	flow.combineSkillPointMap = make(map[uint32]bool)

	// 记录回放
	flow.scene.PushAction(&Proto.CombineSkillPoint{
		Point: 0,
	})
}

// CombineSkillReadyCast 合体技准备释放
func (flow *_CombineSkillFlow) CombineSkillReadyCast() {
	// 检测释放条件
	if len(flow.combineSkillReadyMembers) > 0 || (len(flow.combineSkillPointMap) < flow.combineSkillPointLimit &&
		(flow.combineSkillCountdownEndTime <= 0 || flow.scene.NowTime < flow.combineSkillCountdownEndTime)) {
		return
	}

	// 记录准备参与合体技成员
	for pawnID := range flow.combineSkillPointMap {
		flow.combineSkillReadyMembers = append(flow.combineSkillReadyMembers, pawnID)
	}

	// 单人无法释放合体超必杀
	if len(flow.combineSkillReadyMembers) <= 1 {
		flow.combineSkillReadyMembers = nil
	}

	// 取消倒计时
	flow.CancelCombineSkillCountdown()

	// 所有人合体必杀技灭灯
	flow.ClearCombineSkillPoint()
}

// CastCombineSkill 释放合体技
func (flow *_CombineSkillFlow) CastCombineSkill() {
	// 检测释放条件
	if len(flow.combineSkillReadyMembers) <= 0 || flow.CastUltimateSkillIsLocked() || flow.combineSkillCaster == nil ||
		len(flow.combineSkillCaster.combineSkillList) <= 0 {
		return
	}

	// 合体技
	combineSkill := flow.combineSkillCaster.combineSkillList[0]
	if len(combineSkill.Config.TemplateConfig.AttackConfs) <= 0 {
		return
	}

	// 首个伤害体配置
	atkConf := combineSkill.Config.TemplateConfig.AttackConfs[0]

	var target *Pawn

	// 查找释放目标
	for _, pawn := range flow.scene.pawnList {
		switch atkConf.TargetCategory {
		case conf.AttackTargetCategory_Enemy:
			if !pawn.State.CantBeEnemySelect && pawn.GetCamp() != flow.combineSkillCaster.GetCamp() {
				target = pawn
				break
			}
		case conf.AttackTargetCategory_Friend:
			if !pawn.State.CantBeFriendlySelect && pawn.GetCamp() == flow.combineSkillCaster.GetCamp() {
				target = pawn
				break
			}
		}
	}

	if target == nil {
		return
	}

	// 释放合体技
	flow.combineSkillCaster.UseSkill(combineSkill, target.GetPos(), []*Pawn{target})

	// 清除准备释放合体技的成员列表
	flow.combineSkillReadyMembers = nil
}

// GetCombineSkillReadyMembers 获取合体技准备释放的成员列表
func (flow *_CombineSkillFlow) GetCombineSkillReadyMembers() []uint32 {
	return flow.combineSkillReadyMembers
}
