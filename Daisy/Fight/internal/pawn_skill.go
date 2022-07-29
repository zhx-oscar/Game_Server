package internal

import (
	"Cinder/Base/linemath"
	"Daisy/Fight/internal/conf"
	"Daisy/Proto"
	"fmt"
	"math/rand"
)

// _PawnSkill pawn技能栏模块
type _PawnSkill struct {
	pawn *Pawn

	// 普攻相关
	normalAttackList          []*_SkillItem // 普攻列表
	overDriveNormalAttackList []*_SkillItem // 超载普攻列表
	comboNormalAtkCastOrder   []int         // 普攻连段释放顺序
	curComboNormalAtkIndex    int           // 当前普攻连段序号
	addComboNormalAttack      *_SkillItem   // 普攻追加连段
	breakNormalAtkCombo       bool          // 打断当前普攻连段
	insertNormalAtkSuperSkill *_SkillItem   // 插入普攻连段的超能技
	normalAttackTarget        *Pawn         // 普攻目标

	// 超能技相关
	superSkillList          []*_SkillItem // 超能技列表
	overDriveSuperSkillList []*_SkillItem // 超载超能技列表

	// 必杀技相关
	ultimateSkillList []*_SkillItem // 必杀技表

	// 合体技相关
	combineSkillList []*_SkillItem // 合体技表

	// 当前技能相关
	curSkill *Skill // 当前技能

	// 其他
	publicCDEndTime uint32 // 公共CD结束时间
}

// init 初始化
func (pawnSkill *_PawnSkill) init(pawn *Pawn) error {
	pawnSkill.pawn = pawn

	// 初始化普攻
	for _, normalAttackID := range pawn.Info.NormalAtkList {
		normalAttack := &_SkillItem{}
		if err := normalAttack.init(pawn, normalAttackID); err != nil {
			return fmt.Errorf("pawn type %d config %d init NormalAttack item %d failed, %s", pawn.Info.Type, pawn.Info.ConfigId, normalAttackID, err.Error())
		}

		pawn.normalAttackList = append(pawn.normalAttackList, normalAttack)
	}

	// 初始化超载普攻
	for _, overDriveNormalAttackID := range pawn.Info.OverDriveNormalAttackList {
		overDriveNormalAttack := &_SkillItem{}
		if err := overDriveNormalAttack.init(pawn, overDriveNormalAttackID); err != nil {
			return fmt.Errorf("pawn type %d config %d init OverDriveNormalAttack item %d failed, %s", pawn.Info.Type, pawn.Info.ConfigId, overDriveNormalAttackID, err.Error())
		}

		pawn.overDriveNormalAttackList = append(pawn.overDriveNormalAttackList, overDriveNormalAttack)
	}

	// 初始化追加普攻连段
	if pawn.Info.AddComboAttack > 0 {
		addComboAttack := &_SkillItem{}
		if err := addComboAttack.init(pawn, pawn.Info.AddComboAttack); err != nil {
			return fmt.Errorf("pawn type %d config %d init AddComboAttack item %d failed, %s", pawn.Info.Type, pawn.Info.ConfigId, pawn.Info.AddComboAttack, err.Error())
		}

		pawn.addComboNormalAttack = addComboAttack
	}

	// 初始化超能技
	for _, skillVID := range pawn.Info.SuperSkillList {
		newSkill := &_SkillItem{}
		if err := newSkill.init(pawn, skillVID); err != nil {
			return fmt.Errorf("pawn type %d config %d init SuperSkill item %d failed, %s", pawn.Info.Type, pawn.Info.ConfigId, skillVID, err.Error())
		}

		pawn.superSkillList = append(pawn.superSkillList, newSkill)
	}

	// 初始化超载超能技
	for _, overDriveSkillID := range pawn.Info.OverDriveSuperSkillList {
		newOverDriveSkill := &_SkillItem{}
		if err := newOverDriveSkill.init(pawn, overDriveSkillID); err != nil {
			return fmt.Errorf("pawn type %d config %d init OverDriveSuperSkill item %d failed, %s", pawn.Info.Type, pawn.Info.ConfigId, overDriveSkillID, err.Error())
		}

		pawn.overDriveSuperSkillList = append(pawn.overDriveSuperSkillList, newOverDriveSkill)
	}

	// 初始化必杀技
	for _, skillVID := range pawn.Info.UltimateSkillList {
		newSkill := &_SkillItem{}
		if err := newSkill.init(pawn, skillVID); err != nil {
			return fmt.Errorf("pawn type %d config %d init UltimateSkill item %d failed, %s", pawn.Info.Type, pawn.Info.ConfigId, skillVID, err.Error())
		}

		pawn.ultimateSkillList = append(pawn.ultimateSkillList, newSkill)
	}

	// 初始化合体技
	for _, skillVID := range pawn.Info.CombineSkillList {
		newSkill := &_SkillItem{}
		if err := newSkill.init(pawn, skillVID); err != nil {
			return fmt.Errorf("pawn type %d config %d init CombineSkill item %d failed, %s", pawn.Info.Type, pawn.Info.ConfigId, skillVID, err.Error())
		}

		pawn.combineSkillList = append(pawn.combineSkillList, newSkill)
	}

	if pawn.Scene.Info.Inherit != nil && pawn.IsRole() {
		if pawnInherit, ok := pawn.Scene.Info.Inherit.PawnInheritMap[pawn.Info.Role.RoleId]; ok {
			pawn.Attr.ChangeUltimateSkillPower(pawnInherit.UltimateSkillPower)
		}
	}

	return nil
}

