package internal

import (
	"Daisy/Fight/internal/conf"
	"Daisy/Proto"
)

// ISkillEffect 技能效果
type ISkillEffect interface {
	Init(skill *Skill) error
}

// IBuffEffect buff效果
type IBuffEffect interface {
	Init(buff *Buff) error
}

// IEffectCallback 效果回调
type IEffectCallback interface {
	// 技能回调
	OnSkillInReady(skill *Skill)
	OnSkillReadyFinish(skill *Skill)
	OnSkillInShowTime(skill *Skill)
	OnSkillInDashing(skill *Skill)
	OnSkillDashingFinish(skill *Skill)
	OnSkillInStart(skill *Skill)
	OnSkillInBefore(skill *Skill)
	OnSkillInAttak(skill *Skill)
	OnSkillInLater(skill *Skill)
	OnSkillInEnd(skill *Skill, lastStat Proto.SkillState_Enum, skillEndReason Proto.SkillEndReason_Enum, breakCaster *Pawn)

	// buff回调
	OnBuffAdd(buff *Buff)
	OnBuffRemove(buff *Buff, clear bool)
	OnBuffChangeDuration(buff *Buff, delta int32)
	OnBuffUpdate(buff *Buff)
	OnOtherBuffAdd(other *Buff)
	OnOtherBuffRemove(other *Buff, clear bool)
	OnOtherBuffChangeDuration(other *Buff, delta int32)

	// 伤害体回调
	OnAttackInit(attack *Attack)
	OnAttackBeforeHitAll(attack *Attack, hitTimes uint32)
	OnAttackHitTarget(attack *Attack, target *Pawn, hitTimes uint32)
	OnAttackAfterHitAll(attack *Attack, hitTimes uint32)
	OnAttackDestroy(attack *Attack, isBreak bool, breakCaster *Pawn)

	// 主动攻击流程回调
	OnDamageTarget(attack *Attack, target *Pawn, damageKind conf.DamageKind, damageBit Bits, damageValue, damageHP, damageHPShield int64)
	OnHitTarget(attack *Attack, target *Pawn, damageBit Bits, hitType Proto.HitType_Enum)
	OnKillTarget(attack *Attack, target *Pawn)
	OnBreakTargetSkill(target *Pawn, skill *Skill)

	// 被动攻击流程回调
	OnBeDamage(attack *Attack, damageKind conf.DamageKind, damageBit Bits, damageValue, damageHP, damageHPShield int64)
	OnBeHit(attack *Attack, damageBit Bits, hitType Proto.HitType_Enum)
	OnDead(attack *Attack)
	OnBeBreakSkill(caster *Pawn, skill *Skill)
	OnDecHPShield(attack *Attack, shield *HPShield, oldShieldHP int64)
	OnShieldBroken(attack *Attack, shield *HPShield)

	// 主动伤害运算回调
	OnBeforeDamageStep(damageCtx *DamageContext, step conf.DamageStep)
	OnAfterDamageStep(damageCtx *DamageContext, step conf.DamageStep)
	OnStepWaterfallJudgePassJudge(damageCtx *DamageContext, damageJudgeRv conf.DamageJudgeRv)

	// 被动伤害运算回调
	OnBeBeforeDamageStep(damageCtx *DamageContext, step conf.DamageStep)
	OnBeAfterDamageStep(damageCtx *DamageContext, step conf.DamageStep)
	OnBeStepWaterfallJudgePassJudge(damageCtx *DamageContext, damageJudgeRv conf.DamageJudgeRv)

	// 自定义回调
	OnOverDrive()
	OnAddHaloMember(halo *Halo, pawn *Pawn)
	OnRemoveHaloMember(halo *Halo, pawn *Pawn)
}
