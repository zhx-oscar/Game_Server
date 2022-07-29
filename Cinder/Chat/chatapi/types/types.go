// types 定义结构，避免包循环引用
package types

import (
	"time"
)

type ChatMessage struct {
	From       string // 发送人ID
	FromNick   string // 发送人昵称
	FromData   []byte // 发送人自定义数据，在Login时注册
	SendTime   int64  // 发送时间，Unix 时戳，1970年以来的秒数
	MsgContent []byte // 内容
}

// OfflineMessage 离线消息
type OfflineMessage struct {
	P2PMessages   []*ChatMessage            // 私聊
	GroupMessages map[string][]*ChatMessage // 群聊，以群ID为键
}

type FriendInfo struct {
	ID   string
	Nick string
	Data []byte // Login 时注册的主动数据(active data)

	IsOnline       bool      // 是否在线
	OfflineTime    time.Time // 离线时间
	FollowerNumber int       // 粉丝数
}

type FollowerInfo struct {
	FriendInfo *FriendInfo
	FollowTime time.Time // 关注时间，加粉丝的时间
}

// Target 用于 RPC_ChatRecvGroupMessageV2 的参数，
// 及 RPC_ChatUpdateUser 时的目标
type Target struct {
	ID   string // 目标ID
	Data []byte // 目标主动数据，Login()时发送的数据
}

type UserUpdateType int

const (
	UUT_ILLEGAL                  UserUpdateType = 0
	UUT_LOGOUT                   UserUpdateType = 2
	UUT_NICK                     UserUpdateType = 3
	UUT_DATA                     UserUpdateType = 4
	UUT_LOGIN_WITH_NICK_AND_DATA UserUpdateType = 5
)

type UpdateUserMsg struct {
	Targets []Target

	UserID string         // 变更的用户
	Type   UserUpdateType // 变更类型
	Nick   string
	Data   []byte // 仅主动数据(active data)
}
