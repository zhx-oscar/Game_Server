package user

import (
	"Cinder/Chat/chatapi"
	"Cinder/Chat/rpcproc/logic/bc"
	"Cinder/Chat/rpcproc/logic/rpc"
	"Cinder/Chat/rpcproc/logic/types"
	"Cinder/Chat/rpcproc/logic/user/internal/blklst"
	"Cinder/Chat/rpcproc/logic/user/internal/dbutil"
	"Cinder/Chat/rpcproc/logic/user/internal/follow"
	"Cinder/Chat/rpcproc/logic/user/internal/friend"
	"Cinder/Chat/rpcproc/logic/user/internal/oflmsg"
	"Cinder/Chat/rpcproc/logic/user/internal/userinfo"
	"Cinder/Chat/rpcproc/logic/user/internal/usrgrp"
	"fmt"

	assert "github.com/arl/assertgo"
	log "github.com/cihub/seelog"
)

// User 代表一个玩家
type User struct {
	// User 本身无锁，需要各成员协程安全

	id UserID

	// 用户数据，Login() 所注册的数据, FriendInfo, 用锁保护
	userInfo *userinfo.UserInfo

	// 玩家的所有群
	userGroupMgr *usrgrp.UserGroupMgr

	// 禁止自己拉黑自己，关注自己，加自己好友，防止锁重入死锁
	blacklist *blklst.Blacklist // 黑名单
	followMgr *follow.FollowMgr // 关注和粉丝
	friendMgr *friend.FriendMgr // 好友管理器
}

// NewUser 创建玩家
func NewUser(id UserID, srvID string) (*User, error) {
	// 先按输入创建 UserInfo, 再从DB加载
	userInfo := userinfo.NewUserInfo(id, srvID)
	u := &User{
		id:       id,
		userInfo: userInfo,

		userGroupMgr: usrgrp.NewUserGroupMgr(id),

		blacklist: blklst.NewBlacklist(id),
		followMgr: follow.NewFollowMgr(id, userInfo),
		friendMgr: friend.NewFriendMgr(id, userInfo),
	}

	if err := u.load(); err != nil {
		return nil, fmt.Errorf("load user: %w", err)
	}
	return u, nil
}

func (user *User) String() string {
	return fmt.Sprintf("User:%s", user.id)
}

func (user *User) GetFriendMgr() types.IFriendMgr {
	return user.friendMgr
}

func (user *User) GetUserGroupMgr() types.IUserGroupMgr {
	return user.userGroupMgr
}

func (user *User) GetBlacklist() types.IBlacklist {
	return user.blacklist
}

func (user *User) GetFollowMgr() types.IFollowMgr {
	return user.followMgr
}

// SendChatMessage 发送私聊消息到 to.
func (user *User) SendChatMessage(to UserID, msgContent []byte) {
	fromNick, fromData := user.GetNickAndData()
	targetUser := userMgr.GetUser(to)
	if targetUser != nil {
		targetUser.RecvChatMessage(user.id, fromNick, fromData, msgContent)
		return
	}

	oflmsg.MsgCache.Add(user.id, fromNick, fromData, to, msgContent)
}

// RecvChatMessage 从from接收私聊消息.
func (user *User) RecvChatMessage(from UserID, fromNick string, fromData []byte, msgContent []byte) {
	activeData := user.userInfo.GetActiveData()
	srvID := user.userInfo.GetSrvID()
	go func() {
		// 回调尽量新开协程
		ret := rpc.Rpc(srvID, "RPC_ChatRecvP2PMessage", string(user.id), activeData, string(from), fromNick, fromData, msgContent)
		if ret.Err != nil {
			log.Errorf("RPC_ChatRecvP2PMessage: %s", ret.Err)
		}
	}()
}

// LoadOfflineMessage 从DB加载离线消息。
// 仅私聊消息需要加载。群聊消息已在群加载时加载。
func (user *User) LoadOfflineMessage() (*chatapi.OfflineMessage, error) {
	// 各群的 sequenceID 已加载为 user.initGroupToSeq
	// TODO: 需保证实时消息发送前，离线消息已经发送完成，不然会有重叠

	// 加载离线私聊
	p2pMessages, err := user.loadP2POfflineMessages()
	if err != nil {
		return nil, fmt.Errorf("load p2p offline messages error: %w", err)
	}
	// 加载后即从DB删除
	if err := user.removeDBP2POfflineMessages(); err != nil {
		return nil, fmt.Errorf("remove p2p offline messages error: %w", err)
	}
	// 需要添加内存中的离线消息
	p2pMessages = append(p2pMessages, oflmsg.MsgCache.Pop(user.id)...)

	return &chatapi.OfflineMessage{
		P2PMessages: p2pMessages,
		// 加载各群的离线数据
		GroupMessages: user.userGroupMgr.GetGroupOfflineMessageMap(),
	}, nil
}

