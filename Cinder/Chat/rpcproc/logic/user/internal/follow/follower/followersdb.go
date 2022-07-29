package follower

import (
	"Cinder/Chat/rpcproc/logic/user/internal/follow/follower/dbutil"
)

// 粉丝写DB与Followers是分离的，因为 Followers 仅对在线用户实例化。

// AddFollowerDB DB添加粉丝，followerID 关注 userID 成为其粉丝。
func AddFollowerInDB(userID, followerID UserID) error {
	if err := dbutil.UsersFollowersUtil(userID).Add(followerID); err != nil {
		return err
	}
	_ = dbutil.UsersUtil(userID).IncreaseFollowerNumber()
	// TODO: 如果出错，粉丝数与实际不符，需要上线时检查并修复
	return nil
}

// DeleteFollowerInDB DB删除 userID 的粉丝 followerID。
func DeleteFollowerInDB(userID, followerID UserID) error {
	if err := dbutil.UsersFollowersUtil(userID).Remove(followerID); err != nil {
		return err
	}
	_ = dbutil.UsersUtil(userID).DecreaseFollowerNumber()
	// TODO: 如果出错，粉丝数与实际不符，需要上线时检查并修复
	return nil
}
