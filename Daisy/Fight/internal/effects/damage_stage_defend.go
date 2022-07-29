package effects

import (
	. "Daisy/Fight/internal"
	. "Daisy/Fight/internal/conf"
	"Daisy/Fight/internal/log"
	"Daisy/Proto"
	"fmt"
	"math"
)

// stepDefendDeductHPShield 步骤：防御扣除护盾
func (damage *Damage) stepDefendDeductHPShield(damageCtx *DamageContext) bool {
	damageCtx.PushDebugInfo(func() string {
		damageCtx.StepIndex++
		return fmt.Sprintf("【%d：防御扣除护盾】\n", damageCtx.StepIndex)
	})

	defer func() {
		damageCtx.PushDebugInfo(func() string {
			info := fmt.Sprintf("\t\t总伤害：%f\n", damageCtx.DamageValue())
			return info
		})
	}()

	casterEvents := &damageCtx.Attack.Caster.Events

	// 发送事件
	casterEvents.EmitBeforeDamageStep(damageCtx, DamageStep_DefendDeductHPShield)

	// 扣除护盾
	damage.decHPShield(damageCtx)

	// 发送事件
	casterEvents.EmitAfterDamageStep(damageCtx, DamageStep_DefendDeductHPShield)

	// 检测总伤害值
	if damageCtx.DamageValue() <= 0 {
		return false
	}

	return true
}

// stepResistanceDamageValue 步骤：抗性修正
func (damage *Damage) stepResistanceDamageValue(damageCtx *DamageContext) bool {
	damageCtx.PushDebugInfo(func() string {
		damageCtx.StepIndex++
		return fmt.Sprintf("【%d：抗性修正】\n", damageCtx.StepIndex)
	})

	casterEvents := &damageCtx.Attack.Caster.Events
	casterAttr := &damageCtx.Attack.CasterSnapshot.Attr
	targetAttr := &damageCtx.Target.Attr

	defer func() {
		damageCtx.PushDebugInfo(func() string {
			info := fmt.Sprintf("\t\t总伤害：%f\n", damageCtx.DamageValue())
			return info
		})
	}()

	// 发送事件
	casterEvents.EmitBeforeDamageStep(damageCtx, DamageStep_ResistanceDamageValue)

	for i := DamageValueKind_Begin; i < DamageValueKind_End; i++ {
		// 抗性修正
		var resistance float32

		// 抗性修正
		if i == DamageValueKind_Normal {
			resistance = float32(Min(Max(targetAttr.Armor/(casterAttr.Attack+math.Abs(targetAttr.Armor)), -1), 0.75))

			damageCtx.PushDebugInfo(func() string {
				return fmt.Sprintf("\t\t%s抗性：Min(Max(%f / (%f + Abs(%f))), -1), 0.75) = %f\n",
					getDamageValueKindText(i),
					targetAttr.Armor,
					casterAttr.Attack,
					targetAttr.Armor,
					resistance)
			})

		} else {
			resistance = targetAttr.GetResistance(i)

			damageCtx.PushDebugInfo(func() string {
				return fmt.Sprintf("\t\t%s抗性：%f\n",
					getDamageValueKindText(i),
					resistance)
			})
		}

		// 基础伤害
		baseDamage := damageCtx.DamageValueTab[i]

		// 抗性修正
		damageCtx.DamageValueTab[i] = Ceil(Max(baseDamage*(1-float64(resistance)), 0))

		damageCtx.PushDebugInfo(func() string {
			return fmt.Sprintf("\t\t%s抗性修正：Ceil(Max(%f * (1 - %f), 0)) = %f\n",
				getDamageValueKindText(i),
				baseDamage,
				resistance,
				damageCtx.DamageValueTab[i])
		})
	}

	// 发送事件
	casterEvents.EmitAfterDamageStep(damageCtx, DamageStep_ResistanceDamageValue)

	return true
}

