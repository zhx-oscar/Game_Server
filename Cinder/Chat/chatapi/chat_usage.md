# 聊天服使用说明

本文最新版本见: 
http://gitlab.ztgame.com/Cinder/Server/Cinder/blob/master/Chat/chatapi/chat_usage.md

本文说明如何接入聊天服。聊天服的功能和部署见：[../README.md](../README.md)

请使用 `Cinder/Chat/chatapi` 包接入聊天服。仅支持 go 语言并应用 Cinder 框架开发的游戏。

chatapi 包是 Chat 服的SDK包.
chatapi 由Game，Space或其他服使用，只需调用其导出接口，而不必直接调用Chat服的RPC接口。

## 使用示例
```
	user := chatapi.NewUser(userID)
	if err := user.Login(userID); err != nil {
		panic(err)
	}
	msg := user.GetOfflineMessage()
	errSnd := user.SendMessage(target, msgContent)
```

## 回调接口

要求调用服注册以下几个 RPC 处理函数, 供Chat服回调：

* `RPC_ChatRecvP2PMessage(targetID string, targetData []byte, fromID string, fromNick string, fromData []byte, msgContent []byte)`
	+ 说明：接收私聊消息，接收者回调
	+ 参数
		- targetID: 目标ID
		- targetData: 目标数据，Login()时发送的数据
		- fromID: 来自ID
		- fromNick: 来自昵称
		- fromData: 来自数据，消息发送者Login()时的数据
		- msgContent: 消息内容
* `RPC_ChatRecvGroupMessageV2(groupID string, targetsJson []byte, fromID string, fromNick string, fromData []byte, msgContent []byte)`
	+ 说明：接收群聊消息，接收者回调
	+ 参数
		- groupID: 聊天群ID
		- targetsJson: []Target 的 Json 打包，每个 Target 包含
			- TargetID: 目标ID
			- TargetData: 目标数据，Login()时发送的数据
		- fromID: 来自ID
		- fromNick: 来自昵称
		- fromData: 来自数据，消息发送者Login()时的数据
		- msgContent: 消息内容
* `RPC_AddFriendReq(targetID string, fromID string, reqInfo []byte)`
	+ 说明：请求加好友，被加方回调
	+ 参数
		- targetID: 目标ID
		- fromID: 来自ID
		- reqInfo: 请求信息
* `RPC_AddFriendRet(targetID string, fromID string, ok bool)`
	+ 说明：加好友结果返回，请求方回调
	+ 参数
		- targetID: 目标ID
		- fromID: 来自ID
		- ok: 是否成功，拒绝则为否
* `RPC_AddDelFollower(targetID string, targetData []byte, followerID string, isAdd bool)`
	+ 说明：followerID 关注或取消关注 targetID 时，通知 targetID 的回调
	+ 参数
		- targetID: 目标ID
		- targetData: 目标数据，Login()时发送的数据
		- followerID: 粉丝ID
		- isAdd: 是否添加，取消关注则为否
* `RPC_ChatUpdateUser(updateUserMsgJson []byte)`
	+ 说明：用户上下线，更改昵称或数据时，通知所有好友和群
		- 不是对每个好友和群成员发送，而是对所在服发送一次，消息内包含了多个待通知的目标
	+ 参数：updataUserMsgJson 是 UpdateUserMsg 的 json 打包
* `RPC_FriendDeleted(targetID string, friendID string)`
	+ 说明：friendID 删除好友 targetID 时，通知 targetID 的回调
	+ 参数
		- targetID: 目标ID
		- frientID: 已删好友ID

## 函数

* `Login(id string, nick string, data []byte) (IUser, error)` 登录
* `Logout(id string) error` 登出
* `CreateGroup(groupID string, members []string) error` 创建聊天群
* `DeleteGroup(groupID string) error` 删除聊天群
* `GetOrCreateGroup(groupID string) (IGroup, error)` 获取或创建聊天群
* `GetGroupMemberCounts(groupIDs []string) map[string]int` 获取多个聊天群的成员数
* `GetFriendInfos(userIDs []string) []FriendInfo` 获取一批好友信息


## 接口

### `IUser`

使用 `NewUser(id)` 获取 `IUser` 接口。有以下功能：

```
	// Login 玩家登录时通知聊天服，只有 Login() 之后才能收到聊天服的各种回调。
	Login() error
	LoginWithNickData(nick string, activeData []byte) error
	// Logout 退出，之后就不会有聊天服的回调了
	Logout() error

	// 昵称和自定义数据读写, 自定义数据分为主动数据和被动数据。
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
	AddFriendReq(friendID string, reqInfo []byte) error
	ApplyAddFriendReq(fromID string, ok bool) error

	// 黑名单功能。黑名单是单向关系。
	AddFriendToBlacklist(friendID string) error
	RemoveFriendFromBlacklist(friendID string) error

	DeleteFriend(friendID string) error

	GetFriendList() []FriendInfo
	GetFriendBlacklist() []FriendInfo
```

### `IGroup`

通过 `NewGroup(id)` 获取 `IGroup` 接口。有以下功能：
```
	AddIntoGroup(member string) error
	KickFromGroup(member string) error
	GetGroupMembers() ([]string, error)
	SendGroupMessage(fromID string, msgContent []byte) error
	// 获取群历史消息，取最近 count 条
	GetHistoryMessages(count uint16) []*ChatMessage
```

## 数据类型

* `ChatMessage struct` 聊天消息
* `OfflineMessage struct` 离线消息
* `FriendInfo struct` 好友信息
* `FollowerInfo struct` 粉丝信息
* `UpdateUserMsg struct` 用户更新通知消息
