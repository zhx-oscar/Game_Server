package user

import (
	"Cinder/Chat/chatapi"
	"Cinder/Chat/rpcproc/logic/user/internal/oflinfos"
)

func GetFriendInfos(userIDs []UserID) ([]*chatapi.FriendInfo, error) {
	return oflinfos.GetFriendInfos(userIDs)
}
