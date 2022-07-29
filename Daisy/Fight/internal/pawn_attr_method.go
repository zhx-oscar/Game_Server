package internal

import (
	"Daisy/Fight/internal/conf"
	"Daisy/Proto"
	"fmt"
)

// _HPShieldChange 护盾变化
type _HPShieldChange struct {
	HPShield
	oldShieldHP int64
}

// BeDamageChangeHPShield 受到伤害调整HP护盾
func (attr *FightAttr) BeDamageChangeHPShield(attack *Attack, damageBit Bits, damageValue int64) (damageHPShield int64, exemptionDamage bool) {
	if !attr.inFight || attr.pawn.isSnapshot || !attr.pawn.IsAlive() || damageValue <= 0 {
		return
	}

	// 护盾只能扣除
	if !damageBit.Test(int32(Proto.DamageType_Damage)) {
		return
	}

	// 检测免伤
	if attr.pawn.State.CantBeDamage {
		exemptionDamage = true
		return
	}

	// 记录护盾总值
	oldAllHPShield := attr.AllHPShield

	// 护盾变化表
	var shieldChangeList []_HPShieldChange

	// 剩余伤害值
	leftDamageValue := damageValue

	// 扣除每个护盾
	for i := len(attr.HPShieldList) - 1; i >= 0; i-- {
		shield := attr.HPShieldList[i]

		// 记录护盾旧值
		oldShieldHP := shield.ShieldHP

		// 计算需要扣除值
		decHP := leftDamageValue
		if decHP > oldShieldHP {
			decHP = oldShieldHP
		}

		// 扣除护盾值
		shield.ShieldHP -= decHP

		// 扣除总护盾值
		attr.AllHPShield -= decHP

		// 记录护盾变化
		shieldChangeList = append(shieldChangeList, _HPShieldChange{
			HPShield:    *shield,
			oldShieldHP: oldShieldHP,
		})

		// 检测删除护盾
		if shield.ShieldHP <= 0 {
			attr.HPShieldList = append(attr.HPShieldList[:i], attr.HPShieldList[i+1:]...)
		}

		// 计算剩余伤害值
		leftDamageValue -= decHP

		// 伤害值扣完
		if leftDamageValue <= 0 {
			break
		}
	}

	// 记录护盾与血量扣除数
	damageHPShield = damageValue - leftDamageValue

	// 记录回放
	attr.AttrSyncToClient(Proto.AttrType_HPShield, float64(oldAllHPShield), float64(attr.AllHPShield))

	// 发送护盾变化事件
	for i := range shieldChangeList {
		shieldChange := &shieldChangeList[i]

		// 发送护盾血量扣除事件
		attr.pawn.Events.EmitDecHPShield(attack, &shieldChange.HPShield, shieldChange.oldShieldHP)

		// 发送护盾破碎事件
		if shieldChange.ShieldHP <= 0 {
			attr.pawn.Scene.PushDebugInfo(func() string {
				return fmt.Sprintf("${PawnID:%d}的护盾%d损毁", attr.pawn.UID, shieldChange.UID)
			})

			attr.pawn.Events.EmitShieldBroken(attack, &shieldChange.HPShield)
		}
	}

	return
}

// BeDamageChangeHP 受到伤害调整HP
func (attr *FightAttr) BeDamageChangeHP(attack *Attack, damageKind conf.DamageKind, damageBit Bits, damageValue, damageHPShield int64, damageNotDead bool) (damageHP int64, exemptionDamage bool) {
	if !attr.inFight || attr.pawn.isSnapshot || !attr.pawn.IsAlive() || damageValue <= 0 || damageHPShield < 0 {
		return
	}

	// 选择是伤害还是治疗
	if damageBit.Test(int32(Proto.DamageType_Damage)) {
		// 检测免伤
		if attr.pawn.State.CantBeDamage {
			exemptionDamage = true
			return
		}

		oldHP := attr.CurHP

		// 调整HP
		if damageBit.Test(int32(Proto.DamageType_Invert)) {
			attr.ChangeHP(attr.CurHP + damageValue)
		} else {
			attr.ChangeHP(attr.CurHP - damageValue)
		}

		// 是否伤害不致死
		if attr.CurHP <= 0 && damageNotDead {
			attr.ChangeHP(1)
		}

		// 成功扣除的HP
		damageHP = IntAbs(oldHP - attr.CurHP)

	} else if damageBit.Test(int32(Proto.DamageType_RecoverHP)) {
		oldHP := attr.CurHP

		// 调整HP
		attr.ChangeHP(attr.CurHP + damageValue)

		// 成功恢复的HP
		damageHP = attr.CurHP - oldHP
	}

	// 发送伤害目标事件
	attack.Caster.Events.EmitDamageTarget(attack, attr.pawn, damageKind, damageBit, damageValue+damageHPShield, damageHP, damageHPShield)

	// 发送受到伤害事件
	attr.pawn.Events.EmitBeDamage(attack, damageKind, damageBit, damageValue+damageHPShield, damageHP, damageHPShield)

	return
}

// ChangeHP 调整HP
func (attr *FightAttr) ChangeHP(hp int64) {
	if !attr.inFight || attr.pawn.isSnapshot {
		return
	}

	if hp < 0 {
		hp = 0
	} else if hp > attr.MaxHP {
		hp = attr.MaxHP
	}

	old := attr.CurHP

	attr.CurHP = hp

	attr.AttrSyncToClient(Proto.AttrType_CurHP, float64(old), float64(attr.CurHP))
}

// ChangeUltimateSkillPower 调整必杀技能量
func (attr *FightAttr) ChangeUltimateSkillPower(power int32) {
	if !attr.inFight || attr.pawn.isSnapshot {
		return
	}

	if power < 0 {
		power = 0
	} else if power > attr.SkillPowerLimit {
		power = attr.SkillPowerLimit
	}

	old := attr.UltimateSkillPower

	attr.UltimateSkillPower = power

	attr.AttrSyncToClient(Proto.AttrType_UltimateSkillPower, float64(old), float64(attr.UltimateSkillPower))
}

// UltimateSkillPowerIsFull 必杀技能量满
func (attr *FightAttr) UltimateSkillPowerIsFull() bool {
	return attr.UltimateSkillPower >= attr.SkillPowerLimit
}
