package internal

import (
	"Daisy/Proto"
)

// ExtendAttackTime 延长攻击时间
func (skill *Skill) ExtendAttackTime(time uint32) {
	if skill.Stat != Proto.SkillState_Attack {
		return
	}

	skill.attackExtendTime += time
}

// SaveCombineSkillReadyMembers 记录合体技准备释放的成员列表
func (skill *Skill) SaveCombineSkillReadyMembers(members []uint32) {
	skill.combineSkillReadyMembers = members
}

// Equal 技能对象是否相同
func (skill *Skill) Equal(other *Skill) bool {
	if skill == nil || other == nil {
		return skill == other
	}

	return skill.UID == other.UID
}

//fixTurnTime 修正 turn开始时间+持续时间数据
func (skill *Skill) fixTurnTime() {
	if skill.Config.TimeLineConfig.Turn == nil {
		return
	}

	skill.turnBeginTime = skill.Config.TimeLineConfig.Turn.Begin
	if skill.skipBefore {
		//第一种情况 turn开始时间处于前摇范围内
		if skill.turnBeginTime <= skill.Config.TimeLineConfig.ShowTime+skill.Config.TimeLineConfig.BeforeTime {
			skill.turnBeginTime = skill.Config.TimeLineConfig.ShowTime
		} else {
			skill.turnBeginTime -= skill.Config.TimeLineConfig.BeforeTime
		}
	}
}
