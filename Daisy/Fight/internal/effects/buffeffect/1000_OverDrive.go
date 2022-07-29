package buffeffect

import (
	. "Daisy/Fight/internal"
	"Daisy/Fight/internal/conf"
	"Daisy/Fight/internal/effects"
	"Daisy/Proto"
	"math"
)

// OverDriveStage 超载阶段
type OverDriveStage uint32

const (
	OverDriveStage_Idle      OverDriveStage = iota // 0: 空闲阶段
	OverDriveStage_OverDrive                       // 1: 超载阶段
	OverDriveStage_Weak                            // 2: 虚弱阶段
)

type _1000_OverDrive struct {
	effects.Blank
	buff                 *Buff
	CurOverDriveValue    int32          // 当前超载值
	OverDriveStateTime   uint32         // 超载持续时间
	OverDriveRestoreTime uint32         // 超载恢复时间
	OverDriveCGEndTime   uint32         // 超载开始cg结束时间
	WeakEndTime          uint32         // 虚弱状态结束时间
	WeakDebuff           uint32         // 虚弱buff
	OverDriveStage       OverDriveStage // 超载阶段
}

func (effect *_1000_OverDrive) Init(buff *Buff) error {
	effect.buff = buff
	pawn := effect.buff.Pawn

	effect.CurOverDriveValue = 0
	effect.OverDriveStage = OverDriveStage_Idle
	if weakDebuff, ok := pawn.Scene.GetConstConIntValue(conf.ConstExcel_WeakDebuff); ok {
		effect.WeakDebuff = uint32(weakDebuff)
	}

	return nil
}

// OnBuffAdd 施加buff后（buff自身能收到）
func (effect *_1000_OverDrive) OnBuffAdd(buff *Buff) {
	pawn := effect.buff.Pawn
	if !pawn.IsBoss() || pawn.Attr.OverDriveLimit <= 0 {
		buff.Destroy()
	}
}

// 状态更新

// OnBuffUpdate buff状态更新
func (effect *_1000_OverDrive) OnBuffUpdate(buff *Buff) {
	effect.updateStage()
}

// updateOverDriveStage 更新
func (effect *_1000_OverDrive) updateStage() {
	switch effect.OverDriveStage {
	case OverDriveStage_Idle:
		effect.UpdateIdleStage()
	case OverDriveStage_OverDrive:
		effect.UpdateOverDriveStage()
	case OverDriveStage_Weak:
		effect.UpdateWeakStage()
	}
}

// UpdateIdleStage 更新空闲阶段
func (effect *_1000_OverDrive) UpdateIdleStage() {
	pawn := effect.buff.Pawn
	if pawn.State.BeHitStat != 0 {
		return
	}

	if pawn.IsSkillRunning() {
		return
	}

	if effect.CurOverDriveValue < pawn.Attr.OverDriveLimit {
		return
	}

	effect.changeStage(OverDriveStage_OverDrive)
}

// UpdateOverDriveStage 更新超载阶段
func (effect *_1000_OverDrive) UpdateOverDriveStage() {
	pawn := effect.buff.Pawn
	if pawn.IsBackground() {
		return
	}

	if pawn.State.OverDriveStartCG && effect.OverDriveCGEndTime <= pawn.Scene.NowTime {
		pawn.State.ChangeStat(Stat_OverDriveStartCG, false)
	}

	if effect.OverDriveRestoreTime > pawn.Scene.NowTime {
		return
	}

	if pawn.IsSkillRunning() {
		return
	}

	effect.changeStage(OverDriveStage_Weak)
}

// UpdateWeakStage 更新虚弱状态
func (effect *_1000_OverDrive) UpdateWeakStage() {
	pawn := effect.buff.Pawn
	if pawn.IsBackground() {
		return
	}

	if effect.WeakEndTime > pawn.Scene.NowTime {
		return
	}

	effect.changeStage(OverDriveStage_Idle)
}

// changeStage 修改阶段
func (effect *_1000_OverDrive) changeStage(stage OverDriveStage) {
	if effect.OverDriveStage == stage {
		return
	}

	// 调整阶段
	effect.OverDriveStage = stage

	// 发送事件
	switch stage {
	case OverDriveStage_Idle:
		effect.OnWeakEnd()
	case OverDriveStage_OverDrive:
		effect.OnOverDriveStart()
	case OverDriveStage_Weak:
		effect.OnWeakStart()
	}
}

// 事件回调

