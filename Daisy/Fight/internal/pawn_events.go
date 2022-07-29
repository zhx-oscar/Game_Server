package internal

import (
	"Daisy/Fight/internal/conf"
	"Daisy/Proto"
)

// 事件递归深度
const eventDeepLimit = uint32(32)

// _EventSolt 事件槽
type _EventSolt struct {
	UID     uint32
	Disable bool
	Effects []IEffectCallback
}

// FightEvents 事件源
type FightEvents struct {
	pawn         *Pawn
	eventSoltTab []*_EventSolt
}

// init 初始化
func (e *FightEvents) init(pawn *Pawn) {
	e.pawn = pawn
}

// checkDeepLimit 检测递归深度限制
func (e *FightEvents) checkDeepLimit() (bool, func()) {
	if e.pawn == nil {
		return false, nil
	}

	scene := e.pawn.Scene

	if scene.eventDeep > eventDeepLimit {
		return false, nil
	}

	scene.eventDeep++

	return true, func() {
		scene.eventDeep--
	}
}

// HookEvent 绑定事件
func (e *FightEvents) HookEvent(uid uint32, effects []IEffectCallback) {
	for _, solt := range e.eventSoltTab {
		if solt.UID == uid {
			return
		}
	}

	e.eventSoltTab = append(e.eventSoltTab, &_EventSolt{
		UID:     uid,
		Effects: effects,
	})
}

// UnhookEvent 解绑事件
func (e *FightEvents) UnhookEvent(uid uint32) {
	for i, solt := range e.eventSoltTab {
		if solt.UID == uid {
			solt.Disable = true
			e.eventSoltTab = append(e.eventSoltTab[:i], e.eventSoltTab[i+1:]...)
			return
		}
	}
}

// sendEventToEffectTab 向指定效果表发送事件
func (e *FightEvents) sendEventToEffectTab(fun func(effect IEffectCallback) bool, effectTab []IEffectCallback) {
	for _, effect := range effectTab {
		if effect != nil {
			if !fun(effect) {
				break
			}
		}
	}
}

// sendEventToSoltTab 向事件槽发送事件
func (e *FightEvents) sendEventToSoltTab(fun func(effect IEffectCallback) bool) {
	for _, solt := range e.eventSoltTab {
		if solt.Disable {
			continue
		}

		for _, effect := range solt.Effects {
			if solt.Disable {
				break
			}

			if effect != nil {
				if !fun(effect) {
					break
				}
			}
		}
	}
}

// sendEventToSoltTabEx 向事件槽发送事件
func (e *FightEvents) sendEventToSoltTabEx(fun func(soltUID uint32, effect IEffectCallback) bool) {
	for _, solt := range e.eventSoltTab {
		if solt.Disable {
			continue
		}

		for _, effect := range solt.Effects {
			if solt.Disable {
				break
			}

			if !fun(solt.UID, effect) {
				break
			}
		}
	}
}

// EmitSkillInReady 发送技能进入准备阶段事件
func (e *FightEvents) EmitSkillInReady(skill *Skill) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	fun := func(effect IEffectCallback) bool {
		if skill.Stat == Proto.SkillState_Wait {
			return false
		}

		effect.OnSkillInReady(skill)

		return true
	}

	e.sendEventToSoltTab(fun)
	e.sendEventToEffectTab(fun, skill.effectTab)
}

// EmitSkillReadyFinish 发送技能准备释放事件
func (e *FightEvents) EmitSkillReadyFinish(skill *Skill) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	fun := func(effect IEffectCallback) bool {
		if skill.Stat == Proto.SkillState_Wait {
			return false
		}

		effect.OnSkillReadyFinish(skill)

		return true
	}

	e.sendEventToSoltTab(fun)
	e.sendEventToEffectTab(fun, skill.effectTab)
}

// EmitSkillInShowTime 发送技能进入特写阶段事件
func (e *FightEvents) EmitSkillInShowTime(skill *Skill) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	fun := func(effect IEffectCallback) bool {
		if skill.Stat == Proto.SkillState_Wait {
			return false
		}

		effect.OnSkillInShowTime(skill)

		return true
	}

	e.sendEventToSoltTab(fun)
	e.sendEventToEffectTab(fun, skill.effectTab)
}

