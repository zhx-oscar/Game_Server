package chatapi

import (
	"Cinder/Chat/chatapi/internal"
)

// _ChatGroup 实现了 IGroup 接口
// 由 GetOrCreateGroup(id string) 创建
type _ChatGroup struct {
	id string
}

func (c *_ChatGroup) AddIntoGroup(member string) error {
	// 同 IUser.JoinGroup()
	return internal.MemberJoinGroup(member, c.id)
}

func (c *_ChatGroup) KickFromGroup(member string) error {
	// 同 IUser.LeaveGroup()
	return internal.MemberLeaveGroup(member, c.id)
}

func (c *_ChatGroup) GetGroupMembers() ([]string, error) {
	return internal.GetGroupMembers(c.id)
}

func (c *_ChatGroup) SendGroupMessage(fromID string, msgContent []byte) error {
	return internal.SendGroupMessage(fromID, c.id, msgContent)
}

func (c *_ChatGroup) GetHistoryMessages(count uint16) []*ChatMessage {
	return internal.GetGroupHistoryMessages(c.id, count)
}
