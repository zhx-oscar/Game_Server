package main

import (
	"Cinder/Base/Util"
	"Cinder/Space"
	"Daisy/Battle/sceneconfigdata"
	"Daisy/Data"
	"Daisy/Fight"
	"Daisy/Proto"
	"errors"
	"fmt"
)

var (
	errCanfFindProgressConfig = errors.New("找不到Raid进度配置")
	errCantFindNPCConfig      = errors.New("找不到怪物配置")
)

func (team *_Team) GetBattleFieldCfg(progress uint32, battleFieldID uint32) (*sceneconfigdata.BattleFieldCfg, error) {
	battleAreaCfg, ok := Data.GetSceneConfig().BattleArea_ConfigItems[progress]
	if !ok {
		return nil, errCanfFindProgressConfig
	}

	mapID := battleAreaCfg.MapID
	battleFieldPath := fmt.Sprintf("../res/MapData/%d/Battle_%d.json", mapID, battleFieldID)
	return sceneconfigdata.LoadBattleField(battleFieldPath)
}

// BuildRaidSceneInfo 构建跑图中战场，包括小怪战和BOSS战
func (team *_Team) BuildRaidSceneInfo(progress uint32, isBoss bool, battleFieldID uint32) (*Fight.SceneInfo, error) {
	battleFieldCfg, err := team.GetBattleFieldCfg(progress, battleFieldID)
	if err != nil {
		return nil, err
	}

	info := &Fight.SceneInfo{
		Formation:       make([]*Fight.FormationInfo, 2),
		MaxMilliseconds: Fight.BattleMaxMilliseconds,
		BoundaryPoints:  battleFieldCfg.AreaPoints,
	}

	// 己方站位
	info.Formation[Proto.Camp_Red], err = team.buildRaidOwnFormation(battleFieldCfg, isBoss)
	if err != nil {
		return nil, err
	}

	// 对手站位
	info.Formation[Proto.Camp_Blue], err = team.buildRaidEnemyFormation(battleFieldID, battleFieldCfg, progress, isBoss)
	if err != nil {
		return nil, err
	}

	//加载继承的战斗属性
	team.loadInheritAttr(info)

	return info, nil
}

func (team *_Team) loadInheritAttr(info *Fight.SceneInfo) {
	var inheritAttr = &Proto.FightInherit{
		PawnInheritMap: make(map[string]*Proto.PawnInherit),
	}
	team.TraversalActor(func(ia Space.IActor) {
		role := ia.(*_Role)
		inheritAttr.PawnInheritMap[role.GetID()] = role.prop.Data.InheritAttr
	})
	info.Inherit = inheritAttr
}

func (team *_Team) saveInheritAttr(result *Proto.FightResult) {
	for key, value := range result.Inherit.PawnInheritMap {
		ia, err := team.GetActor(key)
		if err != nil {
			continue
		}

		role := ia.(*_Role)
		role.prop.SyncSetInheritAttr(value)
	}
}

func (team *_Team) buildRaidOwnFormation(battleFieldCfg *sceneconfigdata.BattleFieldCfg, isBoss bool) (*Fight.FormationInfo, error) {

	formation := &Fight.FormationInfo{
		PawnInfos: make([]*Fight.PawnInfo, 0),
		RageTime:  40000,
	}

	for _, val := range battleFieldCfg.PlayerPoints {
		formation.BornPoints = append(formation.BornPoints, &Fight.BornPoint{Point: Proto.Position{
			X: val.X,
			Y: val.Y,
		}, Angle: val.Angle})
	}

	idx := 0
	team.TraversalActor(func(actor Space.IActor) {
		role := actor.(*_Role)
		fightingBuild, ok := role.prop.Data.BuildMap[role.prop.Data.FightingBuildID]
		if !ok {
			return
		}

		specialAgent, ok := role.prop.Data.SpecialAgentList[fightingBuild.SpecialAgentID]
		if !ok {
			return
		}

		specialAgentCfg, ok := Data.GetSpecialAgentConfig().SpecialAgent_ConfigItems[fightingBuild.SpecialAgentID]

		//装备属性
		equipAttr := role.CalcEquip(fightingBuild.BuildID)

		//获取出战build内技能相关
		ultimateSkills, _, superSkills := role.getBuildBattleSkills()

		//todo 合体技系统还没有，先临时获取
		formation.CombineSkills = specialAgentCfg.CombineSkills

		idx = idx % len(battleFieldCfg.PlayerPoints)
		pawn := &Fight.PawnInfo{
			PawnInfo: &Proto.PawnInfo{
				Id:       uint32(idx + 1),
				ConfigId: specialAgent.Base.ConfigID,
				Type:     Proto.PawnType_Role,
				Role: &Proto.FightRoleInfo{
					RoleId: role.GetID(),
					Name:   role.prop.Data.Base.Name,
					Equips: role.GetFightingBDEquips(),
				},
				Camp:  Proto.Camp_Red,
				Level: int32(specialAgent.Base.Level),
				BornPos: &Proto.Position{
					X: battleFieldCfg.PlayerPoints[idx].X,
					Y: battleFieldCfg.PlayerPoints[idx].Y,
				},
				BornAngle: battleFieldCfg.PlayerPoints[idx].Angle,
			},
			NormalAtkList:     specialAgentCfg.NormalAttack,
			AddComboAttack:    specialAgentCfg.AddComboAttack,
			SuperSkillList:    superSkills,
			UltimateSkillList: ultimateSkills,
			BornBuffs:         equipAttr.BornBuffs,
			AttrAffixList:     equipAttr.AttrAffixList,
		}

		//debug 词缀嵌入
		pawn.AttrAffixList = append(pawn.AttrAffixList, role.debugAttraffix...)

		formation.PawnInfos = append(formation.PawnInfos, pawn)
		buffs := role.getCurSpecialAgentTalentBuff()
		for _, v := range buffs {
			pawn.BornBuffs = append(pawn.BornBuffs, v)
		}

		idx++
	})

	return formation, nil
}