// EmitSkillInDashing 发送技能进入冲刺阶段事件
func (e *FightEvents) EmitSkillInDashing(skill *Skill) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	fun := func(effect IEffectCallback) bool {
		if skill.Stat == Proto.SkillState_Wait {
			return false
		}

		effect.OnSkillInDashing(skill)

		return true
	}

	e.sendEventToSoltTab(fun)
	e.sendEventToEffectTab(fun, skill.effectTab)
}

// EmitSkillDashingFinish 发送技能冲刺结束事件
func (e *FightEvents) EmitSkillDashingFinish(skill *Skill) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	fun := func(effect IEffectCallback) bool {
		if skill.Stat == Proto.SkillState_Wait {
			return false
		}

		effect.OnSkillDashingFinish(skill)

		return true
	}

	e.sendEventToSoltTab(fun)
	e.sendEventToEffectTab(fun, skill.effectTab)
}

// EmitSkillInStart 发送技能开始事件
func (e *FightEvents) EmitSkillInStart(skill *Skill) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	fun := func(effect IEffectCallback) bool {
		if skill.Stat == Proto.SkillState_Wait {
			return false
		}

		effect.OnSkillInStart(skill)

		return true
	}

	e.sendEventToSoltTab(fun)
	e.sendEventToEffectTab(fun, skill.effectTab)
}

// EmitSkillInBefore 发送技能进入前摇阶段事件
func (e *FightEvents) EmitSkillInBefore(skill *Skill) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	fun := func(effect IEffectCallback) bool {
		if skill.Stat == Proto.SkillState_Wait {
			return false
		}

		effect.OnSkillInBefore(skill)

		return true
	}

	e.sendEventToSoltTab(fun)
	e.sendEventToEffectTab(fun, skill.effectTab)
}

// EmitSkillInAttack 发送技能进入攻击阶段事件
func (e *FightEvents) EmitSkillInAttack(skill *Skill) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	fun := func(effect IEffectCallback) bool {
		if skill.Stat == Proto.SkillState_Wait {
			return false
		}

		effect.OnSkillInAttak(skill)

		return true
	}

	e.sendEventToSoltTab(fun)
	e.sendEventToEffectTab(fun, skill.effectTab)
}

// EmitSkillInLater 发送技能进入后摇阶段事件
func (e *FightEvents) EmitSkillInLater(skill *Skill) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	fun := func(effect IEffectCallback) bool {
		if skill.Stat == Proto.SkillState_Wait {
			return false
		}

		effect.OnSkillInLater(skill)

		return true
	}

	e.sendEventToSoltTab(fun)
	e.sendEventToEffectTab(fun, skill.effectTab)
}

// EmitSkillInEnd 发送技能结束事件
func (e *FightEvents) EmitSkillInEnd(skill *Skill, lastStat Proto.SkillState_Enum, skillEndReason Proto.SkillEndReason_Enum, breakCaster *Pawn) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	fun := func(effect IEffectCallback) bool {
		if skill.Stat != Proto.SkillState_Wait {
			return false
		}

		effect.OnSkillInEnd(skill, lastStat, skillEndReason, breakCaster)

		return true
	}

	e.sendEventToSoltTab(fun)
	e.sendEventToEffectTab(fun, skill.effectTab)
}

// EmitBuffAdd 发送施加buff事件
func (e *FightEvents) EmitBuffAdd(buff *Buff) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	e.sendEventToEffectTab(func(effect IEffectCallback) bool {
		if buff.IsDestroy {
			return false
		}

		effect.OnBuffAdd(buff)

		return true

	}, buff.effectTab)

	e.sendEventToSoltTabEx(func(soltUID uint32, effect IEffectCallback) bool {
		if buff.IsDestroy || soltUID == buff.BuffKey.UID {
			return false
		}

		effect.OnOtherBuffAdd(buff)

		return true
	})
}

// EmitBuffRemove 发送删除buff事件
func (e *FightEvents) EmitBuffRemove(buff *Buff, clear bool) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	e.sendEventToEffectTab(func(effect IEffectCallback) bool {
		if !buff.IsDestroy {
			return false
		}

		effect.OnBuffRemove(buff, clear)

		return true

	}, buff.effectTab)

	e.sendEventToSoltTabEx(func(soltUID uint32, effect IEffectCallback) bool {
		if !buff.IsDestroy || soltUID == buff.BuffKey.UID {
			return false
		}

		effect.OnOtherBuffRemove(buff, clear)

		return true
	})
}

