package internal

import (
	"Cinder/Base/linemath"
	"Daisy/Fight/internal/conf"
	"Daisy/Proto"
)

// Destroy 销毁自身
func (attack *Attack) Destroy(isBreak bool, breakCaster *Pawn) {
	attack.Caster.Scene.removeAttack(attack, isBreak, breakCaster)
}

// Src 伤害体来源
func (attack *Attack) Src() Proto.AttackSrc_Enum {
	if attack.Skill != nil {
		return Proto.AttackSrc_Skill
	}

	if attack.Buff != nil {
		return Proto.AttackSrc_Buff
	}

	return Proto.AttackSrc_Custom
}

// getTargetPos 查询目标位置
func (attack *Attack) getTargetPos() linemath.Vector2 {
	switch attack.Config.Type {
	case conf.AttackType_Single:
		if len(attack.castTargets) > 0 {
			return attack.castTargets[0].GetPos()
		}
	case conf.AttackType_Aoe:
		return attack.castAoePos
	}
	return attack.Caster.GetPos()
}

// ExtendHitTime 延长hit时间
func (attack *Attack) ExtendHitTime(time uint32) {
	attack.hitExtendTime += time
}

// SearchTargets 查询伤害体目标列表
func (attack *Attack) SearchTargets() []*Pawn {
	return attack.Caster.Scene.searchAttackTargets(attack.Config.AttackArgs, attack.Pos, attack.Angle, attack.autoExtendValue, attack.Caster.GetCamp(), attack.Scale)
}

// Equal 伤害体是否相同
func (attack *Attack) Equal(other *Attack) bool {
	if attack == nil || other == nil {
		return attack == other
	}

	return attack.UID == other.UID
}

// ZoomAttackTime 缩放技能攻击时间
func (attack *Attack) ZoomAttackTime(time uint32) uint32 {
	switch attack.Config.MoveMode {
	case conf.AttackMoveMode_None:
	default:
		return 1
	}

	switch attack.Config.DestroyType {
	case conf.AttackDestroyType_LifeTime:
	default:
		return 1
	}

	switch attack.Src() {
	case Proto.AttackSrc_Skill:
		return attack.Skill.ZoomAttackTime(time)
	default:
		return 1
	}
}
