package internal

import (
	"Cinder/Base/linemath"
	"Daisy/Proto"
	"errors"
	"fmt"
	"reflect"
)

// Skill 技能
type Skill struct {
	*_SkillItem                                    // 技能道具
	UID                      uint32                // 唯一ID
	effectTab                []IEffectCallback     // 技能效果表
	TargetList               []*Pawn               // 目标列表
	TargetPos                linemath.Vector2      // 目标点
	Stat                     Proto.SkillState_Enum // 技能状态
	inStatTime               uint32                // 进入技能阶段时间（单位ms）
	attackBegin              Bits                  // 已创建的伤害体
	attackExtendTime         uint32                // 延长攻击阶段时间
	skipBefore               bool                  // 是否跳过前摇
	actionUseSKill           *Proto.UseSkill       // 记录使用技能帧
	combineSkillReadyMembers []uint32              // 合体技准备释放的成员列表
	scale                    float32               // 缩放比例
	isAlreadyTurned          bool                  //是否已经转向修正过
	beginTime                uint32                //技能开始时间
	beginDashingTime         uint32                //冲刺开始时间
	endDashingTime           uint32                //冲刺结束时间
	turnBeginTime            uint32                //后处理turn开始时间
}

// init 初始化
func (skill *Skill) init(skillItem *_SkillItem) error {
	if skillItem == nil {
		return errors.New("args invalid")
	}

	skill._SkillItem = skillItem
	skill.UID = skill.Caster.Scene.generateUID()
	skill.scale = skill.Caster.Attr.Scale

	effect, err := createEffect(skill.Config.SkillKind, reflect.Value{})
	if err != nil {
		return fmt.Errorf("skill %d create effect failed, %s", skill.Config.ValueID(), err.Error())
	}

	if skillEff, ok := effect.(ISkillEffect); ok {
		if err := skillEff.Init(skill); err != nil {
			return fmt.Errorf("skill %d effect %d init error, %s", skill.Config.ValueID(), skill.Config.SkillKind, err.Error())
		}
	}

	skill.effectTab = append(skill.effectTab, effect)
	skill.beginTime = skill.Caster.Scene.NowTime
	skill.fixTurnTime()

	return nil
}
