package skilleffect

import (
	. "Daisy/Fight/internal"
	"Daisy/Fight/internal/conf"
	"Daisy/Fight/internal/effects"
	"Daisy/Proto"
)

// _1_SuperSkill 超能技模板
type _1_SuperSkill struct {
	effects.Blank  // 继承回调
	effects.Damage // 伤害模板
	effects.Hit    // 打击模板
}

// Init 初始化
func (effect *_1_SuperSkill) Init(skill *Skill) error {
	// 初始化伤害模板
	effect.Damage.DamageKind = conf.DamageKind_Hurt
	effect.Damage.DamageFlow = conf.DamageFlow_CastSkill
	effect.Damage.DamageValueTab[conf.DamageValueKind_Normal] = skill.Config.NormalAttack
	effect.Damage.DamageValueTab[conf.DamageValueKind_Fire] = skill.Config.FireAttack
	effect.Damage.DamageValueTab[conf.DamageValueKind_Cold] = skill.Config.ColdAttack
	effect.Damage.DamageValueTab[conf.DamageValueKind_Poison] = skill.Config.PoisonAttack
	effect.Damage.DamageValueTab[conf.DamageValueKind_Lightning] = skill.Config.LightningAttack

	return nil
}

// OnSkillReadyFinish 技能准备结束（此时已经记录使用技能帧，当前技能与所有buff能收到）
func (effect *_1_SuperSkill) OnSkillReadyFinish(skill *Skill) {
	skill.Caster.AddBuff(skill.Caster, skill.Config.OwnBuff, 0)
}

// OnSkillInStart 技能开始（当前技能与所有buff能收到）
func (effect *_1_SuperSkill) OnSkillInStart(skill *Skill) {
	// 增加必杀技能量
	skill.Caster.Attr.ChangeUltimateSkillPower(skill.Caster.Attr.UltimateSkillPower + skill.Config.CastAddUltimateSkillPower)
}

// OnAttackHitTarget 伤害体打击目标时（创建伤害体的技能或buff能收到）
func (effect *_1_SuperSkill) OnAttackHitTarget(attack *Attack, target *Pawn, hitTimes uint32) {
	// 伤害目标
	damageBit, _, _, _ := effect.DamageTarget(attack, target, 0)

	// 伤害命中
	if !damageBit.Any(BitsStick(int32(Proto.DamageType_Miss), int32(Proto.DamageType_Dodge), int32(Proto.DamageType_ExemptionDamage))) {
		// 增加必杀技能量
		attack.Caster.Attr.ChangeUltimateSkillPower(attack.Caster.Attr.UltimateSkillPower + attack.Skill.Config.HitAddUltimateSkillPower)

		// 清除buff组
		for _, clearGroup := range attack.Skill.Config.ClearGroup {
			target.ClearBuff(clearGroup)
		}

		// 添加命中buff
		target.AddBuff(attack.Skill.Caster, attack.Skill.Config.TargetBuff, 0)
	}

	// 打击目标
	effect.Hit.Hit(attack, target, damageBit)

}

// OnSkillInEnd 技能结束（包含被打断和正常结束，当前技能与所有buff能收到）
func (effect *_1_SuperSkill) OnSkillInEnd(skill *Skill, lastStat Proto.SkillState_Enum, skillEndReason Proto.SkillEndReason_Enum, breakCaster *Pawn) {
	// 重置公共cd结束时间
	if lastStat != Proto.SkillState_Ready {
		skill.Caster.ResetPublicCDTime()
	}
}
