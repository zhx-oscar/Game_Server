package skilleffect

import (
	. "Daisy/Fight/internal"
	"Daisy/Fight/internal/conf"
	"Daisy/Fight/internal/effects"
	"Daisy/Proto"
)

// _2_NormalSkill 普攻模板
type _2_NormalSkill struct {
	effects.Blank  // 继承回调
	effects.Damage // 伤害模板
	effects.Hit    // 打击模板
}

// Init 初始化
func (effect *_2_NormalSkill) Init(skill *Skill) error {
	// 初始化伤害模板
	effect.Damage.DamageKind = conf.DamageKind_Hurt
	effect.Damage.DamageFlow = conf.DamageFlow_Attack

	return nil
}

// OnSkillReadyFinish 技能准备结束（此时已经记录使用技能帧，当前技能与所有buff能收到）
func (effect *_2_NormalSkill) OnSkillReadyFinish(skill *Skill) {
	skill.Caster.AddBuff(skill.Caster, skill.Config.OwnBuff, 0)
}

// OnAttackHitTarget 伤害体打击目标时（创建伤害体的技能或buff能收到）
func (effect *_2_NormalSkill) OnAttackHitTarget(attack *Attack, target *Pawn, hitTimes uint32) {
	// 伤害目标
	damageBit, _, _, _ := effect.DamageTarget(attack, target, 0)

	// 伤害命中
	if !damageBit.Any(BitsStick(int32(Proto.DamageType_Miss), int32(Proto.DamageType_Dodge), int32(Proto.DamageType_ExemptionDamage))) {
		// 增加必杀技能量
		attack.Caster.Attr.ChangeUltimateSkillPower(attack.Caster.Attr.UltimateSkillPower + attack.Skill.Config.HitAddUltimateSkillPower)

		// 添加命中buff
		target.AddBuff(attack.Skill.Caster, attack.Skill.Config.TargetBuff, 0)
	}

	// 打击目标
	effect.Hit.Hit(attack, target, damageBit)

}

// OnSkillInLater 进入技能后摇阶段（当前技能与所有buff能收到）
func (effect *_2_NormalSkill) OnSkillInLater(skill *Skill) {
	skill.Caster.ContinueComboAttack()
}