// stepDecrDamageValue 步骤：减伤修正
func (damage *Damage) stepDecrDamageValue(damageCtx *DamageContext) bool {
	damageCtx.PushDebugInfo(func() string {
		damageCtx.StepIndex++
		return fmt.Sprintf("【%d：减伤修正】\n", damageCtx.StepIndex)
	})

	casterEvents := &damageCtx.Attack.Caster.Events
	targetAttr := &damageCtx.Target.Attr

	defer func() {
		damageCtx.PushDebugInfo(func() string {
			info := fmt.Sprintf("\t\t总伤害：%f\n", damageCtx.DamageValue())
			return info
		})
	}()

	// 发送事件
	casterEvents.EmitBeforeDamageStep(damageCtx, DamageStep_DecrDamageValue)

	for i := DamageValueKind_Begin; i < DamageValueKind_End; i++ {
		// 基础伤害
		baseDamage := damageCtx.DamageValueTab[i]

		// 修正伤害，后续用于处理伤害反转
		modifyDamage := baseDamage*(1+targetAttr.GetBeDamageAdd(i)+targetAttr.GetBeDamageDecr(i)) - targetAttr.GetBeDamageDeduct(i)
		damageNegative := false
		if modifyDamage < 0 {
			modifyDamage = Floor(modifyDamage)
			damageNegative = true
		} else {
			modifyDamage = Ceil(modifyDamage)
		}

		// 减伤修正
		damageCtx.DamageValueTab[i] = modifyDamage

		damageCtx.PushDebugInfo(func() string {
			return fmt.Sprintf("\t\t%s减伤修正：%s(%f * (1 + %f + %f) - %f) = %f\n",
				getDamageValueKindText(i),
				func() string {
					if damageNegative {
						return "Floor"
					}
					return "Ceil"
				}(),
				baseDamage,
				targetAttr.GetBeDamageAdd(i),
				targetAttr.GetBeDamageDecr(i),
				targetAttr.GetBeDamageDeduct(i),
				damageCtx.DamageValueTab[i])
		})
	}

	// 发送事件
	casterEvents.EmitAfterDamageStep(damageCtx, DamageStep_DecrDamageValue)

	return true
}

// _SputteringDamage 分裂伤害效果
type _SputteringDamage struct {
	Blank
	skipTarget  *Pawn // 跳过的目标
	damageValue int64 // 伤害值
}

// OnAttackHitTarget 伤害体打击目标时（创建伤害体的技能或buff能收到）
func (effect *_SputteringDamage) OnAttackHitTarget(attack *Attack, target *Pawn, hitTimes uint32) {
	if effect.skipTarget.Equal(target) {
		return
	}

	damageValue := effect.damageValue
	if damageValue <= 0 {
		return
	}

	attack.Caster.Scene.PushDebugInfo(func() string {
		return fmt.Sprintf("分裂伤害，伤害体ID：%d，对目标${PawnID:%d}产生${DamageValue:%d}点伤害", attack.UID, target.UID, damageValue)
	})

	// 扣除护盾
	damageHPShield, _ := target.Attr.BeDamageChangeHPShield(attack, BitsStick(int32(Proto.DamageType_Damage)), effect.damageValue)

	// 剩余伤害值
	damageValue -= damageHPShield

	if damageValue > 0 {
		// 扣除HP
		target.Attr.BeDamageChangeHP(attack, DamageKind_Sputtering, BitsStick(int32(Proto.DamageType_Damage)),
			damageValue, damageHPShield, false)

		// 血量为0时死亡
		if target.Attr.CurHP <= 0 {
			target.State.Dead(attack)
		}
	}
}

// stepSputtering 步骤：分裂伤害
func (damage *Damage) stepSputtering(damageCtx *DamageContext) bool {
	damageCtx.PushDebugInfo(func() string {
		damageCtx.StepIndex++
		return fmt.Sprintf("【%d：分裂伤害】\n", damageCtx.StepIndex)
	})

	caster := damageCtx.Attack.Caster
	casterSnapshot := damageCtx.Attack.CasterSnapshot
	casterEvents := &caster.Events
	scene := caster.Scene

	// 发送事件
	casterEvents.EmitBeforeDamageStep(damageCtx, DamageStep_Sputtering)

	// 计算分裂伤害值
	sputteringDamageValue := int64(damageCtx.DamageValue() * float64(casterSnapshot.Attr.AttackSputteringRate))

	damageCtx.PushDebugInfo(func() string {
		return fmt.Sprintf("\t\t分裂伤害：Floor(%f * %f) = %d\n", damageCtx.DamageValue(), casterSnapshot.Attr.AttackSputteringRate, sputteringDamageValue)
	})

	if sputteringDamageValue <= 0 {
		damageCtx.PushDebugInfo(func() string {
			return fmt.Sprintf("\t\t分裂伤害为0，无需执行\n")
		})

		// 发送事件
		casterEvents.EmitAfterDamageStep(damageCtx, DamageStep_Sputtering)

		return true
	}

	// 分裂伤害效果
	sputteringDamage := &_SputteringDamage{
		skipTarget:  damageCtx.Target,
		damageValue: sputteringDamageValue,
	}

	// 创建伤害体
	if sputteringAttack, err := scene.CreateCustomAttack(InnerAttackID_SputteringAttack, caster, casterSnapshot, nil, damageCtx.Target.GetPos(), sputteringDamage, damageCtx.Attack.Scale); err != nil {
		damageCtx.PushDebugInfo(func() string {
			return fmt.Sprintf("\t\t创建分裂伤害体失败，%s\n", err.Error())
		})
		log.Error(err)
	} else {
		damageCtx.PushDebugInfo(func() string {
			return fmt.Sprintf("\t\t创建分裂伤害体成功，伤害体ID：%d\n", sputteringAttack.UID)
		})
	}

	// 发送事件
	casterEvents.EmitAfterDamageStep(damageCtx, DamageStep_Sputtering)

	return true
}
