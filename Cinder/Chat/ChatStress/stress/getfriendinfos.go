package stress

import (
	"Cinder/Chat/chatapi"
)

func init() {
	register(getFriendInfos)
	register(getFriendInfosNil)
}

func getFriendInfos(iGo GoroutineIndex, i _RunIndex) {
	roleIDs := make([]string, 0, 10)
	for j := _RunIndex(0); j < 10; j++ {
		roleIDs = append(roleIDs, getRoleID(iGo, j))
	}
	if infos := chatapi.GetFriendInfos(roleIDs); infos == nil {
		flushAndPanic("failed to get friend infos")
	}
}

func getFriendInfosNil(iGo GoroutineIndex, i _RunIndex) {
	if infos := chatapi.GetFriendInfos(nil); infos == nil {
		flushAndPanic("failed to get empty friend infos")
	}
}
