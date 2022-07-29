package room

import (
	"Cinder/Matcher/matchapi/mtypes"
	"Cinder/Matcher/matcherlib/ltypes"
	"sync"

	assert "github.com/arl/assertgo"
	log "github.com/cihub/seelog"
)

// 匹配房间
type _Room struct {
	mtx sync.Mutex

	info *mtypes.RoomInfo
}

var (
	// 仅在 _Room 内部使用
	_roomEvtHdlr ltypes.IRoomEventHandler
)

func newRoom(createRoomInfo mtypes.CreateRoomInfo) *_Room {
	assert.True(_roomEvtHdlr != nil)
	return &_Room{
		info: mtypes.NewRoomInfo(createRoomInfo),
	}
}

func SetRoomEventHandler(hdlr ltypes.IRoomEventHandler) {
	assert.True(hdlr != nil)
	_roomEvtHdlr = hdlr
}

func (r *_Room) ID() mtypes.RoomID {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	return r.info.RoomID
}

func (r *_Room) Info() mtypes.RoomInfo {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	return r.info.Copy()
}

func (r *_Room) IsFull() bool {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	return r.info.IsFull
}

func (r *_Room) IsHidden() bool {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	return r.info.IsHidden
}

func (r *_Room) MatchMode() mtypes.MatchMode {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	return r.info.MatchMode
}

// 没用到
// AddRole 添加角色。不符合条件就返回false.
// func (r *_Room) AddRole(roleInfo mtypes.RoleInfo) bool {
// 	r.mtx.Lock()
// 	defer r.mtx.Unlock()
// 	return r.addRoles(mtypes.RoleMap{roleInfo.RoleID: &roleInfo})
// }

// AddRoles 添加一组角色。不符合条件就返回false.
func (r *_Room) AddRoles(roles mtypes.RoleMap) bool {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	return r.addRoles(roles)
}

// AddRoles 添加一组角色。不符合条件就返回false.
func (r *_Room) addRoles(roles mtypes.RoleMap) bool {
	if r.info.IsDeleting || r.info.IsFull {
		return false // 删除中或人已满，不可加
	}

	// 执行匹配算法: 如判断角色等级，积分
	if false == _roomEvtHdlr.OnAddingRoles(r.info, roles) {
		// 不可加人
		return false
	}

	joinedRoleIDs := make([]mtypes.RoleID, 0, len(roles))
	for roleID, role := range roles {
		// 必须复制 RoleInfo, 使房间内数据不受队伍影响
		roleCopy := *role
		r.info.Roles[roleID] = &roleCopy
		joinedRoleIDs = append(joinedRoleIDs, roleID)
	}

	_roomEvtHdlr.OnAddedRoles(r.info, roles)
	r.newNotifier().NotifyRolesJoinRoom(mtypes.RolesJoinRoomMsg{
		RoomInfo:      *r.info,
		JoinedRoleIDs: joinedRoleIDs,
	})
	return true
}

func (r *_Room) DelRole(roleID mtypes.RoleID) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.delRole(roleID)
}

func (r *_Room) delRole(roleID mtypes.RoleID) {
	if _, ok := r.info.Roles[roleID]; !ok {
		return // 不存在
	}

	_roomEvtHdlr.OnDeletingRole(r.info, roleID)
	r.newNotifier().NotifyRoleLeaveRoom(mtypes.RoleLeaveRoomMsg{
		RoomInfo:      *r.info,
		LeavingRoleID: roleID,
	})
	delete(r.info.Roles, roleID)
}

func (r *_Room) HasRole(roleID mtypes.RoleID) bool {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	_, ok := r.info.Roles[roleID]
	return ok
}

func (r *_Room) OnCreate() error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	return _roomEvtHdlr.OnCreate(r.info)
}

func (r *_Room) OnDestroy() {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	// 先删除所有人
	for roleID := range r.info.Roles {
		r.delRole(roleID)
	}
	_roomEvtHdlr.OnDestroy(r.info)
}

func (r *_Room) IsDeleting() bool {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	return r.info.IsDeleting
}

func (r *_Room) Broadcast(msg interface{}) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.newNotifier().NotifyBroadcastRoom(mtypes.BroadcastRoomMsg{Msg: msg})
	return nil
}

func (r *_Room) newNotifier() *_RoomNotifier {
	return newRoomNotifier(r.info.GetSrvIDToRoleIDs())
}

func (r *_Room) SetRoomData(key string, data interface{}) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.info.SetData(key, data)
	// 没有 _roomEvtHdlr 事件触发
	r.newNotifier().NotifySetRoomData(mtypes.SetRoomDataMsg{
		Key:  key,
		Data: data,
	})
}

func (r *_Room) UpdateRoomRoleData(update mtypes.RoomRoleDataUpdate) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if update.RoomID != r.info.RoomID {
		// log.Error("room ID diffs")
		return
	}

	role, ok := r.info.Roles[update.RoleID]
	if !ok {
		// log.Debugf("no such role in room: %s", update.RoleID)
		return
	}

	log.Debugf("room update role data: %v", update)
	for key, f := range update.FloatData {
		role.SetFloatData(key, f)
	}
	for key, s := range update.StringData {
		role.SetStringData(key, s)
	}
	for _, tag := range update.AddTags {
		role.AddTag(tag)
	}
	for _, tag := range update.DelTags {
		role.DelTag(tag)
	}

	r.newNotifier().NotifyUpdateRoomRoleData(mtypes.UpdateRoomRoleDataMsg{
		Update: update,
	})
}