// IsSkillRunning 是否正在使用技能（所有技能）
func (pawnSkill *_PawnSkill) IsSkillRunning() bool {
	pawn := pawnSkill.pawn

	return pawn.Scene.isSkillRunning(pawn)
}

// IsNormalAttackRunning 是否正在使用普攻
func (pawnSkill *_PawnSkill) IsNormalAttackRunning() bool {
	pawn := pawnSkill.pawn

	if pawn.curSkill == nil {
		return false
	}

	return pawn.Scene.isSkillRunning(pawn) && pawn.curSkill.IsNormalAttack()
}

// IsSuperSkillRunning 是否正在使用超能技
func (pawnSkill *_PawnSkill) IsSuperSkillRunning() bool {
	pawn := pawnSkill.pawn

	if pawn.curSkill == nil {
		return false
	}

	return pawn.Scene.isSkillRunning(pawn) && pawnSkill.curSkill.IsSuperSkill()
}

// IsUltimateSkillRunning 是否正在使用必杀技
func (pawnSkill *_PawnSkill) IsUltimateSkillRunning() bool {
	pawn := pawnSkill.pawn

	if pawn.curSkill == nil {
		return false
	}

	return pawn.Scene.isSkillRunning(pawn) && pawn.curSkill.IsUltimateSkill()
}

// IsCombineSkillRunning 是否正在使用合体技
func (pawnSkill *_PawnSkill) IsCombineSkillRunning() bool {
	pawn := pawnSkill.pawn

	if pawn.curSkill == nil {
		return false
	}

	return pawn.Scene.isSkillRunning(pawn) && pawn.curSkill.IsCombineSkill()
}

// CanUseSkill 能否使用指定技能
func (pawnSkill *_PawnSkill) CanUseSkill(skillItem *_SkillItem) bool {
	pawn := pawnSkill.pawn

	if skillItem == nil || !pawn.IsAlive() {
		return false
	}

	// 检测能否使用技能
	switch skillItem.Config.SkillKind {
	case conf.SkillKind_Super:
		if pawn.State.CantUseSuperSkill {
			return false
		}

		if pawn.publicCDEndTime > pawn.Scene.NowTime || skillItem.cdEndTime > pawn.Scene.NowTime {
			return false
		}

	case conf.SkillKind_NormalAtk:
		if pawn.State.CantUseNormalAtk {
			return false
		}

		if skillItem.cdEndTime > pawn.Scene.NowTime {
			return false
		}

	case conf.SkillKind_Ultimate:
		if pawn.State.CantUseUltimateSkill {
			return false
		}
	}

	// 玩家技能可以打断普攻
	if pawnSkill.IsSkillRunning() {
		if !pawn.IsRole() {
			return false
		}

		if !pawnSkill.IsNormalAttackRunning() {
			return false
		}
	}

	return true
}

// UseSkill 使用技能
func (pawnSkill *_PawnSkill) UseSkill(skillItem *_SkillItem, targetPos linemath.Vector2, targetList []*Pawn) bool {
	pawn := pawnSkill.pawn

	if skillItem == nil {
		return false
	}

	if pawn.IsRole() {
		switch skillItem.Config.SkillKind {
		case conf.SkillKind_NormalAtk:
			// 玩家使用普攻连段
			return pawnSkill.useComboNormalAtk(skillItem, targetPos, targetList)

		case conf.SkillKind_Super:
			// 普攻连段中插入超能技
			if pawnSkill.curSkill != nil && pawnSkill.curSkill.IsNormalAttack() {
				if pawnSkill.curSkill.Stat != Proto.SkillState_Later {
					// 当前技能是普攻命中阶段，则插入超能技
					pawnSkill.insertNormalAtkSuperSkill = skillItem
					return true
				}
			}
		}
	}

	return pawn.Scene.useSkill(pawn, skillItem, targetPos, targetList, false)
}