// OnSkillInEnd 技能结束（包含被打断和正常结束，当前技能与所有buff能收到）
func (effect *_1000_OverDrive) OnSkillInEnd(skill *Skill, lastStat Proto.SkillState_Enum, skillEndReason Proto.SkillEndReason_Enum, breakCaster *Pawn) {
	if effect.OverDriveStage != OverDriveStage_OverDrive {
		return
	}

	// 防止技能连放无法结束超载状态
	effect.UpdateOverDriveStage()
}

// OnBeDamage 受击
func (effect *_1000_OverDrive) OnBeDamage(attack *Attack, damageKind conf.DamageKind, damageBit Bits, damageValue, damageHP, damageHPShield int64) {
	pawn := effect.buff.Pawn

	if pawn.IsBackground() {
		return
	}

	if damageValue <= 0 {
		return
	}

	if !pawn.IsBoss() {
		return
	}

	if damageBit.Test(int32(Proto.DamageType_Miss)) || damageBit.Test(int32(Proto.DamageType_ExemptionDamage)) || damageBit.Test(int32(Proto.DamageType_Dodge)) {
		return
	}

	if !damageBit.Test(int32(Proto.DamageType_Damage)) {
		return
	}

	if effect.OverDriveStage != OverDriveStage_Idle {
		return
	}

	// 每伤害1%的血量，增加OverDriveAddEfficiency点超载值
	overDriveValueAdd := int32(float64(pawn.Attr.OverDriveAddEfficiency) * math.Abs(float64(damageValue)) / (0.01 * float64(pawn.Attr.MaxHP)))
	if overDriveValueAdd == 0 {
		overDriveValueAdd = 1
	}

	oldValue := effect.CurOverDriveValue

	effect.CurOverDriveValue += overDriveValueAdd
	if effect.CurOverDriveValue >= pawn.Attr.OverDriveLimit {
		effect.CurOverDriveValue = pawn.Attr.OverDriveLimit
	}

	pawn.Attr.AttrSyncToClient(Proto.AttrType_OverDrivePower, float64(oldValue), float64(effect.CurOverDriveValue))

	if effect.CurOverDriveValue < pawn.Attr.OverDriveLimit {
		return
	}

	if pawn.State.BeHitStat != 0 {
		return
	}

	if pawn.IsSkillRunning() {
		return
	}

	effect.changeStage(OverDriveStage_OverDrive)
}

// OnOverDriveStart 进入超载状态(超载buff能收到)
func (effect *_1000_OverDrive) OnOverDriveStart() {
	pawn := effect.buff.Pawn

	if pawn.Info.ActConf.OverDrive.Time > 0 {
		effect.OverDriveCGEndTime = pawn.Scene.NowTime + pawn.Info.ActConf.OverDrive.Time
		pawn.State.ChangeStat(Stat_OverDriveStartCG, true)
	}

	pawn.Stop()
	effect.OverDriveRestoreTime = effect.OverDriveCGEndTime + pawn.Attr.OverDriveTime
	pawn.State.ChangeStat(Stat_OverDrive, true)

	pawn.Events.EmitOverDrive()
}

// OnWeakStart 进入虚弱状态(超载buff能受到)
func (effect *_1000_OverDrive) OnWeakStart() {
	pawn := effect.buff.Pawn

	pawn.Stop()
	pawn.State.ChangeStat(Stat_OverDrive, false)
	pawn.State.ChangeStat(Stat_Weak, true)

	oldValue := effect.CurOverDriveValue
	effect.CurOverDriveValue = 0
	pawn.Attr.AttrSyncToClient(Proto.AttrType_OverDrivePower, float64(oldValue), float64(effect.CurOverDriveValue))

	pawn.AddBuff(pawn, effect.WeakDebuff, 0)
	weakTime := pawn.Attr.WeakTime
	if weakTime < pawn.Info.ActConf.WeakBegin.Time+pawn.Info.ActConf.WeakEnd.Time {
		weakTime = pawn.Info.ActConf.WeakBegin.Time + pawn.Info.ActConf.WeakBegin.Time
	}

	effect.WeakEndTime = pawn.Scene.NowTime + weakTime
}

// OnWeakEnd 虚弱状态结束(超载buff能受到)
func (effect *_1000_OverDrive) OnWeakEnd() {
	pawn := effect.buff.Pawn
	pawn.State.ChangeStat(Stat_Weak, false)

	if weakDebuffs, ok := pawn.GetSameBuffs(effect.WeakDebuff); ok {
		for _, weakDebuff := range weakDebuffs {
			pawn.RemoveBuff(weakDebuff.BuffKey)
		}
	}
}
