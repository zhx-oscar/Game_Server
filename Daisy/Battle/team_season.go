package main

import (
	cConst "Cinder/Base/Const"
	"Cinder/Space"
	"Daisy/ActivityTimer"
	"Daisy/Const"
	"Daisy/DB"
	"Daisy/Data"
	"Daisy/DataTables"
	"Daisy/Prop"
	"Daisy/Proto"
	"math"
	"time"
)

const (
	StartFirstSeason             = 3
	SeasonLast                   = 4
	SeasonInterval               = 5
	SeasonLevelAwardMailTile     = 6
	SeasonLevelAwardMailContent  = 7
	SeasonStageAwardMailTile     = 8
	SeasonStageAwardMailContent  = 9
	SeasonLegendAwardMailTile    = 10
	SeasonLegendAwardMailContent = 11
)

func init() {
	ActivityTimer.GetActivities().RegisterActivity(&SeasonActivity{})
}

func getSeasonGeneralString(id uint32) string {
	config, ok := Data.GetSeasonConfig().SeasonGeneral_ConfigItems[id]
	if ok == false {
		return ""
	}
	return config.Characters
}

type SeasonActivity struct {
}

func (act *SeasonActivity) Init() error {
	return nil
}
func (act *SeasonActivity) Timer() {

}
func (act *SeasonActivity) Start() {
}
func (act *SeasonActivity) End() {
}
func (act *SeasonActivity) GetKey() string {
	return Const.SEASON_KEY
}

func (act *SeasonActivity) GetID() uint32 {
	return Const.SEASON
}

func (act *SeasonActivity) GetStartTime() int64 {
	config, ok := Data.GetSeasonConfig().SeasonGeneral_ConfigItems[StartFirstSeason]
	if ok == false {
		return 0
	}

	return int64(config.Value)*3600 + Const.GetZoneOpenTime().Unix() //加上开区时间
}

func (act *SeasonActivity) GetEndTime() int64 {
	return 0
}

func (act *SeasonActivity) GetLast() uint32 {
	config, ok := Data.GetSeasonConfig().SeasonGeneral_ConfigItems[SeasonLast]
	if ok == false {
		return 0
	}
	return config.Value * 3600
}
func (act *SeasonActivity) GetInterval() uint32 {
	config, ok := Data.GetSeasonConfig().SeasonGeneral_ConfigItems[SeasonInterval]
	if ok == false {
		return 0
	}
	return config.Value * 3600
}
func (act *SeasonActivity) GetLoop() uint32 {
	return uint32(len(Data.GetSeasonConfig().SeasonID_ConfigItems))
}

//CheckSeasonChange 检测赛季状态变化
func (team *_Team) CheckSeasonChange() {
	activity := ActivityTimer.GetActivities().GetActivity(Const.SEASON)
	if activity == nil {
		return
	}
	if team.prop.Data.SeasonInfo.SeasonID == 0 {
		//第一次初始化
		team.prop.Data.SeasonInfo.SeasonID = activity.GetStep()
		team.prop.Data.SeasonInfo.SeasonState = activity.IsActive()
		team.prop.SyncUpdateSeasonInfo(team.prop.Data.SeasonInfo)
	}
	//赛季状态变化
	if team.prop.Data.SeasonInfo.SeasonState == true && (activity.IsActive() == false || team.prop.Data.SeasonInfo.SeasonID != activity.GetStep()) {
		team.Infof("赛季:队伍进入赛季结束清算：seasonID=%d, state=%d", team.prop.Data.SeasonInfo.SeasonID, team.prop.Data.SeasonInfo.SeasonState)
		team.SeasonEnd()
	}

	if activity.IsActive() == true && team.prop.Data.SeasonInfo.SeasonState == false {
		team.Infof("赛季:队伍进入赛季开始状态：seasonID=%d, state=%d", team.prop.Data.SeasonInfo.SeasonID, team.prop.Data.SeasonInfo.SeasonState)
		team.SeasonStart()
	}
	//team.Infof("赛季:seasonID=%d,endTime=%d,now=%d, state=%d", activity.GetStep(), activity.GetEndTime(), time.Now().Unix(), activity.IsActive())
}
func (team *_Team) SeasonEnd() {
	Season := ActivityTimer.GetActivities().GetActivity(Const.SEASON).GetStep()
	team.prop.SyncChangeSeasonState(false)

	team.TraversalActor(func(actor Space.IActor) {
		role := actor.(*_Role)
		role.SeasonEnd(Season)
	})
	team.Info("赛季结束，id=", team.prop.Data.SeasonInfo.SeasonID)
}

