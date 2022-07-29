package effects

import (
	. "Daisy/Fight/internal"
	"Daisy/Proto"
)

// execCasterSkillFlow 执行施法流程
func (damage *Damage) execCasterSkillFlow(damageCtx *DamageContext) {
	damageCtx.DamageBit.TurnOn(int32(Proto.DamageType_Damage))

	// 统计施法伤害
	if !damage.stepCountCastSkillDamageValue(damageCtx) {
		return
	}

	// 施法增伤修正
	if !damage.stepIncrCastSkillDamageValue(damageCtx) {
		return
	}

	// 防御扣除护盾
	if !damage.stepDefendDeductHPShield(damageCtx) {
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

	// 伤害扣血
	if !damage.stepDeductHP(damageCtx) {
		return
	}
}
