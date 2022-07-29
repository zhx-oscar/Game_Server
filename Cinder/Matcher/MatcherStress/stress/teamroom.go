package stress

func TeamCreateDeleteRoom(iGo GoroutineIndex, i RunIndex) {
	firstCreateTeam(iGo, i)
	teamID := getTeamID(iGo)

	roomInfo, errCrt := rs.TeamCreateRoom(teamID, "2v2")
	panicIfError(errCrt)
	errDel := rs.DeleteRoom(roomInfo.RoomID)
	panicIfError(errDel)
}

func TeamJoinRoom(iGo GoroutineIndex, i RunIndex) {
	firstCreateTeam(iGo, i)
	teamID := getTeamID(iGo)

	// 创建房间
	roomInfo, errCrt := rs.RoleCreateRoom(getRoleInfo(iGo, i), "2v2")
	panicIfError(errCrt)
	// 加入房间
	_, errJoin := rs.TeamJoinRoom(teamID, roomInfo.RoomID)
	panicIfError(errJoin)
	// 删除房间
	errDel := rs.DeleteRoom(roomInfo.RoomID)
	panicIfError(errDel)
}

func TeamJoinRandomRoom(iGo GoroutineIndex, i RunIndex) {
	firstCreateTeam(iGo, i)
	teamID := getTeamID(iGo)
	RoomInfo, errJoin := rs.TeamJoinRandomRoom(teamID, "2v2")
	panicIfError(errJoin)
	errDel := rs.DeleteRoom(RoomInfo.RoomID)
	panicIfError(errDel)
}
