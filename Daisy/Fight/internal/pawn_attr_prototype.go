package internal

import (
	"Daisy/DataTables"
	"Daisy/Fight/internal/conf"
	"math"
)

// PawnAttr 基本属性
type PawnAttr = DataTables.PropValue_Config

// ExtendAttr 扩展属性
type ExtendAttr struct {
	CollisionRadius float32 // 碰撞半径
	Mass            uint32  // 质量id
	Scale           float32 // 缩放比例

	BlockUnbalanceDegrade int32 // 格挡失衡值每秒扣除值
	CrackUp_Block         int32 // 格挡失衡值上限
}

// _Attr 属性
type _Attr struct {
	*PawnAttr   // 基础属性
	*ExtendAttr // 扩展属性
}

// copy 拷贝
func (attr _Attr) copy() *_Attr {
	tPawnAttr := &PawnAttr{}
	*tPawnAttr = *attr.PawnAttr
	attr.PawnAttr = tPawnAttr

	tEquipAttr := &ExtendAttr{}
	*tEquipAttr = *attr.ExtendAttr
	attr.ExtendAttr = tEquipAttr

	return &attr
}

// growUpFirstProp 成长一级属性数值
func (attr *_Attr) growUpFirstProp(level int32) {
	// 等级系数
	levelCoe := float64(level - 1)
	if levelCoe < 0 {
		levelCoe = 0
	}

	// 一级属性系数
	firstCoe := attr.Strength + attr.Agility + attr.Intelligence + attr.Vitality
	if firstCoe <= 0 {
		return
	}

	attr.Strength = attr.Strength + (levelCoe*firstCoe*0.15+math.Pow(levelCoe, 2.5)*0.1)*attr.Strength/firstCoe
	attr.Agility = attr.Agility + (levelCoe*firstCoe*0.15+levelCoe*0.1)*attr.Agility/firstCoe
	attr.Intelligence = attr.Intelligence + (levelCoe*firstCoe*0.15+math.Pow(levelCoe, 2.5)*0.1)*attr.Intelligence/firstCoe
	attr.Vitality = attr.Vitality + (levelCoe*firstCoe*0.15+math.Pow(levelCoe, 2.5)*0.1)*attr.Vitality/firstCoe
}

// GetEquipAttackPlus 获取武器附加攻击力
func (attr *_Attr) GetEquipAttackPlus(damageValueKind conf.DamageValueKind) float64 {
	switch damageValueKind {
	case conf.DamageValueKind_Normal:
		return attr.EquipAttackNormalPlus
	case conf.DamageValueKind_Fire:
		return attr.EquipAttackFirePlus
	case conf.DamageValueKind_Cold:
		return attr.EquipAttackColdPlus
	case conf.DamageValueKind_Poison:
		return attr.EquipAttackPoisonPlus
	case conf.DamageValueKind_Lightning:
		return attr.EquipAttackLightningPlus
	}
	return 0
}

// GetAttackPlus 获取非武器附加攻击力
func (attr *_Attr) GetAttackPlus(damageValueKind conf.DamageValueKind) float64 {
	switch damageValueKind {
	case conf.DamageValueKind_Normal:
		return attr.AttackNormalPlus
	case conf.DamageValueKind_Fire:
		return attr.AttackFirePlus
	case conf.DamageValueKind_Cold:
		return attr.AttackColdPlus
	case conf.DamageValueKind_Poison:
		return attr.AttackPoisonPlus
	case conf.DamageValueKind_Lightning:
		return attr.AttackLightningPlus
	}
	return 0
}

// GetAttackAdd 获取非武器伤害加深
func (attr *_Attr) GetAttackAdd(damageValueKind conf.DamageValueKind) float64 {
	switch damageValueKind {
	case conf.DamageValueKind_Normal:
		return attr.AttackNoramlAdd
	case conf.DamageValueKind_Fire:
		return attr.AttackFireAdd
	case conf.DamageValueKind_Cold:
		return attr.AttackColdAdd
	case conf.DamageValueKind_Poison:
		return attr.AttackPoisonAdd
	case conf.DamageValueKind_Lightning:
		return attr.AttackLightningAdd
	}
	return 0
}

// GetAttackDecr 获取非武器伤害降低
func (attr *_Attr) GetAttackDecr(damageValueKind conf.DamageValueKind) float64 {
	switch damageValueKind {
	case conf.DamageValueKind_Normal:
		return attr.AttackNoramlDecr
	case conf.DamageValueKind_Fire:
		return attr.AttackFireDecr
	case conf.DamageValueKind_Cold:
		return attr.AttackColdDecr
	case conf.DamageValueKind_Poison:
		return attr.AttackPoisonDecr
	case conf.DamageValueKind_Lightning:
		return attr.AttackLightningDecr
	}
	return 0
}

// GetResistance 获取抗性
func (attr *_Attr) GetResistance(damageValueKind conf.DamageValueKind) float32 {
	switch damageValueKind {
	case conf.DamageValueKind_Fire:
		return attr.ResistanceFire
	case conf.DamageValueKind_Cold:
		return attr.ResistanceCold
	case conf.DamageValueKind_Poison:
		return attr.ResistancePoison
	case conf.DamageValueKind_Lightning:
		return attr.ResistanceLightning
	}
	return 0
}

// GetBeDamageAdd 获取受击伤害加深
func (attr *_Attr) GetBeDamageAdd(damageValueKind conf.DamageValueKind) float64 {
	switch damageValueKind {
	case conf.DamageValueKind_Normal:
		return attr.BeDamageNormalAdd
	case conf.DamageValueKind_Fire:
		return attr.BeDamageFireAdd
	case conf.DamageValueKind_Cold:
		return attr.BeDamageColdAdd
	case conf.DamageValueKind_Poison:
		return attr.BeDamagePoisonAdd
	case conf.DamageValueKind_Lightning:
		return attr.BeDamageLightningAdd
	}
	return 0
}

// GetBeDamageDecr 获取受击伤害降低
func (attr *_Attr) GetBeDamageDecr(damageValueKind conf.DamageValueKind) float64 {
	switch damageValueKind {
	case conf.DamageValueKind_Normal:
		return attr.BeDamageNormalDecr
	case conf.DamageValueKind_Fire:
		return attr.BeDamageFireDecr
	case conf.DamageValueKind_Cold:
		return attr.BeDamageColdDecr
	case conf.DamageValueKind_Poison:
		return attr.BeDamagePoisonDecr
	case conf.DamageValueKind_Lightning:
		return attr.BeDamageLightningDecr
	}
	return 0
}

// GetBeDamageDeduct 获取受击伤害抵挡
func (attr *_Attr) GetBeDamageDeduct(damageValueKind conf.DamageValueKind) float64 {
	switch damageValueKind {
	case conf.DamageValueKind_Normal:
		return attr.BeDamageNormalDeduct
	}
	return 0
}
