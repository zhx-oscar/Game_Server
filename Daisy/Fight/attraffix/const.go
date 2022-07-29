package attraffix

// Field 属性字段
type Field uint32

// 人物属性
const (
	Field_Prop_Begin Field = 1

	Field_Strength Field = iota
	Field_Agility
	Field_Intelligence
	Field_Vitality
	Field_MaxHP
	Field_SkillPowerLimit
	Field_Attack
	Field_AttackNoramlAdd
	Field_AttackNoramlDecr
	Field_AttackFireAdd
	Field_AttackFireDecr
	Field_AttackColdAdd
	Field_AttackColdDecr
	Field_AttackPoisonAdd
	Field_AttackPoisonDecr
	Field_AttackLightningAdd
	Field_AttackLightningDecr
	Field_HitRate
	Field_DodgeRate
	Field_CritRate
	Field_BlockRate
	Field_BlockValue
	Field_CritDamageRate
	Field_Armor
	Field_ResistanceFire
	Field_ResistanceCold
	Field_ResistancePoison
	Field_ResistanceLightning
	Field_BeDamageNormalAdd
	Field_BeDamageNormalDecr
	Field_BeDamageFireAdd
	Field_BeDamageFireDecr
	Field_BeDamageColdAdd
	Field_BeDamageColdDecr
	Field_BeDamagePoisonAdd
	Field_BeDamagePoisonDecr
	Field_BeDamageLightningAdd
	Field_BeDamageLightningDecr
	Field_Lucky
	Field_NormalAttackSpeed
	Field_AttackLucky
	Field_AttackBloodsuckerRate
	Field_AttackSputteringRate
	Field_AttackStealUltimateSkillPower
	Field_BeDamageNormalDeduct
	Field_BeDamageThorns
	Field_RecoverHP
	Field_RecoverUltimateSkillPowerRate
	Field_PowerShieldHP
	Field_PowerShieldRecoverSpeed
	Field_TopAttackElementPlus
	Field_OverDriveLimit
	Field_OverDriveAddEfficiency
	Field_OverDriveTime
	Field_WeakTime
	Field_BreakValueLimit
	Field_BreakStateTime

	Field_Prop_End
)

// 装备属性
const (
	Field_Equip_Begin Field = 1001

	Field_EquipMinAttack Field = iota + 1000
	Field_EquipMaxAttack
	Field_EquipAttackRate
	Field_EquipAttackNormalPlus
	Field_EquipAttackFirePlus
	Field_EquipAttackColdPlus
	Field_EquipAttackPoisonPlus
	Field_EquipAttackLightningPlus
	Field_AttackNormalPlus
	Field_AttackFirePlus
	Field_AttackColdPlus
	Field_AttackPoisonPlus
	Field_AttackLightningPlus

	Field_Equip_End
)

// 逻辑属性
const (
	Field_Logic_Begin Field = 2001

	Field_CollisionRadius Field = iota + 2000
	Field_Mass
	Field_Scale

	Field_Logic_End
)

// AttrAffix 属性词缀
type AttrAffix struct {
	Field Field   // 属性字段
	ParaA float32 // 参数a
	ParaB float64 // 参数b
}
