package effects

import (
	. "Daisy/Fight/internal"
	"Daisy/Fight/internal/conf"
	"Daisy/Proto"
)

type Blank struct{}

// OnSkillInReady 进入技能准备阶段（当前技能与所有buff能收到）
func (*Blank) OnSkillInReady(skill *Skill) {}

// OnSkillReadyFinish 技能准备结束（此时已经记录使用技能帧，当前技能与所有buff能收到）
func (*Blank) OnSkillReadyFinish(skill *Skill) {}

// OnSkillInShowTime 进入技能特写阶段（有特写时会触发，当前技能与所有buff能收到）
func (*Blank) OnSkillInShowTime(skill *Skill) {}

// OnSkillInDashing 进入技能冲刺阶段（有冲刺时会触发，当前技能与所有buff能收到）
func (*Blank) OnSkillInDashing(skill *Skill) {}

// OnSkillDashingFinish 技能冲刺结束（冲刺结束时会触发，冲刺中技能被打断不会触发，当前技能与所有buff能收到）
func (*Blank) OnSkillDashingFinish(skill *Skill) {}

// OnSkillInStart 技能开始（进入前摇或攻击阶段前触发，当前技能与所有buff能收到）
func (*Blank) OnSkillInStart(skill *Skill) {}

// OnSkillInBefore 进入技能前摇阶段（当前技能与所有buff能收到）
func (*Blank) OnSkillInBefore(skill *Skill) {}

// OnSkillInAttak 进入技能攻击阶段（当前技能与所有buff能收到）
func (*Blank) OnSkillInAttak(skill *Skill) {}

// OnSkillInLater 进入技能后摇阶段（当前技能与所有buff能收到）
func (*Blank) OnSkillInLater(skill *Skill) {}

// OnSkillInEnd 技能结束（包含被打断和正常结束，当前技能与所有buff能收到）
func (*Blank) OnSkillInEnd(skill *Skill, lastStat Proto.SkillState_Enum, skillEndReason Proto.SkillEndReason_Enum, breakCaster *Pawn) {
}

// OnBuffAdd 施加buff后（buff自身能收到）
func (*Blank) OnBuffAdd(buff *Buff) {}

// OnBuffRemove 移除buff后（buff自身能收到）
func (*Blank) OnBuffRemove(buff *Buff, clear bool) {}

// OnBuffChangeDuration 修改buff时长后（buff自身能收到）
func (*Blank) OnBuffChangeDuration(buff *Buff, delta int32) {}

// OnBuffUpdate buff帧更新（buff自身能收到）
func (*Blank) OnBuffUpdate(buff *Buff) {}

// OnOtherBuffAdd 其他buff施加后（其他buff能收到）
func (*Blank) OnOtherBuffAdd(other *Buff) {}

// OnOtherBuffRemove 其他buff移除后（其他buff能收到）
func (*Blank) OnOtherBuffRemove(other *Buff, clear bool) {}

// OnOtherBuffChangeDuration 其他buff修改buff时长后（其他buff能收到）
func (*Blank) OnOtherBuffChangeDuration(other *Buff, delta int32) {}

// OnAttackInit 伤害体创建后（创建伤害体的技能或buff能收到）
func (*Blank) OnAttackInit(attack *Attack) {}

// OnAttackBeforeHitAll 伤害体打击所有目标前（创建伤害体的技能或buff能收到）
func (*Blank) OnAttackBeforeHitAll(attack *Attack, hitTimes uint32) {}

// OnAttackHitTarget 伤害体打击目标时（创建伤害体的技能或buff能收到）
func (*Blank) OnAttackHitTarget(attack *Attack, target *Pawn, hitTimes uint32) {}

// OnAttackAfterHitAll 伤害体打击所有目标后（创建伤害体的技能或buff能收到）
func (*Blank) OnAttackAfterHitAll(attack *Attack, hitTimes uint32) {}