// EmitBuffChangeDuration 发送修改buff时长事件
func (e *FightEvents) EmitBuffChangeDuration(buff *Buff, delta int32) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	e.sendEventToEffectTab(func(effect IEffectCallback) bool {
		if buff.IsDestroy {
			return false
		}

		effect.OnBuffChangeDuration(buff, delta)

		return true

	}, buff.effectTab)

	e.sendEventToSoltTabEx(func(soltUID uint32, effect IEffectCallback) bool {
		if buff.IsDestroy || soltUID == buff.BuffKey.UID {
			return false
		}

		effect.OnOtherBuffChangeDuration(buff, delta)

		return true
	})
}

// EmitBuffUpdate 发送帧更新buff事件
func (e *FightEvents) EmitBuffUpdate(buff *Buff) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	e.sendEventToEffectTab(func(effect IEffectCallback) bool {
		if buff.IsDestroy {
			return false
		}

		effect.OnBuffUpdate(buff)

		return true

	}, buff.effectTab)
}

// EmitAttackInit 发送伤害体创建事件
func (e *FightEvents) EmitAttackInit(attack *Attack) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	fun := func(effect IEffectCallback) bool {
		if attack.IsDestroy {
			return false
		}

		effect.OnAttackInit(attack)

		return true
	}

	switch attack.Src() {
	case Proto.AttackSrc_Skill:
		e.sendEventToEffectTab(fun, attack.Skill.effectTab)
	case Proto.AttackSrc_Buff:
		e.sendEventToEffectTab(fun, attack.Buff.effectTab)
	case Proto.AttackSrc_Custom:
		e.sendEventToEffectTab(fun, []IEffectCallback{attack.effectCallback})
	}
}

// EmitAttackBeforeHitAll 伤害体打击所有目标前事件
func (e *FightEvents) EmitAttackBeforeHitAll(attack *Attack) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	fun := func(effect IEffectCallback) bool {
		if attack.IsDestroy {
			return false
		}

		effect.OnAttackBeforeHitAll(attack, attack.HitTimes)

		return true
	}

	switch attack.Src() {
	case Proto.AttackSrc_Skill:
		e.sendEventToEffectTab(fun, attack.Skill.effectTab)
	case Proto.AttackSrc_Buff:
		e.sendEventToEffectTab(fun, attack.Buff.effectTab)
	case Proto.AttackSrc_Custom:
		e.sendEventToEffectTab(fun, []IEffectCallback{attack.effectCallback})
	}
}

// EmitAttackHitTarget 发送伤害体打击目标时事件
func (e *FightEvents) EmitAttackHitTarget(attack *Attack, target *Pawn) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	fun := func(effect IEffectCallback) bool {
		if attack.IsDestroy {
			return false
		}

		effect.OnAttackHitTarget(attack, target, attack.HitTimes)

		return true
	}

	switch attack.Src() {
	case Proto.AttackSrc_Skill:
		e.sendEventToEffectTab(fun, attack.Skill.effectTab)
	case Proto.AttackSrc_Buff:
		e.sendEventToEffectTab(fun, attack.Buff.effectTab)
	case Proto.AttackSrc_Custom:
		e.sendEventToEffectTab(fun, []IEffectCallback{attack.effectCallback})
	}
}

// EmitAttackAfterHitAll 发送伤害体打击所有目标后事件
func (e *FightEvents) EmitAttackAfterHitAll(attack *Attack) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	fun := func(effect IEffectCallback) bool {
		if attack.IsDestroy {
			return false
		}

		effect.OnAttackAfterHitAll(attack, attack.HitTimes)

		return true
	}

	switch attack.Src() {
	case Proto.AttackSrc_Skill:
		e.sendEventToEffectTab(fun, attack.Skill.effectTab)
	case Proto.AttackSrc_Buff:
		e.sendEventToEffectTab(fun, attack.Buff.effectTab)
	case Proto.AttackSrc_Custom:
		e.sendEventToEffectTab(fun, []IEffectCallback{attack.effectCallback})
	}
}

