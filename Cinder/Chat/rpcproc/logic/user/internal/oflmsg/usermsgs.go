package oflmsg

import (
	"Cinder/Chat/chatapi"
	"Cinder/Chat/rpcproc/logic/chatmsg"
	"Cinder/Chat/rpcproc/logic/user/internal/oflmsg/dbutil"
	"sync"
	"time"

	log "github.com/cihub/seelog"
)

// 单个User的离线消息
type _UserOfflineMsgs struct {
	mtx sync.Mutex

	userID UserID
	msgs   []*chatapi.ChatMessage

	// 控制离线消息保存, 每100条 或者 5min 就保存一次
	nextSaveTime time.Time
}

func newUserOfflineMsgs(userID UserID) *_UserOfflineMsgs {
	res := &_UserOfflineMsgs{
		userID: userID,
	}
	res.reset()
	return res
}

func (u *_UserOfflineMsgs) Append(from UserID, fromNick string, fromData []byte, msgContent []byte) {
	u.mtx.Lock()
	defer u.mtx.Unlock()
	u.msgs = append(u.msgs, chatmsg.NewChatMessage(from, fromNick, fromData, msgContent))
	// log.Debugf("user '%v' cached offline message: %d", u.userID, len(u.msgs))

	// 触发DB保存, 保存后清空
	u.tryToSave()
}

func (u *_UserOfflineMsgs) CopyMsgs() []*chatapi.ChatMessage {
	u.mtx.Lock()
	defer u.mtx.Unlock()

	// log.Debugf("user '%v' poped offline messages: %d", u.userID, len(u.msgs))
	result := make([]*chatapi.ChatMessage, len(u.msgs))
	copy(result, u.msgs)
	return result
}

// reset 重置，用于初始化和保存后.
func (u *_UserOfflineMsgs) reset() {
	if len(u.msgs) != 0 {
		log.Debugf("user '%v' reset offline messages: %d -> 0", u.userID, len(u.msgs))
	}
	u.msgs = []*chatapi.ChatMessage{}
	u.nextSaveTime = time.Now().Add(5 * time.Minute)
}

func (u *_UserOfflineMsgs) tryToSave() {
	if !u.canSave() {
		return
	}
	if err := u.save(); err != nil {
		log.Errorf("failed to save user offline messages: %v", err)
		return
	}
	u.reset()
}

func (u *_UserOfflineMsgs) canSave() bool {
	return len(u.msgs) > 100 || time.Now().After(u.nextSaveTime)
}

func (u *_UserOfflineMsgs) save() error {
	// log.Debugf("user '%v' is saving offline messages to DB: %d", u.userID, len(u.msgs))
	return dbutil.UserOflnMsgUtil(u.userID).Insert(u.msgs)
}
