package group

import (
	"Cinder/Chat/chatapi"
	"Cinder/Chat/rpcproc/logic/group/internal"
	"Cinder/Chat/rpcproc/logic/usermgr"
	"sync"

	log "github.com/cihub/seelog"
	"github.com/pkg/errors"
)

// Group 用来管理群组数据
type Group struct {
	mtx sync.Mutex

	id      GroupID
	members *internal.GroupMembers
	msgs    *internal.MsgCache
}

// LoadGroup 创建对象并从DB加载
func LoadGroup(groupID GroupID) (*Group, error) {
	group := &Group{
		id:      groupID,
		members: internal.NewGroupMembers(groupID),
		msgs:    internal.NewMsgCache(groupID),
	}

	group.mtx.Lock()
	defer group.mtx.Unlock()

	// 加载最近消息和成员列表
	if err := group.load(); err != nil {
		log.Debugf("failed to load group '%s': %v", groupID, err)
		return nil, errors.Wrap(err, "load group")
	}

	return group, nil
}

func (group *Group) Save() {
	group.mtx.Lock()
	defer group.mtx.Unlock()

	// 仅需保存离线消息
	group.msgs.Save()
}

// SendGroupMessage 发送消息.
// 会触发定时保存。
func (group *Group) SendGroupMessage(fromRoleID UserID, msgContent []byte) {
	group.mtx.Lock()
	defer group.mtx.Unlock()

	// log.Debugf("send group message: '%s'->group(%s) msg: %v", fromRoleID, group.id, msgContent)
	if !group.members.HasMember(fromRoleID) {
		log.Debugf("non-member '%v' can not post in group '%v'", fromRoleID, group.id)
		return
	}

	fromUser := usermgr.GetUserMgr().GetUser(fromRoleID)
	if fromUser == nil {
		// 先Login(), 然后才能说话，不然没有 customData
		log.Debugf("(*Group).SendGroupMessage(): '%s' is not online", fromRoleID)
		return
	}
	fromNick, fromData := fromUser.GetNickAndData()

	// 群聊消息不管是否全部在线，所有消息都保存
	sequenceID := group.msgs.Add(fromRoleID, fromNick, fromData, msgContent)

	// 群发消息较耗时，异步发送，不然请求会超时
	go internal.GroupBroadcastMessage(group.id, group.CopyOnlineMemberIDs(), fromRoleID, fromNick, fromData, msgContent, sequenceID)
}

// LoginMember 设置成员在线
func (group *Group) LoginMember(userID UserID) {
	group.mtx.Lock()
	defer group.mtx.Unlock()

	group.members.Login(userID)
}

// LogoutMember 设置成员离线
func (group *Group) LogoutMember(userID UserID) {
	group.mtx.Lock()
	defer group.mtx.Unlock()

	group.members.Logout(userID)
}

// AddMember 加成员入群组, 立即存DB.
// 允许添加离线成员。
// 因为成员可能未上线，所以没有 User 对象，需直接在 DB 中为 User 添加群
func (group *Group) AddMember(userID UserID) error {
	group.mtx.Lock()
	defer group.mtx.Unlock()
	return group.addMembers([]UserID{userID})
}

// AddMembers 加成员入群组, 立即存DB.
// 允许添加离线成员。
// 因为成员可能未上线，所以没有 User 对象，需直接在 DB 中为 User 添加群
func (group *Group) AddMembers(userIDs []UserID) error {
	group.mtx.Lock()
	defer group.mtx.Unlock()
	return group.addMembers(userIDs)
}

// addMembers 加成员入群组, 立即存DB.
// 允许添加离线成员。
// 因为成员可能未上线，所以没有 User 对象，需直接在 DB 中为 User 添加群
func (group *Group) addMembers(members []UserID) error {
	// log.Debugf("add members to group '%s': %v", group.id, members)
	return group.members.AddMembers(members, group.msgs.GetMaxSeqID())
}

// KickMember 让成员离开群组, 立即存DB
func (group *Group) KickMember(userID UserID) {
	group.mtx.Lock()
	defer group.mtx.Unlock()
	group.members.DeleteMember(userID)
}

// GetOfflineMessagesAfter 获取离线消息。
// seqID 是已读取序号, 应该返回 seqID+1 及后面的消息。
func (group *Group) GetOfflineMessagesAfter(seqID SequenceID) []*chatapi.ChatMessage {
	group.mtx.Lock()
	defer group.mtx.Unlock()
	return group.msgs.GetMsgsAfter(seqID)
}

// load 加载最近消息，成员列表, 同时更新在线列表。
func (group *Group) load() error {
	if err := group.msgs.Load(); err != nil { // 包括当前序号
		return errors.Wrap(err, "load offline messages")
	}
	if err := group.members.Load(); err != nil {
		return errors.Wrap(err, "load members")
	}
	return nil
}

// GetGroupMemberIDs 获取成员ID列表。
// 以复制方式获取。
func (group *Group) GetGroupMemberIDs() []UserID {
	group.mtx.Lock()
	defer group.mtx.Unlock()
	return group.members.GetMemberIDs()
}

// GetGroupMemberCount 获取成员个数
func (group *Group) GetGroupMemberCount() int {
	group.mtx.Lock()
	defer group.mtx.Unlock()
	return group.members.GetMemberCount()
}

// GetHistoryMessages 获取最近 count 条历史消息
func (group *Group) GetHistoryMessages(count uint16) []*chatapi.ChatMessage {
	group.mtx.Lock()
	defer group.mtx.Unlock()
	return group.msgs.GetHistoryMessages(count)
}

func (group *Group) CopyOnlineMemberIDs() map[UserID]bool {
	return group.members.CopyOnlineMemberIDs()
}
