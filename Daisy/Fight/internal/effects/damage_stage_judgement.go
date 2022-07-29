package effects

import (
	. "Daisy/Fight/internal"
	. "Daisy/Fight/internal/conf"
	"Daisy/Proto"
	"fmt"
	"math/rand"
)

// stepRoundTableJudge 步骤：圆桌判定
func (damage *Damage) stepRoundTableJudge(damageCtx *DamageContext) bool {
	damageCtx.PushDebugInfo(func() string {
		damageCtx.StepIndex++
		return fmt.Sprintf("【%d：圆桌判定】\n", damageCtx.StepIndex)
	})

	defer func() {
		damageCtx.PushDebugInfo(func() string {
			info := fmt.Sprintf("\t\t总伤害：%f\n", damageCtx.DamageValue())
			return info
		})
	}()

	casterEvents := &damageCtx.Attack.Caster.Events
	casterAttr := &damageCtx.Attack.CasterSnapshot.Attr
	targetAttr := &damageCtx.Target.Attr

	// 发送事件
	casterEvents.EmitBeforeDamageStep(damageCtx, DamageStep_RoundTableJudge)

	// 非首次单体伤害
	if damageCtx.Attack.HitTimes > 0 && damageCtx.Attack.Config.Type == AttackType_Single {
		// 首次命中状态
		firstHitDamageBit := GetAttackLogicData(damageCtx.Attack).FirstHitDamageBit

		if firstHitDamageBit.Test(int32(Proto.DamageType_Dodge)) {
			damageCtx.DamageBit.TurnOn(int32(Proto.DamageType_Dodge))

			damageCtx.PushDebugInfo(func() string {
				return fmt.Sprintf("\t\t继承首次判定：闪避\n")
			})

			// 发送事件
			casterEvents.EmitAfterDamageStep(damageCtx, DamageStep_RoundTableJudge)

			return false
		}

		if firstHitDamageBit.Test(int32(Proto.DamageType_Miss)) {
			damageCtx.DamageBit.TurnOn(int32(Proto.DamageType_Miss))

			damageCtx.PushDebugInfo(func() string {
				return fmt.Sprintf("\t\t继承首次判定：丢失\n")
			})

			// 发送事件
			casterEvents.EmitAfterDamageStep(damageCtx, DamageStep_RoundTableJudge)

			return false
		}
	}

	// 圆桌表
	var roundTab _RoundTable

	// 设置数据
	roundTab[DamageJudgeRv_Miss].Max = roundFloat2Int(float32(Max(float64(1-casterAttr.HitRate), 0)))
	roundTab[DamageJudgeRv_Dodge].Max = roundFloat2Int(targetAttr.DodgeRate)
	roundTab[DamageJudgeRv_Crit].Max = roundFloat2Int(casterAttr.CritRate)
	roundTab[DamageJudgeRv_Block].Max = roundFloat2Int(targetAttr.BlockRate)
	roundTab[DamageJudgeRv_Hit].Max = roundFloat2Int(1)

	// 随机结果
	roundRv, n := roundTab.random()

	damageCtx.PushDebugInfo(func() string {
		info := fmt.Sprintf("\t\t圆桌表：\n")
		for i := 0; i < len(roundTab); i++ {
			var lowBound, highBound int32

			if i > 0 {
				lowBound = roundTab[i-1].Max
			}
			highBound = roundTab[i].Max - 1

			if lowBound <= highBound {
				info += fmt.Sprintf("\t\t\t\t=>%-7s%-7d%-7f - %-7f\n", getRoundRvText(DamageJudgeRv(i)), i, roundInt2Float(lowBound), roundInt2Float(highBound))
			} else {
				info += fmt.Sprintf("\t\t\t\t=>%-7s%-7d无\n", getRoundRvText(DamageJudgeRv(i)), i)
			}
		}

		info += fmt.Sprintf("\t\t随机点数：%f\n", roundInt2Float(n)) +
			fmt.Sprintf("\t\t判定结果：%s\n", getRoundRvText(roundRv))

		return info
	})

	// 非首次单体伤害
	if damageCtx.Attack.HitTimes > 0 && damageCtx.Attack.Config.Type == AttackType_Single {
		// 修改圆桌判定结果
		switch roundRv {
		case DamageJudgeRv_Miss:
			fallthrough
		case DamageJudgeRv_Dodge:
			damageCtx.PushDebugInfo(func() string {
				return fmt.Sprintf("\t\t修正非首次判定结果：%s => %s\n", getRoundRvText(roundRv), getRoundRvText(DamageJudgeRv_Hit))
			})
			roundRv = DamageJudgeRv_Hit
		}
	}

	// 执行圆桌判定结果
	switch roundRv {
	case DamageJudgeRv_Miss:
		damageCtx.DamageBit.TurnOn(int32(Proto.DamageType_Miss))

		// 发送事件
		casterEvents.EmitAfterDamageStep(damageCtx, DamageStep_RoundTableJudge)

		return false

	case DamageJudgeRv_Dodge:
		damageCtx.DamageBit.TurnOn(int32(Proto.DamageType_Dodge))

		// 发送事件
		casterEvents.EmitAfterDamageStep(damageCtx, DamageStep_RoundTableJudge)

		return false

	case DamageJudgeRv_Crit:
		damageCtx.DamageBit.TurnOn(int32(Proto.DamageType_Crit))

		for i := DamageValueKind_Begin; i < DamageValueKind_End; i++ {
			// 基础伤害
			baseDamage := damageCtx.DamageValueTab[i]

			// 暴击修正
			damageCtx.DamageValueTab[i] = Ceil(Max(baseDamage*float64(casterAttr.CritDamageRate), 0))

			damageCtx.PushDebugInfo(func() string {
				return fmt.Sprintf("\t\t%s暴击修正：Ceil(Max(%f * %f, 0)) = %f\n",
					getDamageValueKindText(i),
					baseDamage,
					casterAttr.CritDamageRate,
					damageCtx.DamageValueTab[i])
			})
		}

		// 扣除护盾
		damage.decHPShield(damageCtx)

	case DamageJudgeRv_Block:
		damageCtx.DamageBit.TurnOn(int32(Proto.DamageType_Block))

		// 扣除护盾
		damage.decHPShield(damageCtx)

		// 检测总伤害值
		if damageCtx.DamageValue() <= 0 {
			// 发送事件
			casterEvents.EmitAfterDamageStep(damageCtx, DamageStep_RoundTableJudge)

			return false
		}

		for i := DamageValueKind_Begin; i < DamageValueKind_End; i++ {
			// 基础伤害
			baseDamage := damageCtx.DamageValueTab[i]

			// 格挡修正
			damageCtx.DamageValueTab[i] = Ceil(Max(baseDamage-targetAttr.BlockValue, 0))

			damageCtx.PushDebugInfo(func() string {
				return fmt.Sprintf("\t\t%s格挡修正：Ceil(Max(%f - %f, 0)) = %f\n",
					getDamageValueKindText(i),
					baseDamage,
					targetAttr.BlockValue,
					damageCtx.DamageValueTab[i])
			})
		}

	case DamageJudgeRv_Hit:
		// 扣除护盾
		damage.decHPShield(damageCtx)
	}

	// 发送事件
	casterEvents.EmitAfterDamageStep(damageCtx, DamageStep_RoundTableJudge)

	// 检测总伤害值
	if damageCtx.DamageValue() <= 0 {
		return false
	}

	return true
}

