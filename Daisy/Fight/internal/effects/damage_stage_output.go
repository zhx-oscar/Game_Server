package effects

import (
	. "Daisy/Fight/internal"
	. "Daisy/Fight/internal/conf"
	"fmt"
	"math/rand"
)

// stepAttackLuckyJudge 步骤：攻击幸运判定
func (damage *Damage) stepAttackLuckyJudge(damageCtx *DamageContext) bool {
	damageCtx.PushDebugInfo(func() string {
		damageCtx.StepIndex++
		return fmt.Sprintf("【%d：攻击幸运判定】\n", damageCtx.StepIndex)
	})

	casterEvents := &damageCtx.Attack.Caster.Events
	casterAttr := &damageCtx.Attack.CasterSnapshot.Attr

	// 发送事件
	casterEvents.EmitBeforeDamageStep(damageCtx, DamageStep_AttackLuckyJudge)

	// 幸运判定
	damageCtx.PushDebugInfo(func() string {
		return fmt.Sprintf("\t\t幸运判定：\n")
	})

	damageCtx.AttackLucky = false

	randNum := rand.Float32()

	if randNum < casterAttr.AttackLucky {
		damageCtx.PushDebugInfo(func() string {
			return fmt.Sprintf("\t\t\t\t随机点数：%f\n", randNum) +
				fmt.Sprintf("\t\t\t\t幸运值：%f\n", casterAttr.AttackLucky) +
				fmt.Sprintf("\t\t\t\t幸运结果：%f < %f = 成功\n", randNum, casterAttr.AttackLucky) +
				fmt.Sprintf("\t\t判定结果：幸运\n")
		})

		damageCtx.AttackLucky = true

	} else {
		damageCtx.PushDebugInfo(func() string {
			return fmt.Sprintf("\t\t\t\t随机点数：%f\n", randNum) +
				fmt.Sprintf("\t\t\t\t幸运值：%f\n", casterAttr.AttackLucky) +
				fmt.Sprintf("\t\t\t\t幸运结果：%f >= %f = 失败\n", randNum, casterAttr.AttackLucky) +
				fmt.Sprintf("\t\t判定结果：正常\n")
		})
	}

	// 发送事件
	casterEvents.EmitAfterDamageStep(damageCtx, DamageStep_AttackLuckyJudge)

	return true
}

// stepCountAttackDamageValue 步骤：统计攻击伤害
func (damage *Damage) stepCountAttackDamageValue(damageCtx *DamageContext) bool {
	damageCtx.PushDebugInfo(func() string {
		damageCtx.StepIndex++
		return fmt.Sprintf("【%d：统计攻击伤害】\n", damageCtx.StepIndex)
	})

	defer func() {
		damageCtx.PushDebugInfo(func() string {
			info := fmt.Sprintf("\t\t总伤害：%f\n", damageCtx.DamageValue())
			return info
		})
	}()

	casterEvents := &damageCtx.Attack.Caster.Events
	casterAttr := &damageCtx.Attack.CasterSnapshot.Attr

	// 发送事件
	casterEvents.EmitBeforeDamageStep(damageCtx, DamageStep_CountAttackDamageValue)

	// 装备攻击力
	var equipAttack float64

	if damageCtx.AttackLucky {
		equipAttack = casterAttr.EquipMaxAttack

		damageCtx.PushDebugInfo(func() string {
			return fmt.Sprintf("\t\t武器攻击力（幸运）：%f\n",
				equipAttack)
		})

	} else {
		equipAttack = Max(casterAttr.EquipMinAttack+rand.Float64()*(casterAttr.EquipMaxAttack-casterAttr.EquipMinAttack), 0)

		damageCtx.PushDebugInfo(func() string {
			return fmt.Sprintf("\t\t武器攻击力（正常）：Max(Randbetween(%f, %f), 0) = %f\n",
				casterAttr.EquipMaxAttack,
				casterAttr.EquipMinAttack,
				equipAttack)
		})
	}

	// 统计基础攻击力
	damageCtx.DamageValueTab[DamageValueKind_Normal] = Ceil(Max(casterAttr.Attack+equipAttack*(1+casterAttr.EquipAttackRate), 0))

	damageCtx.PushDebugInfo(func() string {
		return fmt.Sprintf("\t\t基础攻击力：Ceil(Max(%f + %f * (1 + %f), 0)) = %f\n",
			casterAttr.Attack,
			equipAttack,
			casterAttr.EquipAttackRate,
			damageCtx.DamageValueTab[DamageValueKind_Normal])
	})

	// 发送事件
	casterEvents.EmitAfterDamageStep(damageCtx, DamageStep_CountAttackDamageValue)

	return true
}

