package Const

const (
	// ChatGroupTypeWorld 世界聊天频道
	ChatGroupTypeWorld = "world"

	// GroupHistoryMessageCount 世界、工会频道的历史消息数
	GroupHistoryMessageCount = 25

	//ChatGroupTypeTeam 队伍聊天频道
	ChatGroupTypeTeam = "team"

	//ChatTeamHistoryMessageCount 队伍频道历史消息数量限制
	ChatTeamHistoryMessageCount = 100

	//ChatMsgMaxCount 聊天内容长度限制 40个汉字
	ChatMsgMaxCount = 40 * 2
)

type MessageType uint8

const (
	ChatTypeNormal       MessageType = iota + 1 // 1:普通类型消息
	ChatTypeRequestSkill                        //  2:乞求技能消息
)
