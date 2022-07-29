package internal

import (
	"Cinder/Chat/chatapi"
	"Cinder/Chat/rpcproc/logic/chatmsg"
	"Cinder/Chat/rpcproc/logic/group/internal/dbutil"
	"Cinder/Chat/rpcproc/logic/types"
	"time"

	assert "github.com/arl/assertgo"
	log "github.com/cihub/seelog"
)

type GroupID = types.GroupID
type SequenceID = types.SequenceID

// 每个群在内存中仅保存1000条，防止内存占用太大。
// 因为离线消息仅读自内存，所以也受此限制。
const kMaxCacheCount = 1000

// 群聊消息缓存。
// 不必加锁，因为Group中已有锁。
type MsgCache struct {
	groupID GroupID

	// 当前保存的最小序号，用于删除
	minSeqID SequenceID
	// 当前记录序号, msgs中的最大序号
	maxSeqID SequenceID

	// 控制离线消息保存, 每100条 或者 5min 就保存一次
	prevSaveSeq  SequenceID
	nextSaveTime time.Time

	// 群聊消息缓存
	msgs map[SequenceID]*chatapi.ChatMessage
}

func NewMsgCache(groupID GroupID) *MsgCache {
	return &MsgCache{
		groupID:      groupID,
		nextSaveTime: time.Now(), // 避免 nil
		msgs:         make(map[SequenceID]*chatapi.ChatMessage),
	}
}

// updateSavePos 记录 save() 的时间和数值
func (m *MsgCache) updateSavePos() {
	m.prevSaveSeq = m.maxSeqID
	m.nextSaveTime = time.Now().Add(5 * time.Minute)
}

// canSave 控制是否可以 save()。
// 100条消息 或者 5min 就保存一次
func (m *MsgCache) canSave() bool {
	return m.prevSaveSeq+100 > m.maxSeqID ||
		m.nextSaveTime.After(time.Now())
}

// TryToSave 如果条件满足就保存群数据
func (m *MsgCache) tryToSave() {
	if m.canSave() {
		m.Save()
	}
}

// Save DB保存离线消息.
func (m *MsgCache) Save() {
	prevSaveSeq := m.prevSaveSeq
	m.updateSavePos() // 5min后再保存
	if prevSaveSeq == m.maxSeqID {
		// 上次保存以来没变过
		return
	}

	// 已保存过的不保存
	db := dbutil.GroupsMessagessUtil(m.groupID)
	if err := db.Insert(m.msgs, prevSaveSeq+1, m.maxSeqID); err != nil {
		log.Errorf("failed to save offline messages, groupID=%s sequenceID=%d..%d, error=%v",
			m.groupID, prevSaveSeq+1, m.maxSeqID, err)
		return
	}
	// 需要删除旧消息吗？暂不删除。只是读取时会有Limit.
}

// Add 添加一条消息，返回序号
func (m *MsgCache) Add(fromID types.UserID, fromNick string, fromData []byte, msgContent []byte) SequenceID {
	msg := chatmsg.NewChatMessage(fromID, fromNick, fromData, msgContent)
	m.maxSeqID += 1
	m.msgs[m.maxSeqID] = msg
	// 限量 kMaxCacheCount
	m.limitCount()

	// 新消息触发保存
	m.tryToSave()
	return m.maxSeqID
}

// LoadFromDB 从DB加载离线消息
func (m *MsgCache) Load() error {
	var err error
	db := dbutil.GroupsMessagessUtil(m.groupID)
	m.msgs, err = db.Load()

	m.updateSequenceIDs() // 加载后立即更新 maxSeqID 为最大ID, minSeqID 为最小ID
	m.updateSavePos()     // 加载等同于刚保存
	return err
}

// GetMsgAfter 获取seqID之后的离线消息。
// seqID 是已读取序号, 应该返回 seqID+1 及后面的消息。
func (m *MsgCache) GetMsgsAfter(seqID SequenceID) []*chatapi.ChatMessage {
	// 消息缓存在内存中
	if seqID >= m.maxSeqID {
		return nil // 已完部读取，无离线消息
	}

	l := int(m.maxSeqID - seqID)
	result := make([]*chatapi.ChatMessage, 0, l)
	for seq := seqID + 1; seq <= m.maxSeqID; seq++ {
		result = append(result, m.msgs[seq])
	}
	return result
}

// updateSequenceIDs 用消息缓存中的最大最小序号更新当前最大最小序号
func (m *MsgCache) updateSequenceIDs() {
	// 先随便设个初值
	for seq, _ := range m.msgs {
		m.maxSeqID = seq
		m.minSeqID = seq
		break
	}
	for seq, _ := range m.msgs {
		if seq > m.maxSeqID {
			m.maxSeqID = seq
		} else if seq < m.minSeqID {
			m.minSeqID = seq
		}
		assert.True(m.minSeqID <= m.maxSeqID)
	}
	// log.Debugf("msg sequenceID updated to (%d, %d)", m.minSeqID, m.maxSeqID)
}

// limitCount 限制缓存消息数
func (m *MsgCache) limitCount() {
	assert.True(m.minSeqID <= m.maxSeqID)
	for ; m.minSeqID < m.maxSeqID; m.minSeqID++ {
		if len(m.msgs) < kMaxCacheCount {
			return
		}
		delete(m.msgs, m.minSeqID)
	}
}

// GetMaxSeqID 获取最大序号，即当前序号，最新一条消息的序号。
func (m *MsgCache) GetMaxSeqID() SequenceID {
	return m.maxSeqID
}

// GetHistoryMessages 获取最近 count 条历史消息
func (m *MsgCache) GetHistoryMessages(count uint16) []*chatapi.ChatMessage {
	if count > kMaxCacheCount {
		count = kMaxCacheCount
	}
	cached := uint16(m.maxSeqID - m.minSeqID + 1)
	if count > cached {
		count = cached
	}

	result := make([]*chatapi.ChatMessage, 0, count)
	startSeq := (m.maxSeqID + 1) - SequenceID(count)
	for seq := startSeq; seq <= m.maxSeqID; seq++ {
		result = append(result, m.msgs[seq])
	}
	return result
}
