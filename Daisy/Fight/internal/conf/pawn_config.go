package conf

import (
	"Daisy/DataTables"
	"encoding/json"
	"fmt"
)

// PawnConfig pawn配置
type PawnConfig struct {
	PropValue       *DataTables.PropValue_Config // 属性数值
	CollisionRadius float32                      // 碰撞体积
	Mass            uint32                       // 质量id
	Scale           float32                      // 缩放比例
	PublicCD        uint32                       // 公共cd时间
	BeHitConf       *BeHitConfig                 // 受击配置
	ActConf         *ActConfig                   // 表演配置
	BornPauseAllAI  bool                         // 出生暂停所有AI
	InnerBornBuffs  []uint32                     // 内部出生buff

	WalkSpeed           float32 //走路速度
	RunSpeed            float32 //跑步速度
	FastSpeed           float32 //冲刺速度
	LookAtSpeed         float32 //迂回速度
	LookAtBackSpeed     float32 //迂回后退速度
	MoveMaxSpeed        float32 //基础移动最大速度
	MoveMaxAcceleration float32 //基础移动最大加速度
	TurnSpeed           float32

	EnergyShieldCapability uint32 // 能量护盾是否生效

	DodgeDist float32 // 闪避距离
	DodgeTime uint32  // 闪避时间
	FastDist  float32 // 冲刺距离

	AIID              uint32                 // ai id
	BlackBoardKeys    string                 // ai黑板所有key原始数据
	BlackBoardKeyData map[string]interface{} //ai黑板所有key
}

// loadRolePawnConfig 加载role配置
func loadRolePawnConfig(specialAgentExcelConf *DataTables.SpecialAgent_Config_Data, propValueExcelConf *DataTables.Prop_Config_Data, playerBeHitConfs map[uint32]*BeHitConfig) (map[uint32]*PawnConfig, error) {
	config := map[uint32]*PawnConfig{}
	var ok bool
	for id, val := range specialAgentExcelConf.SpecialAgent_ConfigItems {
		conf := &PawnConfig{
			PropValue:       nil,
			CollisionRadius: val.CollisionRadius,
			Mass:            val.Mass,
			Scale:           1,
			PublicCD:        val.PublicCD,
			BeHitConf:       nil,
			ActConf:         &ActConfig{},
			InnerBornBuffs: []uint32{
				InnerBuffID_RecoverHP,
				InnerBuffID_RecoverUltimateSkillPower,
				InnerBuffID_Thorns,
				InnerBuffID_StealUltimateSkillPower,
				InnerBuffID_EnergyShield,
			},
			WalkSpeed:              val.WalkSpeed,
			RunSpeed:               val.RunSpeed,
			FastSpeed:              val.FastSpeed,
			TurnSpeed:              val.TurnSpeed,
			LookAtSpeed:            val.LookAtSpeed,
			LookAtBackSpeed:        val.LookAtBackSpeed,
			MoveMaxSpeed:           val.MaxSpeed,
			MoveMaxAcceleration:    0,
			DodgeDist:              val.DodgeDist,
			DodgeTime:              val.DodgeTime,
			FastDist:               val.FastDist,
			AIID:                   val.AIID,
			BlackBoardKeys:         val.BlackBoardKeys,
			BlackBoardKeyData:      map[string]interface{}{},
			EnergyShieldCapability: val.ShieldCapability,
		}

		if len(val.BlackBoardKeys) > 0 {
			err := json.Unmarshal([]byte(val.BlackBoardKeys), &conf.BlackBoardKeyData)
			if err != nil {
				return nil, fmt.Errorf("specialAgentExcelConf BlackBoardKeys Unmarshal fail id: %v, val: %v", id, val.BlackBoardKeys)
			}
		}

		conf.PropValue, ok = propValueExcelConf.PropValue_ConfigItems[val.PropValueID]
		if !ok {
			return nil, fmt.Errorf("job属性配置失败不存在：%v", val.PropValueID)
		}

		conf.BeHitConf, ok = playerBeHitConfs[val.ID]
		if !ok {
			return nil, fmt.Errorf("job受击配置失败不存在：%v", val.ID)
		}

		config[id] = conf
	}

	return config, nil
}

