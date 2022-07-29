package stress

import (
	"Cinder/Chat/chatapi"
)

func init() {
	register(addFriendReqApply)
	register(addFriendToBlacklistAndRemove)
	register(deleteFriend)
	register(getFriendList)
	register(getFriendBlacklist)
}

func addFriendReqApply(iGo GoroutineIndex, i _RunIndex) {
	type FromTo struct {
		From chatapi.IUser
		To   chatapi.IUser
	}
	fromRoleID := getRoleID(iGo, 0)
	toRoleID := fromRoleID + "_to"
	if i == 0 {
		userFrom, errFrom := chatapi.Login(fromRoleID, "nickFrom", []byte("dataFrom"))
		panicIfError(errFrom)
		userTo, errTo := chatapi.Login(toRoleID, "nickTo", []byte("dataTo"))
		panicIfError(errTo)
		goroutineMap.Store(iGo, &FromTo{From: userFrom, To: userTo})
	}

	iData, ok := goroutineMap.Load(iGo)
	if !ok {
		flushAndPanic("can not get 2 users")
	}
	fromTo := iData.(*FromTo)
	errAdd := fromTo.From.AddFriendReq(toRoleID, []byte("reqInfo"))
	panicIfError(errAdd)
	errApply := fromTo.To.ApplyAddFriendReq(fromRoleID, i%2 == 0)
	panicIfError(errApply)
}

func addFriendToBlacklistAndRemove(iGo GoroutineIndex, i _RunIndex) {
	user := firstLoginUser(iGo, i)
	errAdd := user.AddFriendToBlacklist("friend")
	panicIfError(errAdd)
	errRm := user.RemoveFriendFromBlacklist("friend")
	panicIfError(errRm)
}

func deleteFriend(iGo GoroutineIndex, i _RunIndex) {
	user := firstLoginUser(iGo, i)
	user.DeleteFriend("friendID")
}

func getFriendList(iGo GoroutineIndex, i _RunIndex) {
	user := firstLoginUser(iGo, i)
	if user.GetFriendList() == nil {
		flushAndPanic("GetFriendList() returns nil")
	}
}

func getFriendBlacklist(iGo GoroutineIndex, i _RunIndex) {
	user := firstLoginUser(iGo, i)
	if user.GetFriendBlacklist() == nil {
		flushAndPanic("GetFriendBlacklist() returns nil")
	}
}
