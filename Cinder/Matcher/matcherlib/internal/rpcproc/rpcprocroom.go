package rpcproc

import (
	"Cinder/Matcher/matchapi/mtypes"
	"Cinder/Matcher/matcherlib/internal/rpcproc/room"
	"Cinder/Matcher/rpcmsg"
	"encoding/json"

	log "github.com/cihub/seelog"
)

type _RPCProcRoom struct {
}

// RPC_RoleJoinRandomRoom 进入随机房间, 触发房间广播 RPC_MatchNotifyJoinRoom
func (r *_RPCProcRoom) RPC_RoleJoinRandomRoom(reqJson []byte) (rspJson []byte, errStr string) {
	log.Debugf("RPC_RoleJoinRandomRoom(sreqJson=%d bytes)", len(reqJson))
	req := rpcmsg.RoleJoinRandomRoomReq{}
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err.Error()
	}
	roomInfo, errJoin := room.GetMgr().RoleJoinRandomRoom(req.RoleInfo, req.CreateInfo)
	if errJoin != nil {
		return nil, errJoin.Error()
	}
	return rpcmsg.RPCResponse{
		RoleJoinRandomRoomRsp: &rpcmsg.RoleJoinRandomRoomRsp{
			RoomInfo: roomInfo,
		},
	}.Marshal()
}

// RPC_RoleCreateRoom 创建房间
func (r *_RPCProcRoom) RPC_RoleCreateRoom(reqJson []byte) (rspJson []byte, errStr string) {
	log.Debugf("RPC_RoleCreateRoom(reqJson=%d bytes)", len(reqJson))
	req := rpcmsg.RoleCreateRoomReq{}
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err.Error()
	}
	roomInfo, errCrt := room.GetMgr().RoleCreateRoom(req.CreatorInfo, req.CreateInfo)
	if errCrt != nil {
		return nil, errCrt.Error()
	}
	return rpcmsg.RPCResponse{
		RoleCreateRoomRsp: &rpcmsg.RoleCreateRoomRsp{
			RoomInfo: roomInfo,
		},
	}.Marshal()
}

// RPC_RoleJoinRoom 加入指定房间, 触发房间广播 RPC_MatchNotifyRolesJoinRoom
func (r *_RPCProcRoom) RPC_RoleJoinRoom(reqJson []byte) (rspJson []byte, errStr string) {
	log.Debugf("RPC_RoleJoinRoom(reqJson=%d bytes)", len(reqJson))
	req := rpcmsg.RoleJoinRoomReq{}
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err.Error()
	}
	roomInfo, errJoin := room.GetMgr().RoleJoinRoom(req.RoleInfo, req.RoomID)
	if errJoin != nil {
		return nil, errJoin.Error()
	}
	return rpcmsg.RPCResponse{
		RoleJoinRoomRsp: &rpcmsg.RoleJoinRoomRsp{
			RoomInfo: roomInfo,
		},
	}.Marshal()
}

// RPC_RoleLeaveRoom 自己离开房间, 或踢人, 触发房间广播 RPC_MatchNotifyLeaveRoom
func (r *_RPCProcRoom) RPC_RoleLeaveRoom(reqJson []byte) (rspJson []byte, errStr string) {
	log.Debugf("RPC_LeaveRoom(reqJson=%d bytes)", len(reqJson))
	req := rpcmsg.RoleLeaveRoomReq{}
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err.Error()
	}
	roomInfo, err := room.GetMgr().RoleLeaveRoom(req.RoleID, req.RoomID)
	if err != nil {
		return nil, err.Error()
	}
	return rpcmsg.RPCResponse{
		RoleLeaveRoomRsp: &rpcmsg.RoleLeaveRoomRsp{
			RoomInfo: roomInfo,
		},
	}.Marshal()
}

// RPC_TeamJoinRandomRoom 加入随机房间, 触发房间广播 RPC_MatchNotifyJoinRoom
func (r *_RPCProcRoom) RPC_TeamJoinRandomRoom(reqJson []byte) (rspJson []byte, errStr string) {
	log.Debugf("RPC_TeamJoinRandomRoom(reqJson=%d bytes)", len(reqJson))
	req := rpcmsg.TeamJoinRandomRoomReq{}
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err.Error()
	}
	roomInfo, errJoin := room.GetMgr().TeamJoinRandomRoom(req.TeamID, req.CreateInfo)
	if errJoin != nil {
		return nil, errJoin.Error()
	}
	return rpcmsg.RPCResponse{
		TeamJoinRandomRoomRsp: &rpcmsg.TeamJoinRandomRoomRsp{
			RoomInfo: roomInfo,
		},
	}.Marshal()
}

// RPC_TeamCreateRoom 创建房间，触发房间广播 RPC_MatchNotifyJoinRoom
func (r *_RPCProcRoom) RPC_TeamCreateRoom(reqJson []byte) (rspJson []byte, erStr string) {
	log.Debugf("RPC_TeamCreateRoom(reqJson=%d bytes)", len(reqJson))
	req := rpcmsg.TeamCreateRoomReq{}
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err.Error()
	}
	roomInfo, errCrt := room.GetMgr().TeamCreateRoom(req.CreatorTeamID, req.CreateInfo)
	if errCrt != nil {
		return nil, errCrt.Error()
	}
	return rpcmsg.RPCResponse{
		TeamCreateRoomRsp: &rpcmsg.TeamCreateRoomRsp{
			RoomInfo: roomInfo,
		},
	}.Marshal()
}

