// rpcmsg 定义Matcher RPC 的请求与应答消息
package rpcmsg

import (
	"Cinder/Matcher/matchapi/mtypes"
	"encoding/json"
)

// Team service messages
type (
	CreateTeamReq struct {
		CreatorInfo mtypes.RoleInfo
		TeamInfo    mtypes.TeamInfo
	}
	CreateTeamRsp struct {
		TeamInfo mtypes.TeamInfo
	}

	JoinTeamReq struct {
		RoleInfo mtypes.RoleInfo
		TeamID   mtypes.TeamID
		Password string
	}
	JoinTeamRsp struct {
		TeamInfo mtypes.TeamInfo
	}

	SetTeamDataReq struct {
		TeamID mtypes.TeamID
		Key    string
		Data   interface{}
	}
	SetTeamDataRsp struct{}

	BroadcastTeamReq struct {
		TeamID mtypes.TeamID
		Msg    interface{}
	}
	BroadcastTeamRsp struct{}
)

// Room service messages
type (
	RoleJoinRandomRoomReq struct {
		RoleInfo   mtypes.RoleInfo
		CreateInfo mtypes.CreateRoomInfo
	}
	RoleJoinRandomRoomRsp struct {
		RoomInfo mtypes.RoomInfo
	}

	RoleCreateRoomReq struct {
		CreatorInfo mtypes.RoleInfo
		CreateInfo  mtypes.CreateRoomInfo
	}
	RoleCreateRoomRsp struct {
		RoomInfo mtypes.RoomInfo
	}

	RoleJoinRoomReq struct {
		RoleInfo mtypes.RoleInfo
		RoomID   mtypes.RoomID
	}
	RoleJoinRoomRsp struct {
		RoomInfo mtypes.RoomInfo
	}

	RoleLeaveRoomReq struct {
		RoleID mtypes.RoleID
		RoomID mtypes.RoomID
	}
	RoleLeaveRoomRsp struct {
		RoomInfo mtypes.RoomInfo
	}

	TeamJoinRandomRoomReq struct {
		TeamID     mtypes.TeamID
		CreateInfo mtypes.CreateRoomInfo
	}
	TeamJoinRandomRoomRsp struct {
		RoomInfo mtypes.RoomInfo
	}

	TeamCreateRoomReq struct {
		CreatorTeamID mtypes.TeamID
		CreateInfo    mtypes.CreateRoomInfo
	}
	TeamCreateRoomRsp struct {
		RoomInfo mtypes.RoomInfo
	}

	TeamJoinRoomReq struct {
		TeamID mtypes.TeamID
		RoomID mtypes.RoomID
	}
	TeamJoinRoomRsp struct {
		RoomInfo mtypes.RoomInfo
	}

	BroadcastRoomReq struct {
		RoomID mtypes.RoomID
		Msg    interface{}
	}
	BroadcastRoomRsp struct{}

	GetRoomListReq struct {
		MatchMode mtypes.MatchMode
	}
	GetRoomListRsp struct {
		RoomIDs []mtypes.RoomID
	}

	ListRoomsReq struct {
		MatchModes []mtypes.MatchMode
	}
	ListRoomsRsp struct {
		Rooms map[mtypes.MatchMode][]mtypes.RoomInfo
	}

	GetRoomInfosReq struct {
		RoomIDs []mtypes.RoomID
	}
	GetRoomInfosRsp struct {
		RoomInfos map[mtypes.RoomID]mtypes.RoomInfo
	}

	SetRoomDataReq struct {
		RoomID mtypes.RoomID
		Key    string
		Data   interface{}
	}
	SetRoomDataRsp struct {
	}

	UpdateRoomRoleDataReq struct {
		Update mtypes.RoomRoleDataUpdate
	}
	UpdateRoomRoleDataRsp struct{}
)

type RPCResponse struct {
	// Team service response
	CreateTeamRsp *CreateTeamRsp `json:"createTeamRsp,omitempty"`
	JoinTeamRsp   *JoinTeamRsp   `json:"joinTeamRsp,omitempty"`

	// Room service responses

	RoleJoinRandomRoomRsp *RoleJoinRandomRoomRsp `json:"roleJoinRandomRoomRsp,omitempty"`
	RoleCreateRoomRsp     *RoleCreateRoomRsp     `json:"roleCreateRoomRsp,omitempty"`
	RoleJoinRoomRsp       *RoleJoinRoomRsp       `json:"roleJoinRoomRsp,omitempty"`
	RoleLeaveRoomRsp      *RoleLeaveRoomRsp      `json:"roleLeaveRoomRsp,omitempty"`

	TeamJoinRandomRoomRsp *TeamJoinRandomRoomRsp `json:"teamJoinRandomRoomRsp,omitempty"`
	TeamCreateRoomRsp     *TeamCreateRoomRsp     `json:"teamCreateRoomRsp,omitempty"`
	TeamJoinRoomRsp       *TeamJoinRoomRsp       `json:"teamJoinRoomRsp,omitempty"`

	GetRoomListRsp        *GetRoomListRsp        `json:"getRoomListRsp,omitempty"`
	ListRoomsRsp          *ListRoomsRsp          `json:"listRoomsRsp,omitempty"`
	GetRoomInfosRsp       *GetRoomInfosRsp       `json:"getRoomInfos,omitempty"`
	SetRoomDataRsp        *SetRoomDataRsp        `json:"setRoomDataRsp,omitempty"`
	UpdateRoomRoleDataRsp *UpdateRoomRoleDataRsp `json:"setRoomRoleDataRsp,omitempty"`
}

