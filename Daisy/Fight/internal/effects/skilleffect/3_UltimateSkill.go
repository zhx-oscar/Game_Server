package skilleffect

import (
	. "Daisy/Fight/internal"
	"Daisy/Fight/internal/conf"
	"Daisy/Fight/internal/effects"
	"Daisy/Proto"
)

// _3_UltimateSkill 必杀技模板
type _3_UltimateSkill struct {
	effects.Blank           // 继承回调
	effects.Damage          // 伤害模板
	effects.Hit             // 打击模板
	isAddInvincible bool    // 是否已设置无敌
	isAllAIPause    bool    // 是否已暂停所有AI
	isTargetAIPause bool    // 是否已暂停目标AI
	pauseAITargets  []*Pawn // 暂停AI的目标列表
}

// Init 初始化
func (effect *_3_UltimateSkill) Init(skill *Skill) error {
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
func (effect *_3_UltimateSkill) OnSkillInReady(skill *Skill) {
	effect.isAddInvincible = false
	effect.isAllAIPause = false
	effect.isTargetAIPause = false
	effect.pauseAITargets = nil

	// 抢占使用必杀技
	if !skill.Caster.Scene.LockCastUltimateSkill(skill.Caster) {
		// 抢占失败终止必杀技
		skill.Caster.BreakCurSkill(skill.Caster, Proto.SkillBreakReason_Normal)
		return
	}
}

// OnSkillReadyFinish 技能准备释放（在进入技能准备阶段之后可以释放时，当前技能与所有buff能收到）
func (effect *_3_UltimateSkill) OnSkillReadyFinish(skill *Skill) {
	// 重置必杀技能量
	skill.Caster.Attr.ChangeUltimateSkillPower(0)

	// 添加无敌
	if !effect.isAddInvincible {
		skill.Caster.State.ChangeStat(Stat_Invincible, true)
		effect.isAddInvincible = true
	}

	// 暂停所有AI
	if !effect.isAllAIPause {
		skill.Caster.Scene.AllAIPauseSkipOne(skill.Caster, true)
		effect.isAllAIPause = true
	}

	skill.Caster.AddBuff(skill.Caster, skill.Config.OwnBuff, 0)
}

// OnSkillInStart 技能开始（当前技能与所有buff能收到）
func (effect *_3_UltimateSkill) OnSkillInStart(skill *Skill) {
	// 恢复所有AI
	if effect.isAllAIPause {
		skill.Caster.Scene.AllAIPauseSkipOne(skill.Caster, false)
		effect.isAllAIPause = false
	}

	// 暂停目标AI
	if !effect.isTargetAIPause {
		if skill.GetAttackType() == conf.AttackType_Aoe {
			effect.pauseAITargets = skill.SearchTargets(skill.TargetPos)
		} else {
			effect.pauseAITargets = append([]*Pawn{}, skill.TargetList...)
		}

		for _, target := range effect.pauseAITargets {
			if target.Equal(skill.Caster) {
				continue
			}
			target.AIPause(true)
		}

		effect.isTargetAIPause = true
	}

	// 合体必杀技亮灯
	skill.Caster.Scene.TurnOnCombineSkillPoint(skill.Caster)

	// 合体必杀技开始计时
	skill.Caster.Scene.StartCombineSkillCountdown()
}

// OnAttackInit 伤害体创建后（创建伤害体的技能或buff能收到）
func (effect *_3_UltimateSkill) OnAttackInit(attack *Attack) {
	for _, target := range effect.pauseAITargets {
		if target.Equal(attack.Skill.Caster) {
			continue
		}
		target.AIPause(true)
	}
}

// OnAttackDestroy 伤害体销毁后（创建伤害体的技能或buff能收到）
func (effect *_3_UltimateSkill) OnAttackDestroy(attack *Attack, isBreak bool, breakCaster *Pawn) {
	for _, target := range effect.pauseAITargets {
		if target.Equal(attack.Skill.Caster) {
			continue
		}
		target.AIPause(false)
	}
}

// OnAttackHitTarget 伤害体打击目标时（创建伤害体的技能或buff能收到）
func (effect *_3_UltimateSkill) OnAttackHitTarget(attack *Attack, target *Pawn, hitTimes uint32) {
	// 伤害目标
	damageBit, _, _, _ := effect.DamageTarget(attack, target, 0)

	// 命中增加必杀技能量
	if !damageBit.Any(BitsStick(int32(Proto.DamageType_Miss), int32(Proto.DamageType_Dodge), int32(Proto.DamageType_ExemptionDamage))) {
		// 添加命中buff
		target.AddBuff(attack.Skill.Caster, attack.Skill.Config.TargetBuff, 0)
	}

	// 打击目标
	effect.Hit.Hit(attack, target, damageBit)
}

// OnSkillInLater 进入技能后摇阶段（当前技能与所有buff能收到）
func (effect *_3_UltimateSkill) OnSkillInLater(skill *Skill) {
	// 解锁必杀技
	skill.Caster.Scene.UnlockCastUltimateSkill(skill.Caster)

	// 合体必杀技准备释放
	skill.Caster.Scene.CombineSkillReadyCast()
}

// OnSkillInEnd 技能结束（包含被打断和正常结束，当前技能与所有buff能收到）
func (effect *_3_UltimateSkill) OnSkillInEnd(skill *Skill, lastStat Proto.SkillState_Enum, skillEndReason Proto.SkillEndReason_Enum, breakCaster *Pawn) {
	// 重置公共cd结束时间
	if lastStat != Proto.SkillState_Ready {
		skill.Caster.ResetPublicCDTime()
	}

	// 解锁必杀技
	skill.Caster.Scene.UnlockCastUltimateSkill(skill.Caster)

	// 准备使用合体必杀技
	if lastStat != Proto.SkillState_Ready {
		skill.Caster.Scene.CombineSkillReadyCast()
	}

	// 解除无敌
	if effect.isAddInvincible {
		skill.Caster.State.ChangeStat(Stat_Invincible, false)
		effect.isAddInvincible = false
	}

	// 恢复所有AI
	if effect.isAllAIPause {
		skill.Caster.Scene.AllAIPauseSkipOne(skill.Caster, false)
		effect.isAllAIPause = false
	}

	// 恢复目标AI
	if effect.isTargetAIPause {
		for _, target := range effect.pauseAITargets {
			target.AIPause(false)
		}
		effect.isTargetAIPause = false
	}
}
