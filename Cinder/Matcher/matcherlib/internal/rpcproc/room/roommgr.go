package room

import (
	"Cinder/Matcher/matchapi/mtypes"
	"errors"
)

type _RoomMgr struct {
	rooms *_RoomsSafe // 协程安全
}

var mgr = &_RoomMgr{
	rooms: newRoomsSafe(),
}

func GetMgr() *_RoomMgr {
	return mgr
}

// JoinRandomRoom 加入随机房间。
// 如果没有空房，则创建一个.
func (r *_RoomMgr) RoleJoinRandomRoom(roleInfo mtypes.RoleInfo, createInfo mtypes.CreateRoomInfo) (mtypes.RoomInfo, error) {
	roleInfo.TeamID = "" // 强制清空TeamID
	return r.rolesJoinRandomRoom(mtypes.RoleMap{roleInfo.RoleID: &roleInfo}, createInfo)
}

// RoleCreateRoom 创建房间
func (r *_RoomMgr) RoleCreateRoom(creatorInfo mtypes.RoleInfo, createRoomInfo mtypes.CreateRoomInfo) (mtypes.RoomInfo, error) {
	return r.rolesCreateRoom(mtypes.RoleMap{creatorInfo.RoleID: &creatorInfo}, createRoomInfo)
}

func (r *_RoomMgr) RoleJoinRoom(roleInfo mtypes.RoleInfo, roomID mtypes.RoomID) (mtypes.RoomInfo, error) {
	return r.rolesJoinRoom(mtypes.RoleMap{roleInfo.RoleID: &roleInfo}, roomID)
}

// RoleLeaveRoom 离开房间
func (r *_RoomMgr) RoleLeaveRoom(roleID mtypes.RoleID, roomID mtypes.RoomID) (mtypes.RoomInfo, error) {
	room := r.rooms.GetRoom(roomID)
	if room == nil {
		return mtypes.RoomInfo{}, errors.New("room does not exist")
	}
	if !room.HasRole(roleID) {
		return mtypes.RoomInfo{}, errors.New("role does not exist") // 不存在该 role
	}
	room.DelRole(roleID)                   // 触发事件并房间广播
	r.rooms.AdjustOnRoomChanged(room.ID()) // 可能需要移入未满队列
	return room.Info(), nil
}

// TeamJoinRandomRoom 组队加入随机房间
func (r *_RoomMgr) TeamJoinRandomRoom(teamID mtypes.TeamID, createInfo mtypes.CreateRoomInfo) (mtypes.RoomInfo, error) {
	teamRoles, errGet := getTeamRoles(teamID)
	if errGet != nil {
		return mtypes.RoomInfo{}, errGet
	}
	return r.rolesJoinRandomRoom(teamRoles, createInfo)
}

func (r *_RoomMgr) rolesJoinRandomRoom(roles mtypes.RoleMap, createInfo mtypes.CreateRoomInfo) (mtypes.RoomInfo, error) {
	roomIDs := r.rooms.GetRoomIDsNotFull(createInfo.MatchMode)
	for _, id := range roomIDs {
		room := r.rooms.GetRoom(id)
		if room == nil {
			continue
		}

		if !room.AddRoles(roles) {
			continue
		}
		r.rooms.AdjustOnRoomChanged(id) // 可能需要移出未满队列
		return room.Info(), nil
	}

	// 需要新建房间
	return r.rolesCreateRoom(roles, createInfo)
}

// TeamCreateRoom 组队创建房间
func (r *_RoomMgr) TeamCreateRoom(teamID mtypes.TeamID, createInfo mtypes.CreateRoomInfo) (mtypes.RoomInfo, error) {
	teamRoles, errGet := getTeamRoles(teamID)
	if errGet != nil {
		return mtypes.RoomInfo{}, errGet
	}
	return r.rolesCreateRoom(teamRoles, createInfo)
}

