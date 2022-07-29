package stress

import (
	"Cinder/Matcher/matchapi/mtypes"
)

func CreateLeaveTeam(iGo GoroutineIndex, i RunIndex) {
	roleInfo := getRoleInfo(iGo, i)
	teamInfo, errCrt := ts.CreateTeam(roleInfo, getTeamInfo(iGo, i))
	panicIfError(errCrt)
	errLv := ts.LeaveTeam(roleInfo.RoleID, teamInfo.TeamID)
	panicIfError(errLv)
}

func JoinLeaveTeam(iGo GoroutineIndex, i RunIndex) {
	firstCreateTeam(iGo, i)
	if i == 0 {
		return
	}
	roleInfo := getRoleInfo(iGo, i)
	teamID := getTeamID(iGo)
	_, errJoin := ts.JoinTeam(roleInfo, teamID, "")
	panicIfError(errJoin)
	errLeave := ts.LeaveTeam(roleInfo.RoleID, teamID)
	panicIfError(errLeave)
}

func ChangeTeamLeader(iGo GoroutineIndex, i RunIndex) {
	firstCreateTeam(iGo, i)
	teamID := getTeamID(iGo)
	if i == 0 {
		// 加入2人
		roleInfo0 := getRoleInfo(iGo, 10)
		_, err0 := ts.JoinTeam(roleInfo0, teamID, "")
		panicIfError(err0)
		roleInfo1 := getRoleInfo(iGo, 11)
		_, err1 := ts.JoinTeam(roleInfo1, teamID, "")
		panicIfError(err1)
	}

	roleID := getRoleID(iGo, 10+i%2)
	errChg := ts.ChangeTeamLeader(teamID, roleID)
	panicIfError(errChg)
}

func SetTeamData(iGo GoroutineIndex, i RunIndex) {
	firstCreateTeamAndAddMembers(iGo, i)
	teamID := getTeamID(iGo)
	err := ts.SetTeamData(teamID, "testKey", "testData")
	panicIfError(err)
}

func BroadcastTeam(iGo GoroutineIndex, i RunIndex) {
	firstCreateTeamAndAddMembers(iGo, i)
	teamID := getTeamID(iGo)
	err := ts.BroadcastTeam(teamID, "testmsg")
	panicIfError(err)
}

func firstCreateTeam(iGo GoroutineIndex, i RunIndex) {
	if i != 0 {
		return
	}

	roleInfo := getRoleInfo(iGo, i)
	teamInfo, errCrt := ts.CreateTeam(roleInfo, getTeamInfo(iGo, i))
	panicIfError(errCrt)
	goroutineMap.Store(iGo, teamInfo.TeamID)
}

func firstCreateTeamAndAddMembers(iGo GoroutineIndex, i RunIndex) {
	if i != 0 {
		return
	}
	firstCreateTeam(iGo, i)

	// 加入2人
	teamID := getTeamID(iGo)
	prefix := rolePrefix + getRoleSuffix(iGo, i) + "_"
	_, err0 := ts.JoinTeam(mtypes.RoleInfo{RoleID: mtypes.RoleID(prefix + "0")}, teamID, "")
	panicIfError(err0)
	_, err1 := ts.JoinTeam(mtypes.RoleInfo{RoleID: mtypes.RoleID(prefix + "1")}, teamID, "")
	panicIfError(err1)
}
