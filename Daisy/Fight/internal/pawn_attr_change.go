package internal

import (
	. "Daisy/Fight/attraffix"
	"Daisy/Proto"
	"math"
	"sort"
)

// _AttrChanged 属性变化中间值
type _AttrChanged struct {
	Second   float64   // 二级属性附加数值
	ParaA    float32   // 参数a
	ParaB    float64   // 参数b
	Vals     []float64 // 数值表
	Override bool      // 是否覆盖
}

// needSave 是否需要保存
func (attrChanged *_AttrChanged) needSave() bool {
	return attrChanged.Second != 0 || attrChanged.ParaA != 0 || attrChanged.ParaB != 0 || len(attrChanged.Vals) > 0 || attrChanged.Override
}

// overlap 叠加数值
func (attrChanged *_AttrChanged) overlap(overlapCategory *_AttrAffixOverlapCategory, paraA float32, paraB float64, sign, reset bool) {
	switch overlapCategory.Type {
	case _AttrAffixOverlapType_BRF:
		if reset {
			attrChanged.ParaA = 0
			attrChanged.ParaB = 0
		} else {
			if !sign {
				paraA = -paraA
				paraB = -paraB
			}
			attrChanged.ParaA += paraA
			attrChanged.ParaB += paraB
		}

	case _AttrAffixOverlapType_Addition:
		if reset {
			attrChanged.ParaB = 0
		} else {
			if !sign {
				paraB = -paraB
			}
			attrChanged.ParaB += paraB
		}

	case _AttrAffixOverlapType_ReverseMCL:
		if reset {
			attrChanged.Vals = nil
			attrChanged.ParaB = 1
		} else {
			if sign {
				attrChanged.Vals = append(attrChanged.Vals, paraB)
			} else {
				for i := len(attrChanged.Vals) - 1; i >= 0; i-- {
					if FloatEqual(attrChanged.Vals[i], paraB) {
						attrChanged.Vals = append(attrChanged.Vals[:i], attrChanged.Vals[i+1:]...)
						break
					}
				}
			}
			attrChanged.ParaB = 1
			for _, v := range attrChanged.Vals {
				attrChanged.ParaB *= 1 - v/overlapCategory.ConstantTab[0]
			}
		}

	case _AttrAffixOverlapType_Max:
		if reset {
			attrChanged.Vals = nil
		} else {
			if sign {
				attrChanged.Vals = append(attrChanged.Vals, paraB)
				sort.Float64s(attrChanged.Vals)
			} else {
				for i := len(attrChanged.Vals) - 1; i >= 0; i-- {
					if FloatEqual(attrChanged.Vals[i], paraB) {
						attrChanged.Vals = append(attrChanged.Vals[:i], attrChanged.Vals[i+1:]...)
						break
					}
				}
			}
		}

	case _AttrAffixOverlapType_Override:
		if reset {
			attrChanged.Override = false
			attrChanged.ParaB = 0
		} else {
			if sign {
				attrChanged.Override = true
				attrChanged.ParaB = paraB
			} else {
				attrChanged.Override = false
				attrChanged.ParaB = 0
			}
		}
	}
}

// _AttrChangedTab 属性变化中间值表
type _AttrChangedTab map[Field]*_AttrChanged

// ChangeAttr 修改属性数值
func (attr *FightAttr) ChangeAttr(field Field, paraA float32, paraB float64, sign bool) {
	attr.changeAttr(field, paraA, paraB, sign, false)
}

// OverrideAttr 覆盖属性数值
func (attr *FightAttr) OverrideAttr(field Field, value float64) {
	overlapCategory, ok := attrAffixOverlapCategoryTab[field]
	if !ok {
		return
	}

	if overlapCategory.Type != _AttrAffixOverlapType_Override {
		return
	}

	attr.changeAttr(field, 0, value, true, false)
}

// ResetAttr 重置属性数值
func (attr *FightAttr) ResetAttr(field Field) {
	attr.changeAttr(field, 0, 0, false, true)
}

// changeAttr 修改属性数值
func (attr *FightAttr) changeAttr(field Field, paraA float32, paraB float64, sign, reset bool) {
	if attr.inFight && attr.pawn.isSnapshot {
		return
	}

	// 查询属性词缀重叠处理策略
	overlapCategory, ok := attrAffixOverlapCategoryTab[field]
	if !ok {
		return
	}

	// 获取属性变化中间值
	attrChanged := attr.getAttrChanged(field)

	// 叠加属性变化数值
	attrChanged.overlap(overlapCategory, paraA, paraB, sign, reset)

	// 修改属性数值
	attr.changeFirstAttr(field, overlapCategory, attrChanged)

	// 检测中间值是否需要保存
	if !attrChanged.needSave() {
		delete(attr.changedTab, field)
	}
}

