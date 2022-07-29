package usrgrp

import (
	"Cinder/Chat/rpcproc/logic/user/internal/usrgrp/dbutil"
	"sync"
	"time"

	log "github.com/cihub/seelog"
	"github.com/pkg/errors"
)

// 记录聊天群的已收序号
// 协程安全
type GroupToSeq struct {
	mtx    sync.Mutex
	userID UserID
	g2s    map[GroupID]SequenceID

	// 变化 100 次，或 5min 才保存DB
	changeTimes   int
	changedGroups map[GroupID]bool
	nextSaveTime  time.Time
}

// NewGroupToSeq 创建对象
func NewGroupToSeq(userID UserID) *GroupToSeq {
	res := &GroupToSeq{
		userID: userID,
		g2s:    make(map[GroupID]SequenceID),
	}
	res.resetOnSaved()
	return res
}

// resetOnSaved 在保存完成后重置某些成员
func (g *GroupToSeq) resetOnSaved() {
	g.changeTimes = 0
	g.changedGroups = make(map[GroupID]bool)
	g.nextSaveTime = time.Now().Add(5 * time.Minute)
}

// Update 更新序号，须保证序号单调增加
func (g *GroupToSeq) Update(groupID GroupID, seqID SequenceID) {
	g.mtx.Lock()
	defer g.mtx.Unlock()

	oldSeqID := g.g2s[groupID]
	if oldSeqID >= seqID {
		return // 忽略错序的序号
	}
	g.g2s[groupID] = seqID

	g.changeTimes++
	g.changedGroups[groupID] = true
	g.tryToSave()
}

// AddGroupIfNot 添加群。
// 仅内存操作，不写DB。
// 添加删除群写DB无论用户是否在线都会立即执行，但GroupToSeq只有上线用户才有。
func (g *GroupToSeq) AddGroupIfNot(groupID GroupID, seqID SequenceID) {
	g.mtx.Lock()
	defer g.mtx.Unlock()

	// 如果不存在，则初始化为seqID
	if _, ok := g.g2s[groupID]; ok {
		return
	}
	g.g2s[groupID] = seqID
}

// DeleteGroup 删除群。
// 仅内存操作，不写DB
// 添加删除群写DB无论用户是否在线都会立即执行，但GroupToSeq只有上线用户才有。
func (g *GroupToSeq) DeleteGroup(groupID GroupID) {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	delete(g.g2s, groupID)
	delete(g.changedGroups, groupID)
}

// GetGroupIDs 获取群号列表
func (g *GroupToSeq) GetGroupIDs() []GroupID {
	g.mtx.Lock()
	defer g.mtx.Unlock()

	result := make([]GroupID, 0, len(g.g2s))
	for groupID, _ := range g.g2s {
		result = append(result, groupID)
	}
	return result
}

// CopyGroupToSeq 复制群及已读序号
func (g *GroupToSeq) CopyGroupToSeq() map[GroupID]SequenceID {
	g.mtx.Lock()
	defer g.mtx.Unlock()

	result := make(map[GroupID]SequenceID)
	for groupID, seq := range g.g2s {
		result[groupID] = seq
	}
	return result
}

// Load 从DB加载群已读序号
func (g *GroupToSeq) Load() error {
	g.mtx.Lock()
	defer g.mtx.Unlock()

	var err error
	db := dbutil.UsersGroupsSeqUtil(g.userID)
	g.g2s, err = db.LoadGroupToSeq()
	if err != nil {
		log.Debugf("LoadGroups error: %v", err)
		return errors.Wrap(err, "db.LoadGroups")
	}
	// log.Debugf("groupToSeq: %v", g.g2s)
	return nil
}

// Save 强制保存
func (g *GroupToSeq) Save() {
	g.mtx.Lock()
	defer g.mtx.Unlock()

	g.save()
}

// save 强制保存
func (g *GroupToSeq) save() {
	// 判断是否有数据需要保存
	if len(g.changedGroups) <= 0 {
		g.resetOnSaved() // 按保存成功处理
		return
	}

	changedGroupToSeq := g.getChangedGroupToSeq()
	if len(changedGroupToSeq) <= 0 {
		g.resetOnSaved() // 按保存成功处理
		return
	}

	db := dbutil.UsersGroupsSeqUtil(g.userID)
	if err := db.UpdateGroupToSeq(changedGroupToSeq); err != nil {
		log.Errorf("db save groups error: %v", err)
		return
	}
	g.resetOnSaved()
}

// tryToSave 试图保存
func (g *GroupToSeq) tryToSave() {
	if g.canSave() {
		g.save()
	}
}

// canSave 是否可保存。
// 100次变化或5min间隔才保存一次。
func (g *GroupToSeq) canSave() bool {
	return g.changeTimes > 100 || time.Now().After(g.nextSaveTime)
}

// getChangedGroupToSeq 获取有改变的群，及其已读序号
func (g *GroupToSeq) getChangedGroupToSeq() map[GroupID]SequenceID {
	result := make(map[GroupID]SequenceID)
	for groupID, _ := range g.changedGroups {
		if seqID, ok := g.g2s[groupID]; ok { // 群可能已删
			result[groupID] = seqID
		}
	}
	return result
}
