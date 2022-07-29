package chatapi

import (
	"Cinder/Chat/chatapi/internal"
	"Cinder/Chat/chatapi/types"
	"errors"
)

// 接口使用说明见：[chat_usage.md](chat_usage.md)

type ChatMessage = types.ChatMessage
type OfflineMessage = types.OfflineMessage
type FriendInfo = types.FriendInfo
type FollowerInfo = types.FollowerInfo

const (
	SelfReachedMaxFriendCount  = string("you have reached the max friend count")
	PeerReachedMaxFriendCount  = string("the peer has reached the max friend count")
	PeerReachedMaxAddFriendReq = string("the peer has received too many requests to add friend")
)

var (
	// 加好友接口 AddFriendReq(), ReplyAddFriendReq() 会因为本人或对方到达最大好友数而返回以下错误
	ErrSelfReachedMaxFriendCount = errors.New(SelfReachedMaxFriendCount)
	ErrPeerReachedMaxFriendCount = errors.New(PeerReachedMaxFriendCount)
	// AddFriendReq() 会限制对方可接收请求个数，超限则返回以下错误
	ErrPeerReachedMaxAddFriendReq = errors.New(PeerReachedMaxAddFriendReq)
)

// NewUser 返回IUser接口
// 一般需要随后调用 Login().
func NewUser(id string) IUser {
	return newChatUser(id)
}

// Login 玩家登录时创建对象。
// 废弃，应该使用 NewUser(id).LoginWithNickData(nick, data)
func Login(id string, nick string, activeData []byte) (IUser, error) {
	user := NewUser(id)
	return user, user.LoginWithNickData(nick, activeData)
}

// Logout 玩家登出
// 废弃，应该使用 NewUser(id).Logout()
func Logout(id string) error {
	return NewUser(id).Logout()
}

// CreateGroup 创建聊天群。
// 如果已存在，则合并成员，即相当于添加成员。允许重复创建，允许成员列表为nil或空。
func CreateGroup(groupID string, members []string) error {
	return internal.CreateGroup(groupID, members)
}

// DeleteGroup 删除聊天群。
func DeleteGroup(groupID string) error {
	return internal.DeleteGroup(groupID)
}

// GetOrCreateGroup 获取聊天群接口。
// 并不会立即请求聊天服创建聊天群，而是在调用群接口时才创建。
// 废弃，改用 NewGroup(groupID string) IGroup
func GetOrCreateGroup(groupID string) (IGroup, error) {
	// 空群无需创建，调用群接口时聊天服会自动创建
	return NewGroup(groupID), nil
}

// NewGroup 创建聊天群接口。
// 并不会立即请求聊天服创建聊天群，而是在调用群接口时才创建。
func NewGroup(groupID string) IGroup {
	return &_ChatGroup{
		id: groupID,
	}
}

// 获取一批群的成员数，用于判断多个世界群的大小(繁忙程度)
// 返回map, 以群ID为键。返回 nil 表示出错。
func GetGroupMemberCounts(groupIDs []string) map[string]int {
	return internal.GetGroupMemberCounts(groupIDs)
}

// GetFriendInfos 查询一系列玩家的数据
// 返回 nil 表示出错
func GetFriendInfos(userIDs []string) []FriendInfo {
	return internal.GetFriendInfos(userIDs)
}

type IUser interface {
	// Login 玩家登录时通知聊天服，只有 Login() 之后才能收到聊天服的各种回调。
	// 其他接口并不要求 Login().
	// 允许重复多次 Login().
	Login() error
	// activeData 是额外的主动数据，将在消息推送时原样输入到RPC接口。
	LoginWithNickData(nick string, activeData []byte) error
	// Logout 退出，之后就不会有聊天服的回调了
	// Game服宕时，聊天服不会知道，所以不会自动Logout相应的用户，仍认为在线。
	// Game服收到聊天服回调时发现用户已不在该服时，需调用Logout()通知聊天服。
	Logout() error

	// 昵称和自定义数据读写, 自定义数据分为主动数据和被动数据。
	// 登录，昵称和主动数据变化会主动通知相关在线用户，被动数据需要用户自己查询。
	SetNick(nick string) error
	GetNick() (string, error)
	SetActiveData(data []byte) error
	GetActiveData() ([]byte, error)
	SetPassiveData(data []byte) error
	GetPassiveData() ([]byte, error)

	// 获取离线消息
	GetOfflineMessage() *OfflineMessage

	// 发送私聊消息
	SendMessage(target string, msgContent []byte) error

	// 聊天群操作
	// 应该使用自由函数 CreateGroup(), DeleteGroup()
	CreateGroup(groupID string, members []string) error
	DeleteGroup(groupID string) error

	// 应该使用 IGroup.AddIntoGroup()/KickFromGroup()
	JoinGroup(groupID string) error
	LeaveGroup(groupID string) error

	KickFromGroup(groupID string, member string) error
	GetGroupMembers(groupID string) ([]string, error)
	SendGroupMessage(groupID string, msgContent []byte) error

	// 关注功能相关接口。关注是单向关系。
	FollowFriendReq(friendID string) error
	UnFollowFriendReq(friendID string) error
	GetFollowingList() []FriendInfo
	GetFollowerList() []FollowerInfo

	// 好友相关接口。好友是双向关系。
	// 如本人或对方到达最大好友数，分别返回 ErrSelfReachedMaxFriendCount 或 ErrPeerReachedMaxFriendCount
	AddFriendReq(friendID string, reqInfo []byte) error
	// 因名字有歧义而作废，改用 ReplyAddFriendReq()
	ApplyAddFriendReq(fromID string, ok bool) error
	// 同意或拒绝加好友请求
	ReplyAddFriendReq(fromID string, ok bool) error

	// 黑名单功能。黑名单是单向关系。
	AddFriendToBlacklist(friendID string) error
	RemoveFriendFromBacklist(friendID string) error // 作废，待删除
	RemoveFriendFromBlacklist(friendID string) error

	DeleteFriend(friendID string) error

	GetFriendList() []FriendInfo
	GetFriendBlacklist() []FriendInfo
}

type IGroup interface {
	AddIntoGroup(member string) error
	KickFromGroup(member string) error
	GetGroupMembers() ([]string, error)
	SendGroupMessage(fromID string, msgContent []byte) error
	// 获取群历史消息，取最近 count 条
	GetHistoryMessages(count uint16) []*ChatMessage
}
