package internal

import (
	. "Daisy/Fight/attraffix"
)

// assignAttr 赋值属性
func (attr *_Attr) assignAttr(field Field, value float64) {
	switch field {
	// 人物属性
	case Field_Strength:
		attr.Strength = value
	case Field_Agility:
		attr.Agility = value
	case Field_Intelligence:
		attr.Intelligence = value
	case Field_Vitality:
		attr.Vitality = value
	case Field_MaxHP:
		attr.MaxHP = int64(value)
	case Field_SkillPowerLimit:
		attr.SkillPowerLimit = int32(value)
	case Field_Attack:
		attr.Attack = value
	case Field_AttackNoramlAdd:
		attr.AttackNoramlAdd = value
	case Field_AttackNoramlDecr:
		attr.AttackNoramlDecr = value
	case Field_AttackFireAdd:
		attr.AttackFireAdd = value
	case Field_AttackFireDecr:
		attr.AttackFireDecr = value
	case Field_AttackColdAdd:
		attr.AttackColdAdd = value
	case Field_AttackColdDecr:
		attr.AttackColdDecr = value
	case Field_AttackPoisonAdd:
		attr.AttackPoisonAdd = value
	case Field_AttackPoisonDecr:
		attr.AttackPoisonDecr = value
	case Field_AttackLightningAdd:
		attr.AttackLightningAdd = value
	case Field_AttackLightningDecr:
		attr.AttackLightningDecr = value
	case Field_HitRate:
		attr.HitRate = float32(value)
	case Field_DodgeRate:
		attr.DodgeRate = float32(value)
	case Field_CritRate:
		attr.CritRate = float32(value)
	case Field_BlockRate:
		attr.BlockRate = float32(value)
	case Field_BlockValue:
		attr.BlockValue = value
	case Field_CritDamageRate:
		attr.CritDamageRate = float32(value)
	case Field_Armor:
		attr.Armor = value
	case Field_ResistanceFire:
		attr.ResistanceFire = float32(value)
	case Field_ResistanceCold:
		attr.ResistanceCold = float32(value)
	case Field_ResistancePoison:
		attr.ResistancePoison = float32(value)
	case Field_ResistanceLightning:
		attr.ResistanceLightning = float32(value)
	case Field_BeDamageNormalAdd:
		attr.BeDamageNormalAdd = value
	case Field_BeDamageNormalDecr:
		attr.BeDamageNormalDecr = value
	case Field_BeDamageFireAdd:
		attr.BeDamageFireAdd = value
	case Field_BeDamageFireDecr:
		attr.BeDamageFireDecr = value
	case Field_BeDamageColdAdd:
		attr.BeDamageColdAdd = value
	case Field_BeDamageColdDecr:
		attr.BeDamageColdDecr = value
	case Field_BeDamagePoisonAdd:
		attr.BeDamagePoisonAdd = value
	case Field_BeDamagePoisonDecr:
		attr.BeDamagePoisonDecr = value
	case Field_BeDamageLightningAdd:
		attr.BeDamageLightningAdd = value
	case Field_BeDamageLightningDecr:
		attr.BeDamageLightningDecr = value
	case Field_Lucky:
		attr.Lucky = uint32(value)
	case Field_NormalAttackSpeed:
		attr.NormalAttackSpeed = float32(value)
	case Field_AttackLucky:
		attr.AttackLucky = float32(value)
	case Field_AttackBloodsuckerRate:
		attr.AttackBloodsuckerRate = float32(value)
	case Field_AttackSputteringRate:
		attr.AttackSputteringRate = float32(value)
	case Field_AttackStealUltimateSkillPower:
		attr.AttackStealUltimateSkillPower = int32(value)
	case Field_BeDamageNormalDeduct:
		attr.BeDamageNormalDeduct = value
	case Field_BeDamageThorns:
		attr.BeDamageThorns = value
	case Field_RecoverHP:
		attr.RecoverHP = int64(value)
	case Field_RecoverUltimateSkillPowerRate:
		attr.RecoverUltimateSkillPowerRate = float32(value)
	case Field_PowerShieldHP:
		attr.PowerShieldHP = int64(value)
	case Field_PowerShieldRecoverSpeed:
		attr.PowerShieldRecoverSpeed = float32(value)
	case Field_TopAttackElementPlus:
		attr.TopAttackElementPlus = value
	case Field_OverDriveLimit:
		attr.OverDriveLimit = int32(value)
	case Field_OverDriveAddEfficiency:
		attr.OverDriveAddEfficiency = int32(value)
	case Field_OverDriveTime:
		attr.OverDriveTime = uint32(value)
	case Field_WeakTime:
		attr.WeakTime = uint32(value)
	case Field_BreakValueLimit:
		attr.BreakValueLimit = int32(value)
	case Field_BreakStateTime:
		attr.BreakStateTime = uint32(value)

	// 装备属性
	case Field_EquipMinAttack:
		attr.EquipMinAttack = value
	case Field_EquipMaxAttack:
		attr.EquipMaxAttack = value
	case Field_EquipAttackRate:
		attr.EquipAttackRate = value
	case Field_EquipAttackNormalPlus:
		attr.EquipAttackNormalPlus = value
	case Field_EquipAttackFirePlus:
		attr.EquipAttackFirePlus = value
	case Field_EquipAttackColdPlus:
		attr.EquipAttackColdPlus = value
	case Field_EquipAttackPoisonPlus:
		attr.EquipAttackPoisonPlus = value
	case Field_EquipAttackLightningPlus:
		attr.EquipAttackLightningPlus = value
	case Field_AttackNormalPlus:
		attr.AttackNormalPlus = value
	case Field_AttackFirePlus:
		attr.AttackFirePlus = value
	case Field_AttackColdPlus:
		attr.AttackColdPlus = value
	case Field_AttackPoisonPlus:
		attr.AttackPoisonPlus = value
	case Field_AttackLightningPlus:
		attr.AttackLightningPlus = value

	// 逻辑属性
	case Field_CollisionRadius:
		attr.CollisionRadius = float32(value)
	case Field_Mass:
		attr.Mass = uint32(value)
	case Field_Scale:
		attr.Scale = float32(value)
	}
}

