package effects

import (
	. "Daisy/Fight/internal"
	"Daisy/Proto"
)

// execAttackFlow 执行攻击流程
func (damage *Damage) execAttackFlow(damageCtx *DamageContext) {
	damageCtx.DamageBit.TurnOn(int32(Proto.DamageType_Damage))

	// 攻击幸运判定
	if !damage.stepAttackLuckyJudge(damageCtx) {
		return
	}

	// 统计攻击伤害
	if !damage.stepCountAttackDamageValue(damageCtx) {
		return
	}

	// 攻击增伤修正
	if !damage.stepIncrAttackDamageValue(damageCtx) {
		return
	}

	// 瀑布判定
	if !damage.stepWaterfallJudge(damageCtx) {
		return
	}

	// 抗性修正
	if !damage.stepResistanceDamageValue(damageCtx) {
		return
	}

	// 减伤修正
	if !damage.stepDecrDamageValue(damageCtx) {
		return
	}

	// 分裂伤害
	if !damage.stepSputtering(damageCtx) {
		return
	}

	// 伤害扣血
	if !damage.stepDeductHP(damageCtx) {
		return
	}

	// 吸血
	if !damage.stepBloodsucker(damageCtx) {
		return
	}
}
