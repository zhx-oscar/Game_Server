package buffeffect

import (
	. "Daisy/Fight/internal"
	"Daisy/Fight/internal/effects"
	"Daisy/Proto"
)

type _1003_BlockBreak struct {
	effects.Blank
	buff *Buff
}

func (effect *_1003_BlockBreak) Init(buff *Buff) error {
	effect.buff = buff
	pawn := effect.buff.Pawn

	pawn.Attr.BlockUnbalanceDegrade = 5
	if pawn.IsBoss() {
		pawn.Attr.CrackUp_Block = 200
	} else {
		pawn.Attr.CrackUp_Block = 100
	}

	pawn.BeHit.CurBlockBreakValue = 0
	pawn.BeHit.BlockBreakValueDecTime = 0

	return nil
}

func (effect *_1003_BlockBreak) OnBuffUpdate(buff *Buff) {
	pawn := effect.buff.Pawn
	if pawn.State.BeHitStat.Test(int32(Proto.HitType_Block)) || pawn.State.BeHitStat.Test(int32(Proto.HitType_BlockBreak)) {
		return
	}

	if pawn.BeHit.BlockBreakValueDecTime > pawn.Scene.NowTime {
		return
	}

	oldValue := pawn.BeHit.CurBlockBreakValue
	pawn.BeHit.CurBlockBreakValue -= pawn.Attr.BlockUnbalanceDegrade
	if pawn.BeHit.CurBlockBreakValue < 0 {
		pawn.BeHit.CurBlockBreakValue = 0
	}

	if oldValue != pawn.BeHit.CurBlockBreakValue {
		pawn.Scene.PushDebugAction(&Proto.ChangeAttr{
			SelfId:   pawn.UID,
			AttrType: Proto.AttrType_CurBlockBreakValue,
			OldValue: float64(oldValue),
			NewValue: float64(pawn.BeHit.CurBlockBreakValue),
		})
	}

	pawn.BeHit.BlockBreakValueDecTime = pawn.Scene.NowTime + 1000
}

// OnBeHit 受到打击后
func (effect *_1003_BlockBreak) OnBeHit(attack *Attack, damageBit Bits, hitType Proto.HitType_Enum) {
	pawn := effect.buff.Pawn
	if attack.Src() != Proto.AttackSrc_Skill {
		return
	}

	if !damageBit.Test(int32(Proto.DamageType_Block)) {
		return
	}

	if !pawn.State.BeHitStat.Test(int32(Proto.HitType_Block)) {
		return
	}

	oldValue := pawn.BeHit.CurBlockBreakValue
	pawn.BeHit.CurBlockBreakValue += attack.Skill.Config.HitAddBlockBreakValue
	if pawn.BeHit.CurBlockBreakValue >= pawn.Attr.CrackUp_Block {
		pawn.BeHit.CurBlockBreakValue = pawn.Attr.CrackUp_Block
	}

	if oldValue != pawn.BeHit.CurBlockBreakValue {
		pawn.Scene.PushDebugAction(&Proto.ChangeAttr{
			SelfId:   pawn.UID,
			AttrType: Proto.AttrType_CurBlockBreakValue,
			OldValue: float64(oldValue),
			NewValue: float64(pawn.BeHit.CurBlockBreakValue),
		})
	}

	if pawn.BeHit.CurBlockBreakValue == pawn.Attr.CrackUp_Block {
		pawn.State.ChangeBeHitStatBit(Proto.HitType_Block, false)
		pawn.State.ChangeBeHitStatBit(Proto.HitType_BlockBreak, true)
	}
}
