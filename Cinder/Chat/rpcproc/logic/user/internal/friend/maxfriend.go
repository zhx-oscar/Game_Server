package friend

import (
	"Cinder/Chat/chatapi"
	"Cinder/Chat/rpcproc/logic/user/internal/friend/dbutil"

	assert "github.com/arl/assertgo"
	"github.com/spf13/viper"
)

// hasReachedMaxFriendCount 是否已达最大好友数
func hasReachedMaxFriendCount(friendCount int) bool {
	if friendCount <= 0 {
		return false
	}

	max := viper.GetInt("Chat.MaxFriendCount")
	if max < 0 {
		return false
	}
	return friendCount >= max
}

// checkMaxFriendCount 检查自身和对方是否已达最大好友数，没有则返回nil, 否则返回相应错误
// 注意，本函数在 selfID 的 UserMgr 中执行，有锁不可以重入。
func checkMaxFriendCount(selfID UserID, selfFriendCount int, peerID UserID) error {
	assert.True(selfID != peerID)
	max := viper.GetInt("Chat.MaxFriendCount")
	if max < 0 {
		return nil
	}
	if selfFriendCount >= max {
		return chatapi.ErrSelfReachedMaxFriendCount // 自身已达最大好友数
	}

	// 检查对方好友数
	if peerFrndMgr := userMgr.GetUserFriendMgr(peerID); peerFrndMgr != nil {
		if peerFrndMgr.GetFriendCount() >= max {
			return chatapi.ErrPeerReachedMaxFriendCount // 对方已达最大好友数
		}
		return nil
	}
	// 对方不在线，从DB查
	if cnt, err := getFriendCountInDb(peerID); err != nil {
		return err
	} else if cnt >= max {
		return chatapi.ErrPeerReachedMaxFriendCount // 对方已达最大好友数
	}
	return nil
}

func getFriendCountInDb(userID UserID) (count int, e error) {
	return dbutil.UsersUtil(userID).GetFriendCount()
}

// checkPeerMaxAddFriendReq 检查对方是否到达最大被申请数，返回相应错误。
func checkPeerMaxAddFriendReq(peerID UserID) error {
	max := viper.GetInt("Chat.MaxAddFriendReq")
	if max < 0 {
		return nil
	}
	if cnt, err := dbutil.FriendRequestsUtil(peerID).GetCount(); err != nil {
		return err
	} else if cnt >= int64(max) {
		return chatapi.ErrPeerReachedMaxAddFriendReq // 对方已达最大请求数
	}
	return nil
}
