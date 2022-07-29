package internal

import (
	"Cinder/Base/linemath"
	"Daisy/Fight/internal/conf"
	"Daisy/Proto"
	"fmt"
)

// Attack 伤害体
type Attack struct {
	UID              uint32               // 唯一ID
	Config           *conf.AttackConfig   // 伤害体配置
	Skill            *Skill               // 创建伤害体的技能
	Buff             *Buff                // 创建伤害体的buff
	Caster           *Pawn                // 施法者
	effectCallback   IEffectCallback      // 效果回调
	configTab        []*conf.AttackConfig // 伤害体配置表
	index            uint32               // 伤害体索引
	groupID          uint32               // 分组ID
	spawnPawn        *Pawn                // 用于确定出生位置的pawn
	spawnPos         linemath.Vector2     // 出生位置
	Pos              linemath.Vector2     // 当前位置
	Angle            float32              // 当前角度
	autoExtendValue  float32              // 自动调整形状大小数值
	CasterSnapshot   *Pawn                // 施法者数值快照
	castTargets      []*Pawn              // 施法目标列表
	castAoePos       linemath.Vector2     // 施法Aoe位置
	HitTimes         uint32               // hit次数
	HitTargets       []*Pawn              // hit目标列表
	linkTimes        int32                // 链接次数
	linkTargetMap    map[uint32]bool      // 已链接的目标
	createTime       uint32               // 创建时间
	fixTimeMoveEnd   uint32               // 固定时间移动结束时间
	hitExtendTime    uint32               // hit延长时间
	IsDestroy        bool                 // 已删除
	isMoveEnd        bool                 // 是否移动结束
	actionMoveAoe    *Proto.AttackMoveAoe // 上次移动Aoe帧
	LogicData        interface{}          // 逻辑层数据
	Scale            float32              // 缩放比例
	aoeWarnRegionUID uint32               // Aoe警告区域UID
}

