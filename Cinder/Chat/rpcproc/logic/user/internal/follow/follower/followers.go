package follower

import (
	"Cinder/Chat/rpcproc/logic/rpc"
	"Cinder/Chat/rpcproc/logic/types"
	"Cinder/Chat/rpcproc/logic/user/internal/follow/follower/dbutil"
	"fmt"
	"sync"
	"time"

	log "github.com/cihub/seelog"
)

type UserID = types.UserID

type IUserInfo interface {
	GetActiveData() []byte
	GetSrvID() string
}

// 粉丝列表
type Followers struct {
	mtx sync.Mutex

	userID   UserID // 主人ID，主人有哪些粉丝
	userInfo IUserInfo

	// 因为粉丝较多，有可能内存中放不下, 加载时有个数限制
	followers map[UserID]time.Time // 记录加粉丝时间
}

func NewFollowers(userID UserID, userInfo IUserInfo) *Followers {
	return &Followers{
		userID:    userID,
		userInfo:  userInfo,
		followers: make(map[UserID]time.Time),
	}
}

// Add 添加粉丝。
// 不在线就没有 Followers 实例，也不会 Add(). 写DB是调用者处理的，不在线也要写DB.
func (f *Followers) Add(userID UserID) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	if userID == f.userID {
		return // 禁止添加自己
	}
	f.followers[userID] = time.Now()
	isAdd := true
	f.rpcAddDelFollower(userID, isAdd)
}

// rpcAddDelFollower 通知我有粉丝加减。
func (f *Followers) rpcAddDelFollower(userID UserID, isAdd bool) {
	activeData := f.userInfo.GetActiveData()
	srvID := f.userInfo.GetSrvID()
	targetID := string(f.userID)
	// RPC 回调时不能加锁，后台执行
	go func() {
		ret := rpc.Rpc(srvID, "RPC_AddDelFollower", targetID, activeData, string(userID), isAdd)
		if ret.Err != nil {
			log.Errorf("RPC_AddDelFollower error: %v", ret.Err)
		}
	}()
}

// Delete 删除粉丝。
// 不在线就没有 Followers 实例，也不会 Delete(). 写DB是调用者处理的，不在线也要写DB.
func (f *Followers) Delete(userID UserID) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	delete(f.followers, userID)
	isAdd := false
	f.rpcAddDelFollower(userID, isAdd)
}

// Load 从DB加载粉丝
func (f *Followers) Load() error {
	followers, err := dbutil.UsersFollowersUtil(f.userID).Load()
	if err != nil {
		return fmt.Errorf("db load user followers: %w", err)
	}

	f.mtx.Lock()
	defer f.mtx.Unlock()

	f.followers = make(map[UserID]time.Time)
	for _, doc := range followers {
		f.followers[doc.Follower] = doc.OID.Timestamp()
	}
	delete(f.followers, f.userID) // 排除自己
	return nil
}

// GetIDs 获取ID列表
func (f *Followers) GetIDs() []UserID {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	result := make([]UserID, 0, len(f.followers))
	for id, _ := range f.followers {
		result = append(result, id)
	}
	return result
}

// GetFollowers 获取ID-Time列表
func (f *Followers) GetFollowers() map[UserID]time.Time {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	result := make(map[UserID]time.Time, len(f.followers))
	for id, t := range f.followers {
		result[id] = t
	}
	return result
}
