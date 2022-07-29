package main

import (
	"Cinder/Base/linemath"
	"Daisy/Data"
	"Daisy/DataTables"
	"Daisy/Fight"
	"Daisy/Proto"
	"encoding/json"
	"fmt"
	"github.com/sipt/GoJsoner"
	"io/ioutil"
	"strconv"
)

// StandPos 站位信息
type StandPos struct {
	Index    int
	Position linemath.Vector2
	Angle    float32
	Camp     Proto.Camp_Enum
}

// MineDefault 我方默认配置
type MineDefault struct {
	Lv              int32
	MaxHP           int64
	Attack          float64
	ExtendDodgeRate float32
	ExtendBlockRate float32
	ExtendHitRate   float32
	ExtendCritRate  float32
	CombineSkills   []uint32
}

// Attribute 战斗属性
type Attribute struct {
	Lv              *int32
	MaxHP           *int64
	Attack          *float64
	ExtendDodgeRate *float32
	ExtendBlockRate *float32
	ExtendHitRate   *float32
	ExtendCritRate  *float32
}

func (attr *Attribute) Fill(def *MineDefault) {
	if attr.Lv == nil {
		//不允许 出现0级 默认1级
		if def.Lv == 0 {
			def.Lv = 1
		}
		attr.Lv = &def.Lv
	}

	if attr.MaxHP == nil {
		attr.MaxHP = &def.MaxHP
	}

	if attr.Attack == nil {
		attr.Attack = &def.Attack
	}

	if attr.ExtendDodgeRate == nil {
		attr.ExtendDodgeRate = &def.ExtendDodgeRate
	}

	if attr.ExtendBlockRate == nil {
		attr.ExtendBlockRate = &def.ExtendBlockRate
	}

	if attr.ExtendHitRate == nil {
		attr.ExtendHitRate = &def.ExtendHitRate
	}

	if attr.ExtendCritRate == nil {
		attr.ExtendCritRate = &def.ExtendCritRate
	}
}

// Toy 玩具
type Toy struct {
	SuitId uint32
	BuffID uint32
}

// Player 玩家
type Player struct {
	Name           string
	SpecialAgentID uint32
	Attr           *Attribute
	NormalSkills   []uint32
	SuperSkills    []uint32
	UltimateSkills []uint32
	BornBuffs      []uint32
}

func (player *Player) Fill(def *MineDefault) {
	if player.Attr == nil {
		player.Attr = &Attribute{}
	}

	player.Attr.Fill(def)
}

// Npc npc
type Npc struct {
	ConfigId               uint32
	NormalSkills           []uint32
	SuperSkills            []uint32
	BornBuffs              []uint32
	OverDriveNormalAttacks []uint32
	OverDriveSuperSkills   []uint32
}

// Stand 站位
type Stand struct {
	Player *Player
	Npc    *Npc
}

// FightConf 战斗配置
type FightConf struct {
	Battlefield                         string
	MineDefault                         MineDefault
	Enemy, Mine                         map[int]*Stand
	_MineStandPosMap, _EnemyStandPosMap map[int]*StandPos
	BattlefieldData
	BattleAreaID           uint32
	EnemyFormationBuffList []uint32
	MineFormationBuffList  []uint32
}

type bornPoint struct {
	linemath.Vector2
	Angle float32
}

//BattlefieldData 战场数据
type BattlefieldData struct {
	AreaPoints   []linemath.Vector2 `json:areaPoints`
	EnemyPoints  []bornPoint        `json:enemyPoints`
	PlayerPoints []bornPoint        `json:playerPoints`
}

func (conf *FightConf) LoadBattlefield(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("读取对战Json配置出错，%s", err.Error()))
	}

	if err := json.Unmarshal([]byte(data), &conf.BattlefieldData); err != nil {
		panic(fmt.Sprintf("解析BattlefieldJson配置出错，%s", err.Error()))
	}
}

