package usermgr

import (
	"Cinder/Chat/chatapi"
	"Cinder/Chat/rpcproc/logic/types"
	"Cinder/Chat/rpcproc/logic/user"
	"fmt"
	"sync"

	log "github.com/cihub/seelog"
)

// UserManager 聊天管理器, 负责聊天用户管理, 登录登出验证等. 管理所有在线用户。
type UserManager struct {
	userMap *sync.Map
}

// UserMgr 用户管理器单例
var UserMgr = &UserManager{
	userMap: &sync.Map{},
}

// GetUserMgr 获取用户管理器单例
func GetUserMgr() *UserManager {
	return UserMgr
}

// LoginUser 加入用户
func (mgr *UserManager) LoginUser(id UserID, srvID string) (types.IUser, error) {
	// 如果已经登录, 则更新，重新DB加载太费。
	if u := mgr.GetUser(id); u != nil {
		u.SetSrvID(srvID)
		return u, nil
	}

	u, err := user.NewUser(id, srvID)
	if err != nil {
		return nil, fmt.Errorf("NewUser: %w", err)
	}
	mgr.userMap.Store(id, u) // 必须在 OnLogin() 之前加入userMgr
	u.OnLogin()              // callback UUT_LOGIN_WITH_NICK_AND_DATA
	// OfflinedInfoMgr 只保存离线的，userMap添加之后立即删除
	user.GetOfldInfoMgr().Remove(id)
	return u, nil
}

// RemoveUser 移除用户
func (mgr *UserManager) RemoveUser(id UserID) {
	// log.Debug("User Delete: ", id)
	u := mgr.GetUser(id)
	if u == nil {
		return // 无此用户
	}

	// OfflinedInfoMgr 保存离线的，userMap删除前添加
	user.GetOfldInfoMgr().Add(u.GetFriendInfo())
	u.OnLogout()

	// 最后才删
	mgr.userMap.Delete(id)
}

// GetUser 获取用户
func (mgr *UserManager) GetUser(id UserID) types.IUser {
	v, ok := mgr.userMap.Load(id)
	if !ok {
		return nil
	}

	return v.(types.IUser)
}

// GetUserGroupMgr 获取用户 userGroupMgr 指针
func (mgr *UserManager) GetUserGroupMgr(id UserID) types.IUserGroupMgr {
	if user := mgr.GetUser(id); user != nil {
		return user.GetUserGroupMgr()
	}
	return nil
}

// GetUserFriendMgr 获取用户 friendMgr 指针
func (mgr *UserManager) GetUserFriendMgr(id UserID) types.IFriendMgr {
	if user := mgr.GetUser(id); user != nil {
		return user.GetFriendMgr()
	}
	return nil
}

// GetUserFollowMgr 获取用户 followMgr 指针
func (mgr *UserManager) GetUserFollowMgr(id UserID) types.IFollowMgr {
	if user := mgr.GetUser(id); user != nil {
		return user.GetFollowMgr()
	}
	return nil
}

// GetUserFriendInfo 获取用户信息
// nil 表示不在线
func (mgr *UserManager) GetUserFriendInfo(id UserID) *chatapi.FriendInfo {
	if user := mgr.GetUser(id); user != nil {
		info := user.GetFriendInfo()
		return &info
	}
	return nil
}

func (mgr *UserManager) SaveAllToDBOnExit() {
	log.Info("user manager save all to DB...")
	mgr.userMap.Range(func(_, user interface{}) bool {
		user.(types.IUser).SaveOnLogout()
		return true
	})
	log.Info("user manager done")
}

func (mgr *UserManager) IsUserOnline(userID UserID) bool {
	return mgr.GetUser(userID) != nil
}
