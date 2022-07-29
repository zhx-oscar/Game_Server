package internal

import "Daisy/Proto"

// _StatChanged 状态变化中间值
type _StatChanged struct {
	count int16
}

// _StatChangedTab 状态变化中间值表
type _StatChangedTab []_StatChanged

// ChangeStat 修改状态
func (state *FightState) ChangeStat(stat Stat, enable bool) {
	statChanged := &state.changedTab[stat]

	if enable {
		statChanged.count++
	} else {
		statChanged.count--
	}

	switch stat {
	case Stat_CantUseNormalAtk:
		state.CantUseNormalAtk = statChanged.count > 0

	case Stat_CantUseSuperSkill:
		state.CantUseSuperSkill = statChanged.count > 0

	case Stat_CantUseUltimateSkill:
		state.CantUseUltimateSkill = statChanged.count > 0

	case Stat_CantMove:
		state.CantMove = statChanged.count > 0

	case Stat_CantBeDamage:
		state.CantBeDamage = statChanged.count > 0

	case Stat_CantBeEnemySelect:
		state.CantBeEnemySelect = statChanged.count > 0

	case Stat_CantBeHitControl:
		state.CantBeHitControl = statChanged.count > 0

	case Stat_CantBeAddIncrBuff:
		state.CantBeAddIncrBuff = statChanged.count > 0

	case Stat_CantBeAddDecrBuff:
		state.CantBeAddDecrBuff = statChanged.count > 0

	case Stat_CantBeFriendlySelect:
		state.CantBeFriendlySelect = statChanged.count > 0

	case Stat_Invincible:
		state.Invincible = statChanged.count > 0

		// 修改一级状态
		state.ChangeStat(Stat_CantBeDamage, enable)
		state.ChangeStat(Stat_CantBeHitControl, enable)
		state.ChangeStat(Stat_CantBeAddDecrBuff, enable)

	case Stat_Weak:
		old := state.Weak
		state.Weak = statChanged.count > 0
		state.StateSyncToClient(Proto.StatType_Weak, old, state.Weak)

		// 修改一级状态
		state.ChangeStat(Stat_CantMove, enable)
		state.ChangeStat(Stat_CantUseNormalAtk, enable)
		state.ChangeStat(Stat_CantUseSuperSkill, enable)
		state.ChangeStat(Stat_CantUseUltimateSkill, enable)
		state.ChangeStat(Stat_CantBeHitControl, enable)

	case Stat_OverDrive:
		old := state.OverDrive
		state.OverDrive = statChanged.count > 0
		state.StateSyncToClient(Proto.StatType_OverDrive, old, state.OverDrive)

		// 修改一级状态
		state.ChangeStat(Stat_CantBeHitControl, enable)

	case Stat_OverDriveStartCG:
		state.OverDriveStartCG = statChanged.count > 0

		// 修改一级状态
		state.ChangeStat(Stat_CantMove, enable)
		state.ChangeStat(Stat_CantUseNormalAtk, enable)
		state.ChangeStat(Stat_CantUseSuperSkill, enable)
		state.ChangeStat(Stat_CantUseUltimateSkill, enable)
		state.ChangeStat(Stat_CantBeHitControl, enable)

	case Stat_Balance:
		old := state.Balance
		state.Balance = statChanged.count > 0
		state.StateSyncToClient(Proto.StatType_Break, !old, !state.Balance)

		// 修改一级状态
		state.ChangeStat(Stat_CantBeHitControl, enable)

	case Stat_Dodging:
		state.Dodging = statChanged.count > 0

		// 修改一级状态
		state.ChangeStat(Stat_CantMove, enable)
		state.ChangeStat(Stat_CantBeDamage, enable)
		state.ChangeStat(Stat_CantUseNormalAtk, enable)
		state.ChangeStat(Stat_CantUseSuperSkill, enable)
		state.ChangeStat(Stat_CantUseUltimateSkill, enable)
		state.ChangeStat(Stat_CantBeHitControl, enable)

	case Stat_Death:
		old := state.Death
		state.Death = statChanged.count > 0
		state.StateSyncToClient(Proto.StatType_Dead, old, state.Death)

		state.ChangeStat(Stat_CantMove, enable)
		state.ChangeStat(Stat_CantBeDamage, enable)
		state.ChangeStat(Stat_CantBeEnemySelect, enable)
		state.ChangeStat(Stat_CantBeFriendlySelect, enable)
		state.ChangeStat(Stat_CantUseNormalAtk, enable)
		state.ChangeStat(Stat_CantUseSuperSkill, enable)
		state.ChangeStat(Stat_CantUseUltimateSkill, enable)
		state.ChangeStat(Stat_CantBeHitControl, enable)

	case Stat_Raged:
		state.Raged = statChanged.count > 0

	case Stat_EnergyShieldOn:
		state.EnergyShieldOn = statChanged.count > 0

		state.ChangeStat(Stat_CantBeHitControl, enable)

	case Stat_Frozen:
		state.Frozen = statChanged.count > 0

		state.ChangeStat(Stat_CantMove, enable)
	}
}

// ChangeBeHitStatBit 修改受击状态比特位
func (state *FightState) ChangeBeHitStatBit(hitType Proto.HitType_Enum, enable bool) {
	// 不能重复置位
	if state.BeHitStat.Test(int32(hitType)) == enable {
		return
	}

	if hitType == Proto.HitType_None {
		return
	}

	if enable {
		state.BeHitStat.TurnOn(int32(hitType))
	} else {
		state.BeHitStat.TurnOff(int32(hitType))
	}

	// 修改一级状态
	state.ChangeStat(Stat_CantMove, enable)
	state.ChangeStat(Stat_CantUseNormalAtk, enable)
	state.ChangeStat(Stat_CantUseSuperSkill, enable)
	state.ChangeStat(Stat_CantUseUltimateSkill, enable)
}

// StateSyncToClient 状态同步客户端
func (state *FightState) StateSyncToClient(stateType Proto.StatType_Enum, old, new bool) {
	if state.pawn.Attr.inFight && new != old {
		state.pawn.Scene.PushAction(&Proto.ChangeStat{
			SelfId:    state.pawn.UID,
			StatType:  stateType,
			StatValue: new,
		})
	}
}
