package room

import (
	"Cinder/Matcher/matchapi/mtypes"
	"Cinder/Matcher/matcherlib/internal/rpcproc/notify"
)

type _RoomNotifier struct {
	notifier *notify.Notifier
}

func newRoomNotifier(srvIDToRoleIDs mtypes.SrvIDToRoleIDs) *_RoomNotifier {
	return &_RoomNotifier{
		notifier: notify.NewNotifier(srvIDToRoleIDs),
	}
}

func (r *_RoomNotifier) NotifyRolesJoinRoom(msg mtypes.RolesJoinRoomMsg) {
	r.notifier.PostToNotify(mtypes.NotifyMsg{RolesJoinRoomMsg: &msg})
}

func (r *_RoomNotifier) NotifyRoleLeaveRoom(msg mtypes.RoleLeaveRoomMsg) {
	r.notifier.PostToNotify(mtypes.NotifyMsg{RoleLeaveRoomMsg: &msg})
}

func (r *_RoomNotifier) NotifyBroadcastRoom(msg mtypes.BroadcastRoomMsg) {
	r.notifier.PostToNotify(mtypes.NotifyMsg{BroadcastRoomMsg: &msg})
}

func (r *_RoomNotifier) NotifySetRoomData(msg mtypes.SetRoomDataMsg) {
	r.notifier.PostToNotify(mtypes.NotifyMsg{SetRoomDataMsg: &msg})
}

func (r *_RoomNotifier) NotifyUpdateRoomRoleData(msg mtypes.UpdateRoomRoleDataMsg) {
	r.notifier.PostToNotify(mtypes.NotifyMsg{UpdateRoomRoleDataMsg: &msg})
}