// RPC_TeamJoinRoom 加入指定房间, 触发房间广播 RPC_MatchNotifyRolesJoinRoom
func (r *_RPCProcRoom) RPC_TeamJoinRoom(reqJson []byte) (rspJson []byte, errStr string) {
	log.Debugf("RPC_TeamJoinRoom(reqJson=%d bytes)", len(reqJson))
	req := rpcmsg.TeamJoinRoomReq{}
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err.Error()
	}
	roomInfo, errJoin := room.GetMgr().TeamJoinRoom(req.TeamID, req.RoomID)
	if errJoin != nil {
		return nil, errJoin.Error()
	}
	return rpcmsg.RPCResponse{
		TeamJoinRoomRsp: &rpcmsg.TeamJoinRoomRsp{
			RoomInfo: roomInfo,
		},
	}.Marshal()
}

// RPC_BroadcastRoom 房间广播
func (r *_RPCProcRoom) RPC_BroadcastRoom(reqJson []byte) (rspJson []byte, errStr string) {
	log.Debugf("RPC_BroadcastRoom(reqJson=%d bytes)", len(reqJson))
	req := rpcmsg.BroadcastRoomReq{}
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err.Error()
	}
	if err := room.GetMgr().BroadcastRoom(req.RoomID, req.Msg); err != nil {
		return nil, err.Error()
	}
	return rpcmsg.RPCResponse{}.Marshal()
}

// RPC_GetRoomList 获取房间列表
func (r *_RPCProcRoom) RPC_GetRoomList(reqJson []byte) (rspJson []byte, errStr string) {
	log.Debugf("RPC_GetRoomList(reqJson=%d bytes)", len(reqJson))
	req := rpcmsg.GetRoomListReq{}
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err.Error()
	}
	roomIDs := room.GetMgr().GetRoomIDs(req.MatchMode)
	return rpcmsg.RPCResponse{
		GetRoomListRsp: &rpcmsg.GetRoomListRsp{
			RoomIDs: roomIDs,
		},
	}.Marshal()
}

// RPC_ListRooms 批量获取房间列表
func (r *_RPCProcRoom) RPC_ListRooms(reqJson []byte) (rspJson []byte, errStr string) {
	log.Debugf("RPC_ListRooms(reqJson=%d bytes)", len(reqJson))
	req := rpcmsg.ListRoomsReq{}
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err.Error()
	}
	rooms := room.GetMgr().ListRooms(req.MatchModes)
	return rpcmsg.RPCResponse{
		ListRoomsRsp: &rpcmsg.ListRoomsRsp{
			Rooms: rooms,
		},
	}.Marshal()
}

// RPC_GetRoomInfos 获取房间详情
func (r *_RPCProcRoom) RPC_GetRoomInfos(reqJson []byte) (rspJson []byte, errStr string) {
	log.Debugf("RPC_GetRoomInfos(reqJson=%d bytes)", len(reqJson))
	req := rpcmsg.GetRoomInfosReq{}
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err.Error()
	}
	roomInfos := room.GetMgr().GetRoomInfos(req.RoomIDs)
	return rpcmsg.RPCResponse{
		GetRoomInfosRsp: &rpcmsg.GetRoomInfosRsp{
			RoomInfos: roomInfos,
		},
	}.Marshal()
}

// RPC_SetRoomData 设置房间数据
func (r *_RPCProcRoom) RPC_SetRoomData(reqJson []byte) (rspJson []byte, errStr string) {
	log.Debugf("RPC_SetRoomData(reqJson=%d bytes)", len(reqJson))
	req := rpcmsg.SetRoomDataReq{}
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err.Error()
	}
	room.GetMgr().SetRoomData(req.RoomID, req.Key, req.Data)
	return rpcmsg.RPCResponse{
		SetRoomDataRsp: &rpcmsg.SetRoomDataRsp{},
	}.Marshal()
}

// RPC_UpdateRoomRoleData 更新房间内角色数据
func (r *_RPCProcRoom) RPC_UpdateRoomRoleData(reqJson []byte) (rspJson []byte, errStr string) {
	log.Debugf("RPC_UpdateRoomRoleData(reqJson=%d bytes)", len(reqJson))
	req := rpcmsg.UpdateRoomRoleDataReq{}
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err.Error()
	}
	room.GetMgr().UpdateRoomRoleData(req.Update)
	return rpcmsg.RPCResponse{
		UpdateRoomRoleDataRsp: &rpcmsg.UpdateRoomRoleDataRsp{},
	}.Marshal()
}

// RPC_DeleteRoom 删除房间
func (r *_RPCProcRoom) RPC_DeleteRoom(roomID string) (errStr string) {
	log.Debugf("RPC_DeleteRoom(roomID=%s)", roomID)
	room.GetMgr().DeleteRoom(mtypes.RoomID(roomID))
	return ""
}
