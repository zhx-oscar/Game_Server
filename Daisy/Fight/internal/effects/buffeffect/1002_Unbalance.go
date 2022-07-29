package buffeffect

import (
	. "Daisy/Fight/internal"
	"Daisy/Fight/internal/conf"
	"Daisy/Fight/internal/effects"
	"Daisy/Proto"
)

type _1002_Unbalance struct {
	effects.Blank
	buff *Buff

	CurBreakValue         int32  //当前崩坏值
	BreakValueRestoreTime uint32 //崩坏值恢复时间
}

func (effect *_1002_Unbalance) Init(buff *Buff) error {
	effect.buff = buff
	pawn := effect.buff.Pawn

	effect.CurBreakValue = 0
	effect.BreakValueRestoreTime = 0

	if pawn.Attr.BreakValueLimit > 0 {
		pawn.State.ChangeStat(Stat_Balance, true)
	}

	return nil
}

// 受击
func (effect *_1002_Unbalance) OnBeDamage(attack *Attack, damageKind conf.DamageKind, damageBit Bits, damageValue, damageHP, damageHPShield int64) {
	pawn := effect.buff.Pawn
	if attack.Src() != Proto.AttackSrc_Skill {
		return
	}

	if damageValue < 0 {
		return
	}

	if !damageBit.Test(int32(Proto.DamageType_Damage)) {
		return
	}

	if pawn.Info.Type != Proto.PawnType_Npc {
		return
	}

	if pawn.Attr.BreakValueLimit <= 0 {
		return
	}

	if attack.Skill.Config.HitAddBreakValue <= 0 {
		return
	}

	if effect.CurBreakValue >= pawn.Attr.BreakValueLimit {
		return
	}

	if pawn.State.CantBeHitControl {
		return
	}

	if effect.CurBreakValue+attack.Skill.Config.HitAddBreakValue >= pawn.Attr.BreakValueLimit {
		effect.CurBreakValue = pawn.Attr.BreakValueLimit
		effect.BreakValueRestoreTime = pawn.Scene.NowTime + pawn.Attr.BreakStateTime
		// 失衡开始
		pawn.State.ChangeStat(Stat_Balance, false)
	} else {
		effect.CurBreakValue += attack.Skill.Config.HitAddBreakValue
	}
}

func (effect *_1002_Unbalance) OnBuffUpdate(buff *Buff) {
	pawn := effect.buff.Pawn
	if pawn.Info.Type == Proto.PawnType_BG {
		return
	}

	if pawn.Attr.BreakValueLimit <= 0 {
		return
	}

	if effect.CurBreakValue != pawn.Attr.BreakValueLimit {
		return
	}

	if effect.BreakValueRestoreTime > pawn.Scene.NowTime {
		return
	}

	effect.CurBreakValue = 0
	// 失衡结束
	pawn.State.ChangeStat(Stat_Balance, true)
}

func (effect *_1002_Unbalance) OnOverDrive() {
	pawn := effect.buff.Pawn
	// 解除失衡状态
	if !pawn.State.Balance {
		effect.CurBreakValue = 0
		effect.BreakValueRestoreTime = 0
		// 失衡结束
		if pawn.Attr.BreakValueLimit > 0 {
			pawn.State.ChangeStat(Stat_Balance, true)
		}
	}
}
