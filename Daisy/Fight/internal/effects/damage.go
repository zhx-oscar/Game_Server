package effects

import (
	. "Daisy/Fight/internal"
	. "Daisy/Fight/internal/conf"
	"Daisy/Fight/internal/log"
	"Daisy/Proto"
	"fmt"
)

// Damage 伤害模板
type Damage struct {
	DamageKind     DamageKind                   // 伤害类型
	DamageFlow     DamageFlow                   // 伤害流程
	DamageValueTab [DamageValueKind_End]float64 // 伤害值表
}

// DamageTarget 伤害目标
func (damage *Damage) DamageTarget(attack *Attack, target *Pawn, extendValue float64) (damageBit Bits, damageValue, damageHP, damageHPShield int64) {
	// 伤害上下文
	damageCtx := &DamageContext{}

	// 初始化伤害上下文
	if err := damageCtx.Init(damage.DamageKind, damage.DamageFlow, attack, target, extendValue); err != nil {
		log.Errorf("damage context init failed, %s", err.Error())
		return
	}

	// 执行伤害流程
	damage.execDamageFlow(damageCtx)

	// 记录伤害结果
	damageBit = damageCtx.DamageBit
	damageValue = int64(damageCtx.DamageValue())
	damageHP = damageCtx.DamageHP
	damageHPShield = damageCtx.DamageHPShield

	damageCtx.PushDebugInfo(func() string {
		return "【伤害结果】\n" +
			fmt.Sprintf("\t\t总伤害值：${DamageValue:%d}\n", damageValue+damageHPShield) +
			fmt.Sprintf("\t\t扣除HP值：${DamageHP:%d}\n", damageHP) +
			fmt.Sprintf("\t\t扣除HP护盾值：${DamageHPShield:%d}\n", damageHPShield)
	})

	// 刷新debug信息
	damageCtx.FlushDebugInfo()

	// 首次单体伤害
	if damageCtx.Attack.HitTimes <= 0 && damageCtx.Attack.Config.Type == AttackType_Single {
		// 记录首次hit伤害类型
		GetAttackLogicData(damageCtx.Attack).SetFirstHitDamageBit(damageCtx.DamageBit)
	}

	floatWordType := Proto.DamageFloatWordType_Normal
	if attack.Skill != nil && attack.Skill.Config.SkillKind == SkillKind_Combine {
		floatWordType = Proto.DamageFloatWordType_CombineSkill
	}

	// 记录回放
	attack.Caster.Scene.PushAction(&Proto.AttackHit{
		TargetId:            target.UID,
		AttackId:            attack.UID,
		HitTimes:            attack.HitTimes,
		DamageBit:           uint64(damageBit),
		DamageHP:            damageHP,
		DamageHPShield:      damageHPShield,
		DamageValue:         damageValue,
		DamageFloatWordType: floatWordType,
	})

	attack.Caster.Scene.PushDebugInfo(func() string {
		damageBitTextFun := func() string {
			info := ""
			if damageBit.Test(int32(Proto.DamageType_Damage)) {
				info += "伤害"
			}
			if damageBit.Test(int32(Proto.DamageType_RecoverHP)) {
				if info != "" {
					info += "|"
				}
				info += "恢复"
			}
			if damageBit.Test(int32(Proto.DamageType_Crit)) {
				if info != "" {
					info += "|"
				}
				info += "暴击"
			}
			if damageBit.Test(int32(Proto.DamageType_Dodge)) {
				if info != "" {
					info += "|"
				}
				info += "闪避"
			}
			if damageBit.Test(int32(Proto.DamageType_Block)) {
				if info != "" {
					info += "|"
				}
				info += "格挡"
			}
			if damageBit.Test(int32(Proto.DamageType_ExemptionDamage)) {
				if info != "" {
					info += "|"
				}
				info += "免伤"
			}
			if damageBit.Test(int32(Proto.DamageType_Miss)) {
				if info != "" {
					info += "|"
				}
				info += "未命中"
			}
			if damageBit.Test(int32(Proto.DamageType_ExemptionDamage)) {
				if info != "" {
					info += "|"
				}
				info += "伤害反转"
			}
			return info
		}

		switch attack.Src() {
		case Proto.AttackSrc_Skill:
			return fmt.Sprintf("${PawnID:%d}的技能${SkillID:%d}，技能流水：%d，对目标${PawnID:%d}产生伤害，Hit段数：%d，伤害标记：[%s]，伤害值：${DamageValue:%d}，扣除HP：${DamageHP:%d}，扣除HP护盾：${DamageHPShield:%d}，目标剩余HP：%d",
				attack.Caster.UID,
				attack.Skill.Config.ValueID(),
				attack.Skill.UID,
				target.UID,
				attack.HitTimes+1,
				damageBitTextFun(),
				damageValue+damageHPShield,
				damageHP,
				damageHPShield,
				target.Attr.CurHP)
		case Proto.AttackSrc_Buff:
			return fmt.Sprintf("${PawnID:%d}对目标${PawnID:%d}施加的Buff${BuffID:%d}产生伤害，Hit段数：%d，伤害标记：%s，伤害值：${DamageValue:%d}，扣除HP：${DamageHP:%d}，扣除HP护盾：${DamageHPShield:%d}，目标剩余HP：%d",
				attack.Buff.Caster.UID,
				target.UID,
				attack.Buff.Config.MainID(),
				attack.HitTimes+1,
				damageBitTextFun(),
				damageValue+damageHPShield,
				damageHP,
				damageHPShield,
				target.Attr.CurHP)
		case Proto.AttackSrc_Custom:
			return fmt.Sprintf("${PawnID:%d}对目标${PawnID:%d}产生自定义伤害，Hit段数：%d，伤害标记：%s，伤害值：${DamageValue:%d}，扣除HP：${DamageHP:%d}，扣除HP护盾：${DamageHPShield:%d}，目标剩余HP：%d",
				attack.Buff.Caster.UID,
				target.UID,
				attack.HitTimes+1,
				damageBitTextFun(),
				damageValue+damageHPShield,
				damageHP,
				damageHPShield,
				target.Attr.CurHP)
		}
		return ""
	})

	// 血量为0时死亡
	if target.Attr.CurHP <= 0 {
		target.State.Dead(attack)
	}

	return
}

