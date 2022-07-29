package stress

import (
	"Cinder/Chat/chatapi"
)

func init() {
	register(createDeleteGroup)
	register(createSameGroup)
}

func createDeleteGroup(iGo GoroutineIndex, i _RunIndex) {
	groupID := getGroupID(iGo, i)
	errCrt := chatapi.CreateGroup(groupID, nil)
	panicIfError(errCrt)
	errDel := chatapi.DeleteGroup(groupID)
	panicIfError(errDel)
}

func createSameGroup(iGo GoroutineIndex, _ _RunIndex) {
	groupID := getGroupID(iGo, 0)
	errCrt := chatapi.CreateGroup(groupID, nil)
	panicIfError(errCrt)
}