// LoadStandPos 加载站位
func (conf *FightConf) LoadStandPos() {
	for idx, val := range conf.EnemyPoints {
		conf._EnemyStandPosMap[idx+1] = &StandPos{
			Index:    idx + 1,
			Position: val.Vector2,
			Angle:    val.Angle,
			Camp:     1,
		}
	}

	for idx, val := range conf.PlayerPoints {
		conf._MineStandPosMap[idx+1] = &StandPos{
			Index:    idx + 1,
			Position: val.Vector2,
			Angle:    val.Angle,
		}
	}
}

// LoadConfig 加载配置
func (conf *FightConf) LoadConfig(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("读取对战Json配置出错，%s", err.Error()))
	}

	// utf-8 bom头部
	if data[0] == 239 && data[1] == 187 && data[2] == 191 {
		data = data[3:]
	}

	dataStr, err := GoJsoner.Discard(string(data))
	if err != nil {
		panic(fmt.Sprintf("读取对战Json配置出错，%s", err.Error()))
	}

	if err := json.Unmarshal([]byte(dataStr), conf); err != nil {
		panic(fmt.Sprintf("解析对战Json配置出错，%s", err.Error()))
	}

	parseStands := func(stands map[int]*Stand) error {
		for _, stand := range stands {
			if stand.Player != nil {
				stand.Player.Fill(&conf.MineDefault)
			}
		}

		return nil
	}

	if err := parseStands(conf.Mine); err != nil {
		panic(fmt.Sprintf("解析Mine站位出错，%s", err.Error()))
	}

	if err := parseStands(conf.Enemy); err != nil {
		panic(fmt.Sprintf("解析Enemy站位出错，%s", err.Error()))
	}
}

// BuildSceneInfo 创建场景信息
func (conf *FightConf) BuildSceneInfo() *Fight.SceneInfo {
	info := &Fight.SceneInfo{
		Formation: []*Fight.FormationInfo{
			conf.buildFormation(&conf.Mine),
			conf.buildFormation(&conf.Enemy),
		},
		BoundaryPoints:  conf.BattlefieldData.AreaPoints,
		SimulatorMode:   true,
		TestMode:        true,
		MaxMilliseconds: Fight.BattleMaxMilliseconds,
	}

	for _, val := range conf.BattlefieldData.PlayerPoints {
		info.Formation[int(Proto.Camp_Red)].BornPoints = append(info.Formation[int(Proto.Camp_Red)].BornPoints, &Fight.BornPoint{
			Point: Proto.Position{X: val.X, Y: val.Y},
			Angle: val.Angle,
		})
	}

	for _, val := range conf.BattlefieldData.EnemyPoints {
		info.Formation[int(Proto.Camp_Blue)].BornPoints = append(info.Formation[int(Proto.Camp_Blue)].BornPoints, &Fight.BornPoint{
			Point: Proto.Position{X: val.X, Y: val.Y},
			Angle: val.Angle,
		})
	}
	info.Formation[int(Proto.Camp_Red)].FormationBuffList = conf.MineFormationBuffList
	info.Formation[int(Proto.Camp_Blue)].FormationBuffList = conf.EnemyFormationBuffList

	return info
}

// buildFormation 创建阵型信息
func (conf *FightConf) buildFormation(stands *map[int]*Stand) *Fight.FormationInfo {
	return &Fight.FormationInfo{
		PawnInfos: func() (pawnInfos []*Fight.PawnInfo) {
			//for idx, stand := range *stands {
			for idx := 1; idx <= len(*stands); idx++ {
				stand := (*stands)[idx]
				if stand == nil {
					continue
				}
				pawnInfos = append(pawnInfos, conf.buildPawns(stand, func() *StandPos {
					if stands == &conf.Enemy {
						pos, ok := conf._EnemyStandPosMap[idx]
						if !ok {
							panic(fmt.Sprintf("站位点 %d 配置错误", idx))
						}
						return pos
					} else {
						pos, ok := conf._MineStandPosMap[idx]
						if !ok {
							panic(fmt.Sprintf("站位点 %d 配置错误", idx))
						}
						return pos
					}
				}())...)
			}
			return
		}(),
		RageTime:      40000, //毫秒
		CombineSkills: conf.MineDefault.CombineSkills,
	}
}