// execDamageFlow 执行伤害流程
func (damage *Damage) execDamageFlow(damageCtx *DamageContext) {
	damageCtx.PushDebugInfo(func() string {
		caster := damageCtx.Attack.Caster
		target := damageCtx.Target

		info := "【技能伤害】\n"

		switch damage.DamageFlow {
		case DamageFlow_Attack:
			info += "【攻击流程】\n"
		case DamageFlow_CastSkill:
			info += "【施法流程】\n"
		}

		info += "【攻击方】\n" +
			fmt.Sprintf("\t\t名称：${PawnID:%d}\n", caster.UID) +
			fmt.Sprintf("\t\t等级：%d\n", caster.Info.Level) +
			fmt.Sprintf("\t\t属性ID：%d\n", caster.Attr.ID) +
			"【防守方】\n" +
			fmt.Sprintf("\t\t名称：${PawnID:%d}\n", target.UID) +
			fmt.Sprintf("\t\t等级：%d\n", target.Info.Level) +
			fmt.Sprintf("\t\t属性ID：%d\n", target.Attr.ID)

		info += "【伤害来源】\n"

		switch damageCtx.Attack.Src() {
		case Proto.AttackSrc_Skill:
			info += fmt.Sprintf("\t\t技能：${SkillID:%d}\n", damageCtx.Attack.Skill.Config.ValueID()) +
				fmt.Sprintf("\t\t技能ID：%d\n", damageCtx.Attack.Skill.Config.ValueID()) +
				fmt.Sprintf("\t\t技能流水：%d\n", damageCtx.Attack.Skill.UID)
		case Proto.AttackSrc_Buff:
			info += fmt.Sprintf("\t\tBuff：${BuffID:%d}\n", damageCtx.Attack.Buff.Config.MainID()) +
				fmt.Sprintf("\t\tBuffID：%d\n", damageCtx.Attack.Buff.Config.MainID())
		case Proto.AttackSrc_Custom:
			info += "\t\t自定义\n"
		}

		info += fmt.Sprintf("\t\tHit段数：%d\n", damageCtx.Attack.HitTimes+1)

		return info
	})

	switch damage.DamageFlow {
	case DamageFlow_Attack:
		damage.execAttackFlow(damageCtx)
	case DamageFlow_CastSkill:
		damage.execCasterSkillFlow(damageCtx)
	}
}
