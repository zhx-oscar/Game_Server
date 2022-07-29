package buffeffect

import (
	. "Daisy/Fight/internal"
	. "Daisy/Fight/internal/conf"
	. "Daisy/Fight/internal/effects"
	"Daisy/Proto"
	"fmt"
)

type _1008_StealUltimateSkillPower struct {
	Blank
}

// OnStepWaterfallJudgePassJudge 在 步骤：瀑布判定 中通过指定阶段（所有buff能收到）
func (effect *_1008_StealUltimateSkillPower) OnStepWaterfallJudgePassJudge(damageCtx *DamageContext, damageJudgeRv DamageJudgeRv) {
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
	target := damageCtx.Target

	// 检测必杀技能量偷取值
	stealPower := damageCtx.Attack.CasterSnapshot.Attr.AttackStealUltimateSkillPower
	if stealPower <= 0 {
		return
	}

	damageCtx.PushDebugInfo(func() string {
		return fmt.Sprintf("\t\t执行偷取必杀技能量：\n") +
			fmt.Sprintf("\t\t\t\t偷取值：%d\n", stealPower)
	})

	targetOldPower := target.Attr.UltimateSkillPower

	// 目标扣除必杀技能量
	if target.IsAlive() {
		target.Attr.ChangeUltimateSkillPower(target.Attr.UltimateSkillPower - stealPower)
	}

	damageCtx.PushDebugInfo(func() string {
		return fmt.Sprintf("\t\t\t\t目标当前能量：%d\n", targetOldPower) +
			fmt.Sprintf("\t\t\t\t目标准备扣除能量：%d\n", stealPower) +
			fmt.Sprintf("\t\t\t\t目标实际扣除能量：%d\n", targetOldPower-target.Attr.UltimateSkillPower) +
			fmt.Sprintf("\t\t\t\t目标死亡状态：%t\n", !target.IsAlive())
	})

	casterOldPower := caster.Attr.UltimateSkillPower

	// 施法者增加必杀技能量
	if caster.IsAlive() {
		caster.Attr.ChangeUltimateSkillPower(caster.Attr.UltimateSkillPower + stealPower)
	}

	damageCtx.PushDebugInfo(func() string {
		return fmt.Sprintf("\t\t\t\t施法者当前能量：%d\n", casterOldPower) +
			fmt.Sprintf("\t\t\t\t施法者准备增加能量：%d\n", stealPower) +
			fmt.Sprintf("\t\t\t\t施法者实际增加能量：%d\n", caster.Attr.UltimateSkillPower-casterOldPower) +
			fmt.Sprintf("\t\t\t\t施法者死亡状态：%t\n", !target.IsAlive())
	})
}