// GetAttr 查询属性
func (attr *_Attr) GetAttr(field Field) float64 {
	switch field {
	// 人物属性
	case Field_Strength:
		return attr.Strength
	case Field_Agility:
		return attr.Agility
	case Field_Intelligence:
		return attr.Intelligence
	case Field_Vitality:
		return attr.Vitality
	case Field_MaxHP:
		return float64(attr.MaxHP)
	case Field_SkillPowerLimit:
		return float64(attr.SkillPowerLimit)
	case Field_Attack:
		return attr.Attack
	case Field_AttackNoramlAdd:
		return attr.AttackNoramlAdd
	case Field_AttackNoramlDecr:
		return attr.AttackNoramlDecr
	case Field_AttackFireAdd:
		return attr.AttackFireAdd
	case Field_AttackFireDecr:
		return attr.AttackFireDecr
	case Field_AttackColdAdd:
		return attr.AttackColdAdd
	case Field_AttackColdDecr:
		return attr.AttackColdDecr
	case Field_AttackPoisonAdd:
		return attr.AttackPoisonAdd
	case Field_AttackPoisonDecr:
		return attr.AttackPoisonDecr
	case Field_AttackLightningAdd:
		return attr.AttackLightningAdd
	case Field_AttackLightningDecr:
		return attr.AttackLightningDecr
	case Field_HitRate:
		return float64(attr.HitRate)
	case Field_DodgeRate:
		return float64(attr.DodgeRate)
	case Field_CritRate:
		return float64(attr.CritRate)
	case Field_BlockRate:
		return float64(attr.BlockRate)
	case Field_BlockValue:
		return attr.BlockValue
	case Field_CritDamageRate:
		return float64(attr.CritDamageRate)
	case Field_Armor:
		return attr.Armor
	case Field_ResistanceFire:
		return float64(attr.ResistanceFire)
	case Field_ResistanceCold:
		return float64(attr.ResistanceCold)
	case Field_ResistancePoison:
		return float64(attr.ResistancePoison)
	case Field_ResistanceLightning:
		return float64(attr.ResistanceLightning)
	case Field_BeDamageNormalAdd:
		return attr.BeDamageNormalAdd
	case Field_BeDamageNormalDecr:
		return attr.BeDamageNormalDecr
	case Field_BeDamageFireAdd:
		return attr.BeDamageFireAdd
	case Field_BeDamageFireDecr:
		return attr.BeDamageFireDecr
	case Field_BeDamageColdAdd:
		return attr.BeDamageColdAdd
	case Field_BeDamageColdDecr:
		return attr.BeDamageColdDecr
	case Field_BeDamagePoisonAdd:
		return attr.BeDamagePoisonAdd
	case Field_BeDamagePoisonDecr:
		return attr.BeDamagePoisonDecr
	case Field_BeDamageLightningAdd:
		return attr.BeDamageLightningAdd
	case Field_BeDamageLightningDecr:
		return attr.BeDamageLightningDecr
	case Field_Lucky:
		return float64(attr.Lucky)
	case Field_NormalAttackSpeed:
		return float64(attr.NormalAttackSpeed)
	case Field_AttackLucky:
		return float64(attr.AttackLucky)
	case Field_AttackBloodsuckerRate:
		return float64(attr.AttackBloodsuckerRate)
	case Field_AttackSputteringRate:
		return float64(attr.AttackSputteringRate)
	case Field_AttackStealUltimateSkillPower:
		return float64(attr.AttackStealUltimateSkillPower)
	case Field_BeDamageNormalDeduct:
		return attr.BeDamageNormalDeduct
	case Field_BeDamageThorns:
		return attr.BeDamageThorns
	case Field_RecoverHP:
		return float64(attr.RecoverHP)
	case Field_RecoverUltimateSkillPowerRate:
		return float64(attr.RecoverUltimateSkillPowerRate)
	case Field_PowerShieldHP:
		return float64(attr.PowerShieldHP)
	case Field_PowerShieldRecoverSpeed:
		return float64(attr.PowerShieldRecoverSpeed)
	case Field_TopAttackElementPlus:
		return attr.TopAttackElementPlus
	case Field_OverDriveLimit:
		return float64(attr.OverDriveLimit)
	case Field_OverDriveAddEfficiency:
		return float64(attr.OverDriveAddEfficiency)
	case Field_OverDriveTime:
		return float64(attr.OverDriveTime)
	case Field_WeakTime:
		return float64(attr.WeakTime)
	case Field_BreakValueLimit:
		return float64(attr.BreakValueLimit)
	case Field_BreakStateTime:
		return float64(attr.BreakStateTime)

	// 装备属性
	case Field_EquipMinAttack:
		return attr.EquipMinAttack
	case Field_EquipMaxAttack:
		return attr.EquipMaxAttack
	case Field_EquipAttackRate:
		return attr.EquipAttackRate
	case Field_EquipAttackNormalPlus:
		return attr.EquipAttackNormalPlus
	case Field_EquipAttackFirePlus:
		return attr.EquipAttackFirePlus
	case Field_EquipAttackColdPlus:
		return attr.EquipAttackColdPlus
	case Field_EquipAttackPoisonPlus:
		return attr.EquipAttackPoisonPlus
	case Field_EquipAttackLightningPlus:
		return attr.EquipAttackLightningPlus
	case Field_AttackNormalPlus:
		return attr.AttackNormalPlus
	case Field_AttackFirePlus:
		return attr.AttackFirePlus
	case Field_AttackColdPlus:
		return attr.AttackColdPlus
	case Field_AttackPoisonPlus:
		return attr.AttackPoisonPlus
	case Field_AttackLightningPlus:
		return attr.AttackLightningPlus

	// 逻辑属性
	case Field_CollisionRadius:
		return float64(attr.CollisionRadius)
	case Field_Mass:
		return float64(attr.Mass)
	case Field_Scale:
		return float64(attr.Scale)
	}

	return 0
}
