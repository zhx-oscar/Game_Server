package internal

import (
	. "Daisy/Fight/attraffix"
	"Daisy/Proto"
)

// _AttrAffixOverlapType 属性词缀重叠处理方式
type _AttrAffixOverlapType uint32

const (
	_AttrAffixOverlapType_BRF        _AttrAffixOverlapType = iota // BRF公式
	_AttrAffixOverlapType_Addition                                // 加法
	_AttrAffixOverlapType_ReverseMCL                              // 反向乘法
	_AttrAffixOverlapType_Max                                     // 最高生效
	_AttrAffixOverlapType_Override                                // 覆盖
)

// AttrMaxValue float64能精确表示的最大值 2的52次方减1
const AttrMaxValue float64 = 4503599627370495

// _AttrAffixOverlapCategory 属性词缀重叠处理策略
type _AttrAffixOverlapCategory struct {
	Type        _AttrAffixOverlapType                   // 重叠处理方式
	ConstantTab []float64                               // 常量表
	Precision   float64                                 // 精度
	Max         float64                                 // 属性最大值
	Min         float64                                 // 属性最小值
	PostFun     func(attr *FightAttr, old, new float64) // 后处理函数
}

// attrAffixOverlapCategoryTab 属性词缀重叠处理方式表
var attrAffixOverlapCategoryTab = map[Field]*_AttrAffixOverlapCategory{
	// 人物属性
	Field_Strength:     {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},
	Field_Agility:      {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},
	Field_Intelligence: {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},
	Field_Vitality:     {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},
	Field_MaxHP: {
		Type: _AttrAffixOverlapType_BRF,
		Max:  AttrMaxValue,
		Min:  1,
		PostFun: func(attr *FightAttr, old, new float64) {
			attr.AttrSyncToClient(Proto.AttrType_MaxHP, old, new)
		},
	},
	Field_SkillPowerLimit:       {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},
	Field_Attack:                {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},
	Field_AttackNoramlAdd:       {Type: _AttrAffixOverlapType_Addition, Precision: 10000, Max: AttrMaxValue, Min: 0},
	Field_AttackNoramlDecr:      {Type: _AttrAffixOverlapType_Addition, Precision: 10000, Max: 0, Min: -AttrMaxValue},
	Field_AttackFireAdd:         {Type: _AttrAffixOverlapType_Addition, Precision: 10000, Max: AttrMaxValue, Min: 0},
	Field_AttackFireDecr:        {Type: _AttrAffixOverlapType_Addition, Precision: 10000, Max: 0, Min: -AttrMaxValue},
	Field_AttackColdAdd:         {Type: _AttrAffixOverlapType_Addition, Precision: 10000, Max: AttrMaxValue, Min: 0},
	Field_AttackColdDecr:        {Type: _AttrAffixOverlapType_Addition, Precision: 10000, Max: 0, Min: -AttrMaxValue},
	Field_AttackPoisonAdd:       {Type: _AttrAffixOverlapType_Addition, Precision: 10000, Max: AttrMaxValue, Min: 0},
	Field_AttackPoisonDecr:      {Type: _AttrAffixOverlapType_Addition, Precision: 10000, Max: 0, Min: -AttrMaxValue},
	Field_AttackLightningAdd:    {Type: _AttrAffixOverlapType_Addition, Precision: 10000, Max: AttrMaxValue, Min: 0},
	Field_AttackLightningDecr:   {Type: _AttrAffixOverlapType_Addition, Precision: 10000, Max: 0, Min: -AttrMaxValue},
	Field_HitRate:               {Type: _AttrAffixOverlapType_ReverseMCL, ConstantTab: []float64{1}, Precision: 10000, Max: 1, Min: 0},
	Field_DodgeRate:             {Type: _AttrAffixOverlapType_ReverseMCL, ConstantTab: []float64{1}, Precision: 10000, Max: 1, Min: 0},
	Field_CritRate:              {Type: _AttrAffixOverlapType_ReverseMCL, ConstantTab: []float64{1}, Precision: 10000, Max: 1, Min: 0},
	Field_BlockRate:             {Type: _AttrAffixOverlapType_ReverseMCL, ConstantTab: []float64{1}, Precision: 10000, Max: 1, Min: 0},
	Field_BlockValue:            {Type: _AttrAffixOverlapType_Max, Max: AttrMaxValue, Min: 0},
	Field_CritDamageRate:        {Type: _AttrAffixOverlapType_Addition, Precision: 10000, Max: AttrMaxValue, Min: 1},
	Field_Armor:                 {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},
	Field_ResistanceFire:        {Type: _AttrAffixOverlapType_ReverseMCL, ConstantTab: []float64{1}, Precision: 10000, Max: 1, Min: 0},
	Field_ResistanceCold:        {Type: _AttrAffixOverlapType_ReverseMCL, ConstantTab: []float64{1}, Precision: 10000, Max: 1, Min: 0},
	Field_ResistancePoison:      {Type: _AttrAffixOverlapType_ReverseMCL, ConstantTab: []float64{1}, Precision: 10000, Max: 1, Min: 0},
	Field_ResistanceLightning:   {Type: _AttrAffixOverlapType_ReverseMCL, ConstantTab: []float64{1}, Precision: 10000, Max: 1, Min: 0},
	Field_BeDamageNormalAdd:     {Type: _AttrAffixOverlapType_Addition, Precision: 10000, Max: AttrMaxValue, Min: 0},
	Field_BeDamageNormalDecr:    {Type: _AttrAffixOverlapType_Addition, Precision: 10000, Max: 0, Min: -AttrMaxValue},
	Field_BeDamageFireAdd:       {Type: _AttrAffixOverlapType_Addition, Precision: 10000, Max: AttrMaxValue, Min: 0},
	Field_BeDamageFireDecr:      {Type: _AttrAffixOverlapType_Addition, Precision: 10000, Max: 0, Min: -AttrMaxValue},
	Field_BeDamageColdAdd:       {Type: _AttrAffixOverlapType_Addition, Precision: 10000, Max: AttrMaxValue, Min: 0},
	Field_BeDamageColdDecr:      {Type: _AttrAffixOverlapType_Addition, Precision: 10000, Max: 0, Min: -AttrMaxValue},
	Field_BeDamagePoisonAdd:     {Type: _AttrAffixOverlapType_Addition, Precision: 10000, Max: AttrMaxValue, Min: 0},
	Field_BeDamagePoisonDecr:    {Type: _AttrAffixOverlapType_Addition, Precision: 10000, Max: 0, Min: -AttrMaxValue},
	Field_BeDamageLightningAdd:  {Type: _AttrAffixOverlapType_Addition, Precision: 10000, Max: AttrMaxValue, Min: 0},
	Field_BeDamageLightningDecr: {Type: _AttrAffixOverlapType_Addition, Precision: 10000, Max: 0, Min: -AttrMaxValue},
	Field_Lucky:                 {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},
	Field_NormalAttackSpeed: {
		Type:        _AttrAffixOverlapType_ReverseMCL,
		ConstantTab: []float64{2},
		Precision:   100,
		Max:         2,
		Min:         0.1,
		PostFun: func(attr *FightAttr, old, new float64) {
			attr.AttrSyncToClient(Proto.AttrType_NormalAttackSpeed, old, new)
		},
	},
	Field_AttackLucky:                   {Type: _AttrAffixOverlapType_Addition, Precision: 10000, Max: 1, Min: 0},
	Field_AttackBloodsuckerRate:         {Type: _AttrAffixOverlapType_Addition, Precision: 10000, Max: AttrMaxValue, Min: 0},
	Field_AttackSputteringRate:          {Type: _AttrAffixOverlapType_Addition, Precision: 10000, Max: AttrMaxValue, Min: 0},
	Field_AttackStealUltimateSkillPower: {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},
	Field_BeDamageNormalDeduct:          {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},
	Field_BeDamageThorns:                {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},
	Field_RecoverHP:                     {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},
	Field_RecoverUltimateSkillPowerRate: {Type: _AttrAffixOverlapType_Addition, Precision: 10000, Max: 1, Min: 0},
	Field_PowerShieldHP:                 {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},
	Field_PowerShieldRecoverSpeed:       {Type: _AttrAffixOverlapType_Addition, Precision: 10000, Max: 1, Min: 0},
	Field_TopAttackElementPlus:          {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},
	Field_OverDriveLimit:                {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},
	Field_OverDriveAddEfficiency:        {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},
	Field_OverDriveTime:                 {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},
	Field_WeakTime:                      {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},
	Field_BreakValueLimit:               {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},
	Field_BreakStateTime:                {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},

	// 装备属性
	Field_EquipMinAttack:           {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},
	Field_EquipMaxAttack:           {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},
	Field_EquipAttackRate:          {Type: _AttrAffixOverlapType_BRF, Precision: 10000, Max: AttrMaxValue, Min: 0},
	Field_EquipAttackNormalPlus:    {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},
	Field_EquipAttackFirePlus:      {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},
	Field_EquipAttackColdPlus:      {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},
	Field_EquipAttackPoisonPlus:    {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},
	Field_EquipAttackLightningPlus: {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},
	Field_AttackNormalPlus:         {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},
	Field_AttackFirePlus:           {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},
	Field_AttackColdPlus:           {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},
	Field_AttackPoisonPlus:         {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},
	Field_AttackLightningPlus:      {Type: _AttrAffixOverlapType_BRF, Max: AttrMaxValue, Min: 0},

	// 逻辑属性
	Field_CollisionRadius: {
		Type:      _AttrAffixOverlapType_Override,
		Precision: 100,
		Max:       AttrMaxValue,
		Min:       0.1,
		PostFun: func(attr *FightAttr, old, new float64) {
			if attr.inFight {
				attr.pawn.Scene.setPawnShapeRadius(attr.pawn.UID, new)
			}
			attr.AttrSyncToClientDebug(Proto.AttrType_CollisionRadius, old, new)
		},
	},
	Field_Mass: {Type: _AttrAffixOverlapType_Override, Max: AttrMaxValue, Min: 0},
	Field_Scale: {
		Type:      _AttrAffixOverlapType_Addition,
		Precision: 100,
		Max:       AttrMaxValue,
		Min:       0.1,
		PostFun: func(attr *FightAttr, old, new float64) {
			attr.AttrSyncToClient(Proto.AttrType_Scale, old, new)
		},
	},
}
