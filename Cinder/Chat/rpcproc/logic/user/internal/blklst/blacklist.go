package blklst

import (
	"Cinder/Chat/chatapi"
	"Cinder/Chat/rpcproc/logic/types"
	"Cinder/Chat/rpcproc/logic/user/internal/blklst/dbutil"
	"Cinder/Chat/rpcproc/logic/user/internal/oflinfos"
	"fmt"
	"sync"

	log "github.com/cihub/seelog"
	"github.com/pkg/errors"
)

// 限制黑名单长度
const kMaxBlacklistLen = 1000

type UserID = types.UserID

// Blacklist 处理黑名单功能
type Blacklist struct {
	mtx sync.Mutex

	userID UserID // 主人ID, 主人拉黑了一批人

	members map[UserID]bool // 黑名单集合
}

func NewBlacklist(userID UserID) *Blacklist {
	return &Blacklist{
		userID:  userID,
		members: make(map[UserID]bool),
	}
}

// Add 添加黑名单，将某人拉黑
func (b *Blacklist) Add(userID UserID) error {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	if userID == b.userID {
		return nil // 禁止拉黑自己
	}
	if _, ok := b.members[userID]; ok {
		return nil // 已存在
	}
	if len(b.members) > kMaxBlacklistLen {
		return fmt.Errorf("blacklist is longer than %d", kMaxBlacklistLen)
	}

	// 立即写DB
	if err := dbutil.UsersBlacklistUtil(b.userID).AddToBlacklist(userID); err != nil {
		return errors.Wrap(err, "failed to add to blacklist")
	}

	b.members[userID] = true
	return nil
}

func (b *Blacklist) Remove(userID UserID) error {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	if err := dbutil.UsersBlacklistUtil(b.userID).RemoveFromBlacklist(userID); err != nil {
		return fmt.Errorf("failed to remove from blacklist: %w", err)
	}
	return nil
}

// Set 设置黑名单
func (b *Blacklist) Set(ids []UserID) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	b.members = make(map[UserID]bool)
	for _, id := range ids {
		b.members[id] = true
	}
	delete(b.members, b.userID) // 排除自己
}

// GetBlacklistInfos 获取黑名单
func (b *Blacklist) GetBlacklistInfos() []*chatapi.FriendInfo {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	result := make([]*chatapi.FriendInfo, 0, len(b.members))
	for id, _ := range b.members {
		info, err := oflinfos.GetFriendInfo(id)
		if err != nil {
			log.Errorf("failed get blacklist of user '%v': %v", b.userID, err)
			return nil
		}
		result = append(result, info)
	}
	return result
}