// init 初始化伤害体
func (attack *Attack) init(caster *Pawn, casterSnapshot *Pawn, configTab []*conf.AttackConfig, index uint32,
	targetList []*Pawn, targetPos linemath.Vector2, effectCallback IEffectCallback, frontAttack *Attack, scale float32) error {
	if caster == nil {
		return fmt.Errorf("nil caster")
	}

	if casterSnapshot == nil {
		return fmt.Errorf("nil caster snapshot")
	}

	// 读取伤害体配置
	if index >= uint32(len(configTab)) {
		return fmt.Errorf("init attack failed, attack index %d great equal AttackConfs len %d",
			index, len(configTab))
	}

	attack.configTab = configTab
	attack.Config = configTab[index]

	// 初始化参数
	attack.UID = caster.Scene.generateUID()
	attack.Caster = caster
	attack.effectCallback = effectCallback
	attack.index = index
	attack.CasterSnapshot = casterSnapshot
	attack.Scale = scale

	// 初始化目标信息
	if frontAttack != nil {
		if attack.Config.IsLink() {
			if !attack.Config.RepeatLink {
				attack.linkTargetMap = frontAttack.linkTargetMap

				if attack.linkTargetMap == nil {
					attack.linkTargetMap = map[uint32]bool{}
				}

				if len(frontAttack.castTargets) > 0 {
					attack.linkTargetMap[frontAttack.castTargets[0].UID] = true
				}
			}

			castTargets := frontAttack.SearchTargets()

			for i := len(castTargets) - 1; i >= 0; i-- {
				target := castTargets[i]

				if len(frontAttack.castTargets) > 0 {
					if target.Equal(frontAttack.castTargets[0]) {
						continue
					}
				}

				if !attack.Config.RepeatLink {
					if _, ok := attack.linkTargetMap[target.UID]; ok {
						continue
					}
				}

				attack.castTargets = append(attack.castTargets, target)
			}

		} else {
			// 继承目标信息
			switch frontAttack.Config.Type {
			case conf.AttackType_Single:
				switch attack.Config.Type {
				case conf.AttackType_Single:
					if frontAttack.Config.TargetCategory != attack.Config.TargetCategory {
						return fmt.Errorf("init attack failed, attack index %d inherit TargetCategory was different", index)
					}
					attack.castTargets = append([]*Pawn{}, frontAttack.castTargets...)

				case conf.AttackType_Aoe:
					attack.castAoePos = frontAttack.getTargetPos()
				}
			case conf.AttackType_Aoe:
				switch attack.Config.Type {
				case conf.AttackType_Single:
					if frontAttack.Config.TargetCategory != attack.Config.TargetCategory {
						return fmt.Errorf("init attack failed, attack index %d inherit TargetCategory was different", index)
					}
					attack.castTargets = frontAttack.SearchTargets()

				case conf.AttackType_Aoe:
					attack.castAoePos = frontAttack.getTargetPos()
				}
			}
		}
	} else {
		// 目标信息
		switch attack.Config.Type {
		case conf.AttackType_Single:
			attack.castTargets = append([]*Pawn{}, targetList...)

		case conf.AttackType_Aoe:
			attack.castAoePos = targetPos
		}
	}

	// 出生位置信息
	switch attack.Config.Spawn.Pos {
	case conf.AttackSpawnPos_Caster:
		attack.spawnPos = caster.GetPos()
		attack.Angle = caster.GetAngle()

	case conf.AttackSpawnPos_Target:
		if attack.Config.Type == conf.AttackType_Single {
			if len(attack.castTargets) > 0 {
				attack.spawnPawn = attack.castTargets[0]
			}
		}
		attack.spawnPos = attack.getTargetPos()
		attack.Angle = caster.GetAngle()

	case conf.AttackSpawnPos_Inherit:
		if frontAttack == nil {
			return fmt.Errorf("init attack failed, attack index %d inherit SpawnPos with nil frontAttack", index)
		}

		if attack.Config.Type == conf.AttackType_Single {
			if len(attack.castTargets) > 0 {
				attack.spawnPawn = attack.castTargets[0]
			}
		}
		attack.spawnPos = frontAttack.getTargetPos()
		attack.Angle = frontAttack.Angle
	}

	// 分组信息
	if attack.Config.Grouping {
		if frontAttack != nil {
			attack.groupID = frontAttack.groupID
		} else {
			attack.groupID = caster.Scene.generateUID()
		}
	}

	// 链式伤害体增加链接次数
	if attack.Config.IsLink() {
		if frontAttack != nil {
			attack.linkTimes = frontAttack.linkTimes + 1
		}
	}

	return nil
}

// initNextAttack 初始化继续攻击的伤害体
func (attack *Attack) initNextAttack(frontAttack *Attack, index uint32) error {
	// 调用普通初始化
	if err := attack.init(frontAttack.Caster, frontAttack.CasterSnapshot, frontAttack.configTab, index,
		nil, linemath.Vector2{}, frontAttack.effectCallback, frontAttack, frontAttack.Scale); err != nil {
		return fmt.Errorf("init next attack %s, %s", func() string {
			switch attack.Src() {
			case Proto.AttackSrc_Skill:
				return fmt.Sprintf(", skill %d", attack.Skill.Config.ValueID())
			case Proto.AttackSrc_Buff:
				return fmt.Sprintf(", buff %d", attack.Buff.Config.MainID())
			}
			return ""
		}(), err.Error())
	}

	// 设置伤害体来源
	attack.Skill = frontAttack.Skill
	attack.Buff = frontAttack.Buff

	return nil
}

// initSkillAttack 初始化技能伤害体
func (attack *Attack) initSkillAttack(skill *Skill, index uint32) error {
	// 调用普通初始化
	if err := attack.init(skill.Caster, skill.Caster.Snapshot(), skill.Config.TemplateConfig.AttackConfs, index,
		skill.TargetList, skill.TargetPos, nil, nil, skill.scale); err != nil {
		return fmt.Errorf("skill %d, %s", skill.Config.ValueID(), err.Error())
	}

	// 设置伤害体来源
	attack.Skill = skill

	return nil
}
