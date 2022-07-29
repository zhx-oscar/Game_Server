package oflmsg

import (
	"Cinder/Chat/chatapi"
	"Cinder/Chat/rpcproc/logic/types"
	"sync"
)

type UserID = types.UserID

// P2P 私聊消息缓存
type _OfflineMessageCache struct {
	userToMsgs sync.Map // 接收者UserID -> UserOfflineMsgs
}

var MsgCache = &_OfflineMessageCache{}

// Add 添加一条缓存消息
func (o *_OfflineMessageCache) Add(from UserID, fromNick string, fromData []byte, to UserID, msgContent []byte) {
	msgs, _ := o.userToMsgs.LoadOrStore(to, newUserOfflineMsgs(to))
	msgs.(*_UserOfflineMsgs).Append(from, fromNick, fromData, msgContent)
}

// Pop 弹出玩家的所有缓存离线消息
// 玩家上线时，离线消息除了从DB加载，还需要添加缓存中的离线消息
func (o *_OfflineMessageCache) Pop(userID UserID) []*chatapi.ChatMessage {
	msgs, ok := o.userToMsgs.Load(userID)
	if !ok {
		return nil
	}

	o.userToMsgs.Delete(userID)
	return msgs.(*_UserOfflineMsgs).CopyMsgs()
}
