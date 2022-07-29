package skilleffect

import (
	. "Daisy/Fight/internal"
	"Daisy/Fight/internal/conf"
	"Daisy/Fight/internal/effects"
	"Daisy/Proto"
)

// _4_CombineSkill 合体技模板
type _4_CombineSkill struct {
	effects.Blank        // 继承回调
	effects.Damage       // 伤害模板
	effects.Hit          // 打击模板
	isAddInvincible bool // 是否已设置无敌
}

// Init 初始化
func (effect *_4_CombineSkill) Init(skill *Skill) error {
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

// OnSkillInReady 进入技能准备阶段（当前技能与所有buff能收到）
func (effect *_4_CombineSkill) OnSkillInReady(skill *Skill) {
	effect.isAddInvincible = false

	// 抢占使用合体技
	if !skill.Caster.Scene.LockCastUltimateSkill(skill.Caster) {
		// 抢占失败终止合体技
		skill.Caster.BreakCurSkill(skill.Caster, Proto.SkillBreakReason_Normal)
		return
	}

	// 记录合体技准备释放的成员列表
	skill.SaveCombineSkillReadyMembers(skill.Caster.Scene.GetCombineSkillReadyMembers())
}

// OnSkillReadyFinish 技能准备释放（在进入技能准备阶段之后可以释放时，当前技能与所有buff能收到）
func (effect *_4_CombineSkill) OnSkillReadyFinish(skill *Skill) {
	// 重置合体技能量
	skill.Caster.Attr.ChangeUltimateSkillPower(0)

	// 添加无敌
	if !effect.isAddInvincible {
		skill.Caster.State.ChangeStat(Stat_Invincible, true)
		effect.isAddInvincible = true
	}
}

// OnSkillInStart 技能开始（当前技能与所有buff能收到）
func (effect *_4_CombineSkill) OnSkillInStart(skill *Skill) {
	// 合体合体技亮灯
	skill.Caster.Scene.TurnOnCombineSkillPoint(skill.Caster)

	// 合体合体技开始计时
	skill.Caster.Scene.StartCombineSkillCountdown()
}

// OnAttackHitTarget 伤害体打击目标时（创建伤害体的技能或buff能收到）
func (effect *_4_CombineSkill) OnAttackHitTarget(attack *Attack, target *Pawn, hitTimes uint32) {
	// 伤害目标
	damageBit, _, _, _ := effect.DamageTarget(attack, target, 0)

	// 打击目标
	effect.Hit.Hit(attack, target, damageBit)
}

// OnSkillInLater 进入技能后摇阶段（当前技能与所有buff能收到）
func (effect *_4_CombineSkill) OnSkillInLater(skill *Skill) {
	// 解锁合体技
	skill.Caster.Scene.UnlockCastUltimateSkill(skill.Caster)

	// 合体合体技准备释放
	skill.Caster.Scene.CombineSkillReadyCast()
}

// OnSkillInEnd 技能结束（包含被打断和正常结束，当前技能与所有buff能收到）
func (effect *_4_CombineSkill) OnSkillInEnd(skill *Skill, lastStat Proto.SkillState_Enum, skillEndReason Proto.SkillEndReason_Enum, breakCaster *Pawn) {
	// 解锁合体技
	skill.Caster.Scene.UnlockCastUltimateSkill(skill.Caster)

	// 准备使用合体合体技
	if lastStat != Proto.SkillState_Ready {
		skill.Caster.Scene.CombineSkillReadyCast()
	}

	// 解除无敌
	if effect.isAddInvincible {
		skill.Caster.State.ChangeStat(Stat_Invincible, false)
		effect.isAddInvincible = false
	}
}