func IsRequest(i interface{}) bool {
	switch i.(type) {
	case BroadcastRoomReq,
		*BroadcastRoomReq,
		BroadcastTeamReq,
		*BroadcastTeamReq,
		CreateTeamReq,
		*CreateTeamReq,
		GetRoomInfosReq,
		*GetRoomInfosReq,
		GetRoomListReq,
		*GetRoomListReq,
		ListRoomsReq,
		*ListRoomsReq,
		JoinTeamReq,
		*JoinTeamReq,
		RoleCreateRoomReq,
		*RoleCreateRoomReq,
		RoleJoinRandomRoomReq,
		*RoleJoinRandomRoomReq,
		RoleJoinRoomReq,
		*RoleJoinRoomReq,
		RoleLeaveRoomReq,
		*RoleLeaveRoomReq,
		SetRoomDataReq,
		*SetRoomDataReq,
		SetTeamDataReq,
		*SetTeamDataReq,
		UpdateRoomRoleDataReq,
		*UpdateRoomRoleDataReq,
		TeamCreateRoomReq,
		*TeamCreateRoomReq,
		TeamJoinRandomRoomReq,
		*TeamJoinRandomRoomReq,
		TeamJoinRoomReq,
		*TeamJoinRoomReq:
		return true
	}
	return false
}

func (r RPCResponse) Marshal() ([]byte, string) {
	buf, err := json.Marshal(r)
	if err != nil {
		return nil, err.Error()
	}
	return buf, ""
}

func (r *RPCResponse) Unmarshal(buf []byte) error {
	if err := json.Unmarshal(buf, r); err != nil {
		return err
	}
	return nil
}

func (r *RPCResponse) GetCreateTeamRsp() CreateTeamRsp {
	if r != nil && r.CreateTeamRsp != nil {
		return *r.CreateTeamRsp
	}
	return CreateTeamRsp{}
}

func (r *RPCResponse) GetJoinTeamRsp() JoinTeamRsp {
	if r != nil && r.JoinTeamRsp != nil {
		return *r.JoinTeamRsp
	}
	return JoinTeamRsp{}
}

func (r *RPCResponse) GetRoleJoinRandomRoomRsp() RoleJoinRandomRoomRsp {
	if r != nil && r.RoleJoinRandomRoomRsp != nil {
		return *r.RoleJoinRandomRoomRsp
	}
	return RoleJoinRandomRoomRsp{}
}

func (r *RPCResponse) GetRoleCreateRoomRsp() RoleCreateRoomRsp {
	if r != nil && r.RoleCreateRoomRsp != nil {
		return *r.RoleCreateRoomRsp
	}
	return RoleCreateRoomRsp{}
}

func (r *RPCResponse) GetRoleJoinRoomRsp() RoleJoinRoomRsp {
	if r != nil && r.RoleJoinRoomRsp != nil {
		return *r.RoleJoinRoomRsp
	}
	return RoleJoinRoomRsp{}
}

func (r *RPCResponse) GetRoleLeaveRoomRsp() RoleLeaveRoomRsp {
	if r != nil && r.RoleLeaveRoomRsp != nil {
		return *r.RoleLeaveRoomRsp
	}
	return RoleLeaveRoomRsp{}
}

func (r *RPCResponse) GetTeamJoinRandomRoomRsp() TeamJoinRandomRoomRsp {
	if r != nil && r.TeamJoinRandomRoomRsp != nil {
		return *r.TeamJoinRandomRoomRsp
	}
	return TeamJoinRandomRoomRsp{}
}

func (r *RPCResponse) GetTeamCreateRoomRsp() TeamCreateRoomRsp {
	if r != nil && r.TeamCreateRoomRsp != nil {
		return *r.TeamCreateRoomRsp
	}
	return TeamCreateRoomRsp{}
}

func (r *RPCResponse) GetTeamJoinRoomRsp() TeamJoinRoomRsp {
	if r != nil && r.TeamJoinRoomRsp != nil {
		return *r.TeamJoinRoomRsp
	}
	return TeamJoinRoomRsp{}
}

func (r *RPCResponse) GetGetRoomListRsp() GetRoomListRsp {
	if r != nil && r.GetRoomListRsp != nil {
		return *r.GetRoomListRsp
	}
	return GetRoomListRsp{}
}

func (r *RPCResponse) GetListRoomsRsp() ListRoomsRsp {
	if r != nil && r.ListRoomsRsp != nil {
		return *r.ListRoomsRsp
	}
	return ListRoomsRsp{}
}

func (r *RPCResponse) GetGetRoomInfosRsp() GetRoomInfosRsp {
	if r != nil && r.GetRoomInfosRsp != nil {
		return *r.GetRoomInfosRsp
	}
	return GetRoomInfosRsp{}
}

func (r *RPCResponse) GetSetRoomDataRsp() SetRoomDataRsp {
	if r != nil && r.SetRoomDataRsp != nil {
		return *r.SetRoomDataRsp
	}
	return SetRoomDataRsp{}
}

func (r *RPCResponse) GetUpdateRoomRoleDataRsp() UpdateRoomRoleDataRsp {
	if r != nil && r.UpdateRoomRoleDataRsp != nil {
		return *r.UpdateRoomRoleDataRsp
	}
	return UpdateRoomRoleDataRsp{}
}