// buildPawns 创建pawn
func (conf *FightConf) buildPawns(stand *Stand, standPos *StandPos) []*Fight.PawnInfo {
	if stand.Player != nil {
		var jobConf *DataTables.SpecialAgent_Config

		jobConfs := Data.GetSpecialAgentConfig().SpecialAgent_ConfigItems
		for _, t := range jobConfs {
			if t.ID == stand.Player.SpecialAgentID {
				jobConf = t
				break
			}
		}

		// 获取职业配置
		if jobConf == nil {
			panic(fmt.Sprintf("获取Job %d 配置失败", stand.Player.SpecialAgentID))
		}

		playerInfo := &Fight.PawnInfo{
			PawnInfo: &Proto.PawnInfo{
				ConfigId: jobConf.ID,
				Type:     Proto.PawnType_Role,
				Role: &Proto.FightRoleInfo{
					RoleId: strconv.Itoa(int(standPos.Camp)*10 + standPos.Index),
					Name:   stand.Player.Name,
				},
				Camp:  standPos.Camp,
				Level: *stand.Player.Attr.Lv,
				BornPos: &Proto.Position{
					X: standPos.Position.X,
					Y: standPos.Position.Y,
				},
				BornAngle: standPos.Angle,
			},
			NormalAtkList:     stand.Player.NormalSkills,
			AddComboAttack:    jobConf.AddComboAttack,
			SuperSkillList:    stand.Player.SuperSkills,
			UltimateSkillList: stand.Player.UltimateSkills,
			BornBuffs:         stand.Player.BornBuffs,
			SimulatorModeInfo: &Fight.SimulatorModeInfo{
				MaxHP:           *stand.Player.Attr.MaxHP,
				Attack:          *stand.Player.Attr.Attack,
				ExtendDodgeRate: *stand.Player.Attr.ExtendDodgeRate,
				ExtendBlockRate: *stand.Player.Attr.ExtendBlockRate,
				ExtendHitRate:   *stand.Player.Attr.ExtendHitRate,
				ExtendCritRate:  *stand.Player.Attr.ExtendCritRate,
			},
		}

		// 创建玩家信息
		pawnList := append([]*Fight.PawnInfo{}, playerInfo)

		return pawnList

	} else if stand.Npc != nil {
		// 获取NPC配置
		var npcConf *DataTables.Logic_Config
		var ok bool
		npcConf, ok = Data.GetMonsterConfig().Logic_ConfigItems[stand.Npc.ConfigId]
		if !ok {
			panic(fmt.Sprintf("获取NPC %d 配置失败", stand.Npc.ConfigId))
		}

		npcInfo := &Fight.PawnInfo{
			PawnInfo: &Proto.PawnInfo{
				ConfigId: stand.Npc.ConfigId,
				Type:     Proto.PawnType_Npc,
				Npc: &Proto.FightNpcInfo{
					IsBoss: npcConf.Difficulty != 0,
				},
				Camp:  standPos.Camp,
				Level: int32(npcConf.Level),
				BornPos: &Proto.Position{
					X: standPos.Position.X,
					Y: standPos.Position.Y,
				},
				BornAngle: standPos.Angle,
			},
			NormalAtkList:             stand.Npc.NormalSkills,
			SuperSkillList:            stand.Npc.SuperSkills,
			OverDriveNormalAttackList: stand.Npc.OverDriveNormalAttacks,
			OverDriveSuperSkillList:   stand.Npc.OverDriveSuperSkills,
			BornBuffs:                 stand.Npc.BornBuffs,
		}

		return []*Fight.PawnInfo{
			npcInfo,
		}
	}

	return nil
}

// clipInt32 限制取值范围
func clipInt32(num, left, right int32) int32 {
	if num < left {
		return left
	}

	if num > right {
		return right
	}

	return num
}
