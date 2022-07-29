package buffeffect

import (
	. "Daisy/Fight/internal"
	"Daisy/Fight/internal/conf"
	"Daisy/Fight/internal/effects"
)

type _2002_HPShield struct {
	effects.Blank
	Args   *conf.HPShieldArgs
	buff   *Buff
	shield *HPShield
}

// Init 初始化
func (effect *_2002_HPShield) Init(buff *Buff) error {
	effect.buff = buff
	return nil
}

// OnBuffAdd 施加buff后（buff自身能收到）
func (effect *_2002_HPShield) OnBuffAdd(buff *Buff) {
	caster := effect.buff.Caster
	pawn := effect.buff.Pawn

	var value float64

	// 取出原始数值
	switch effect.Args.HPShieldSource {
	case conf.ShieldValueSrc_CasterMaxHP:
		value = float64(caster.Attr.MaxHP)
	case conf.ShieldValueSrc_ExtValue:
		value = float64(effect.buff.ExtValue)
	default:
		effect.buff.Destroy()
		return
	}

	// 计算护盾HP
	shieldHP := int64(Max(value*effect.Args.HPShieldAddRate+float64(effect.Args.HPShieldAddFix), 1))

	// 增加护盾
	effect.shield = pawn.Attr.AddHPShield(shieldHP, 100, effect)
	if effect.shield == nil {
		effect.buff.Destroy()
		return
	}
}

// OnBuffRemove 移除buff后（buff自身能收到）
func (effect *_2002_HPShield) OnBuffRemove(buff *Buff, clear bool) {
	if effect.shield == nil {
		return
	}

	buff.Pawn.Attr.RemoveHPShield(effect.shield.UID)
}

// OnShieldBroken 护盾破碎后（护盾绑定的effect能收到）
func (effect *_2002_HPShield) OnShieldBroken(attack *Attack, shield *HPShield) {
	effect.buff.Destroy()
}