// useComboNormalAtk 使用连段普攻
func (pawnSkill *_PawnSkill) useComboNormalAtk(normalAtk *_SkillItem, targetPos linemath.Vector2, targetList []*Pawn) bool {
	pawn := pawnSkill.pawn

	if !normalAtk.IsNormalAttack() || len(targetList) <= 0 {
		return false
	}

	// 随机普攻序列
	pawnSkill.comboNormalAtkCastOrder = rand.Perm(len(pawnSkill.normalAttackList))
	if len(pawnSkill.comboNormalAtkCastOrder) <= 0 {
		return false
	}

	// 优先释放AI指定普攻
	for i, idx := range pawnSkill.comboNormalAtkCastOrder {
		if pawnSkill.normalAttackList[idx] == normalAtk {
			t := pawnSkill.comboNormalAtkCastOrder[0]
			pawnSkill.comboNormalAtkCastOrder[0], pawnSkill.comboNormalAtkCastOrder[i] = idx, t
		}
	}

	// 取出首个普攻
	pawnSkill.curComboNormalAtkIndex = 0
	normalSkill := pawnSkill.normalAttackList[pawnSkill.comboNormalAtkCastOrder[pawnSkill.curComboNormalAtkIndex]]

	// 设置连段参数
	pawnSkill.breakNormalAtkCombo = false
	pawnSkill.insertNormalAtkSuperSkill = nil
	pawnSkill.normalAttackTarget = targetList[0]

	// 使用普攻
	return pawn.Scene.useSkill(pawn, normalSkill, targetPos, targetList, false)
}

// ContinueComboAttack 继续连段攻击
func (pawnSkill *_PawnSkill) ContinueComboAttack() {
	// 检测连段普攻条件
	if !pawnSkill.pawn.IsRole() {
		return
	}
	if pawnSkill.curSkill == nil || !pawnSkill.curSkill.IsNormalAttack() || pawnSkill.curSkill.Stat != Proto.SkillState_Later {
		return
	}
	if pawnSkill.normalAttackTarget == nil || pawnSkill.normalAttackTarget.State.CantBeEnemySelect {
		return
	}

	scene := pawnSkill.pawn.Scene

	// 是否有插入普攻连段的超能技
	if pawnSkill.insertNormalAtkSuperSkill != nil {
		// 取出下个技能
		nextSkill := pawnSkill.insertNormalAtkSuperSkill
		pawnSkill.insertNormalAtkSuperSkill = nil

		// 当前技能跳过后摇
		if pawnSkill.curSkill.actionUseSKill != nil {
			pawnSkill.curSkill.actionUseSKill.SkipLater = true
		}

		// 打断当前技能并释放下个技能
		scene.breakCurSkill(pawnSkill.pawn, pawnSkill.pawn, Proto.SkillBreakReason_Combo)
		scene.useSkill(pawnSkill.pawn, nextSkill, pawnSkill.normalAttackTarget.GetPos(), []*Pawn{pawnSkill.normalAttackTarget}, true)

	} else {
		// 检测重置中断连段普攻
		if pawnSkill.breakNormalAtkCombo {
			return
		}

		// 下个普攻索引
		pawnSkill.curComboNormalAtkIndex++

		if pawnSkill.curComboNormalAtkIndex < len(pawnSkill.comboNormalAtkCastOrder) {
			// 取出下个普攻
			normalAttack := pawnSkill.normalAttackList[pawnSkill.comboNormalAtkCastOrder[pawnSkill.curComboNormalAtkIndex]]

			// 当前技能跳过后摇
			if pawnSkill.curSkill.actionUseSKill != nil {
				pawnSkill.curSkill.actionUseSKill.SkipLater = true
			}

			// 打断当前技能并释放下个普攻
			scene.breakCurSkill(pawnSkill.pawn, pawnSkill.pawn, Proto.SkillBreakReason_Combo)
			scene.useSkill(pawnSkill.pawn, normalAttack, pawnSkill.normalAttackTarget.GetPos(), []*Pawn{pawnSkill.normalAttackTarget}, true)

		} else {
			// 普攻追加连段
			if pawnSkill.addComboNormalAttack != nil && pawnSkill.addComboNormalAttack != pawnSkill.curSkill._SkillItem {
				// 当前技能跳过后摇
				if pawnSkill.curSkill.actionUseSKill != nil {
					pawnSkill.curSkill.actionUseSKill.SkipLater = true
				}

				// 打断当前技能并释放下个普攻
				scene.breakCurSkill(pawnSkill.pawn, pawnSkill.pawn, Proto.SkillBreakReason_Combo)
				scene.useSkill(pawnSkill.pawn, pawnSkill.addComboNormalAttack, pawnSkill.normalAttackTarget.GetPos(), []*Pawn{pawnSkill.normalAttackTarget}, true)
			}
		}
	}
}

