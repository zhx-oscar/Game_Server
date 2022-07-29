package internal

import (
	. "Daisy/Fight/attraffix"
	"Daisy/Fight/internal/conf"
	"errors"
)

// FightAttr 战斗属性
type FightAttr struct {
	pawn *Pawn

	// 标记
	inFight bool // 在战斗中，不是用于属性计算

	// 属性
	base       *_Attr          // 基础属性值
	changedTab _AttrChangedTab // 属性变化中间值表
	*_Attr                     // 当前属性值

	// 血量与护盾
	CurHP        int64       // 当前血量
	AllHPShield  int64       // 血量护盾
	HPShieldList []*HPShield // 血量护盾列表

	// 必杀技能量
	UltimateSkillPower int32 // 必杀技能量
}

// init 初始化
func (attr *FightAttr) init(pawn *Pawn) error {
	attr.pawn = pawn
	attr.inFight = true

	// 初始化属性
	if err := attr.initAttr(pawn.Info.PawnConfig, pawn.Info.Level, pawn.Info.AttrAffixList); err != nil {
		return err
	}

	// 初始化HP
	attr.ChangeHP(attr.MaxHP)

	return nil
}

// initAttr 初始化属性
func (attr *FightAttr) initAttr(pawnConf *conf.PawnConfig, level int32, attrAffixList []AttrAffix) error {
	if pawnConf == nil {
		return errors.New("nil pawnConf")
	}

	// 属性基础值
	attr.base = &_Attr{
		PawnAttr: &PawnAttr{},
		ExtendAttr: &ExtendAttr{
			CollisionRadius: pawnConf.CollisionRadius,
			Mass:            pawnConf.Mass,
			Scale:           pawnConf.Scale,
		},
	}
	*attr.base.PawnAttr = *pawnConf.PropValue

	// 成长一级属性数值
	attr.base.growUpFirstProp(level)

	// 模拟器模式下 基础属性处理
	attr.updateBaseAttrInSimulatorMode()

	// 属性变化中间值表
	attr.changedTab = _AttrChangedTab{}

	// 属性当前值
	attr._Attr = &_Attr{
		PawnAttr: &PawnAttr{
			ID:   attr.base.ID,
			Type: attr.base.Type,
		},
		ExtendAttr: &ExtendAttr{},
	}

	// 初始化属性当前值
	for i := Field_Prop_Begin; i < Field_Prop_End; i++ {
		attr.ResetAttr(i)
	}

	for i := Field_Equip_Begin; i < Field_Equip_End; i++ {
		attr.ResetAttr(i)
	}

	for i := Field_Logic_Begin; i < Field_Logic_End; i++ {
		attr.ResetAttr(i)
	}

	// 初始化词缀附加属性
	for _, v := range attrAffixList {
		attr.ChangeAttr(v.Field, v.ParaA, v.ParaB, true)
	}

	return nil
}

//updateBaseAttrInSimulatorMode 模拟器下 基础属性处理
func (attr FightAttr) updateBaseAttrInSimulatorMode() {
	if !attr.inFight || !attr.pawn.Scene.SimulatorMode() || attr.pawn.Info.SimulatorModeInfo == nil {
		return
	}

	if attr.pawn.Info.SimulatorModeInfo.MaxHP != 0 {
		attr.base.MaxHP = attr.pawn.Info.SimulatorModeInfo.MaxHP
	}

	if attr.pawn.Info.SimulatorModeInfo.Attack != 0 {
		attr.base.Attack = attr.pawn.Info.SimulatorModeInfo.Attack
	}

	attr.base.DodgeRate += attr.pawn.Info.SimulatorModeInfo.ExtendDodgeRate
	attr.base.BlockRate += attr.pawn.Info.SimulatorModeInfo.ExtendBlockRate
	attr.base.HitRate += attr.pawn.Info.SimulatorModeInfo.ExtendHitRate
	attr.base.CritRate += attr.pawn.Info.SimulatorModeInfo.ExtendCritRate
}

// copy 拷贝
func (attr FightAttr) copy(pawn *Pawn) *FightAttr {
	attr.pawn = pawn

	// 清除基础值
	attr.base = nil

	// 清除属性变化中间值
	attr.changedTab = nil

	// 拷贝当前属性值
	attr._Attr = attr._Attr.copy()

	// 清除护盾数据
	attr.HPShieldList = nil

	return &attr
}
