package internal

import (
	"Daisy/Proto"
	"errors"
)

// FormationInfo 战斗阵形信息
type FormationInfo struct {
	PawnInfos         []*PawnInfo  // 所有pawn信息
	CombineSkills     []uint32     // 合体必杀技列表
	RageTime          uint32       // 狂暴时间(毫秒)
	FormationBuffList []uint32     //阵营buff 施加到背景NPC上
	BornPoints        []*BornPoint //出生点信息
}

type BornPoint struct {
	Point Proto.Position
	Angle float32
}

// _FormationInfo 战斗阵形信息
type _FormationInfo struct {
	*FormationInfo
	Camp Proto.Camp_Enum
}

// Formation 战斗阵形
type Formation struct {
	scene          *Scene
	Info           _FormationInfo
	PawnList       []*Pawn
	BackgroundPawn *Pawn
	RageBeginTime  uint32 // 狂暴时间点 (毫秒)
	Raged          bool   // 是否已狂暴
}

// init 阵形初始化
func (formation *Formation) init(scene *Scene, formationInfo *FormationInfo, camp Proto.Camp_Enum) error {
	if scene == nil || formationInfo == nil {
		return errors.New("args invalid")
	}

	formation.scene = scene
	formation.Info.FormationInfo = formationInfo
	formation.Info.Camp = camp

	if len(formation.Info.PawnInfos) <= 0 {
		return errors.New("no pawn in formation")
	}

	return nil
}

// putPawns 放置所有pawn
func (formation *Formation) putBackGroundPawns() error {
	info := &PawnInfo{
		PawnInfo: &Proto.PawnInfo{
			Type:    Proto.PawnType_BG,
			Camp:    formation.Info.Camp,
			BornPos: &Proto.Position{},
		},
		CombineSkillList: formation.Info.CombineSkills,
		BornBuffs:        formation.Info.FormationBuffList,
	}

	pawn, err := formation.scene.Summon(info)
	if err != nil {
		return err
	}

	//背景NPC 先天自带
	pawn.State.ChangeStat(Stat_CantBeDamage, true)
	pawn.State.ChangeStat(Stat_CantBeHitControl, true)
	pawn.State.ChangeStat(Stat_CantBeAddDecrBuff, true)
	pawn.State.ChangeStat(Stat_CantBeEnemySelect, true)
	pawn.State.ChangeStat(Stat_CantBeAddIncrBuff, true)
	pawn.State.ChangeStat(Stat_CantBeFriendlySelect, true)

	formation.BackgroundPawn = pawn

	return nil
}

// putPawns 放置所有pawn
func (formation *Formation) putPawns() error {
	for _, pawnInfo := range formation.Info.PawnInfos {
		// 强制阵营一致
		pawnInfo.Camp = formation.Info.Camp

		// 构造pawn
		pawn := &Pawn{}
		if err := pawn.init(formation.scene, pawnInfo); err != nil {
			return err
		}

		// 放置pwan
		if err := formation.scene.putPawn(pawn); err != nil {
			return err
		}
	}

	return nil
}
