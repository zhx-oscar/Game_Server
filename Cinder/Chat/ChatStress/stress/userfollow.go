package stress

func init() {
	register(followUnfollowFriendReq)
	register(getFollowingList)
	register(getFollowerList)
}

func followUnfollowFriendReq(iGo GoroutineIndex, i _RunIndex) {
	user := firstLoginUser(iGo, i)
	user.FollowFriendReq("friend")
	user.UnFollowFriendReq("friend")
}

func getFollowingList(iGo GoroutineIndex, i _RunIndex) {
	user := firstLoginUser(iGo, i)
	user.GetFollowingList()
}

func getFollowerList(iGo GoroutineIndex, i _RunIndex) {
	user := firstLoginUser(iGo, i)
	user.GetFollowerList()
}
