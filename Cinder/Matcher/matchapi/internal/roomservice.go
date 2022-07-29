package internal

import (
	"Cinder/Matcher/matchapi/mtypes"
	"Cinder/Matcher/rpcmsg"
)

type RoomService struct {
	rpcCaller *_RpcCaller
}

func NewRoomService(matcherAreaID string, matcherServerID string) *RoomService {
	return &RoomService{
		rpcCaller: newRpcCaller(matcherAreaID, matcherServerID),
	}
}

// JoinRandomRoom 加入随机房间, 触发房间广播 RPC_MatchNotifyJoinRoom
func (r *RoomService) RoleJoinRandomRoom(roleInfo mtypes.RoleInfo, matchMode mtypes.MatchMode) (mtypes.RoomInfo, error) {
	return r.RoleJoinRandomRoomWithInfo(roleInfo, mtypes.CreateRoomInfo{
		MatchMode: matchMode,
	})
}

// 加入随机房间, 带初始信息，触发房间广播 JoinRoomMsg
func (r *RoomService) RoleJoinRandomRoomWithInfo(roleInfo mtypes.RoleInfo, createInfo mtypes.CreateRoomInfo) (mtypes.RoomInfo, error) {
	roleInfo.SrvID = getServiceID()
	ret, errJoin := r.rpcCaller.rpc2("RPC_RoleJoinRandomRoom", rpcmsg.RoleJoinRandomRoomReq{
		RoleInfo:   roleInfo,
		CreateInfo: createInfo,
	})
	return ret.GetRoleJoinRandomRoomRsp().RoomInfo, errJoin
}

// CreateRoom 创建房间
func (r *RoomService) RoleCreateRoom(creatorInfo mtypes.RoleInfo, matchMode mtypes.MatchMode) (mtypes.RoomInfo, error) {
	return r.RoleCreateRoomWithInfo(creatorInfo, mtypes.CreateRoomInfo{
		MatchMode: matchMode,
	})
}

// 创建房间，带初始化信息，触发房间广播 JoinRoomMsg
func (r *RoomService) RoleCreateRoomWithInfo(creatorInfo mtypes.RoleInfo, createInfo mtypes.CreateRoomInfo) (mtypes.RoomInfo, error) {
	creatorInfo.SrvID = getServiceID()
	ret, err := r.rpcCaller.rpc2("RPC_RoleCreateRoom", rpcmsg.RoleCreateRoomReq{
		CreatorInfo: creatorInfo,
		CreateInfo:  createInfo,
	})
	return ret.GetRoleCreateRoomRsp().RoomInfo, err
}

// JoinRoom 加入指定房间, 触发房间广播 RPC_MatchNotifyJoinRoom
func (r *RoomService) RoleJoinRoom(roleInfo mtypes.RoleInfo, roomID mtypes.RoomID) (mtypes.RoomInfo, error) {
	roleInfo.SrvID = getServiceID()
	ret, err := r.rpcCaller.rpc2("RPC_RoleJoinRoom", rpcmsg.RoleJoinRoomReq{
		RoleInfo: roleInfo,
		RoomID:   roomID,
	})
	return ret.GetRoleJoinRoomRsp().RoomInfo, err
}

// LeaveRoom 自己离开房间, 或踢人, 触发房间广播 RPC_MatchNotifyLeaveRoom
// 如果有组队，则同时离开队伍
func (r *RoomService) RoleLeaveRoom(roleID mtypes.RoleID, roomID mtypes.RoomID) (mtypes.RoomInfo, error) {
	ret, err := r.rpcCaller.rpc2("RPC_RoleLeaveRoom", rpcmsg.RoleLeaveRoomReq{
		RoleID: roleID,
		RoomID: roomID,
	})
	if err != nil {
		return mtypes.RoomInfo{}, err
	}
	return ret.GetRoleLeaveRoomRsp().RoomInfo, nil
}

// 加入随机房间, 触发房间广播 RPC_MatchNotifyJoinRoom
func (r *RoomService) TeamJoinRandomRoom(teamID mtypes.TeamID, matchMode mtypes.MatchMode) (mtypes.RoomInfo, error) {
	return r.TeamJoinRandomRoomWithInfo(teamID, mtypes.CreateRoomInfo{
		MatchMode: matchMode,
	})
}

// 加入随机房间, 触发房间广播 RPC_MatchNotifyJoinRoom, 如果需新建房间则设置房间数据
func (r *RoomService) TeamJoinRandomRoomWithInfo(teamID mtypes.TeamID, createInfo mtypes.CreateRoomInfo) (mtypes.RoomInfo, error) {
	ret, err := r.rpcCaller.rpc2("RPC_TeamJoinRandomRoom", rpcmsg.TeamJoinRandomRoomReq{
		TeamID:     teamID,
		CreateInfo: createInfo,
	})
	return ret.GetTeamJoinRandomRoomRsp().RoomInfo, err
}

// 创建房间，队长成为 Leader
func (r *RoomService) TeamCreateRoom(creatorTeamID mtypes.TeamID, matchMode mtypes.MatchMode) (mtypes.RoomInfo, error) {
	return r.TeamCreateRoomWithInfo(creatorTeamID, mtypes.CreateRoomInfo{
		MatchMode: matchMode,
	})
}

// 创建房间, 带初始化信息
func (r *RoomService) TeamCreateRoomWithInfo(creatorTeamID mtypes.TeamID, createInfo mtypes.CreateRoomInfo) (mtypes.RoomInfo, error) {
	ret, err := r.rpcCaller.rpc2("RPC_TeamCreateRoom", rpcmsg.TeamCreateRoomReq{
		CreatorTeamID: creatorTeamID,
		CreateInfo:    createInfo,
	})
	return ret.GetTeamCreateRoomRsp().RoomInfo, err
}

