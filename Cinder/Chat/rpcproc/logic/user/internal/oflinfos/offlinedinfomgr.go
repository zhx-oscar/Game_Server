package oflinfos

import (
	"Cinder/Chat/chatapi"
	"Cinder/Chat/rpcproc/logic/types"
	"Cinder/Chat/rpcproc/logic/user/internal/dbutil"
	"fmt"
	"time"

	"github.com/hashicorp/golang-lru"
)

// _OfflinedInfoMgr 管理离线用户的 FriendInfo。
type _OfflinedInfoMgr struct {
	// least-recent-used列表保存只读 *chatapi.FriendInfo
	cache *lru.TwoQueueCache
}

var offlinedInfoMgr *_OfflinedInfoMgr

func init() {
	cache, err := lru.New2Q(1024 * 1024)
	if err != nil {
		panic(fmt.Errorf("lru new cache failed: %w", err))
	}
	offlinedInfoMgr = &_OfflinedInfoMgr{
		cache: cache,
	}
}

func GetOfldInfoMgr() *_OfflinedInfoMgr {
	return offlinedInfoMgr
}

func (i *_OfflinedInfoMgr) Add(info chatapi.FriendInfo) {
	info.IsOnline = false
	info.OfflineTime = time.Now()
	i.cache.Add(info.ID, &info)
}

func (i *_OfflinedInfoMgr) Remove(userID types.UserID) {
	i.cache.Remove(userID)
}

// _Get 获取一个ID对应的信息。
// 仅供包内部调用。
func (i *_OfflinedInfoMgr) _Get(userID types.UserID) (*chatapi.FriendInfo, error) {
	// 从Cache获取一个ID对应的信息。
	if cached, ok := i.cache.Get(userID); ok && cached != nil {
		return cached.(*chatapi.FriendInfo), nil
	}

	// 从DB加载
	// log.Debugf("user info is missing in cache, load from DB: %v", userID)
	doc, err := dbutil.UsersUtil(userID).Load()
	if err != nil {
		return nil, fmt.Errorf("load user info: %w", err)
	}

	// 加载后加入cache
	info := userDocToFriendInfo(*doc)
	i.cache.Add(userID, info)
	return info, nil
}

func userDocToFriendInfo(doc dbutil.UserDoc) *chatapi.FriendInfo {
	return &chatapi.FriendInfo{
		ID:   string(doc.UserID),
		Nick: doc.Nick,
		Data: doc.ActiveData,

		IsOnline:       false,
		OfflineTime:    doc.OfflineTime,
		FollowerNumber: doc.FollowerNumber,
	}
}
