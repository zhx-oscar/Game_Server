package stress

import (
	"Cinder/Chat/chatapi"
)

func init() {
	register(getGroupMemberCounts)
}

func getGroupMemberCounts(iGo GoroutineIndex, i _RunIndex) {
	// 创建几个 group
	if i <= 10 {
		groupID := getGroupID(iGo, i)
		errCrt := chatapi.CreateGroup(groupID, nil)
		panicIfError(errCrt)
	}
	groupIDs := make([]string, 0, 10)
	for j := _RunIndex(0); j < 10; j++ {
		groupIDs = append(groupIDs, getGroupID(iGo, j))
	}
	if counts := chatapi.GetGroupMemberCounts(groupIDs); counts == nil {
		flushAndPanic("failed to get group member counts")
	}
}