// 加入指定房间, 触发房间广播 RPC_MatchNotifyJoinRoom
func (r *RoomService) TeamJoinRoom(teamID mtypes.TeamID, roomID mtypes.RoomID) (mtypes.RoomInfo, error) {
	ret, err := r.rpcCaller.rpc2("RPC_TeamJoinRoom", rpcmsg.TeamJoinRoomReq{
		TeamID: teamID,
		RoomID: roomID,
	})
	return ret.GetTeamJoinRoomRsp().RoomInfo, err
}

// BroadcastRoom 房间广播。触发房间广播 RPC_MatchNotifyBroadcastRoom
func (r *RoomService) BroadcastRoom(roomID mtypes.RoomID, msg interface{}) error {
	_, err := r.rpcCaller.rpc2("RPC_BroadcastRoom", rpcmsg.BroadcastRoomReq{
		RoomID: roomID,
		Msg:    msg,
	})
	return err
}

// GetRoomList 获取房间列表。批量获取多个模式请用 ListRooms()。
func (r *RoomService) GetRoomList(matchMode mtypes.MatchMode) (roomIDs []mtypes.RoomID, err error) {
	ret, err := r.rpcCaller.rpc2("RPC_GetRoomList", rpcmsg.GetRoomListReq{
		MatchMode: matchMode,
	})
	return ret.GetGetRoomListRsp().RoomIDs, err
}

// 获取房间详情。批量获取多个房间请用 GetRoomInfos()
func (r *RoomService) GetRoomInfo(roomID mtypes.RoomID) (mtypes.RoomInfo, error) {
	roomInfos, err := r.GetRoomInfos([]mtypes.RoomID{roomID})
	if err != nil {
		return mtypes.RoomInfo{}, err
	}
	if roomInfos != nil {
		return roomInfos[roomID], nil
	}
	return mtypes.RoomInfo{}, nil
}

// 获取房间列表，不可见房间不会列出
func (r *RoomService) ListRooms(matchModes []mtypes.MatchMode) (rooms map[mtypes.MatchMode][]mtypes.RoomInfo, err error) {
	ret, err := r.rpcCaller.rpc2("RPC_ListRooms", rpcmsg.ListRoomsReq{
		MatchModes: matchModes,
	})
	return ret.GetListRoomsRsp().Rooms, err
}

// 获取房间详情
func (r *RoomService) GetRoomInfos(roomIDs []mtypes.RoomID) (map[mtypes.RoomID]mtypes.RoomInfo, error) {
	ret, err := r.rpcCaller.rpc2("RPC_GetRoomInfos", rpcmsg.GetRoomInfosReq{
		RoomIDs: roomIDs,
	})
	return ret.GetGetRoomInfosRsp().RoomInfos, err
}

// SetRoomData 设置房间属性, 触发房间内广播 RPC_MatchNotifySetRoomData
func (r *RoomService) SetRoomData(roomID mtypes.RoomID, key string, data interface{}) error {
	_, err := r.rpcCaller.rpc2("RPC_SetRoomData", rpcmsg.SetRoomDataReq{
		RoomID: roomID,
		Key:    key,
		Data:   data,
	})
	return err
}

// 设置角色数据, 触发房间内广播 RPC_MatchNotifyUpdateRoomRoleData
func (r *RoomService) SetRoomRoleFloatData(roomID mtypes.RoomID, roleID mtypes.RoleID, key string, data float64) error {
	update := makeRoomRoleDataUpdate(roomID, roleID)
	update.FloatData = map[string]float64{key: data}
	return r.updateRoomRoleData(update)
}

// 设置角色数据, 触发房间内广播 RPC_MatchNotifyUpdateRoomRoleData
func (r *RoomService) SetRoomRoleStringData(roomID mtypes.RoomID, roleID mtypes.RoleID, key string, data string) error {
	update := makeRoomRoleDataUpdate(roomID, roleID)
	update.StringData = map[string]string{key: data}
	return r.updateRoomRoleData(update)
}

// 设置角色数据, 触发房间内广播 RPC_MatchNotifyUpdateRoomRoleData
func (r *RoomService) AddRoomRoleTag(roomID mtypes.RoomID, roleID mtypes.RoleID, tag string) error {
	update := makeRoomRoleDataUpdate(roomID, roleID)
	update.AddTags = []string{tag}
	return r.updateRoomRoleData(update)
}

// 设置角色数据, 触发房间内广播 RPC_MatchNotifyUpdateRoomRoleData
func (r *RoomService) DelRoomRoleTag(roomID mtypes.RoomID, roleID mtypes.RoleID, tag string) error {
	update := makeRoomRoleDataUpdate(roomID, roleID)
	update.DelTags = []string{tag}
	return r.updateRoomRoleData(update)
}

func makeRoomRoleDataUpdate(roomID mtypes.RoomID, roleID mtypes.RoleID) mtypes.RoomRoleDataUpdate {
	return mtypes.RoomRoleDataUpdate{
		RoomID: roomID,
		RoleID: roleID,
	}
}

// updateRoomRoleData 更新角色属性, 触发房间内广播 RPC_MatchNotifyUpdateRoomRoleData
func (r *RoomService) updateRoomRoleData(update mtypes.RoomRoleDataUpdate) error {
	_, err := r.rpcCaller.rpc2("RPC_UpdateRoomRoleData", rpcmsg.UpdateRoomRoleDataReq{
		Update: update,
	})
	return err
}

// DeleteRoom 删除房间
func (r *RoomService) DeleteRoom(roomID mtypes.RoomID) error {
	return r.rpcCaller.rpc("RPC_DeleteRoom", string(roomID))
}
