package conf

import (
	"Daisy/DataTables"
	"Daisy/Proto"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
)

// AttackTemplArgs 伤害体模板参数
type AttackTemplArgs struct {
	SpawnPos       AttackSpawnPos       // 出生位置
	AOE            uint32               // AOE配置ID
	TargetCategory AttackTargetCategory // 目标选取策略
	MaxTarget      uint32               // 最大目标人数
	MaxLinkTarget  uint32               // 最大连接数量
	LinkDistance   float32              // 最大连接距离
	RepeatLink     bool                 // 能否重复连接
}

// AttackConfig 伤害体配置
type AttackConfig struct {
	*AttackArgs
	*AttackTimeLine
}

// AttackTemplate 伤害体模板配置
type AttackTemplate struct {
	ID          AttackTemplID   // 伤害体模板类型
	AttackConfs []*AttackConfig // 伤害体配置列表
}

// loadAttackTemplateConfig 加载伤害体模板配置
func loadAttackTemplateConfig(attackExcelConf *DataTables.Attack_Config_Data, templID uint32, templArgs string, attackTimeLineConfs []*AttackTimeLine) (*AttackTemplate, error) {
	templConf, ok := attackExcelConf.AttackTemplate_ConfigItems[templID]
	if !ok {
		return nil, fmt.Errorf("在工作簿【AttackTemplate_Config】未找到模板【%d】的配置信息", templID)
	}

	args := &AttackTemplArgs{}

	if err := json.Unmarshal([]byte(templConf.Args), args); err != nil {
		return nil, fmt.Errorf("解析工作簿【AttackTemplate_Config】中模板Args错误，%s", err.Error())
	}

	if templArgs != "" {
		if err := json.Unmarshal([]byte(templArgs), args); err != nil {
			return nil, fmt.Errorf("解析外部的配置参数TemplateArgs错误，%s", err.Error())
		}
	}

	if templID >= AttackTemplID_End {
		return nil, fmt.Errorf("伤害体模板【%d】暂不支持", templID)
	}

	aoeShape := &AttackShape{}

	if args.AOE > 0 {
		aoeConf, ok := attackExcelConf.AoeTemplate_ConfigItems[args.AOE]
		if !ok {
			return nil, fmt.Errorf("在工作簿【AoeTemplate_Config】未找到AOE模板【%d】的配置信息", args.AOE)
		}

		if err := json.Unmarshal([]byte(aoeConf.Args), aoeShape); err != nil {
			return nil, fmt.Errorf("解析工作簿【AoeTemplate_Config】中AOE模板【%d】Args错误，%s", args.AOE, err.Error())
		}

		aoeShape.FanAngle = aoeShape.FanAngle / 180 * math.Pi
	}

	templ := &AttackTemplate{
		ID: templID,
	}

	switch templID {
	case AttackTemplID_GravityGun:
		var attackArgsList []*AttackArgs

		if err := json.Unmarshal([]byte(templConf.Attacks), &attackArgsList); err != nil {
			return nil, fmt.Errorf("解析工作簿【AttackTemplate_Config】中模板Attacks错误，%s", err.Error())
		}

		if len(attackArgsList) != len(attackTimeLineConfs) {
			return nil, errors.New("解析工作簿【AttackTemplate_Config】中模板Attacks错误，与美术配置的Attack时间轴数量不一致")
		}

		for i := range attackTimeLineConfs {
			templ.AttackConfs = append(templ.AttackConfs, &AttackConfig{
				AttackArgs:     attackArgsList[i],
				AttackTimeLine: attackTimeLineConfs[i],
			})
		}

		// 初始化第一段
		templ.AttackConfs[0].TargetCategory = args.TargetCategory
		templ.AttackConfs[0].Spawn.Rotate = templ.AttackConfs[0].Spawn.Rotate / 180 * math.Pi

		// 初始化第二段
		templ.AttackConfs[1].TargetCategory = args.TargetCategory
		templ.AttackConfs[1].Shape = *aoeShape
		templ.AttackConfs[1].MaxHitTarget = int32(args.MaxTarget)
		templ.AttackConfs[1].Spawn.Rotate = templ.AttackConfs[1].Spawn.Rotate / 180 * math.Pi
		templ.AttackConfs[1].CastRange = true

	default:
		attackArgs := &AttackArgs{}

		if err := json.Unmarshal([]byte(templConf.Attacks), attackArgs); err != nil {
			return nil, fmt.Errorf("解析工作簿【AttackTemplate_Config】中模板Attacks错误，%s", err.Error())
		}

		attackArgs.Spawn.Pos = args.SpawnPos
		attackArgs.Spawn.Rotate = attackArgs.Spawn.Rotate / 180 * math.Pi
		attackArgs.TargetCategory = args.TargetCategory
		attackArgs.Shape = *aoeShape
		attackArgs.MaxHitTarget = int32(args.MaxTarget)
		attackArgs.MaxLinkTarget = int32(args.MaxLinkTarget)
		attackArgs.LinkDistance = args.LinkDistance
		attackArgs.RepeatLink = args.RepeatLink

		for _, atkTimeLine := range attackTimeLineConfs {
			tAttackArgs := &AttackArgs{}
			*tAttackArgs = *attackArgs

			templ.AttackConfs = append(templ.AttackConfs, &AttackConfig{
				AttackArgs:     tAttackArgs,
				AttackTimeLine: atkTimeLine,
			})
		}

		if len(templ.AttackConfs) > 0 {
			templ.AttackConfs[0].CastRange = true
		}

		for _, conf := range templ.AttackConfs {
			if conf.IsLink() {
				if conf.Type != AttackType_Single {
					return nil, fmt.Errorf("解析工作簿【AttackTemplate_Config】中模板【%d】Args错误，链式伤害体只能配置为单体伤害", templID)
				}
				conf.Shape.Type = Proto.AttackShapeType_Circle
				conf.Shape.Radius = conf.LinkDistance
			}
		}
	}

	for _, v := range templ.AttackConfs {
		if v.Type != AttackType_Aoe {
			continue
		}

		switch v.Shape.Type {
		case Proto.AttackShapeType_Rect:
			if v.Shape.Extend.X <= 0 {
				return nil, fmt.Errorf("解析工作簿【AoeTemplate_Config】中AOE模板【%d】Args错误，参数Extend.X不能小于等于0", args.AOE)
			}

			if !v.Spawn.AutoExtend {
				if v.Shape.Extend.Y <= 0 {
					return nil, fmt.Errorf("解析工作簿【AoeTemplate_Config】中AOE模板【%d】Args错误，对于不能修改形状的伤害体模板【%d】，参数Extend.X不能小于等于0", templID, args.AOE)
				}
			}

		case Proto.AttackShapeType_Circle:
			if !v.Spawn.AutoExtend {
				if v.Shape.Radius <= 0 {
					return nil, fmt.Errorf("解析工作簿【AoeTemplate_Config】中AOE模板【%d】Args错误，对于不能修改形状的伤害体模板【%d】，参数Radius不能小于等于0", templID, args.AOE)
				}
			}

		case Proto.AttackShapeType_Fan:
			if !v.Spawn.AutoExtend {
				if v.Shape.Radius <= 0 {
					return nil, fmt.Errorf("解析工作簿【AoeTemplate_Config】中AOE模板【%d】Args错误，对于不能修改形状的伤害体模板【%d】，参数Radius不能小于等于0", templID, args.AOE)
				}
			}

			if v.Shape.FanAngle <= 0 {
				return nil, fmt.Errorf("解析工作簿【AoeTemplate_Config】中AOE模板【%d】Args错误，参数FanAngle不能小于等于0", args.AOE)
			}
		}
	}

	return templ, nil
}

