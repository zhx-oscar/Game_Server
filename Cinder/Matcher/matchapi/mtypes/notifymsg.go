package mtypes

import (
	"encoding/json"
)

// NotifyMsgsToOneSrv 是发送给某个SrvID的一批通知。合并多个通知减少消息量。
type NotifyMsgsToOneSrv struct {
	Msgs []NotifyMsgToOneSrv
}

// NotifyMsgToOneSrv 是发送给某个SrvID的一个通知。一个通知针对同一服上的多人。
type NotifyMsgToOneSrv struct {
	RoleIDs   []RoleID
	NotifyMsg NotifyMsg
}

type NotifyMsg struct {
	// 组队通知
	JoinTeamMsg         *JoinTeamMsg         `json:"joinTeamMsg,omitempty"`
	LeaveTeamMsg        *LeaveTeamMsg        `json:"leaveTeamMsg,omitempty"`
	ChangeTeamLeaderMsg *ChangeTeamLeaderMsg `json:"changeTeamLeaderMsg,omitempty"`
	SetTeamDataMsg      *SetTeamDataMsg      `json:"setTeamDataMsg,omitempty"`
	BroadcastTeamMsg    *BroadcastTeamMsg    `json:"broadcastTeamMsg,omitempty"`

	// 房间通知
	RolesJoinRoomMsg      *RolesJoinRoomMsg      `json:"rolesJoinRoomMsg,omitempty"`
	RoleLeaveRoomMsg      *RoleLeaveRoomMsg      `json:"roleLeaveRoomMsg,omitempty"`
	BroadcastRoomMsg      *BroadcastRoomMsg      `json:"broadcastRoomMsg,omitempty"`
	SetRoomDataMsg        *SetRoomDataMsg        `json:"setRoomDataMsg,omitempty"`
	UpdateRoomRoleDataMsg *UpdateRoomRoleDataMsg `json:"updateRoomRoleDataMsg,omitempty"`
}

type (

	// 组队通知
	JoinTeamMsg struct {
		RoleInfo RoleInfo
	}
	LeaveTeamMsg struct {
		RoleID RoleID
	}
	ChangeTeamLeaderMsg struct {
		NewLeaderID RoleID
	}
	SetTeamDataMsg struct {
		TeamID TeamID
		Key    string
		Data   interface{}
	}
	BroadcastTeamMsg struct {
		Msg interface{}
	}

	// 房间通知
	// TODO: RolesJoinRoomMsg RolesLeaveRoomMsg 合并成一个，并且连续多个就忽略前面的，只需最后一个。
	RolesJoinRoomMsg struct {
		RoomInfo      RoomInfo
		JoinedRoleIDs []RoleID
	}
	RoleLeaveRoomMsg struct {
		RoomInfo      RoomInfo
		LeavingRoleID RoleID
	}
	BroadcastRoomMsg struct {
		Msg interface{}
	}
	SetRoomDataMsg struct {
		Key  string
		Data interface{}
	}
	UpdateRoomRoleDataMsg struct {
		Update RoomRoleDataUpdate
	}
)

type RoomRoleDataUpdate struct {
	RoomID RoomID
	RoleID RoleID

	FloatData  map[string]float64 // 更改数值
	StringData map[string]string  // 更改字符串
	AddTags    []string           // 添加标签
	DelTags    []string           // 删除标签
}

func UnmarshalNotifyOneSrvMsg(buf []byte) (NotifyMsgsToOneSrv, error) {
	result := NotifyMsgsToOneSrv{}
	if err := json.Unmarshal(buf, &result); err != nil {
		return result, err
	}
	return result, nil
}
