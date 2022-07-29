package buffeffect

import (
	. "Daisy/Fight/internal"
	"Daisy/Fight/internal/effects"
	"Daisy/Proto"
	"fmt"
)

type _1004_BornAct struct {
	effects.Blank
	buff        *Buff
	pawnPause   bool
	buffUpdated bool
	halo        *Halo
}

func (effect *_1004_BornAct) Init(buff *Buff) error {
	effect.buff = buff
	effect.halo = nil

	return nil
}

// OnBuffAdd 施加buff后（buff自身能收到）
func (effect *_1004_BornAct) OnBuffAdd(buff *Buff) {
	pawn := effect.buff.Pawn

	effect.buffUpdated = false

	if pawn.Info.Type == Proto.PawnType_Npc {
		if pawn.Info.Npc == nil {
			effect.buff.Destroy()
			return
		}

		if pawn.Info.Npc.SpawnID != "" {
			effect.buff.Destroy()
			return
		}
	}

	if pawn.Info.ActConf.Born.Time <= 0 {
		effect.buff.Destroy()
		return
	}
}

// OnBuffUpdate buff帧更新（buff自身能收到）
func (effect *_1004_BornAct) OnBuffUpdate(buff *Buff) {
	pawn := buff.Pawn
	if effect.buffUpdated {
		return
	}

	effect.buff.SetDestroyTime(pawn.Scene.NowTime + pawn.Info.ActConf.Born.Time)

	if pawn.Info.BornPauseAllAI {
		effect.halo = pawn.Scene.AddHalo(pawn, effect)
	} else {
		pawn.AIPause(true)
	}

	effect.pawnPause = true
	effect.buffUpdated = true
}

// OnBeHit 受到打击后（所有buff能收到）
func (effect *_1004_BornAct) OnBeHit(attack *Attack, damageBit Bits, hitType Proto.HitType_Enum) {
	if hitType == Proto.HitType_None {
		return
	}

	effect.buff.Destroy()
}

// OnBuffRemove 移除buff后（buff自身能收到）
func (effect *_1004_BornAct) OnBuffRemove(buff *Buff, clear bool) {
	if effect.pawnPause {
		if buff.Pawn.Info.BornPauseAllAI {
			if effect.halo != nil {
				buff.Pawn.Scene.RemoveHalo(effect.halo.UID)
			}
		} else {
			effect.buff.Pawn.AIPause(false)
		}
	}
}

// OnAddHaloMember 增加光环成员(光环绑定的effect能收到)
func (effect *_1004_BornAct) OnAddHaloMember(halo *Halo, pawn *Pawn) {
	if pawn.Info.Type == Proto.PawnType_BG {
		return
	}

	pawn.Scene.PushDebugInfo(func() string {
		return fmt.Sprintf("${PawnID:%d}的Buff${BuffID:%d}创建的光环%d对${PawnID:%d}生效",
			effect.buff.Pawn.UID,
			effect.buff.Config.MainID(),
			halo.UID,
			pawn.UID)
	})

	pawn.AIPause(true)
}

// OnRemoveHaloMember 移除光环成员(光环绑定的effect能收到)
func (effect *_1004_BornAct) OnRemoveHaloMember(halo *Halo, pawn *Pawn) {
	if pawn.Info.Type == Proto.PawnType_BG {
		return
	}

	pawn.Scene.PushDebugInfo(func() string {
		return fmt.Sprintf("${PawnID:%d}的Buff${BuffID:%d}创建的光环%d对${PawnID:%d}失效",
			effect.buff.Pawn.UID,
			effect.buff.Config.MainID(),
			halo.UID,
			pawn.UID)
	})

	pawn.AIPause(false)
}
