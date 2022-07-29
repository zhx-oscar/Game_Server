package effects

import (
	. "Daisy/Fight/internal"
	. "Daisy/Fight/internal/conf"
	"fmt"
)

// decHPShield 扣除护盾
func (damage *Damage) decHPShield(damageCtx *DamageContext) int64 {
	targetAttr := &damageCtx.Target.Attr

	// 基础伤害
	damageValue := damageCtx.DamageValue()

	damageCtx.PushDebugInfo(func() string {
		return fmt.Sprintf("\t\t扣除护盾：\n") +
			fmt.Sprintf("\t\t\t\t当前伤害：%f\n", damageValue) +
			fmt.Sprintf("\t\t\t\t扣除前HP护盾：%d\n", targetAttr.AllHPShield)
	})

	// 扣除护盾
	damageHPShield, exemptionDamage := targetAttr.BeDamageChangeHPShield(damageCtx.Attack, damageCtx.DamageBit, int64(damageValue))

	damageCtx.PushDebugInfo(func() string {
		return fmt.Sprintf("\t\t\t\t扣除HP护盾：%d\n", damageHPShield) +
			fmt.Sprintf("\t\t\t\t免疫伤害：%t\n", exemptionDamage) +
			fmt.Sprintf("\t\t\t\t扣除后HP护盾：%d\n", targetAttr.AllHPShield)
	})

	if damageHPShield <= 0 {
		return 0
	}

	// 记录扣除的HP护盾值
	damageCtx.DamageHPShield += damageHPShield

	// 扣除护盾比例
	deductRate := Max(1-float64(damageHPShield)/damageValue, 0)

	// 按比例降低伤害值
	for i := DamageValueKind_Begin; i < DamageValueKind_End; i++ {
		// 基础伤害
		baseDamage := damageCtx.DamageValueTab[i]

		// 护盾修正
		damageCtx.DamageValueTab[i] = Ceil(Max(baseDamage*deductRate, 0))

		damageCtx.PushDebugInfo(func() string {
			return fmt.Sprintf("\t\t\t\t%s护盾修正：Ceil(Max(%f * %f, 0)) = %f\n",
				getDamageValueKindText(i),
				baseDamage,
				deductRate,
				damageCtx.DamageValueTab[i])
		})
	}

	return damageHPShield
}
