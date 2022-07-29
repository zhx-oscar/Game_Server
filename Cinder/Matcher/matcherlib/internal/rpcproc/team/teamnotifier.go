package team

import (
	"Cinder/Matcher/matchapi/mtypes"
	"Cinder/Matcher/matcherlib/internal/rpcproc/notify"
)

type _TeamNotifier struct {
	notifier *notify.Notifier
}

func newTeamNotifier(srvIDToRoleIDs mtypes.SrvIDToRoleIDs) *_TeamNotifier {
	return &_TeamNotifier{
		notifier: notify.NewNotifier(srvIDToRoleIDs),
	}
}

func (r *_TeamNotifier) NotifyJoinTeam(msg mtypes.JoinTeamMsg) {
	go r.notifier.PostToNotify(mtypes.NotifyMsg{JoinTeamMsg: &msg})
}

func (r *_TeamNotifier) NotifyLeaveTeam(msg mtypes.LeaveTeamMsg) {
	go r.notifier.PostToNotify(mtypes.NotifyMsg{LeaveTeamMsg: &msg})
}

func (r *_TeamNotifier) NotifyChangeTeamLeader(msg mtypes.ChangeTeamLeaderMsg) {
	go r.notifier.PostToNotify(mtypes.NotifyMsg{ChangeTeamLeaderMsg: &msg})
}

func (r *_TeamNotifier) NotifyBroadcastTeam(msg mtypes.BroadcastTeamMsg) {
	go r.notifier.PostToNotify(mtypes.NotifyMsg{BroadcastTeamMsg: &msg})
}

func (r *_TeamNotifier) NotifySetTeamData(msg mtypes.SetTeamDataMsg) {
	go r.notifier.PostToNotify(mtypes.NotifyMsg{SetTeamDataMsg: &msg})
}
