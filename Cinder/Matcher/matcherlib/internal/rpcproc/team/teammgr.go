package team

import (
	"Cinder/Matcher/matchapi/mtypes"
	"errors"
)

type _TeamMgr struct {
	teams *_TeamsSafe
}

var mgr = &_TeamMgr{
	teams: newTeamsSafe(),
}

func GetMgr() *_TeamMgr {
	return mgr
}

func (t *_TeamMgr) CreateTeam(roleInfo mtypes.RoleInfo, teamInfo mtypes.TeamInfo) (mtypes.TeamInfo, error) {
	team := newTeam(teamInfo)
	t.teams.Insert(team)
	team.Add(roleInfo)
	return team.Info(), nil
}

func (t *_TeamMgr) JoinTeam(roleInfo mtypes.RoleInfo, teamID mtypes.TeamID, passwd string) (mtypes.TeamInfo, error) {
	info := mtypes.TeamInfo{}
	team := t.teams.GetTeam(teamID)
	if team == nil {
		return info, errors.New("team does not exist")
	}
	if team.IsFull() {
		return info, errors.New("team is already full")
	}
	if team.Passwd() != passwd {
		return info, errors.New("wrong password")
	}
	team.Add(roleInfo)
	return team.Info(), nil
}

func (t *_TeamMgr) LeaveTeam(roleID mtypes.RoleID, teamID mtypes.TeamID) {
	team := t.teams.GetTeam(teamID)
	if team == nil {
		return
	}
	team.Del(roleID)
	if team.IsEmpty() {
		t.teams.Delete(teamID)
	}
}

func (t *_TeamMgr) GetTeamInfo(teamID mtypes.TeamID) (mtypes.TeamInfo, bool) {
	team := t.teams.GetTeam(teamID)
	if team == nil {
		return mtypes.TeamInfo{}, false
	}
	return team.Info(), true
}

func (t *_TeamMgr) ChangeTeamLeader(teamID mtypes.TeamID, newLeaderID mtypes.RoleID) error {
	team := t.teams.GetTeam(teamID)
	if team == nil {
		return errors.New("team does not exist")
	}
	return team.ChangeLeaderTo(newLeaderID)
}

func (t *_TeamMgr) SetTeamData(teamID mtypes.TeamID, key string, data interface{}) {
	team := t.teams.GetTeam(teamID)
	if team != nil {
		team.SetTeamData(key, data)
	}
}

func (t *_TeamMgr) BroadcastTeam(teamID mtypes.TeamID, msg interface{}) {
	team := t.teams.GetTeam(teamID)
	if team != nil {
		team.Broadcast(msg)
	}
}

func (t *_TeamMgr) GetTeamCount() int {
	return t.teams.GetCount()
}
