package main

import (
	"Cinder/Base/Const"
	"Cinder/Base/User"
	"Cinder/plugin/navmesh"
	"Cinder/plugin/physxgo"
	"Daisy/Battle/sceneconfigdata"
	"Daisy/Fight"
	"Daisy/Proto"
	"errors"
	"math/rand"
	"time"
)

type _RunModel struct {
	query   *navmesh.Query
	pxScene physxgo.IPxScene
	paths   *sceneconfigdata.PathCfg

	positive        bool //正向or反向
	targetRoadIndex uint32

	//GM调试
	debugRoadIndex int //仅用于调试的路点，默认-1

	spawnMonsterInfos map[uint32][]*Fight.PawnInfo
}

func NewRunModel() *_RunModel {
	return &_RunModel{
		positive:          true,
		debugRoadIndex:    -1,
		spawnMonsterInfos: make(map[uint32][]*Fight.PawnInfo),
	}
}

func (run *_RunModel) Destroy() {
	if run.query != nil {
		navmesh.DestroyQuery(run.query)
		run.query = nil
	}
}

func (team *_Team) onTriggerRaidBattle(isBoss bool, battleFieldID uint32) {
	if len(team.PendingLeaveUser) > 0 {
		team.Debug("队伍有成员在转移中，不触发战斗")
		return
	}

	if team.CanSetState(TeamState_Raidbattling) {
		team.SetState(TeamState_Raidbattling, isBoss, battleFieldID)
	}
}

func (team *_Team) onTriggerChest(param *sceneconfigdata.TriggerParamChest) {
	team.EnterRunChest(param)
}

func (team *_Team) onSetLocation(curLoc, targetLoc *Proto.PVector3, velocity float32, targetRoadIndex uint32) {
	team.prop.SetLocation(curLoc, targetLoc, velocity, targetRoadIndex)
}

func (team *_Team) RandomSpawnIdx() uint32 {
	if len(team.paths.FirstStartIdx) == 0 {
		panic("跑图没有配置出生点")
	}

	idx := rand.New(rand.NewSource(time.Now().Unix())).Intn(len(team.paths.FirstStartIdx))
	return team.paths.FirstStartIdx[idx]
}

//GetNextRoadIndexs 根据地图配置的生效方向获取nexts或者previous
func (team *_Team) GetNextRoadIndexs(in uint32) []uint32 {
	if team.positive {
		return team.paths.Path[in].Nexts
	} else {
		return team.paths.Path[in].Previous
	}
}

func (team *_Team) SetRoadByName(name string) error {
	if team.GetState() != TeamState_Running {
		team.Errorf("不在跑图中无法设置路点:%s", name)
		return errors.New("不在跑图中无法设置路点")
	}

	idx := -1
	for _, value := range team.paths.Path {
		if value.Name == name {
			idx = int(value.Idx)
		}
	}

	if idx == -1 {
		team.Errorf("设置路点名:%s 不存在", name)
		return errors.New("设置路点名不存在")
	}

	team.debugRoadIndex = idx
	team.Infof("设置路点 name:%s idx:%d", name, idx)
	return nil
}

func (team *_Team) onTriggerSpawnMonster(param *sceneconfigdata.TriggerParamSpawnMonster) error {
	cfg, err := team.GetBattleFieldCfg(team.prop.Data.Raid.Progress, param.BattleID)
	if err != nil {
		team.Errorf("onTriggerSpawnMonster Progress:%d BattleFieldID:%d err:%s", team.prop.Data.Raid.Progress, param.BattleID, err)
		return err
	}

	enemys, err := team.buildRaidEnemyInfos(cfg, team.prop.Data.Raid.Progress, false, true)
	if err != nil {
		team.Errorf("onTriggerSpawnMonster Progress:%d BattleFieldID:%d err:%s", team.prop.Data.Raid.Progress, param.BattleID, err)
		return err
	}

	team.spawnMonsterInfos[param.BattleID] = enemys

	spawnInfo := &Proto.SpawnMonsterInfo{
		BattleFieldID: param.BattleID,
		PawnInfos:     make([]*Proto.PawnInfo, len(enemys), len(enemys)),
	}
	for i := 0; i < len(enemys); i++ {
		spawnInfo.PawnInfos[i] = enemys[i].PawnInfo
	}
	team.TraversalUser(func(iu User.IUser) bool {
		iu.Rpc(Const.Agent, "RPC_SpawnMonster", spawnInfo)
		return true
	})

	return nil
}
