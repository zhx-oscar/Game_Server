package effects

import (
	. "Daisy/Fight/internal"
	. "Daisy/Fight/internal/conf"
	"Daisy/Proto"
	"fmt"
)

// stepDeductHP 步骤：伤害扣血
func (damage *Damage) stepDeductHP(damageCtx *DamageContext) bool {
	damageCtx.PushDebugInfo(func() string {
		damageCtx.StepIndex++
		return fmt.Sprintf("【%d：伤害扣血】\n", damageCtx.StepIndex)
	})

	casterEvents := &damageCtx.Attack.Caster.Events
	targetAttr := &damageCtx.Target.Attr

	// 发送事件
	casterEvents.EmitBeforeDamageStep(damageCtx, DamageStep_DeductHP)

	targetOldHP := targetAttr.CurHP

	// 统计伤害值
	damageValue, damageInvert := damageCtx.DamageValueWithInvert()

	if damageInvert {
		damageCtx.DamageBit.TurnOn(int32(Proto.DamageType_Invert))
	}

	// 扣除HP
	damageHP, exemptionDamage := targetAttr.BeDamageChangeHP(damageCtx.Attack, damage.DamageKind, damageCtx.DamageBit, int64(damageValue), damageCtx.DamageHPShield, false)

	damageCtx.PushDebugInfo(func() string {
		return fmt.Sprintf("\t\t需要扣除HP：%d\n", int64(damageValue)) +
			fmt.Sprintf("\t\t实际扣除HP：%d\n", damageHP) +
			fmt.Sprintf("\t\t免疫伤害：%t\n", exemptionDamage)
	})

	// 标记免伤
	if exemptionDamage {
		damageCtx.DamageBit.TurnOn(int32(Proto.DamageType_ExemptionDamage))
	}

	if damageInvert {
		damageCtx.PushDebugInfo(func() string {
			return fmt.Sprintf("\t\t目标剩余HP：Min(%d + %d, %d) = %d \n",
				targetOldHP,
				damageHP,
				targetAttr.MaxHP,
				targetAttr.CurHP)
		})
	} else {
		damageCtx.PushDebugInfo(func() string {
			return fmt.Sprintf("\t\t目标剩余HP：Max(%d - %d, 0) = %d \n",
				targetOldHP,
				damageHP,
				targetAttr.CurHP)
		})
	}

	// 记录扣除的HP值
	damageCtx.DamageHP = damageHP

	// 发送事件
	casterEvents.EmitAfterDamageStep(damageCtx, DamageStep_DeductHP)

	return true
}

// stepBloodsucker 步骤：吸血
func (damage *Damage) stepBloodsucker(damageCtx *DamageContext) bool {
	damageCtx.PushDebugInfo(func() string {
		damageCtx.StepIndex++
		return fmt.Sprintf("【%d：吸血】\n", damageCtx.StepIndex)
	})

	casterEvents := &damageCtx.Attack.Caster.Events

	// 发送事件
	casterEvents.EmitBeforeDamageStep(damageCtx, DamageStep_Bloodsucker)

	caster := damageCtx.Attack.Caster
	casterAttr := &damageCtx.Attack.CasterSnapshot.Attr

	// 计算吸血值
	damageCtx.BloodsuckerValue = int64(float64(damageCtx.DamageHP) * float64(casterAttr.AttackBloodsuckerRate))

	damageCtx.PushDebugInfo(func() string {
		return fmt.Sprintf("\t\t吸血值：Floor(%d * %f) = %d\n", damageCtx.DamageHP, casterAttr.AttackBloodsuckerRate, damageCtx.BloodsuckerValue)
	})

	if damageCtx.BloodsuckerValue <= 0 {
		damageCtx.PushDebugInfo(func() string {
			return fmt.Sprintf("\t\t吸血值为0，无需执行\n")
		})

		// 发送事件
		casterEvents.EmitAfterDamageStep(damageCtx, DamageStep_Bloodsucker)

		return true
	}

	damageCtx.PushDebugInfo(func() string {
		return fmt.Sprintf("\t\t执行吸血：\n") +
			fmt.Sprintf("\t\t\t\t吸血前HP：%d\n", caster.Attr.CurHP) +
			fmt.Sprintf("\t\t\t\t死亡状态：%t\n", !caster.IsAlive())
	})

	// 恢复HP
	damageCtx.BloodsuckerRecoverHP, _ = caster.Attr.BeDamageChangeHP(damageCtx.Attack, DamageKind_Bloodsucking, BitsStick(int32(Proto.DamageType_RecoverHP)),
		damageCtx.BloodsuckerValue, 0, false)

	damageCtx.PushDebugInfo(func() string {
		return fmt.Sprintf("\t\t\t\t实际恢复HP：%d\n", damageCtx.BloodsuckerRecoverHP) +
			fmt.Sprintf("\t\t\t\t吸血后HP：%d\n", caster.Attr.CurHP) +
			fmt.Sprintf("\t\t\t\t死亡状态：%t\n", !caster.IsAlive())
	})

	// 发送事件
	casterEvents.EmitAfterDamageStep(damageCtx, DamageStep_Bloodsucker)

	return true
}