// OnAttackDestroy 伤害体销毁后（创建伤害体的技能或buff能收到）
func (*Blank) OnAttackDestroy(attack *Attack, isBreak bool, breakCaster *Pawn) {}

// OnDamageTarget 伤害目标后（创建伤害体的技能和所有buff能收到）
func (*Blank) OnDamageTarget(attack *Attack, target *Pawn, damageKind conf.DamageKind, damageBit Bits, damageValue, damageHP, damageHPShield int64) {
}

// OnHitTarget 打击目标后（创建伤害体的技能和所有buff能收到）
func (*Blank) OnHitTarget(attack *Attack, target *Pawn, damageBit Bits, hitType Proto.HitType_Enum) {
}

// OnKillTarget 击杀目标后（创建伤害体的技能和所有buff能收到）
func (*Blank) OnKillTarget(attack *Attack, target *Pawn) {}

// OnBreakTargetSkill 打断目标施法后（所有buff能收到）
func (*Blank) OnBreakTargetSkill(target *Pawn, skill *Skill) {}

// OnBeDamage 受到伤害后（所有buff能收到）
func (*Blank) OnBeDamage(attack *Attack, damageKind conf.DamageKind, damageBit Bits, damageValue, damageHP, damageHPShield int64) {
}

// OnBeHit 受到打击后（所有buff能收到）
func (*Blank) OnBeHit(attack *Attack, damageBit Bits, hitType Proto.HitType_Enum) {}

// OnDead 死亡后（所有buff能收到，attack为nil表示非伤害造成的死亡）
func (*Blank) OnDead(attack *Attack) {}

// OnBeBreakSkill 受到打断施法后（所有buff能收到）
func (*Blank) OnBeBreakSkill(caster *Pawn, skill *Skill) {}

// OnDecHPShield 护盾血量扣除(护盾绑定的effect收到)
func (*Blank) OnDecHPShield(attack *Attack, shield *HPShield, oldHP int64) {}

// OnShieldBroken 护盾破碎后（护盾绑定的effect能收到）
func (*Blank) OnShieldBroken(attack *Attack, shield *HPShield) {}

// OnBeforeDamageStep 进入步骤（所有buff能收到）
func (*Blank) OnBeforeDamageStep(damageCtx *DamageContext, step conf.DamageStep) {}

// OnAfterDamageStep 离开步骤（所有buff能收到）
func (*Blank) OnAfterDamageStep(damageCtx *DamageContext, step conf.DamageStep) {}

// OnStepWaterfallJudgePassJudge 在 步骤：瀑布判定 中通过指定阶段（所有buff能收到）
func (*Blank) OnStepWaterfallJudgePassJudge(damageCtx *DamageContext, damageJudgeRv conf.DamageJudgeRv) {
}

// OnBeBeforeDamageStep 进入步骤（所有buff能收到）
func (*Blank) OnBeBeforeDamageStep(damageCtx *DamageContext, step conf.DamageStep) {}

// OnBeAfterDamageStep 离开步骤（所有buff能收到）
func (*Blank) OnBeAfterDamageStep(damageCtx *DamageContext, step conf.DamageStep) {}

// OnBeStepWaterfallJudgePassJudge 在 步骤：瀑布判定 中通过指定阶段（所有buff能收到）
func (*Blank) OnBeStepWaterfallJudgePassJudge(damageCtx *DamageContext, damageJudgeRv conf.DamageJudgeRv) {
}

// OnOverDrive 进入超载状态(所有buff能收到)
func (*Blank) OnOverDrive() {}

// OnAddHaloMember 增加光环成员(光环绑定的effect能收到)
func (*Blank) OnAddHaloMember(halo *Halo, pawn *Pawn) {}

// OnRemoveHaloMember 移除光环成员(光环绑定的effect能收到)
func (*Blank) OnRemoveHaloMember(halo *Halo, pawn *Pawn) {}
