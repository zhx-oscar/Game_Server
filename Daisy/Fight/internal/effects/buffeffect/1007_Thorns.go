package buffeffect

import (
	. "Daisy/Fight/internal"
	. "Daisy/Fight/internal/conf"
	. "Daisy/Fight/internal/effects"
	"Daisy/Proto"
	"fmt"
)

type _1007_Thorns struct {
	Blank
}

// OnBeStepWaterfallJudgePassJudge 在 步骤：瀑布判定 中通过指定阶段（所有buff能收到）
func (effect *_1007_Thorns) OnBeStepWaterfallJudgePassJudge(damageCtx *DamageContext, damageJudgeRv DamageJudgeRv) {
	// 检测判定阶段
	if damageJudgeRv != DamageJudgeRv_Dodge {
		return
	}

	// 检测伤害类型
	if !damageCtx.DamageBit.Test(int32(Proto.DamageType_Damage)) ||
		damageCtx.DamageBit.Test(int32(Proto.DamageType_Miss)) ||
		damageCtx.DamageBit.Test(int32(Proto.DamageType_Dodge)) ||
		damageCtx.DamageBit.Test(int32(Proto.DamageType_ExemptionDamage)) {
		return
	}

	caster := damageCtx.Attack.Caster
	self := damageCtx.Target

	// 反伤值
	damageValue := int64(self.Attr.BeDamageThorns)
	if damageValue <= 0 {
		return
	}

	damageCtx.PushDebugInfo(func() string {
		return fmt.Sprintf("\t\t执行反伤：\n") +
			fmt.Sprintf("\t\t\t\t目标反伤值：%d\n", damageValue)
	})

	damageCtx.PushDebugInfo(func() string {
		return fmt.Sprintf("\t\t\t\t受到反伤前HP护盾：%d\n", caster.Attr.AllHPShield) +
			fmt.Sprintf("\t\t\t\t死亡状态：%t\n", !caster.IsAlive())
	})

	// 扣除护盾
	damageHPShield, exemptionDamage := caster.Attr.BeDamageChangeHPShield(damageCtx.Attack, BitsStick(int32(Proto.DamageType_Damage)), damageValue)

	damageCtx.PushDebugInfo(func() string {
		return fmt.Sprintf("\t\t\t\t扣除HP护盾：%d\n", damageHPShield) +
			fmt.Sprintf("\t\t\t\t免疫伤害：%t\n", exemptionDamage) +
			fmt.Sprintf("\t\t\t\t受到反伤后HP护盾：%d\n", caster.Attr.AllHPShield) +
			fmt.Sprintf("\t\t\t\t死亡状态：%t\n", !caster.IsAlive())
	})

	// 剩余伤害值
	damageValue -= damageHPShield

	if damageValue > 0 {
		damageCtx.PushDebugInfo(func() string {
			return fmt.Sprintf("\t\t\t\t受到反伤前HP：%d\n", caster.Attr.CurHP) +
				fmt.Sprintf("\t\t\t\t死亡状态：%t\n", !caster.IsAlive())
		})

		// 扣除HP
		damageHP, exemptionDamage := caster.Attr.BeDamageChangeHP(damageCtx.Attack, DamageKind_Thorns, BitsStick(int32(Proto.DamageType_Damage)),
			damageValue, damageHPShield, false)

		damageCtx.PushDebugInfo(func() string {
			return fmt.Sprintf("\t\t\t\t扣除HP：%d\n", damageHP) +
				fmt.Sprintf("\t\t\t\t免疫伤害：%t\n", exemptionDamage) +
				fmt.Sprintf("\t\t\t\t受到反伤后HP：%d\n", caster.Attr.CurHP) +
				fmt.Sprintf("\t\t\t\t死亡状态：%t\n", !caster.IsAlive())
		})

		// 血量为0时死亡
		if caster.Attr.CurHP <= 0 {
			caster.State.Dead(damageCtx.Attack)
		}
	}
}
