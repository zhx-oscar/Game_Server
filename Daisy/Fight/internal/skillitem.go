package internal

import (
	"Daisy/Fight/internal/conf"
	"errors"
	"fmt"
)

// _SkillItem 技能道具
type _SkillItem struct {
	Config    *conf.SkillConfig // 技能配置
	Caster    *Pawn             // Pawn
	cdEndTime uint32            // 冷却结束时间
}

// init 初始化
func (skillItem *_SkillItem) init(pawn *Pawn, valueID uint32) error {
	if pawn == nil {
		return errors.New("args invalid")
	}

	skillItem.Caster = pawn

	var ok bool
	if skillItem.Config, ok = skillItem.Caster.Scene.GetSkillConfig(valueID); !ok {
		return fmt.Errorf("skill config %d not found", valueID)
	}

	return nil
}

// createSkill 创建技能
func (skillItem *_SkillItem) createSkill() (*Skill, error) {
	skill := &Skill{}
	if err := skill.init(skillItem); err != nil {
		return nil, fmt.Errorf("pawn type %d config %d init skill %d failed, %s", skillItem.Caster.Info.Type,
			skillItem.Caster.Info.ConfigId, skillItem.Config.ValueID(), err.Error())
	}

	return skill, nil
}
