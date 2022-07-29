package buffeffect

import (
	. "Daisy/Fight/attraffix"
	. "Daisy/Fight/internal"
	"Daisy/Fight/internal/conf"
	"Daisy/Fight/internal/effects"
)

type _2000_ChangeAttr struct {
	effects.Blank
	buff *Buff
	Args *conf.ChangeAttrArgs
}

func (effect *_2000_ChangeAttr) Init(buff *Buff) error {
	effect.buff = buff

	return nil
}

func (effect *_2000_ChangeAttr) OnBuffAdd(buff *Buff) {
	for attrName, attrArgs := range effect.Args.ChangeAttr {
		effect.ChangeAttr(attrName, attrArgs.Rate, attrArgs.Fix, true)
	}
}

func (effect *_2000_ChangeAttr) OnBuffRemove(buff *Buff, clear bool) {
	for field, attrArgs := range effect.Args.ChangeAttr {
		effect.ChangeAttr(field, attrArgs.Rate, attrArgs.Fix, false)
	}
}

// ChangeAttr 修改属性数值
func (effect *_2000_ChangeAttr) ChangeAttr(field uint32, rate float32, fix float64, sign bool) {
	effect.buff.Pawn.Attr.ChangeAttr(Field(field), rate, fix, sign)
}