// loadP2POfflineMessages 从DB加载私聊离线消息
func (user *User) loadP2POfflineMessages() ([]*chatapi.ChatMessage, error) {
	return dbutil.UserOflnMsgUtil(user.id).Load()
}

// removeDBP2POfflineMessages 从DB删除私聊离线消息
func (user *User) removeDBP2POfflineMessages() error {
	return dbutil.UserOflnMsgUtil(user.id).Remove()
}

// SaveOnLogout 保存数据到DB.
func (user *User) SaveOnLogout() {
	// 写 OfflineTime
	if err := dbutil.UsersUtil(user.id).SaveOfflineTime(); err != nil {
		log.Errorf("failed to save '%v' offline time: %v", user.id, err)
	}

	// 群已读序号需写DB
	user.userGroupMgr.Save()

	// 其他都是变化时就写DB的。
}

// load 从DB加载玩家相关数据
func (user *User) load() error {
	// chat.users
	if err := user.loadUser(); err != nil {
		return fmt.Errorf("load user: %w", err)
	}

	// chat.users.groups
	// 实际需要的是从DB加载用户的所有群已读记录数，但是为了简化处理，设计为预加载玩家的所有群，详见 README.md
	if err := user.userGroupMgr.LoadGroups(); err != nil {
		return fmt.Errorf("load user groups: %w", err)
	}

	// chat.users.followers
	if err := user.followMgr.LoadFollowers(); err != nil {
		return fmt.Errorf("follow mgr load: %w", err)
	}

	// chat.users.friend_reqeusts, chat.users.friend_responses
	// 加载成功后立即处理加好友请求和应答
	if err := user.friendMgr.HandleReqResp(); err != nil {
		return fmt.Errorf("handle friend requests and responses: %w", err)
	}

	// chat.users.offline_messages 将在请求时读取
	return nil
}

// loadUser 加载 chat.users 集合中的用户文档
func (user *User) loadUser() error {
	userDoc, err := dbutil.UsersUtil(user.id).Load()
	if err != nil {
		return fmt.Errorf("load user: %w", err)
	}
	assert.True(userDoc != nil)

	if err := user.initWithUserDoc(userDoc); err != nil {
		return fmt.Errorf("init with user doc: %w", err)
	}
	return nil
}

// initWithUserDoc 用DB中的用户文档初始化
func (user *User) initWithUserDoc(doc *dbutil.UserDoc) error {
	// FriendInfo 取 doc 中的数据。
	user.userInfo.InitWithUserDoc(doc)

	// 群数据不在 doc 中，userGroupMgr 将另外加载: userGroupMgr.LoadGroups()
	user.blacklist.Set(doc.Blacklist)
	user.followMgr.SetFollowings(doc.FollowIDs)
	// 粉丝列表不在 doc 中，followers 将另外加载: followMgr.LoadFollowers()
	user.friendMgr.SetFriends(doc.FriendIDs)
	return nil
}

func (user *User) SetSrvID(srvID string) {
	user.userInfo.SetSrvID(srvID)
}

func (user *User) SetNick(nick string) (changed bool, e error) {
	return user.userInfo.SetNick(nick)
}

func (user *User) GetNick() string {
	return user.userInfo.GetNick()
}

func (user *User) SetActiveData(data []byte) (changed bool, e error) {
	return user.userInfo.SetActiveData(data)
}

// GetData 获取注册的自定义主动数据
func (user *User) GetActiveData() []byte {
	return user.userInfo.GetActiveData()
}

func (user *User) SetPassiveData(data []byte) error {
	return user.userInfo.SetPassiveData(data)
}

func (user *User) GetPassiveData() []byte {
	return user.userInfo.GetPassiveData()
}

// GetNickAndData 获取Login时注册的Nick和自定义主动数据
func (user *User) GetNickAndData() (nick string, data []byte) {
	return user.userInfo.GetNickAndData()
}

// GetFriendInfo 获取用户信息
func (user *User) GetFriendInfo() chatapi.FriendInfo {
	return user.userInfo.GetFriendInfo()
}

// GetSrvID 获取 srvID
func (user *User) GetSrvID() string {
	return user.userInfo.GetSrvID()
}

// OnLogin 处理登录事件
func (user *User) OnLogin() {
	// 广播更新
	bc.BroadcastUserLogin(user.id)
}

// OnLogout 处理登出事件
func (user *User) OnLogout() {
	// 广播更新
	bc.BroadcastUserLogout(user.id)

	user.GetUserGroupMgr().LogoutAllGroups()
	user.SaveOnLogout()
}