func (team *_Team) SeasonStart() {
	Season := ActivityTimer.GetActivities().GetActivity(Const.SEASON).GetStep()
	team.prop.SyncChangeSeasonState(true)
	team.SeasonChange(Season)
	team.TraversalActor(func(actor Space.IActor) {
		role := actor.(*_Role)
		role.SeasonStart(Season)
	})
	team.Info("赛季开始，id=", team.prop.Data.SeasonInfo.SeasonID)
}

func (team *_Team) SeasonChange(Season uint32) {
	//更新状态
	team.prop.Data.SeasonInfo.SeasonID = Season
	team.prop.Data.SeasonInfo.SeasonTopLevel = team.GetSeasonLevelByScore(Season, team.prop.Data.SeasonInfo.TeamScore)
	team.prop.Data.SeasonInfo.SeasonTopLevel = team.prop.Data.SeasonInfo.SeasonLevel
	team.prop.SyncUpdateSeasonInfo(team.prop.Data.SeasonInfo)
}

//Pass过关
func (team *_Team) Pass(raidPro uint32) {
	//赛季休整期
	if team.prop.Data.SeasonInfo.SeasonState == false {
		return
	}

	//计算分数
	config, ok := Data.GetSceneConfig().BattleArea_ConfigItems[raidPro]
	if ok == false {
		return
	}
	score := config.SeasonSceore

	team.ChangeSeasonScore(team.prop.Data.SeasonInfo.TeamScore + score)

	team.TraversalActor(func(actor Space.IActor) {
		role := actor.(*_Role)
		//重复打不加分
		if role.prop.Data.Base.RaidProgress <= raidPro {
			role.prop.SyncAddSeasonScore(score)
			user := role.GetOwnerUser()
			if user == nil {
				role.Error("[findNewTitle] role's user is nil")
			} else {
				user.Rpc(cConst.Game, "RPC_UpdateChatUserActivateData")
			}
		}
	})
}

func (team *_Team) ChangeSeasonScore(num uint32) uint32 {
	old := team.prop.Data.SeasonInfo.TeamScore
	seasonID := team.prop.Data.SeasonInfo.SeasonID
	oldLevel := team.GetSeasonLevel(seasonID)
	config := team.GetSeasonLevelConfig(seasonID, oldLevel) //等级为0时，会返回nil

	team.prop.SyncChangeSeasonScore(num)
	team.Info("[赛季],队伍积分改变.%d->%d", old, team.prop.Data.SeasonInfo.TeamScore)
	//赛季休整期
	if team.prop.Data.SeasonInfo.SeasonState == false {
		return team.prop.Data.SeasonInfo.TeamScore
	}

	team.prop.Data.SeasonInfo.SeasonLevel = team.GetSeasonLevel(seasonID)
	if team.prop.Data.SeasonInfo.SeasonTopLevel < team.prop.Data.SeasonInfo.SeasonLevel {
		team.prop.Data.SeasonInfo.SeasonTopLevel = team.prop.Data.SeasonInfo.SeasonLevel
	}
	team.prop.SyncUpdateSeasonInfo(team.prop.Data.SeasonInfo)
	//加分
	if team.prop.Data.SeasonInfo.TeamScore > old {
		if team.CanInLegendRank() == true {
			DB.DefaultIRankUtil.UpdateMemberData(DB.SeasonLegendRank, team.prop.Data.SeasonInfo.SeasonID, team.GetID(), team.GetSeasonRankData())
			DB.DefaultIRankUtil.UpdateScore(DB.SeasonLegendRank, team.prop.Data.SeasonInfo.SeasonID, team.GetSeasonRankScore(), team.GetID())
		}
	} else {
		if config != nil && config.LegendRank == true && team.CanInLegendRank() == false {
			// 删除排行榜
			DB.DefaultIRankUtil.Remove(DB.SeasonLegendRank, team.prop.Data.SeasonInfo.SeasonID, team.GetID())
			team.Info("[赛季],积分下降，退出传说排行榜")
		}
	}

	DB.DefaultIRankUtil.UpdateScore(DB.SeasonRank, team.prop.Data.SeasonInfo.SeasonID, team.GetSeasonRankScore(), team.GetID()) //普通排行榜只有teamid
	return team.prop.Data.SeasonInfo.TeamScore
}

func (team *_Team) GetSeasonRankScore() float64 {
	return float64(team.prop.Data.SeasonInfo.TeamScore) + 0.1/(float64(time.Now().Unix()-Const.GetZoneOpenTime().Unix()))
}