func (team *_Team) buildRaidEnemyFormation(battleFieldID uint32, battleFieldCfg *sceneconfigdata.BattleFieldCfg, progress uint32, isBoss bool) (*Fight.FormationInfo, error) {
	battleAreaCfg, ok := Data.GetSceneConfig().BattleArea_ConfigItems[progress]
	if !ok {
		return nil, errCanfFindProgressConfig
	}

	npcs := make([]*Fight.PawnInfo, 0)
	//如果怪提前刷过则直接获取已刷的怪，不重新构建
	if npcs, ok = team.spawnMonsterInfos[battleFieldID]; !ok {
		var err error
		npcs, err = team.buildRaidEnemyInfos(battleFieldCfg, progress, isBoss, false)
		if err != nil {
			return nil, err
		}
	} else {
		delete(team.spawnMonsterInfos, battleFieldID)
	}

	formationBuffList := battleAreaCfg.NormalFormationBuffList
	if isBoss {
		formationBuffList = battleAreaCfg.BossFormationBuffList
	}

	formation := &Fight.FormationInfo{
		PawnInfos:         npcs,
		CombineSkills:     []uint32{},
		RageTime:          uint32(60 * 1000),
		FormationBuffList: formationBuffList,
	}

	for _, val := range battleFieldCfg.EnemyPoints {
		formation.BornPoints = append(formation.BornPoints, &Fight.BornPoint{Point: Proto.Position{
			X: val.X,
			Y: val.Y,
		}, Angle: val.Angle})
	}

	return formation, nil
}

func (team *_Team) buildRaidEnemyInfos(battleFieldCfg *sceneconfigdata.BattleFieldCfg, progress uint32, isBoss bool, preSpawn bool) ([]*Fight.PawnInfo, error) {
	battleAreaCfg, ok := Data.GetSceneConfig().BattleArea_ConfigItems[progress]
	if !ok {
		return nil, errCanfFindProgressConfig
	}

	npcIDs := battleAreaCfg.NpcId
	if isBoss {
		npcIDs = battleAreaCfg.BossNpcId
	}
	npcs := make([]*Fight.PawnInfo, 0)

	for i := 0; i < len(npcIDs); i++ {
		if npcIDs[i] == 0 {
			continue
		}

		npc, err := team.buildNPCInfo(npcIDs[i], battleAreaCfg.MonsterBulidID, battleAreaCfg.BossBulidID)
		if err != nil {
			continue
		}

		if preSpawn {
			npc.PawnInfo.Npc.SpawnID = Util.GetGUID()
		}

		value := battleFieldCfg.EnemyPoints[i%len(battleFieldCfg.EnemyPoints)]
		npc.BornPos = &Proto.Position{
			X: value.X,
			Y: value.Y,
		}
		npc.BornAngle = value.Angle
		npc.Id = uint32(i + 1)
		npcs = append(npcs, npc)
	}

	return npcs, nil
}

func (team *_Team) buildNPCInfo(id uint32, monsterBuildIds, bossBuildIds []uint32) (*Fight.PawnInfo, error) {
	logicCfg, ok := Data.GetMonsterConfig().Logic_ConfigItems[id]
	if !ok {
		return nil, errCantFindNPCConfig
	}

	isBoss := logicCfg.Difficulty != 0

	info := &Fight.PawnInfo{
		PawnInfo: &Proto.PawnInfo{
			Type:     Proto.PawnType_Npc,
			ConfigId: id,
			Npc: &Proto.FightNpcInfo{
				IsBoss: isBoss,
			},
			Camp:  Proto.Camp_Blue,
			Level: int32(logicCfg.Level),
		},
		NormalAtkList:             logicCfg.NormalAttack,
		SuperSkillList:            logicCfg.SuperSkill,
		OverDriveNormalAttackList: logicCfg.OverDriveNormalAttack,
		OverDriveSuperSkillList:   logicCfg.OverDriveSuperSkill,
		BornBuffs:                 logicCfg.BornBuffs,
	}
	return info, nil
}

func (team *_Team) getMonsterSkillsByBD(skills, buildIds []uint32) []uint32 {
	result := make([]uint32, 0)
	for _, skillID := range skills {
		for _, bd := range buildIds {
			newSkillID := Data.GetSkillIDByBuildID(skillID, bd)
			result = append(result, newSkillID)
		}
	}
	return result
}
