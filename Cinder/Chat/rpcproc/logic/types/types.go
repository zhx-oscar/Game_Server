package types

import (
	"Cinder/Chat/chatapi"
)

type UserID string
type GroupID string

// 群内的消息序号，从1开始
type SequenceID uint64

type IBlacklist interface {
	Add(userID UserID) error
	Remove(userID UserID) error
	GetBlacklistInfos() []*chatapi.FriendInfo
}

type IUser interface {
	SetSrvID(srvID string)
	GetSrvID() string
	SetNick(nick string) (changed bool, e error)
	GetNick() string
	SetActiveData(data []byte) (changed bool, e error)
	GetActiveData() []byte
	SetPassiveData(data []byte) error
	GetPassiveData() []byte

	RecvChatMessage(from UserID, fromNick string, fromData []byte, msgContent []byte)
	GetFriendInfo() chatapi.FriendInfo
	OnLogout()
	GetUserGroupMgr() IUserGroupMgr
	GetFriendMgr() IFriendMgr
	GetFollowMgr() IFollowMgr
	SaveOnLogout()
	GetNickAndData() (nick string, data []byte)
	GetBlacklist() IBlacklist
	SendChatMessage(to UserID, msgContent []byte)
	LoadOfflineMessage() (*chatapi.OfflineMessage, error)
}

type IFollowMgr interface {
	AddFollower(userID UserID)
	DeleteFollower(userID UserID)
	Follow(userID UserID) error
	Unfollow(userID UserID) error
	GetFollowingList() ([]*chatapi.FriendInfo, error)
	GetFollowerList() ([]*chatapi.FollowerInfo, error)
}

type IFriendMgr interface {
	RecvRequest(fromID UserID, reqInfo []byte)
	RecvResponse(fromID UserID, ok bool)
	DeleteFriendPassive(friendID UserID)
	AddFriendWithoutDB(responderID UserID)
	SendRequest(friendID UserID, reqInfo []byte) error
	SendResponse(fromID UserID, ok bool) error
	DeleteFriendActive(friendID UserID) error
	GetFriendList() ([]*chatapi.FriendInfo, error)
	GetFriendIDs() []UserID
	GetFriendCount() int
}

type IUserGroupMgr interface {
	AddGroup(groupID GroupID, seqID SequenceID)
	DeleteGroup(groupID GroupID)
	LogoutAllGroups()
	UpdateGroupMsgSeq(groupID GroupID, sequenceID SequenceID)
	GetGroupIDs() []GroupID
}

type IUserManager interface {
	GetUser(id UserID) IUser
	GetUserFollowMgr(userID UserID) IFollowMgr
	GetUserFriendMgr(userID UserID) IFriendMgr
	GetUserFriendInfo(userID UserID) *chatapi.FriendInfo
	GetUserGroupMgr(userID UserID) IUserGroupMgr
}

type IGroup interface {
	CopyOnlineMemberIDs() map[UserID]bool
	AddMembers(userIDs []UserID) error
	GetGroupMemberIDs() []UserID
	LoginMember(userID UserID)
	GetGroupMemberCount() int
	AddMember(userID UserID) error
	KickMember(userID UserID)
	SendGroupMessage(fromRoleID UserID, msgContent []byte)
	GetHistoryMessages(count uint16) []*chatapi.ChatMessage
}