// loadMonsterPawnConfig 加载怪物配置
func loadMonsterPawnConfig(monsterExcelConf *DataTables.Monster_Config_Data, propValueExcelConf *DataTables.Prop_Config_Data,
	monsterBeHitConfs map[uint32]*BeHitConfig, npcActConfs map[uint32]*ActConfig) (map[uint32]*PawnConfig, error) {
	config := map[uint32]*PawnConfig{}
	for id, val := range monsterExcelConf.Logic_ConfigItems {
		modelConfig, ok := monsterExcelConf.Model_ConfigItems[val.ModelID]
		if !ok {
			return nil, fmt.Errorf("Model_ConfigItems不存在：%v", val.ModelID)
		}

		conf := &PawnConfig{
			PropValue:       nil,
			CollisionRadius: modelConfig.CollisionRadius,
			Mass:            val.Mass,
			Scale:           val.ModelScale,
			PublicCD:        val.PublicCD,
			BeHitConf:       nil,
			ActConf:         nil,
			BornPauseAllAI:  val.BornPauseAllAI,
			InnerBornBuffs: []uint32{
				InnerBuffID_RecoverHP,
				InnerBuffID_Thorns,
				InnerBuffID_Unbalance,
				InnerBuffID_BornAct,
			},
			WalkSpeed:           modelConfig.WalkSpeed,
			RunSpeed:            modelConfig.RunSpeed,
			FastSpeed:           modelConfig.FastSpeed,
			TurnSpeed:           modelConfig.TurnSpeed,
			LookAtSpeed:         modelConfig.LookAtSpeed,
			LookAtBackSpeed:     modelConfig.LookAtBackSpeed,
			MoveMaxSpeed:        modelConfig.MaxSpeed,
			MoveMaxAcceleration: 0,
			DodgeDist:           modelConfig.DodgeDist,
			DodgeTime:           modelConfig.DodgeTime,
			FastDist:            0,
			AIID:                val.AIID,
			BlackBoardKeys:      val.BlackBoardKeys,
			BlackBoardKeyData:   map[string]interface{}{},
		}

		if len(val.BlackBoardKeys) > 0 {
			err := json.Unmarshal([]byte(val.BlackBoardKeys), &conf.BlackBoardKeyData)
			if err != nil {
				return nil, fmt.Errorf("monsterExcelConf BlackBoardKeys Unmarshal fail id: %v, val: %v", id, val.BlackBoardKeys)
			}
		}

		if floatLessEqual(float64(conf.Scale), 0) {
			conf.Scale = 1
		}

		conf.PropValue, ok = propValueExcelConf.PropValue_ConfigItems[val.PropValueID]
		if !ok {
			return nil, fmt.Errorf("monster属性配置失败不存在：%v %v", val.PropValueID, val)
		}

		conf.BeHitConf, ok = monsterBeHitConfs[val.NpcID]
		if !ok {
			return nil, fmt.Errorf("monster受击配置失败不存在：%v", val.NpcID)
		}

		conf.ActConf, ok = npcActConfs[val.NpcID]
		if !ok {
			return nil, fmt.Errorf("monster表演配置失败不存在：%v", val.NpcID)
		}

		if val.Difficulty != 0 {
			conf.InnerBornBuffs = append(conf.InnerBornBuffs, InnerBuffID_OverDrive)
		}

		config[id] = conf
	}

	return config, nil
}

// loadBGPawnConfig 加载怪物配置
func loadBGPawnConfig() map[uint32]*PawnConfig {
	config := map[uint32]*PawnConfig{}

	config[0] = &PawnConfig{
		PropValue:           &DataTables.PropValue_Config{},
		CollisionRadius:     0,
		Mass:                0,
		Scale:               1,
		PublicCD:            0,
		BeHitConf:           &BeHitConfig{},
		ActConf:             &ActConfig{},
		InnerBornBuffs:      nil,
		WalkSpeed:           0,
		RunSpeed:            0,
		FastSpeed:           0,
		LookAtSpeed:         0,
		LookAtBackSpeed:     0,
		MoveMaxSpeed:        0,
		MoveMaxAcceleration: 0,
		DodgeDist:           0,
		DodgeTime:           0,
		FastDist:            0,
		AIID:                0,
		BlackBoardKeys:      "",
		BlackBoardKeyData:   map[string]interface{}{},
	}

	return config
}
