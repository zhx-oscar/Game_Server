package room

import (
	"Cinder/Matcher/matchapi/mtypes"
	"sync"

	assert "github.com/arl/assertgo"
)

// goroutine safe rooms
type _RoomsSafe struct {
	mtx sync.Mutex

	allRooms _IDToRoom // 所有房间，不包括删除中的房间

	modedRooms   _ModedRooms // 分模式所有房间
	roomsNotFull _ModedRooms // 允许加人的房间
}

type _IDToRoom map[mtypes.RoomID]*_Room
type _ModedRooms map[mtypes.MatchMode]_IDToRoom

func newRoomsSafe() *_RoomsSafe {
	return &_RoomsSafe{
		allRooms:     _IDToRoom{},
		modedRooms:   _ModedRooms{},
		roomsNotFull: _ModedRooms{},
	}
}

func (r *_RoomsSafe) GetRoomIDsNotFull(matchMode mtypes.MatchMode) []mtypes.RoomID {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	rooms, ok := r.roomsNotFull[matchMode]
	if !ok {
		return nil
	}

	result := make([]mtypes.RoomID, 0, len(rooms))
	for id, _ := range rooms {
		result = append(result, id)
	}
	return result
}

func (r *_RoomsSafe) GetRoom(roomID mtypes.RoomID) *_Room {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if result, ok := r.allRooms[roomID]; ok {
		return result
	}
	return nil
}

// AdjustOnRoomChanged 调整房间未满队列.
// room 有任何变化，都需要调用调整。
func (r *_RoomsSafe) AdjustOnRoomChanged(roomID mtypes.RoomID) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	room, ok := r.allRooms[roomID]
	if !ok {
		return // 已删除
	}

	// 检查是否可删除
	if room.IsDeleting() {
		r.deleteRoom(room)
		return
	}

	// 检查是否满人，调整 roomsNotFull
	r.checkRoomFull(room)
}

// InsertRoom 插入房间
func (r *_RoomsSafe) InsertRoom(room *_Room) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	assert.True(room != nil)
	r.allRooms[room.ID()] = room
	insertToModedRooms(r.modedRooms, room)
	r.checkRoomFull(room)
}

// DeleteRoom 删除房间
func (r *_RoomsSafe) DeleteRoom(roomID mtypes.RoomID) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	room, ok := r.allRooms[roomID]
	if !ok {
		return // 已删除
	}

	r.deleteRoom(room)
}

// checkRoomFull 检查房间是否已满，调整未满队列
func (r *_RoomsSafe) checkRoomFull(room *_Room) {
	if room.IsFull() {
		delete(r.roomsNotFull[room.MatchMode()], room.ID())
		return // 已满
	}

	// 加入未满队列
	insertToModedRooms(r.roomsNotFull, room)
}

func (r *_RoomsSafe) deleteRoom(room *_Room) {
	assert.True(room != nil)
	roomID := room.ID()
	delete(r.allRooms, roomID)
	matchMode := room.MatchMode()
	delete(r.modedRooms[matchMode], roomID)
	delete(r.roomsNotFull[matchMode], roomID)

	// room.OnDestroy() 延时执行，因为房间信息可能还需要广播
	deletingRoomList.Push(room)
}

func insertToModedRooms(rooms _ModedRooms, room *_Room) {
	assert.True(room != nil)
	matchMode := room.MatchMode()
	if rooms[matchMode] == nil {
		rooms[matchMode] = _IDToRoom{}
	}
	rooms[matchMode][room.ID()] = room
}

func (r *_RoomsSafe) GetRoomIDs(matchMode mtypes.MatchMode, limit int) []mtypes.RoomID {
	if limit <= 0 {
		return nil
	}

	r.mtx.Lock()
	defer r.mtx.Unlock()

	rooms, ok := r.modedRooms[matchMode]
	if !ok {
		return nil
	}

	l := len(rooms)
	if l > limit {
		l = limit
	}
	result := make([]mtypes.RoomID, 0, l)
	for id, room := range rooms {
		assert.True(room != nil)
		if !room.IsHidden() {
			result = append(result, id)
		}
		if len(result) >= limit {
			return result
		}
	}
	return result
}

// GetRoomCount 获取房间数
func (r *_RoomsSafe) GetRoomCount() int {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	return len(r.allRooms)
}

// GetRoomModeCount 获取房间模式数
func (r *_RoomsSafe) GetRoomModeCount() int {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	return len(r.modedRooms)
}

// ListRooms 批量列举房间，最大limit
func (r *_RoomsSafe) ListRooms(matchModes []mtypes.MatchMode, limit int) map[mtypes.MatchMode][]mtypes.RoomInfo {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	ret := map[mtypes.MatchMode][]mtypes.RoomInfo{}
	for _, mode := range matchModes {
		roomInfos := r.listRoomsOfMode(mode, limit)
		ret[mode] = roomInfos
		limit -= len(roomInfos)
		if limit <= 0 {
			return ret
		}
	}
	return ret
}

// listRoomsOfMode 列举某一模式的房间，最大limit
func (r *_RoomsSafe) listRoomsOfMode(matchMode mtypes.MatchMode, limit int) []mtypes.RoomInfo {
	rooms, ok := r.modedRooms[matchMode]
	if !ok {
		return nil
	}

	l := len(rooms)
	if l > limit {
		l = limit
	}
	result := make([]mtypes.RoomInfo, 0, l)
	for _, room := range rooms {
		assert.True(room != nil)
		if !room.IsHidden() {
			result = append(result, room.Info())
		}
		if len(result) >= limit {
			return result
		}
	}
	return result
}
