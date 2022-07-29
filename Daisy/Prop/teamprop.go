package Prop

import (
	"Cinder/Base/Message"
	"Cinder/Base/Prop"
	"Daisy/Proto"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type TeamProp struct {
	Prop.Prop

	Data *Proto.Team
}

func (t *TeamProp) Marshal() ([]byte, error) {
	return t.Data.Marshal()
}

func (t *TeamProp) UnMarshal(data []byte) error {
	t.Data = &Proto.Team{}

	if len(data) != 0 {
		if err := t.Data.Unmarshal(data); err != nil {
			return err
		}
	}

	t.fillDefault()

	return nil
}

func (t *TeamProp) MarshalCache() ([]byte, error) {
	cache := &Proto.TeamCache{
		Base: t.Data.Base,
	}

	return cache.Marshal()
}

func (t *TeamProp) UnMarshalCache(data []byte) (interface{}, error) {
	cache := &Proto.TeamCache{}
	if err := cache.Unmarshal(data); err != nil {
		return nil, err
	}

	return cache, nil
}

func (t *TeamProp) MarshalToBson() ([]byte, error) {
	return bson.Marshal(t.Data)
}

func (t *TeamProp) UnMarshalFromBson(data []byte) error {
	t.Data = &Proto.Team{}

	if len(data) != 0 {
		if err := bson.Unmarshal(data, t.Data); err != nil {
			return err
		}
	}

	t.fillDefault()

	return nil
}

func (t *TeamProp) MarshalPart() ([]byte, error) {
	return t.Data.Marshal()
}

func (t *TeamProp) UnMarshalPart(data []byte) error {
	return t.Data.Unmarshal(data)
}

// 填充各种字段
func (t *TeamProp) fillDefault() {
	if t.Data.Base == nil {
		t.Data.Base = &Proto.TeamBase{
			Members: make(map[string]*Proto.TeamMemberInfo),
		}
	}
	if t.Data.Base.Members == nil {
		t.Data.Base.Members = make(map[string]*Proto.TeamMemberInfo)
	}

	if t.Data.Run == nil {
		t.Data.Run = &Proto.TeamRun{
			OriginLoaction: &Proto.PVector3{},
			TargetLocation: &Proto.PVector3{},
		}
	}

	if t.Data.Raid == nil {
		t.Data.Raid = &Proto.RaidInfo{
			Progress: 1,
		}
	}

	if t.Data.Guide == nil {
		t.Data.Guide = &Proto.GuideInfo{
			Step:     1,
			IsFinish: true,
		}
	}

	if t.Data.Applys == nil {
		t.Data.Applys = make(map[string]*Proto.TeamApplyInfo)
	}

	if t.Data.Supply == nil {
		t.Data.Supply = &Proto.Supply{}
	}
	if t.Data.SeasonInfo == nil {
		t.Data.SeasonInfo = &Proto.TeamSeasonInfo{}
	}
}

func (t *TeamProp) SyncInitedOnce() {
	t.Sync("InitedOnce", Message.PackArgs(), true)
}

func (t *TeamProp) InitedOnce() {
	t.Data.InitedOnce = true
}

func (t *TeamProp) AddTeamMember(userID string, status uint32) {
	if _, ok := t.Data.Base.Members[userID]; ok {
		return
	}

	member := &Proto.TeamMemberInfo{
		Status:   status,
		JoinTime: time.Now().Unix(),
	}
	t.Sync("AddTeamMemberImpl", Message.PackArgs(userID, member), true, Prop.Target_All_Clients)
}

func (t *TeamProp) AddTeamMemberImpl(userID string, member *Proto.TeamMemberInfo) {
	t.Data.Base.Members[userID] = member
	t.Data.Base.Num = uint32(len(t.Data.Base.Members))
}

func (t *TeamProp) RemoveTeamMember(userID string) {
	if _, ok := t.Data.Base.Members[userID]; !ok {
		return
	}

	t.Sync("RemoveTeamMemberImpl", Message.PackArgs(userID), true, Prop.Target_All_Clients)
}

func (t *TeamProp) RemoveTeamMemberImpl(userID string) {
	delete(t.Data.Base.Members, userID)
	t.Data.Base.Num = uint32(len(t.Data.Base.Members))
}

func (t *TeamProp) SetLocation(originLocation *Proto.PVector3, targetLocation *Proto.PVector3, velocity float32, idx uint32) {
	t.Sync("SetLocationImpl", Message.PackArgs(originLocation, targetLocation, velocity, idx), false, Prop.Target_All_Clients)
}

func (t *TeamProp) SetLocationImpl(originLocation *Proto.PVector3, targetLocation *Proto.PVector3, velocity float32, idx uint32) {
	t.Data.Run.OriginLoaction = originLocation
	t.Data.Run.TargetLocation = targetLocation
	t.Data.Run.Velocity = velocity
}

func (t *TeamProp) SetGuideStep(step uint32) {
	t.Sync("SetGuideStepImpl", Message.PackArgs(step), true, Prop.Target_All_Clients)
}

func (t *TeamProp) SetGuideStepImpl(step uint32) {
	t.Data.Guide.Step = step
}

func (t *TeamProp) FinishGuide() {
	t.Sync("FinishGuideImpl", Message.PackArgs(), true, Prop.Target_All_Clients)
}

func (t *TeamProp) FinishGuideImpl() {
	t.Data.Guide.IsFinish = true
}

func (t *TeamProp) SyncAddApplyInfo(userID, message string) {
	info := &Proto.TeamApplyInfo{
		Time:    time.Now().Unix(),
		Message: message,
	}
	t.Sync("AddApplyInfo", Message.PackArgs(userID, info), false, Prop.Target_All_Clients)
}
func (t *TeamProp) AddApplyInfo(userID string, info *Proto.TeamApplyInfo) {
	t.Data.Applys[userID] = info
}

func (t *TeamProp) SyncRemoveApplyInfo(userID string) {
	t.Sync("RemoveApplyInfo", Message.PackArgs(userID), false, Prop.Target_All_Clients)
}
func (t *TeamProp) RemoveApplyInfo(userID string) {
	delete(t.Data.Applys, userID)
}

func (t *TeamProp) SyncClearApplyInfo() {
	t.Sync("ClearApplyInfo", Message.PackArgs(), false, Prop.Target_All_Clients)
}
func (t *TeamProp) ClearApplyInfo() {
	t.Data.Applys = make(map[string]*Proto.TeamApplyInfo)
}

func (t *TeamProp) SyncSetRaidProgress(progress uint32) {
	if t.Data.Raid.Progress != progress {
		t.Sync("SetRaidProgress", Message.PackArgs(progress), true, Prop.Target_All_Clients)
	}
}

func (t *TeamProp) SetRaidProgress(progress uint32) {
	t.Data.Raid.Progress = progress
}

func (t *TeamProp) SyncSetOwnTickets(ownTickets uint32) {
	t.Sync("SetOwnTickets", Message.PackArgs(ownTickets), true, Prop.Target_All_Clients)
}

func (t *TeamProp) SetOwnTickets(ownTickets uint32) {
	t.Data.Raid.OwnTickets = ownTickets
}

func (t *TeamProp) SyncSetLastLogoutTime(time int64) {
	t.Sync("SetLastLogoutTime", Message.PackArgs(time), true, Prop.Target_All_Clients)
}

func (t *TeamProp) SetLastLogoutTime(time int64) {
	t.Data.Base.LastDestoryTime = time
}

func (t *TeamProp) SyncSetUID(uid uint64) {
	if t.Data.Base.UID != uid {
		t.Sync("SetUID", Message.PackArgs(uid), true)
	}
}

func (t *TeamProp) SetUID(uid uint64) {
	t.Data.Base.UID = uid
}

func (t *TeamProp) SyncSetName(name string) {
	if t.Data.Base.Name != name {
		t.Sync("SetName", Message.PackArgs(name), true, Prop.Target_All_Clients)
	}
}

func (t *TeamProp) SetName(name string) {
	t.Data.Base.Name = name
}

func (t *TeamProp) SyncPublishTeam(board string, needAuth bool, autoJoinIdx uint32) {
	t.Sync("PublishTeam", Message.PackArgs(board, needAuth, autoJoinIdx), true, Prop.Target_All_Clients)
}

func (t *TeamProp) PublishTeam(board string, needAuth bool, autoJoinIdx uint32) {
	t.Data.Base.Published = true
	t.Data.Base.Board = board
	t.Data.Base.NeedAuth = needAuth
	t.Data.Base.AutoJoinIdx = autoJoinIdx
}

func (t *TeamProp) SyncSetStatus(userID string, status uint32) {
	t.Sync("SetStatus", Message.PackArgs(userID, status), true, Prop.Target_All_Clients)
}

func (t *TeamProp) SetStatus(userID string, status uint32) {
	if member, ok := t.Data.Base.Members[userID]; ok {
		member.Status = status
	}
}

func (t *TeamProp) SyncSetFastBattleResetTimestamp(ts int64) {
	t.Sync("SetFastBattleResetTimestamp", Message.PackArgs(ts), true, Prop.Target_All_Clients)
}

func (t *TeamProp) SetFastBattleResetTimestamp(ts int64) {
	t.Data.Raid.FastBattleResetTimestamp = ts
}

func (t *TeamProp) SyncResetMemberDaily(time int64) {
	t.Sync("ResetMemberDaily", Message.PackArgs(time), true)
}

func (t *TeamProp) ResetMemberDaily(nextTimestamp int64) {
	t.Data.Base.DailyResetTimestamp = nextTimestamp
}

func (t *TeamProp) SyncSetSupplyResetTimestamp(ts int64) {
	t.Sync("SetSupplyResetTimestamp", Message.PackArgs(ts), true, Prop.Target_All_Clients)
}

func (t *TeamProp) SetSupplyResetTimestamp(ts int64) {
		t.Data.Supply.SupplyResetTimestamp = ts
	}

func (t *TeamProp) SyncUpdateSeasonInfo(info *Proto.TeamSeasonInfo) {
	t.Sync("UpdateSeasonInfo", Message.PackArgs(info), true, Prop.Target_All_Clients)
}
func (t *TeamProp) UpdateSeasonInfo(info *Proto.TeamSeasonInfo) {
	t.Data.SeasonInfo = info
}
func (t *TeamProp) SyncChangeSeasonScore(score uint32) {
	t.Sync("ChangeSeasonScore", Message.PackArgs(score), true, Prop.Target_All_Clients)
}
func (t *TeamProp) ChangeSeasonScore(score uint32) {
	t.Data.SeasonInfo.TeamScore = score
}

//SyncChangeSeasonState改变赛季状态
func (t *TeamProp)SyncChangeSeasonState(state bool){
	t.Sync("ChangeSeasonState", Message.PackArgs(state), true, Prop.Target_All_Clients)
}
func (t *TeamProp)ChangeSeasonState(state bool) {
	t.Data.SeasonInfo.SeasonState = state
}
