package room

import (
	"container/list"
	"sync"
	"time"
)

// 房间延时删除。如果房间立即删除，加入房间满人立即删除就会返回空的房间信息, OnDestroy()触发太早。
type _DeletingRoomList struct {
	mtx sync.Mutex

	roomList *list.List // of *_TimedRoom
}

type _TimedRoom struct {
	DeleteTime time.Time
	Room       *_Room
}

var deletingRoomList = newDeletingRoomList()

func init() {
	go runDeleteRooms()
}

func newDeletingRoomList() *_DeletingRoomList {
	return &_DeletingRoomList{
		roomList: list.New(),
	}
}

func newTimedRoom(room *_Room) *_TimedRoom {
	return &_TimedRoom{
		DeleteTime: time.Now().Add(time.Second),
		Room:       room,
	}
}

func runDeleteRooms() {
	for {
		time.Sleep(5 * time.Second)
		deletingRoomList.DeleteRooms()
	}
}

func (d *_DeletingRoomList) DeleteRooms() {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	now := time.Now()
	for e := d.roomList.Front(); e != nil; e = e.Next() {
		timedRoom := e.Value.(*_TimedRoom)
		if now.After(timedRoom.DeleteTime) {
			return // 还不到删除时间
		}
		timedRoom.Room.OnDestroy()
		d.roomList.Remove(e)
	}
}

func (d *_DeletingRoomList) Push(room *_Room) {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	d.roomList.PushBack(newTimedRoom(room))
}
