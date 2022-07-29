package userinfo

import (
	"Cinder/Chat/chatapi"
	"Cinder/Chat/rpcproc/logic/types"
	db "Cinder/Chat/rpcproc/logic/user/internal/dbutil"
	"Cinder/Chat/rpcproc/logic/user/internal/userinfo/dbutil"
	"bytes"
	"sync"
)

type UserID = types.UserID

// UserInfo 是 FriendInfo 加上 srvID, passiveData。
// 未来可能加上其他数据，一起在 Mutex 保护之下。因为 User 本身不加锁，所以各种数据必须在子对象锁之下。
// activeData 在 FriendInfo 中，原来没有 passiveData 时， activeData 也称为 customData, data, targetData, fromData。
type UserInfo struct {
	mtx sync.Mutex

	userID      UserID
	srvID       string // 注册服ID, 用作消息推送的目标服
	friendInfo  chatapi.FriendInfo
	passiveData []byte // 被动数据
}

func NewUserInfo(userID UserID, srvID string) *UserInfo {
	// 其他数据须从DB加载
	return &UserInfo{
		userID: userID,
		srvID:  srvID,
		friendInfo: chatapi.FriendInfo{
			ID:       string(userID),
			IsOnline: true,
		},
	}
}

func (l *UserInfo) GetNickAndData() (nick string, data []byte) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	return l.friendInfo.Nick, l.friendInfo.Data
}

// GetActiveData 获取主动数据
// 注册服发送的主动数据，在消息推送时原样用作参数
func (l *UserInfo) GetActiveData() []byte {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	return l.friendInfo.Data
}

func (l *UserInfo) GetSrvID() string {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	return l.srvID
}

func (l *UserInfo) SetSrvID(srvID string) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	l.srvID = srvID // 仅内存数据，无DB保存
}

func (l *UserInfo) SetNick(nick string) (changed bool, e error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if l.friendInfo.Nick == nick {
		return false, nil // 未改变
	}
	// 有变化则更新并写DB
	l.friendInfo.Nick = nick
	return true, dbutil.UsersUtil(l.userID).UpdateNick(nick)
}

func (l *UserInfo) GetNick() string {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	return l.friendInfo.Nick
}

func (l *UserInfo) SetActiveData(data []byte) (changed bool, e error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if bytes.Equal(l.friendInfo.Data, data) {
		return false, nil
	}
	l.friendInfo.Data = data
	return true, dbutil.UsersUtil(l.userID).UpdateActiveData(data)
}

func (l *UserInfo) SetPassiveData(data []byte) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if bytes.Equal(l.passiveData, data) {
		return nil
	}
	l.passiveData = data
	return dbutil.UsersUtil(l.userID).UpdatePassiveData(data)
}

func (l *UserInfo) GetPassiveData() []byte {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	return l.passiveData
}

func (l *UserInfo) InitWithUserDoc(doc *db.UserDoc) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	l.friendInfo.Nick = doc.Nick
	l.friendInfo.Data = doc.ActiveData
	l.friendInfo.OfflineTime = doc.OfflineTime
	l.friendInfo.FollowerNumber = doc.FollowerNumber
	l.passiveData = doc.PassiveData
}

// GetFriendInfo 获取用户信息
func (l *UserInfo) GetFriendInfo() chatapi.FriendInfo {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	return l.friendInfo
}