// loadInnerAttackTemplate 加载内置伤害体模板
func loadInnerAttackTemplate(fightConstConf *DataTables.FightConst_Config_Data, attackExcelConf *DataTables.Attack_Config_Data) (map[InnerAttackID]*AttackTemplate, error) {
	innerAttackTemplateMap := map[InnerAttackID]*AttackTemplate{}

	// 分裂伤害体
	loadSputteringAttackFun := func() error {
		templIDConf, ok := fightConstConf.FightConst_ConfigItems[ConstExcel_SputteringAttackTmplID]
		if !ok {
			return fmt.Errorf("在工作簿【FightConst_Config】中未找到分裂伤害体模板ID配置")
		}

		templID, err := strconv.Atoi(templIDConf.Value)
		if err != nil {
			return fmt.Errorf("解析工作簿【FightConst_Config】分裂伤害体模板ID失败，%s", err.Error())
		}

		timeLineConf, ok := fightConstConf.FightConst_ConfigItems[ConstExcel_SputteringAttackTimeLine]
		if !ok {
			return fmt.Errorf("在工作簿【FightConst_Config】中未找到分裂伤害体时间轴配置")
		}

		timeLine, err := loadAttackTimelineConfig(timeLineConf.Value)
		if err != nil {
			return fmt.Errorf("加载工作簿【FightConst_Config】中分裂伤害体时间轴配置失败，%s", err.Error())
		}

		templateConf, err := loadAttackTemplateConfig(attackExcelConf, uint32(templID), "", []*AttackTimeLine{timeLine})
		if err != nil {
			return fmt.Errorf("加载工作簿【FightConst_Config】中分裂伤害体模板配置失败，%s", err.Error())
		}

		innerAttackTemplateMap[InnerAttackID_SputteringAttack] = templateConf

		return nil
	}

	if err := loadSputteringAttackFun(); err != nil {
		return nil, err
	}

	return innerAttackTemplateMap, nil
}
