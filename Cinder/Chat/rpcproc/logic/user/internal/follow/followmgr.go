package follow

import (
	"Cinder/Chat/chatapi"
	"Cinder/Chat/rpcproc/logic/types"
	"Cinder/Chat/rpcproc/logic/user/internal/follow/follower"
	"Cinder/Chat/rpcproc/logic/user/internal/follow/following"
	"Cinder/Chat/rpcproc/logic/user/internal/oflinfos"
	"errors"
	"fmt"
)

type UserID = types.UserID

type FollowMgr struct {
	userID UserID

	// 成员协程安全
	followings *following.Followings // 关注列表
	followers  *follower.Followers   // 粉丝列表
}

func NewFollowMgr(userID UserID, userInfo follower.IUserInfo) *FollowMgr {
	return &FollowMgr{
		userID: userID,

		followings: following.NewFollowings(userID),
		followers:  follower.NewFollowers(userID, userInfo),
	}
}

// Follow 关注某人, 立即写DB.
func (f *FollowMgr) Follow(userID UserID) error {
	if userID == f.userID {
		return nil // 自关注会引起 mutex 重入死锁
	}

	// DB和内存添加关注
	if err := f.followings.Add(userID); err != nil {
		if errors.Is(err, following.ErrNotAllowedToAdd) {
			return nil
		}
		return err
	}

	// 对方DB添加粉丝
	if err := follower.AddFollowerInDB(userID, f.userID); err != nil {
		return err // TODO: 有数据不一致
	}

	// 对方（被关注者）内存添加粉丝。如果不在线就不能加了。
	if followMgr := userMgr.GetUserFollowMgr(userID); followMgr != nil {
		followMgr.AddFollower(f.userID) // 仅内存添加
	}
	return nil
}

// Unfollow 取消关注某人, 立即写DB.
func (f *FollowMgr) Unfollow(userID UserID) error {
	// DB和内存取消关注
	if err := f.followings.Delete(userID); err != nil {
		if errors.Is(err, following.ErrFollowingNotExists) {
			return nil // 不存在则忽略请求
		}
		return err
	}

	// 对方DB删粉丝
	if err := follower.DeleteFollowerInDB(userID, f.userID); err != nil {
		return err // TODO: 有数据不一致
	}

	// 对方（被关注者）内存删粉丝。如果不在线就不用。
	if followMgr := userMgr.GetUserFollowMgr(userID); followMgr != nil {
		followMgr.DeleteFollower(f.userID) // 仅内存删除
	}
	return nil
}

// AddFollwer 添加粉丝，不写DB.
// 因为粉丝加我的时候，我可能不在线，所以写DB是在User类外处理的。
func (f *FollowMgr) AddFollower(userID UserID) {
	f.followers.Add(userID)
}

// DeleteFollwer 删除粉丝，不写DB.
// 因为粉丝取消关注我的时候，我可能不在线，所以写DB是在User类外处理的。
func (f *FollowMgr) DeleteFollower(userID UserID) {
	f.followers.Delete(userID)
}

// LoadFollowers 加载粉丝。
// 关注会在User加载时初始化。
func (f *FollowMgr) LoadFollowers() error {
	// chat.users.followers
	if err := f.followers.Load(); err != nil {
		return fmt.Errorf("foad followers: %w", err)
	}
	return nil
}

// SetFollowings 设置关注列表
func (f *FollowMgr) SetFollowings(ids []UserID) {
	f.followings.Set(ids)
}

// GetFollowingList 获取关注列表
func (f *FollowMgr) GetFollowingList() ([]*chatapi.FriendInfo, error) {
	ids := f.followings.GetIDs()
	return oflinfos.GetFriendInfos(ids)
}

// GetFollowerList 获取粉丝列表
func (f *FollowMgr) GetFollowerList() ([]*chatapi.FollowerInfo, error) {
	m := f.followers.GetFollowers()
	result := make([]*chatapi.FollowerInfo, 0, len(m))
	for id, time := range m {
		info, err := oflinfos.GetFriendInfo(id)
		if err != nil {
			return nil, fmt.Errorf("get friend info of '%v': %v", id, err)
		}
		result = append(result, &chatapi.FollowerInfo{
			FriendInfo: info,
			FollowTime: time,
		})
	}
	return result, nil
}