func (team *_Team) GetSeasonRankData() []byte {
	data := &Proto.RankTeamData{}
	data.ID = team.GetID()
	data.Name = team.prop.Data.Base.Name
	data.Score = team.prop.Data.SeasonInfo.TeamScore
	data.Members = make([]*Proto.RankMemberData, 0)
	team.TraversalActor(func(actor Space.IActor) {
		role := actor.(*_Role)
		_data := role.GetSeasonRankData()
		_data.Status = team.prop.Data.Base.Members[_data.ID].Status
		data.Members = append(data.Members, _data)
	})

	byData, err := data.Marshal()
	if err != nil {
		team.Error("GetSeasonRankData Marshal err=", err)
	}
	return byData
}

func (team *_Team) GetRankList() *Proto.RankData {
	data := &Proto.RankData{}

	data.MyPlace = team.GetLegendRankPlace(team.prop.Data.SeasonInfo.SeasonID)
	data.SeasonID = team.prop.Data.SeasonInfo.SeasonID
	if team.prop.Data.SeasonInfo.SeasonState == true {
		act := ActivityTimer.GetActivities().GetActivity(Const.SEASON)
		if act != nil {
			data.LeftTime = uint32(act.GetEndTime() - time.Now().Unix())
		}
	}
	data.LastSeason = make([]*Proto.RankTeamData, 0)
	data.ThisSeason = make([]*Proto.RankTeamData, 0)

	list := DB.DefaultIRankUtil.GetList(DB.SeasonLegendRank, team.prop.Data.SeasonInfo.SeasonID, 100)
	for _, v := range list {
		_data := &Proto.RankTeamData{}
		err := _data.Unmarshal([]byte(v))
		if err == nil {
			data.ThisSeason = append(data.ThisSeason, _data)
		} else {
			team.Debug("GetRankList Unmarshal err=", err)
		}
	}

	if team.prop.Data.SeasonInfo.SeasonID > 1 {
		_list := DB.DefaultIRankUtil.GetList(DB.SeasonLegendRank, team.prop.Data.SeasonInfo.SeasonID-1, 100)
		for _, v := range _list {
			_data := &Proto.RankTeamData{}
			err := _data.Unmarshal([]byte(v))
			if err == nil {
				data.LastSeason = append(data.LastSeason, _data)
			} else {
				team.Debug("GetRankList Unmarshal err=", err)
			}
		}
	}

	return data
}

func (team *_Team) GetSeasonPlace(seasonID uint32) uint32 {
	return DB.DefaultIRankUtil.GetPlacePercentage(DB.SeasonRank, seasonID, team.GetID())
}

func (team *_Team) GetLegendRankPlace(seasonID uint32) uint32 {
	return DB.DefaultIRankUtil.GetPlace(DB.SeasonLegendRank, seasonID, team.GetID())
}

func (team *_Team) GetSeasonLevel(seasonID uint32) uint32 {
	data, ok := Data.SeasonLevelData[seasonID]
	if ok == false {
		return 0
	}
	level := uint32(0)
	for {
		_data, _ok := data[level+1]
		if _ok == false {
			break
		}
		if team.prop.Data.SeasonInfo.TeamScore < _data.SeasonGoalScore {
			return level
		}
		level++
	}
	return level
}

func (team *_Team) CanInLegendRank() bool {
	config := team.GetSeasonLevelConfig(team.prop.Data.SeasonInfo.SeasonID, team.GetSeasonLevel(team.prop.Data.SeasonInfo.SeasonID))
	if config == nil {
		return false
	}
	return config.LegendRank
}

func (team *_Team) GetSeasonLevelConfig(seasionID, level uint32) *DataTables.SeasonLevel_Config {
	data, ok := Data.SeasonLevelData[seasionID]
	if ok == false {
		return nil
	}
	_data, _ok := data[level]
	if _ok == false {
		return nil
	}
	return _data
}

func (team *_Team) GetSeasonLevelByScore(seasonID, score uint32) uint32 {
	data, ok := Data.SeasonLevelData[seasonID]
	if ok == false {
		return 0
	}
	level := uint32(0)
	for {
		_data, _ok := data[level+1]
		if _ok == false {
			break
		}
		if score < _data.SeasonGoalScore {
			return level
		}
		level++
	}
	return level
}

func (team *_Team) UpdateSeasonDataOnMemberChange() {
	team.Debug("UpdateSeasonDataOnMemberChange:", team.GetID())
	if len(team.prop.Data.Base.Members) == 0 {
		// 会有空队伍存在，空队伍不做处理
		return
	}
	minScore := uint32(math.MaxUint32)
	team.TraversalActor(func(ia Space.IActor) {
		if _, ok := team.prop.Data.Base.Members[ia.GetID()]; ok {
			progress := ia.GetProp().(*Prop.RoleProp).Data.SeasonInfo.MyScore
			if minScore > progress {
				minScore = progress
			}
		}
	})
	team.Debug("赛季积分改变:", team.GetID(), minScore)
	team.ChangeSeasonScore(minScore)
}
