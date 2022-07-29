package team

import (
	"Cinder/Base/Util"
	"Cinder/Matcher/matchapi/mtypes"
	"errors"
	"sync"
)

// _Team 组队
type _Team struct {
	mtx  sync.Mutex
	info mtypes.TeamInfo
}

func newTeam(info mtypes.TeamInfo) *_Team {
	info.TeamID = mtypes.TeamID(Util.GetGUID())
	info.Roles = mtypes.RoleMap{}

	return &_Team{
		info: info,
	}
}

func (t *_Team) ID() mtypes.TeamID {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	return t.info.TeamID
}

func (t *_Team) Add(roleInfo mtypes.RoleInfo) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	if t.isFull() {
		return
	}

	roleInfo.TeamID = t.info.TeamID
	t.info.Roles[roleInfo.RoleID] = &roleInfo
	if t.info.LeaderID == "" {
		t.info.LeaderID = roleInfo.RoleID
	}
	t.newNotifier().NotifyJoinTeam(mtypes.JoinTeamMsg{
		RoleInfo: roleInfo,
	})
}

func (t *_Team) Del(roleID mtypes.RoleID) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	// 删除队长时先转移队长
	if roleID == t.info.LeaderID {
		t.changeLeader()
	}

	delete(t.info.Roles, roleID)
	t.newNotifier().NotifyLeaveTeam(mtypes.LeaveTeamMsg{
		RoleID: roleID,
	})
}

func (t *_Team) changeLeader() {
	for id := range t.info.Roles {
		if id != t.info.LeaderID {
			t.setLeaderID(id)
			return
		}
	}
}

func (t *_Team) IsFull() bool {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	return t.isFull()
}

func (t *_Team) isFull() bool {
	return len(t.info.Roles) >= t.info.MaxRole
}

func (t *_Team) IsEmpty() bool {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	return len(t.info.Roles) == 0
}

func (t *_Team) Info() mtypes.TeamInfo {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	return t.info.Copy()
}

func (t *_Team) newNotifier() *_TeamNotifier {
	return newTeamNotifier(t.info.GetSrvIDToRoleIDs())
}

func (t *_Team) Passwd() string {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	return t.info.Passwd
}

func (t *_Team) ChangeLeaderTo(newLeaderID mtypes.RoleID) error {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	if newLeaderID == t.info.LeaderID {
		return errors.New("already leader")
	}
	if _, ok := t.info.Roles[newLeaderID]; !ok {
		return errors.New("new leader is not a member")
	}
	t.setLeaderID(newLeaderID)
	return nil
}

func (t *_Team) setLeaderID(newLeaderID mtypes.RoleID) {
	t.info.LeaderID = newLeaderID
	t.newNotifier().NotifyChangeTeamLeader(mtypes.ChangeTeamLeaderMsg{
		NewLeaderID: newLeaderID,
	})
}

func (t *_Team) SetTeamData(key string, data interface{}) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	t.info.SetData(key, data)
	t.newNotifier().NotifySetTeamData(mtypes.SetTeamDataMsg{
		TeamID: t.info.TeamID,
		Key:    key,
		Data:   data,
	})
}

func (t *_Team) Broadcast(msg interface{}) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	t.newNotifier().NotifyBroadcastTeam(mtypes.BroadcastTeamMsg{
		Msg: msg,
	})
}
