package buffeffect

import (
	. "Daisy/Fight/internal"
	"Daisy/Fight/internal/conf"
	"Daisy/Fight/internal/effects"
)

type _2001_ChangeState struct {
	effects.Blank
	buff *Buff
	Args *conf.ChangeStateArgs
}

func (effect *_2001_ChangeState) Init(buff *Buff) error {
	effect.buff = buff

	return nil
}

func (effect *_2001_ChangeState) OnBuffAdd(buff *Buff) {
	for _, stateName := range effect.Args.ChangeState {
		effect.ChangeState(stateName, true)
	}
}

func (effect *_2001_ChangeState) OnBuffRemove(buff *Buff, clear bool) {
	for _, stateName := range effect.Args.ChangeState {
		effect.ChangeState(stateName, false)
	}
}

// ChangeAttr 修改属性数值
func (effect *_2001_ChangeState) ChangeState(stateType uint32, turnOn bool) {
	effect.buff.Pawn.State.ChangeStat(Stat(stateType), turnOn)
}
