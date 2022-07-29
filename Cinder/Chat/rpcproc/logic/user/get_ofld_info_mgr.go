package user

import (
	"Cinder/Chat/chatapi"
	"Cinder/Chat/rpcproc/logic/user/internal/oflinfos"
)

type IOfflinedInfoMgr interface {
	Remove(userID UserID)
	Add(info chatapi.FriendInfo)
}

func GetOfldInfoMgr() IOfflinedInfoMgr {
	return oflinfos.GetOfldInfoMgr()
}