// EmitAttackDestroy 发送伤害体伤害体销毁事件
func (e *FightEvents) EmitAttackDestroy(attack *Attack, isBreak bool, breakCaster *Pawn) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	fun := func(effect IEffectCallback) bool {
		if !attack.IsDestroy {
			return false
		}

		effect.OnAttackDestroy(attack, isBreak, breakCaster)

		return true
	}

	switch attack.Src() {
	case Proto.AttackSrc_Skill:
		e.sendEventToEffectTab(fun, attack.Skill.effectTab)
	case Proto.AttackSrc_Buff:
		e.sendEventToEffectTab(fun, attack.Buff.effectTab)
	case Proto.AttackSrc_Custom:
		e.sendEventToEffectTab(fun, []IEffectCallback{attack.effectCallback})
	}
}

// EmitDamageTarget 发送伤害目标事件
func (e *FightEvents) EmitDamageTarget(attack *Attack, target *Pawn, damageKind conf.DamageKind, damageBit Bits, damageValue, damageHP, damageHPShield int64) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	fun := func(effect IEffectCallback) bool {
		if attack.IsDestroy {
			return false
		}

		effect.OnDamageTarget(attack, target, damageKind, damageBit, damageValue, damageHP, damageHPShield)

		return true
	}

	switch attack.Src() {
	case Proto.AttackSrc_Skill:
		e.sendEventToEffectTab(fun, attack.Skill.effectTab)
	case Proto.AttackSrc_Buff:
		e.sendEventToEffectTab(fun, attack.Buff.effectTab)
	case Proto.AttackSrc_Custom:
		e.sendEventToEffectTab(fun, []IEffectCallback{attack.effectCallback})
	}

	e.sendEventToSoltTab(fun)
}

// EmitHitTarget 发送打击目标事件
func (e *FightEvents) EmitHitTarget(attack *Attack, target *Pawn, damagBit Bits, hitType Proto.HitType_Enum) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	fun := func(effect IEffectCallback) bool {
		if attack.IsDestroy {
			return false
		}

		effect.OnHitTarget(attack, target, damagBit, hitType)

		return true
	}

	switch attack.Src() {
	case Proto.AttackSrc_Skill:
		e.sendEventToEffectTab(fun, attack.Skill.effectTab)
	case Proto.AttackSrc_Buff:
		e.sendEventToEffectTab(fun, attack.Buff.effectTab)
	case Proto.AttackSrc_Custom:
		e.sendEventToEffectTab(fun, []IEffectCallback{attack.effectCallback})
	}

	e.sendEventToSoltTab(fun)
}

// EmitKillTarget 发送击杀目标成功事件
func (e *FightEvents) EmitKillTarget(attack *Attack, target *Pawn) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	fun := func(effect IEffectCallback) bool {
		if attack.IsDestroy {
			return false
		}

		effect.OnKillTarget(attack, target)

		return true
	}

	switch attack.Src() {
	case Proto.AttackSrc_Skill:
		e.sendEventToEffectTab(fun, attack.Skill.effectTab)
	case Proto.AttackSrc_Buff:
		e.sendEventToEffectTab(fun, attack.Buff.effectTab)
	case Proto.AttackSrc_Custom:
		e.sendEventToEffectTab(fun, []IEffectCallback{attack.effectCallback})
	}

	e.sendEventToSoltTab(fun)
}

// EmitBreakTargetSkill 发送打断目标施法事件
func (e *FightEvents) EmitBreakTargetSkill(target *Pawn, skill *Skill) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	e.sendEventToSoltTab(func(effect IEffectCallback) bool {
		effect.OnBreakTargetSkill(target, skill)
		return true
	})
}

// EmitBeDamage 发送受到伤害事件
func (e *FightEvents) EmitBeDamage(attack *Attack, damageKind conf.DamageKind, damageBit Bits, damageValue, damageHP, damageHPShield int64) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	e.sendEventToSoltTab(func(effect IEffectCallback) bool {
		if attack.IsDestroy {
			return false
		}

		effect.OnBeDamage(attack, damageKind, damageBit, damageValue, damageHP, damageHPShield)

		return true
	})
}

// EmitBeHit 发送受到打击事件
func (e *FightEvents) EmitBeHit(attack *Attack, damagBit Bits, hitType Proto.HitType_Enum) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	e.sendEventToSoltTab(func(effect IEffectCallback) bool {
		if attack.IsDestroy {
			return false
		}

		effect.OnBeHit(attack, damagBit, hitType)

		return true
	})
}

// EmitDead 发送死亡事件
func (e *FightEvents) EmitDead(attack *Attack) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	e.sendEventToSoltTab(func(effect IEffectCallback) bool {
		if attack != nil && attack.IsDestroy {
			return false
		}

		effect.OnDead(attack)

		return true
	})
}

