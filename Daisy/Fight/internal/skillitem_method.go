package internal

import (
	"Cinder/Base/linemath"
	"Daisy/Fight/internal/conf"
)

// IsNormalAttack 是否是普攻
func (skillItem *_SkillItem) IsNormalAttack() bool {
	return skillItem.Config.SkillKind == conf.SkillKind_NormalAtk
}

// IsSuperSkill 是否是超能技
func (skillItem *_SkillItem) IsSuperSkill() bool {
	return skillItem.Config.SkillKind == conf.SkillKind_Super
}

// IsUltimateSkill 是否是必杀技
func (skillItem *_SkillItem) IsUltimateSkill() bool {
	return skillItem.Config.SkillKind == conf.SkillKind_Ultimate
}

// IsCombineSkill 是否是合体必杀技
func (skillItem *_SkillItem) IsCombineSkill() bool {
	return skillItem.Config.SkillKind == conf.SkillKind_Combine
}

// PosInAttackRange 目标点是否在攻击范围内（包含冲刺距离）
func (skillItem *_SkillItem) PosInAttackRange(pos linemath.Vector2) bool {
	if skillItem.GetMaxCastDistance() <= 0 {
		return true
	}

	dis := Distance(skillItem.Caster.GetPos(), pos)
	minDis := skillItem.GetMinCastDistance()
	maxDis := skillItem.GetMaxCastDistance() + skillItem.Caster.Attr.CollisionRadius

	if skillItem.Config.CanDash {
		minDis = skillItem.GetMinCastDistance()
		maxDis = skillItem.GetMaxCastDistance() + skillItem.Caster.Info.FastDist + skillItem.Caster.Attr.CollisionRadius
	}

	return dis >= minDis && dis <= maxDis
}

// TargetInAttackRange 目标对象是否在攻击范围内（包含冲刺距离）
func (skillItem *_SkillItem) TargetInAttackRange(target *Pawn) bool {
	if target == nil {
		return false
	}

	if skillItem.GetMaxCastDistance() <= 0 {
		return true
	}

	//自己选择自己作为目标 不做攻击距离检测
	if skillItem.Caster.UID == target.UID {
		return true
	}

	dis := DistancePawn(skillItem.Caster, target)
	minDis := skillItem.GetMinCastDistance()
	maxDis := skillItem.Caster.Attr.CollisionRadius + target.Attr.CollisionRadius + skillItem.GetMaxCastDistance()

	if skillItem.Config.CanDash {
		minDis = skillItem.GetMinCastDistance()
		maxDis = skillItem.Caster.Attr.CollisionRadius + target.Attr.CollisionRadius + skillItem.GetMaxCastDistance() + skillItem.Caster.Info.FastDist
	}
	var result bool

	result = dis <= maxDis
	if skillItem.GetMinCastDistance() > 0 {
		result = dis >= minDis && dis <= maxDis
	}

	return result
}

// GetAttackRange 获取技能攻击距离 最小攻击距离 最大攻击距离
func (skillItem *_SkillItem) GetAttackRange() (float32, float32) {
	if skillItem.Config.CanDash {
		return skillItem.GetMinCastDistance() + skillItem.Caster.Info.FastDist, skillItem.GetMaxCastDistance() + skillItem.Caster.Info.FastDist
	}

	return skillItem.GetMinCastDistance(), skillItem.GetMaxCastDistance()
}

// TargetInCastDistance 目标是否在施法距离内
func (skillItem *_SkillItem) TargetInCastDistance(target *Pawn) bool {
	if target == nil {
		return false
	}

	if skillItem.GetMaxCastDistance() <= 0 {
		return true
	}

	return DistancePawn(skillItem.Caster, target) <= skillItem.GetMaxCastDistance()+skillItem.Caster.Attr.CollisionRadius+target.Attr.CollisionRadius
}

// PosInCastDistance 目标点是否在施法距离内
func (skillItem *_SkillItem) PosInCastDistance(pos linemath.Vector2) bool {
	if skillItem.GetMaxCastDistance() <= 0 {
		return true
	}

	return Distance(skillItem.Caster.GetPos(), pos) <= skillItem.GetMaxCastDistance()+skillItem.Caster.Attr.CollisionRadius
}

// GetAttackType 获取伤害体类型
func (skillItem *_SkillItem) GetAttackType() conf.AttackType {
	if len(skillItem.Config.TemplateConfig.AttackConfs) <= 0 {
		return conf.AttackType_Single
	}

	return skillItem.Config.TemplateConfig.AttackConfs[0].Type
}

// SearchTargets 查询技能目标列表
func (skillItem *_SkillItem) SearchTargets(targetPos linemath.Vector2) []*Pawn {
	return skillItem.Caster.Scene.searchSkillTargets(skillItem, targetPos)
}

// GetAttackSpeed 获取技能攻击速度
func (skillItem *_SkillItem) GetAttackSpeed() float32 {
	switch skillItem.Config.SkillKind {
	case conf.SkillKind_NormalAtk:
		return skillItem.Caster.Attr.NormalAttackSpeed
	default:
		return 1
	}
}

// ZoomAttackTime 缩放技能攻击时间
func (skillItem *_SkillItem) ZoomAttackTime(time uint32) uint32 {
	return uint32(float32(time) / skillItem.GetAttackSpeed())
}

//GetMinCastDistance 获取最小攻击距离
func (skillItem *_SkillItem) GetMinCastDistance() float32 {
	return skillItem.Config.MinCastDistance * skillItem.Caster.Attr.Scale
}

//GetMaxCastDistance 获取最大攻击距离
func (skillItem *_SkillItem) GetMaxCastDistance() float32 {
	return skillItem.Config.MaxCastDistance * skillItem.Caster.Attr.Scale
}
