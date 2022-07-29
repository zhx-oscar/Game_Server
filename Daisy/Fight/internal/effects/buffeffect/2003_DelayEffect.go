package buffeffect

import (
	. "Daisy/Fight/internal"
	"Daisy/Fight/internal/conf"
	"Daisy/Fight/internal/effects"
)

type _2003_DelayEffect struct {
	effects.Blank
	buff    *Buff
	Args    *conf.DelayEffectArgs
	AddTime uint32
}

func (effect *_2003_DelayEffect) Init(buff *Buff) error {
	effect.buff = buff
	effect.AddTime = buff.Pawn.Scene.NowTime

	return nil
}

// OnBuffUpdate buff状态更新
func (effect *_2003_DelayEffect) OnBuffUpdate(buff *Buff) {
	pawn := buff.Caster
	if pawn.Scene.NowTime < effect.AddTime+effect.Args.Time {
		return
	}

	pawn.BatchAddBuffs(pawn, effect.Args.BuffId, 0)

	buff.Destroy()
}