// stepWaterfallJudge 步骤：瀑布判定
func (damage *Damage) stepWaterfallJudge(damageCtx *DamageContext) bool {
	damageCtx.PushDebugInfo(func() string {
		damageCtx.StepIndex++
		return fmt.Sprintf("【%d：瀑布判定】\n", damageCtx.StepIndex)
	})

	defer func() {
		damageCtx.PushDebugInfo(func() string {
			info := fmt.Sprintf("\t\t总伤害：%f\n", damageCtx.DamageValue())
			return info
		})
	}()

	casterEvents := &damageCtx.Attack.Caster.Events
	casterAttr := &damageCtx.Attack.CasterSnapshot.Attr
	targetAttr := &damageCtx.Target.Attr

	// 发送事件
	casterEvents.EmitBeforeDamageStep(damageCtx, DamageStep_WaterfallJudge)

	// 非首次单体伤害
	if damageCtx.Attack.HitTimes > 0 && damageCtx.Attack.Config.Type == AttackType_Single {
		// 首次命中状态
		firstHitDamageBit := GetAttackLogicData(damageCtx.Attack).FirstHitDamageBit

		if firstHitDamageBit.Test(int32(Proto.DamageType_Dodge)) {
			damageCtx.DamageBit.TurnOn(int32(Proto.DamageType_Dodge))

			damageCtx.PushDebugInfo(func() string {
				return fmt.Sprintf("\t\t继承首次判定：闪避\n")
			})

			// 发送事件
			casterEvents.EmitAfterDamageStep(damageCtx, DamageStep_WaterfallJudge)

			return false
		}

		if firstHitDamageBit.Test(int32(Proto.DamageType_Miss)) {
			damageCtx.DamageBit.TurnOn(int32(Proto.DamageType_Miss))

			damageCtx.PushDebugInfo(func() string {
				return fmt.Sprintf("\t\t继承首次判定：丢失\n")
			})

			// 发送事件
			casterEvents.EmitAfterDamageStep(damageCtx, DamageStep_WaterfallJudge)

			return false
		}
	}

	// 计算命中
	if damageCtx.Attack.HitTimes <= 0 && damageCtx.Attack.Config.Type == AttackType_Single ||
		damageCtx.Attack.Config.Type == AttackType_Aoe {
		damageCtx.PushDebugInfo(func() string {
			return fmt.Sprintf("\t\t计算命中：\n")
		})

		randNum := rand.Float32()

		damageCtx.PushDebugInfo(func() string {
			return fmt.Sprintf("\t\t\t\t随机点数：%f\n", randNum) +
				fmt.Sprintf("\t\t\t\t命中率：%f\n", casterAttr.HitRate)
		})

		if randNum >= casterAttr.HitRate {
			damageCtx.PushDebugInfo(func() string {
				return fmt.Sprintf("\t\t\t\t命中结果：%f >= %f = 失败\n", randNum, casterAttr.HitRate) +
					fmt.Sprintf("\t\t判定结果：未命中\n")
			})

			damageCtx.DamageBit.TurnOn(int32(Proto.DamageType_Miss))

			// 发送事件
			casterEvents.EmitAfterDamageStep(damageCtx, DamageStep_WaterfallJudge)

			return false
		}

		damageCtx.PushDebugInfo(func() string {
			return fmt.Sprintf("\t\t\t\t命中结果：%f < %f = 成功\n", randNum, casterAttr.HitRate)
		})

	} else {
		damageCtx.PushDebugInfo(func() string {
			return fmt.Sprintf("\t\t计算命中：非首次单体伤害，无需计算\n")
		})
	}

	// 发送事件
	casterEvents.EmitStepWaterfallJudgePass(damageCtx, DamageJudgeRv_Miss)

	// 计算闪避
	if damageCtx.Attack.HitTimes <= 0 && damageCtx.Attack.Config.Type == AttackType_Single ||
		damageCtx.Attack.Config.Type == AttackType_Aoe {
		damageCtx.PushDebugInfo(func() string {
			return fmt.Sprintf("\t\t计算闪避：\n")
		})

		randNum := rand.Float32()

		damageCtx.PushDebugInfo(func() string {
			return fmt.Sprintf("\t\t\t\t随机点数：%f\n", randNum) +
				fmt.Sprintf("\t\t\t\t闪避率：%f\n", targetAttr.DodgeRate)
		})

		if randNum < targetAttr.DodgeRate {
			damageCtx.PushDebugInfo(func() string {
				return fmt.Sprintf("\t\t\t\t闪避结果：%f < %f = 成功\n", randNum, targetAttr.DodgeRate) +
					fmt.Sprintf("\t\t判定结果：闪避\n")
			})

			damageCtx.DamageBit.TurnOn(int32(Proto.DamageType_Dodge))

			// 发送事件
			casterEvents.EmitAfterDamageStep(damageCtx, DamageStep_WaterfallJudge)

			return false
		}

		damageCtx.PushDebugInfo(func() string {
			return fmt.Sprintf("\t\t\t\t闪避结果：%f >= %f = 失败\n", randNum, targetAttr.DodgeRate)
		})

	} else {
		damageCtx.PushDebugInfo(func() string {
			return fmt.Sprintf("\t\t计算闪避：非首次单体伤害，无需计算\n")
		})
	}

	// 发送事件
	casterEvents.EmitStepWaterfallJudgePass(damageCtx, DamageJudgeRv_Dodge)

	// 计算暴击
	{
		damageCtx.PushDebugInfo(func() string {
			return fmt.Sprintf("\t\t计算暴击：\n")
		})

		randNum := rand.Float32()

		damageCtx.PushDebugInfo(func() string {
			return fmt.Sprintf("\t\t\t\t随机点数：%f\n", randNum) +
				fmt.Sprintf("\t\t\t\t暴击率：%f\n", casterAttr.CritRate)
		})

		if randNum < casterAttr.CritRate {
			damageCtx.PushDebugInfo(func() string {
				return fmt.Sprintf("\t\t\t\t暴击结果：%f < %f = 成功\n", randNum, casterAttr.CritRate) +
					fmt.Sprintf("\t\t判定结果：暴击\n")
			})

			damageCtx.DamageBit.TurnOn(int32(Proto.DamageType_Crit))

			for i := DamageValueKind_Begin; i < DamageValueKind_End; i++ {
				// 基础伤害
				baseDamage := damageCtx.DamageValueTab[i]

				// 暴击修正
				damageCtx.DamageValueTab[i] = Ceil(Max(baseDamage*float64(casterAttr.CritDamageRate), 0))

				damageCtx.PushDebugInfo(func() string {
					return fmt.Sprintf("\t\t\t\t%s暴击修正：Ceil(Max(%f * %f, 0) = %f\n",
						getDamageValueKindText(i),
						baseDamage,
						casterAttr.CritDamageRate,
						damageCtx.DamageValueTab[i])
				})
			}

			// 扣除护盾
			damage.decHPShield(damageCtx)

			// 发送事件
			casterEvents.EmitAfterDamageStep(damageCtx, DamageStep_WaterfallJudge)

			// 检测总伤害值
			if damageCtx.DamageValue() <= 0 {
				return false
			}

			return true
		}

		damageCtx.PushDebugInfo(func() string {
			return fmt.Sprintf("\t\t\t\t暴击结果：%f >= %f = 失败\n", randNum, targetAttr.DodgeRate)
		})
	}

	// 发送事件
	casterEvents.EmitStepWaterfallJudgePass(damageCtx, DamageJudgeRv_Crit)

	// 扣除护盾
	damage.decHPShield(damageCtx)

	// 检测总伤害值
	if damageCtx.DamageValue() <= 0 {
		// 发送事件
		casterEvents.EmitAfterDamageStep(damageCtx, DamageStep_WaterfallJudge)

		return false
	}

	// 计算格挡
	{
		damageCtx.PushDebugInfo(func() string {
			return fmt.Sprintf("\t\t计算格挡：\n")
		})

		randNum := rand.Float32()

		damageCtx.PushDebugInfo(func() string {
			return fmt.Sprintf("\t\t\t\t随机点数：%f\n", randNum) +
				fmt.Sprintf("\t\t\t\t格挡率：%f\n", targetAttr.BlockRate)
		})

		if randNum < targetAttr.BlockRate {
			damageCtx.PushDebugInfo(func() string {
				return fmt.Sprintf("\t\t\t\t格挡结果：%f < %f = 成功\n", randNum, targetAttr.BlockRate) +
					fmt.Sprintf("\t\t判定结果：格挡\n")
			})

			damageCtx.DamageBit.TurnOn(int32(Proto.DamageType_Block))

			for i := DamageValueKind_Begin; i < DamageValueKind_End; i++ {
				// 基础伤害
				baseDamage := damageCtx.DamageValueTab[i]

				// 格挡修正
				damageCtx.DamageValueTab[i] = Ceil(Max(baseDamage-targetAttr.BlockValue, 0))

				damageCtx.PushDebugInfo(func() string {
					return fmt.Sprintf("\t\t\t\t%s格挡修正：Ceil(Max(%f - %f, 0)) = %f\n",
						getDamageValueKindText(i),
						baseDamage,
						targetAttr.BlockValue,
						damageCtx.DamageValueTab[i])
				})
			}

			// 发送事件
			casterEvents.EmitAfterDamageStep(damageCtx, DamageStep_WaterfallJudge)

			// 检测总伤害值
			if damageCtx.DamageValue() <= 0 {
				return false
			}

			return true
		}

		damageCtx.PushDebugInfo(func() string {
			return fmt.Sprintf("\t\t\t\t格挡结果：%f >= %f = 失败\n", randNum, targetAttr.BlockRate)
		})
	}

	// 发送事件
	casterEvents.EmitStepWaterfallJudgePass(damageCtx, DamageJudgeRv_Block)

	damageCtx.PushDebugInfo(func() string {
		return fmt.Sprintf("\t\t判定结果：正常\n")
	})

	// 发送事件
	casterEvents.EmitAfterDamageStep(damageCtx, DamageStep_WaterfallJudge)

	// 检测总伤害值
	if damageCtx.DamageValue() <= 0 {
		return false
	}

	return true
}