// stepIncrAttackDamageValue 步骤：攻击增伤修正
func (damage *Damage) stepIncrAttackDamageValue(damageCtx *DamageContext) bool {
	damageCtx.PushDebugInfo(func() string {
		damageCtx.StepIndex++
		return fmt.Sprintf("【%d：攻击增伤修正】\n", damageCtx.StepIndex)
	})

	defer func() {
		damageCtx.PushDebugInfo(func() string {
			info := fmt.Sprintf("\t\t总伤害：%f\n", damageCtx.DamageValue())
			return info
		})
	}()

	casterEvents := &damageCtx.Attack.Caster.Events
	casterAttr := &damageCtx.Attack.CasterSnapshot.Attr

	// 发送事件
	casterEvents.EmitBeforeDamageStep(damageCtx, DamageStep_IncrAttackDamageValue)

	// 选择伤害最高元素
	topElementDamageValue := 0.0
	topElement := DamageValueKind_Fire

	for i := DamageValueKind_Fire; i < DamageValueKind_End; i++ {
		if damageCtx.DamageValueTab[i] > topElementDamageValue {
			topElementDamageValue = damageCtx.DamageValueTab[i]
			topElement = i
		}
	}

	for i := DamageValueKind_Begin; i < DamageValueKind_End; i++ {
		// 基础伤害
		baseDamage := damageCtx.DamageValueTab[i]

		// 最高元素伤害附加值
		topElementDamagePlus := 0.0

		if topElement == i {
			topElementDamagePlus = casterAttr.TopAttackElementPlus
		}

		// 增伤修正
		damageCtx.DamageValueTab[i] = Ceil(Max((baseDamage+topElementDamagePlus+casterAttr.GetEquipAttackPlus(i))*(1+casterAttr.GetAttackAdd(i)+casterAttr.GetAttackDecr(i))+casterAttr.GetAttackPlus(i), 0))

		damageCtx.PushDebugInfo(func() string {
			return fmt.Sprintf("\t\t%s增伤修正：Ceil(Max((%f + %f + %f) * (1 + %f + %f) + %f, 0)) = %f\n",
				getDamageValueKindText(i),
				baseDamage,
				topElementDamagePlus,
				casterAttr.GetEquipAttackPlus(i),
				casterAttr.GetAttackAdd(i),
				casterAttr.GetAttackDecr(i),
				casterAttr.GetAttackPlus(i),
				damageCtx.DamageValueTab[i])
		})
	}

	// 发送事件
	casterEvents.EmitAfterDamageStep(damageCtx, DamageStep_IncrAttackDamageValue)

	return true
}

// stepCountCastSkillDamageValue 步骤：统计施法伤害
func (damage *Damage) stepCountCastSkillDamageValue(damageCtx *DamageContext) bool {
	damageCtx.PushDebugInfo(func() string {
		damageCtx.StepIndex++
		return fmt.Sprintf("【%d：统计施法伤害】\n", damageCtx.StepIndex)
	})

	defer func() {
		damageCtx.PushDebugInfo(func() string {
			info := fmt.Sprintf("\t\t总伤害：%f\n", damageCtx.DamageValue())
			return info
		})
	}()

	casterEvents := &damageCtx.Attack.Caster.Events

	// 发送事件
	casterEvents.EmitBeforeDamageStep(damageCtx, DamageStep_CountCastSkillDamageValue)

	// 拷贝技能伤害
	for damageKind := DamageValueKind_Normal; damageKind < DamageValueKind_End; damageKind++ {
		damageCtx.DamageValueTab[damageKind] = Ceil(Max(damage.DamageValueTab[damageKind], 0))
	}

	damageCtx.PushDebugInfo(func() string {
		return fmt.Sprintf("\t\t物理攻击力：%f\n", damageCtx.DamageValueTab[DamageValueKind_Normal]) +
			fmt.Sprintf("\t\t火焰攻击力：%f\n", damageCtx.DamageValueTab[DamageValueKind_Fire]) +
			fmt.Sprintf("\t\t冰霜攻击力：%f\n", damageCtx.DamageValueTab[DamageValueKind_Cold]) +
			fmt.Sprintf("\t\t毒素攻击力：%f\n", damageCtx.DamageValueTab[DamageValueKind_Poison]) +
			fmt.Sprintf("\t\t闪电攻击力：%f\n", damageCtx.DamageValueTab[DamageValueKind_Lightning])
	})

	// 发送事件
	casterEvents.EmitAfterDamageStep(damageCtx, DamageStep_CountCastSkillDamageValue)

	return true
}

// stepIncrCastSkillDamageValue 步骤：施法增伤修正
func (damage *Damage) stepIncrCastSkillDamageValue(damageCtx *DamageContext) bool {
	damageCtx.PushDebugInfo(func() string {
		damageCtx.StepIndex++
		return fmt.Sprintf("【%d：施法增伤修正】\n", damageCtx.StepIndex)
	})

	defer func() {
		damageCtx.PushDebugInfo(func() string {
			info := fmt.Sprintf("\t\t总伤害：%f\n", damageCtx.DamageValue())
			return info
		})
	}()

	casterEvents := &damageCtx.Attack.Caster.Events

	// 发送事件
	casterEvents.EmitBeforeDamageStep(damageCtx, DamageStep_IncrCastSkillDamageValue)

	// 发送事件
	casterEvents.EmitAfterDamageStep(damageCtx, DamageStep_IncrCastSkillDamageValue)

	return true
}