// EmitBeBreakSkill 发送受到打断施法事件
func (e *FightEvents) EmitBeBreakSkill(caster *Pawn, skill *Skill) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	e.sendEventToSoltTab(func(effect IEffectCallback) bool {
		effect.OnBeBreakSkill(caster, skill)
		return true
	})
}

// EmitShieldBroken 护盾破裂
func (e *FightEvents) EmitShieldBroken(attack *Attack, shield *HPShield) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	if shield.effect == nil {
		return
	}

	e.sendEventToEffectTab(func(effect IEffectCallback) bool {
		if attack.IsDestroy {
			return false
		}

		effect.OnShieldBroken(attack, shield)

		return true

	}, []IEffectCallback{shield.effect})
}

// EmitDecHPShield 扣除护盾血量
func (e *FightEvents) EmitDecHPShield(attack *Attack, shield *HPShield, oldShieldHP int64) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	if shield.effect == nil {
		return
	}

	e.sendEventToEffectTab(func(effect IEffectCallback) bool {
		if attack.IsDestroy {
			return false
		}

		effect.OnDecHPShield(attack, shield, oldShieldHP)

		return true

	}, []IEffectCallback{shield.effect})
}

// EmitBeforeDamageStep 进入步骤（所有buff能收到）
func (e *FightEvents) EmitBeforeDamageStep(damageCtx *DamageContext, step conf.DamageStep) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	e.sendEventToSoltTab(func(effect IEffectCallback) bool {
		if damageCtx.Attack.IsDestroy {
			return false
		}

		effect.OnBeforeDamageStep(damageCtx, step)

		return true
	})

	damageCtx.Target.Events.sendEventToSoltTab(func(effect IEffectCallback) bool {
		if damageCtx.Attack.IsDestroy {
			return false
		}

		effect.OnBeBeforeDamageStep(damageCtx, step)

		return true
	})
}

// EmitAfterDamageStep 离开步骤（所有buff能收到）
func (e *FightEvents) EmitAfterDamageStep(damageCtx *DamageContext, step conf.DamageStep) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	e.sendEventToSoltTab(func(effect IEffectCallback) bool {
		if damageCtx.Attack.IsDestroy {
			return false
		}

		effect.OnAfterDamageStep(damageCtx, step)

		return true
	})

	damageCtx.Target.Events.sendEventToSoltTab(func(effect IEffectCallback) bool {
		if damageCtx.Attack.IsDestroy {
			return false
		}

		effect.OnBeAfterDamageStep(damageCtx, step)

		return true
	})
}

// EmitStepWaterfallJudgePass 在 步骤：瀑布判定 中通过指定阶段（所有buff能收到）
func (e *FightEvents) EmitStepWaterfallJudgePass(damageCtx *DamageContext, damageJudgeRv conf.DamageJudgeRv) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	e.sendEventToSoltTab(func(effect IEffectCallback) bool {
		if damageCtx.Attack.IsDestroy {
			return false
		}

		effect.OnStepWaterfallJudgePassJudge(damageCtx, damageJudgeRv)

		return true
	})

	damageCtx.Target.Events.sendEventToSoltTab(func(effect IEffectCallback) bool {
		if damageCtx.Attack.IsDestroy {
			return false
		}

		effect.OnBeStepWaterfallJudgePassJudge(damageCtx, damageJudgeRv)

		return true
	})
}

// EmitOverDrive 发送超载事件
func (e *FightEvents) EmitOverDrive() {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	e.sendEventToSoltTab(func(effect IEffectCallback) bool {
		effect.OnOverDrive()
		return true
	})
}

// EmitHaloAddMember 发送光环增加成员事件
func (e *FightEvents) EmitHaloAddMember(halo *Halo, pawn *Pawn) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	if halo.casterEffect == nil {
		return
	}

	halo.casterEffect.OnAddHaloMember(halo, pawn)
}

// EmitHaloRemoveMember 发送光环移除成员事件
func (e *FightEvents) EmitHaloRemoveMember(halo *Halo, pawn *Pawn) {
	if ok, fun := e.checkDeepLimit(); ok {
		defer fun()
	} else {
		return
	}

	if halo.casterEffect == nil {
		return
	}

	halo.casterEffect.OnRemoveHaloMember(halo, pawn)
}
