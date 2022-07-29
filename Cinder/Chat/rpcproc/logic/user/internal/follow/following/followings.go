package following

import (
	"Cinder/Chat/rpcproc/logic/types"
	"Cinder/Chat/rpcproc/logic/user/internal/follow/following/dbutil"
	"errors"
	"fmt"
	"sync"
)

const kMaxFollowCount = 1000

type UserID = types.UserID

type Followings struct {
	mtx sync.Mutex

	userID UserID // 主人ID, 主人关注哪些人

	followingIDs map[UserID]bool
}

// 添加时已经存在
var ErrNotAllowedToAdd = errors.New("not allowed to add")

// 删除时不存在
var ErrFollowingNotExists = errors.New("following not exists")

func NewFollowings(userID UserID) *Followings {
	return &Followings{
		userID:       userID,
		followingIDs: make(map[UserID]bool),
	}
}

// Add 添加关注.
// 仅修改自身数据，被关注方的粉丝列表由调用者修改。
func (f *Followings) Add(userID UserID) error {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	if userID == f.userID {
		return ErrNotAllowedToAdd // 禁止自关注
	}
	if _, ok := f.followingIDs[userID]; ok {
		return ErrNotAllowedToAdd
	}
	if len(f.followingIDs) > kMaxFollowCount {
		return ErrNotAllowedToAdd
	}

	// DB 添加 followIDs
	if err := dbutil.UsersUtil(f.userID).AddFollowID(userID); err != nil {
		return fmt.Errorf("users add following: %w", err)
	}

	f.followingIDs[userID] = true
	return nil
}

// Delete 删除关注.
// 仅修改自身数据，被关注方的粉丝列表由调用者修改。
func (f *Followings) Delete(userID UserID) error {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	if _, ok := f.followingIDs[userID]; !ok {
		return ErrFollowingNotExists
	}

	// DB 删除 followIDs
	if err := dbutil.UsersUtil(f.userID).DeleteFollowID(userID); err != nil {
		return fmt.Errorf("users delete following: %w", err)
	}

	delete(f.followingIDs, userID)
	return nil
}

// Set 设置关注列表
func (f *Followings) Set(ids []UserID) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	f.followingIDs = make(map[UserID]bool)
	for _, id := range ids {
		f.followingIDs[id] = true
	}
	delete(f.followingIDs, f.userID) // 排除自己
}

// GetIDs 获取ID列表
func (f *Followings) GetIDs() []UserID {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	result := make([]UserID, 0, len(f.followingIDs))
	for id, _ := range f.followingIDs {
		result = append(result, id)
	}
	return result
}
