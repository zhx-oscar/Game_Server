package oflinfos

import (
	"Cinder/Chat/chatapi"
	"Cinder/Chat/rpcproc/logic/types"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	assert "github.com/arl/assertgo"
)

// GetFriendInfos 获取一批ID对应的FriendInfo, 内存查不到就读取DB.
// 如果DB中查不到的，就填nil.
func GetFriendInfos(ids []types.UserID) ([]*chatapi.FriendInfo, error) {
	result := make([]*chatapi.FriendInfo, 0, len(ids))
	for _, id := range ids {
		info, err := GetFriendInfo(id)
		// 测试中会有许多userID没有在DB中，忽略这些ID, 而不是出错停止
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("failed to get friend info of '%v': %w", id, err)
		}
		result = append(result, info)
	}
	return result, nil
}

// GetFriendInfo 获取一个ID对应的FriendInfo, 内存查不到就读取DB.
func GetFriendInfo(id types.UserID) (*chatapi.FriendInfo, error) {
	// 优先查询在线用户
	onlinedInfo := getInfoFromOnlineUsers(id)
	if onlinedInfo != nil {
		return onlinedInfo, nil
	}
	return GetOfldInfoMgr()._Get(id)
}

// getInfoFromOnlineUsers 从在线用户中查找用户信息。
func getInfoFromOnlineUsers(userID types.UserID) *chatapi.FriendInfo {
	assert.True(userMgr != nil)
	return userMgr.GetUserFriendInfo(userID)
}