// changeFirstAttr 修改一级属性数值
func (attr *FightAttr) changeFirstAttr(field Field, overlapCategory *_AttrAffixOverlapCategory, attrChanged *_AttrChanged) {
	if attrChanged == nil {
		return
	}

	oldValue := attr.GetAttr(field)
	attr.assignAttr(field, Min(Max(calcAttrValue(attr.base.GetAttr(field), overlapCategory, attrChanged), overlapCategory.Min), overlapCategory.Max))
	if overlapCategory.PostFun != nil {
		overlapCategory.PostFun(attr, oldValue, attr.GetAttr(field))
	}

	// 修改二级属性
	attr.changeSecondAttr(field, attrChanged)
}

// changeSecondAttr 修改二级属性数值
func (attr *FightAttr) changeSecondAttr(firstField Field, firstAttrChanged *_AttrChanged) {
	changeFun := func(field Field, value float64) {
		attrChanged := attr.getAttrChanged(field)
		attrChanged.Second = value

		overlapCategory, ok := attrAffixOverlapCategoryTab[field]
		if ok {
			attr.changeFirstAttr(field, overlapCategory, attrChanged)
		}
	}

	switch firstField {
	case Field_Strength:
		changeFun(Field_Attack, attr.Strength)
	case Field_Agility:
		changeFun(Field_HitRate, attr.Agility*0.001)
	case Field_Intelligence:
		changeFun(Field_TopAttackElementPlus, attr.Intelligence)
	case Field_Vitality:
		changeFun(Field_MaxHP, attr.Vitality*10)
	case Field_Attack:
		defAttack := attr.base.Attack + firstAttrChanged.Second
		changeFun(Field_Armor, (defAttack-0.5*defAttack)/0.5)
	case Field_Scale:
		attr.ChangeAttr(Field_CollisionRadius, 0, float64(attr.base.CollisionRadius*attr.Scale), true)
	}
}

// getAttrChanged 获取属性变化中间值
func (attr *FightAttr) getAttrChanged(field Field) *_AttrChanged {
	attrChanged, ok := attr.changedTab[field]
	if !ok {
		attrChanged = &_AttrChanged{}
		attr.changedTab[field] = attrChanged
	}
	return attrChanged
}

// AttrSyncToClient 属性同步客户端
func (attr *FightAttr) AttrSyncToClient(attrType Proto.AttrType_Enum, old, new float64) {
	if attr.inFight && !FloatEqual(new, old) {
		attr.pawn.Scene.PushAction(&Proto.ChangeAttr{
			SelfId:   attr.pawn.UID,
			AttrType: attrType,
			OldValue: old,
			NewValue: new,
		})
	}
}

// AttrSyncToClientDebug 属性同步客户端debug
func (attr *FightAttr) AttrSyncToClientDebug(attrType Proto.AttrType_Enum, old, new float64) {
	if attr.inFight && !FloatEqual(new, old) {
		attr.pawn.Scene.PushDebugAction(&Proto.ChangeAttr{
			SelfId:   attr.pawn.UID,
			AttrType: attrType,
			OldValue: old,
			NewValue: new,
		})
	}
}

// calcAttrValue 计算属性最终数值
func calcAttrValue(baseVal float64, overlapCategory *_AttrAffixOverlapCategory, attrChanged *_AttrChanged) float64 {
	value := 0.0

	switch overlapCategory.Type {
	case _AttrAffixOverlapType_BRF:
		value = (baseVal+attrChanged.Second)*float64(1+attrChanged.ParaA) + attrChanged.ParaB

	case _AttrAffixOverlapType_Addition:
		value = baseVal + attrChanged.Second + attrChanged.ParaB

	case _AttrAffixOverlapType_ReverseMCL:
		value = (1 - (1-(baseVal+attrChanged.Second)/overlapCategory.ConstantTab[0])*attrChanged.ParaB) * overlapCategory.ConstantTab[0]

	case _AttrAffixOverlapType_Max:
		value = baseVal

		for _, v := range attrChanged.Vals {
			if v > baseVal {
				value = v
			}
		}

	case _AttrAffixOverlapType_Override:
		if attrChanged.Override {
			value = attrChanged.ParaB
		} else {
			value = baseVal
		}
	}

	if overlapCategory.Precision > 1 {
		value = math.Round(value*overlapCategory.Precision) / overlapCategory.Precision
	} else {
		value = math.Round(value)
	}

	return value
}