// GetUsableNormalAttacks 获取可用普攻列表
func (pawnSkill *_PawnSkill) GetUsableNormalAttacks() []*_SkillItem {
	pawn := pawnSkill.pawn

	// 不能使用普攻
	if pawn.State.CantUseNormalAtk {
		return nil
	}

	// 有技能正在运行
	if pawn.IsSkillRunning() {
		return nil
	}

	// 有插入连段普攻的超能技
	if pawn.insertNormalAtkSuperSkill != nil {
		return nil
	}

	if pawn.State.OverDrive {
		return pawn.overDriveNormalAttackList
	}

	return pawn.normalAttackList
}

// GetUsableSuperSkills 获取可用超能技列表
func (pawnSkill *_PawnSkill) GetUsableSuperSkills() []*_SkillItem {
	pawn := pawnSkill.pawn

	// 不能使用超能技
	if pawn.State.CantUseSuperSkill {
		return nil
	}

	// 玩家超能技可以打断普攻
	if pawnSkill.IsSkillRunning() {
		if !pawn.IsRole() {
			return nil
		}

		if !pawnSkill.IsNormalAttackRunning() {
			return nil
		}
	}

	// 有插入连段普攻的超能技
	if pawn.insertNormalAtkSuperSkill != nil {
		return nil
	}

	// 判断cd
	if pawn.publicCDEndTime > pawn.Scene.NowTime {
		return nil
	}

	var skillList []*_SkillItem

	// 超载状态
	if pawn.State.OverDrive {
		for _, skill := range pawn.overDriveSuperSkillList {
			if skill.cdEndTime > pawn.Scene.NowTime {
				continue
			}

			skillList = append(skillList, skill)
		}

		return skillList
	}

	for _, skill := range pawn.superSkillList {
		if skill.cdEndTime > pawn.Scene.NowTime {
			continue
		}

		skillList = append(skillList, skill)
	}

	return skillList
}

// GetUsableUltimateSkills 获取可用必杀技列表
func (pawnSkill *_PawnSkill) GetUsableUltimateSkills() []*_SkillItem {
	pawn := pawnSkill.pawn

	// 不能使用超能及
	if pawn.State.CantUseUltimateSkill {
		return nil
	}

	// 玩家超能技可以打断普攻
	if pawnSkill.IsSkillRunning() && !pawnSkill.IsNormalAttackRunning() {
		return nil
	}

	// 必杀技能量未满或不能锁定释放必杀技
	if !pawn.Attr.UltimateSkillPowerIsFull() || pawn.Scene.CastUltimateSkillIsLocked() {
		return nil
	}

	return pawn.ultimateSkillList
}

// GetSuperSkillList 获取超能技列表
func (pawnSkill *_PawnSkill) GetSuperSkillList() []*_SkillItem {
	pawn := pawnSkill.pawn

	if pawn.State.OverDrive {
		return pawn.overDriveSuperSkillList
	}

	return pawn.superSkillList
}

// GetUltimateSkillList 获取必杀技列表
func (pawnSkill *_PawnSkill) GetUltimateSkillList() []*_SkillItem {
	pawn := pawnSkill.pawn

	return pawn.ultimateSkillList
}

// BreakCurSkill 打断当前技能
func (pawnSkill *_PawnSkill) BreakCurSkill(caster *Pawn, breakReason Proto.SkillBreakReason_Enum) bool {
	pawn := pawnSkill.pawn

	return pawn.Scene.breakCurSkill(caster, pawn, breakReason)
}

// BreakNormalAttackCombo 打断普攻连段
func (pawnSkill *_PawnSkill) BreakNormalAttackCombo(caster *Pawn) bool {
	if !pawnSkill.IsNormalAttackRunning() {
		return false
	}

	pawnSkill.breakNormalAtkCombo = true

	return true
}

// ResetPublicCDTime 重置公共CD结束时间
func (pawnSkill *_PawnSkill) ResetPublicCDTime() {
	pawnSkill.pawn.publicCDEndTime = pawnSkill.pawn.Scene.NowTime + pawnSkill.pawn.Info.PublicCD
}