// TeamJoinRoom 组队加入房间
func (r *_RoomMgr) TeamJoinRoom(teamID mtypes.TeamID, roomID mtypes.RoomID) (mtypes.RoomInfo, error) {
	teamRoles, errGet := getTeamRoles(teamID)
	if errGet != nil {
		return mtypes.RoomInfo{}, errGet
	}
	return r.rolesJoinRoom(teamRoles, roomID)
}

func (r *_RoomMgr) rolesJoinRoom(roles mtypes.RoleMap, roomID mtypes.RoomID) (mtypes.RoomInfo, error) {
	room := r.rooms.GetRoom(roomID)
	if room == nil {
		return mtypes.RoomInfo{}, errors.New("room does not exist")
	}
	if !room.AddRoles(roles) {
		return mtypes.RoomInfo{}, errors.New("can not add role to room")
	}
	r.rooms.AdjustOnRoomChanged(room.ID()) // 可能需要移入未满队列
	return room.Info(), nil
}

func (r *_RoomMgr) rolesCreateRoom(roles mtypes.RoleMap, createInfo mtypes.CreateRoomInfo) (mtypes.RoomInfo, error) {
	room := newRoom(createInfo)
	if errCreate := room.OnCreate(); errCreate != nil { // 触发Space创建房间
		return mtypes.RoomInfo{}, errCreate
	}

	if !room.AddRoles(roles) {
		room.OnDestroy() // 触发通知销毁
		return mtypes.RoomInfo{}, errors.New("can not add role to new room")
	}

	r.rooms.InsertRoom(room)
	return room.Info(), nil
}

// BroadcastRoom 广播房间
func (r *_RoomMgr) BroadcastRoom(roomID mtypes.RoomID, msg interface{}) error {
	room := r.rooms.GetRoom(roomID)
	if room == nil {
		return nil
	}
	return room.Broadcast(msg)
}

// DeleteRoom 删除房间
func (r *_RoomMgr) DeleteRoom(roomID mtypes.RoomID) {
	r.rooms.DeleteRoom(roomID)
}

// SetRoomData 设置房间数据
func (r *_RoomMgr) SetRoomData(roomID mtypes.RoomID, key string, data interface{}) {
	room := r.rooms.GetRoom(roomID)
	if room == nil {
		return
	}
	room.SetRoomData(key, data)
	// 因为 SetRoomData 不会触发事件回调，所以不用 AdjustOnRoomChanged()
}

// UpdateRoomRoleData 更新房间角色数据
func (r *_RoomMgr) UpdateRoomRoleData(update mtypes.RoomRoleDataUpdate) {
	room := r.rooms.GetRoom(update.RoomID)
	if room != nil {
		room.UpdateRoomRoleData(update)
	}
}

// GetRoomInfos 获取房间信息
func (r *_RoomMgr) GetRoomInfos(roomIDs []mtypes.RoomID) map[mtypes.RoomID]mtypes.RoomInfo {
	ret := make(map[mtypes.RoomID]mtypes.RoomInfo, len(roomIDs))
	for _, roomID := range roomIDs {
		room := r.rooms.GetRoom(roomID)
		if room == nil {
			continue
		}
		ret[roomID] = room.Info()
	}
	return ret
}

// GetRoomList 获取房间列表
func (r *_RoomMgr) GetRoomIDs(matchMode mtypes.MatchMode) []mtypes.RoomID {
	return r.rooms.GetRoomIDs(matchMode, 1000)
}

// ListRooms 批量列举房间，最大100个
func (r *_RoomMgr) ListRooms(matchModes []mtypes.MatchMode) map[mtypes.MatchMode][]mtypes.RoomInfo {
	return r.rooms.ListRooms(matchModes, 100)
}

// GetRoomCount 获取房间数
func (r *_RoomMgr) GetRoomCount() int {
	return r.rooms.GetRoomCount()
}

// GetRoomModeCount 获取房间模式数
func (r *_RoomMgr) GetRoomModeCount() int {
	return r.rooms.GetRoomModeCount()
}
