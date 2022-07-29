package stress

import (
	"Cinder/Base/Util"
	"Cinder/Matcher/matchapi"
	"Cinder/Matcher/matchapi/mtypes"
	"strconv"
	"sync"
)

var (
	rs = matchapi.GetRoomService()
	ts = matchapi.GetTeamService()

	rolePrefix   = Util.GetGUID() + "_"
	goroutineMap sync.Map // [int]interface{}
)

func getRoleSuffix(iGo GoroutineIndex, i RunIndex) string {
	return strconv.Itoa(int(iGo)) + "_" + strconv.Itoa(int(i))
}

func getRoleID(iGo GoroutineIndex, i RunIndex) mtypes.RoleID {
	return mtypes.RoleID(rolePrefix + getRoleSuffix(iGo, i))
}

func getRoleInfo(iGo GoroutineIndex, i RunIndex) mtypes.RoleInfo {
	return mtypes.RoleInfo{RoleID: getRoleID(iGo, i)}
}

func getTeamInfo(iGo GoroutineIndex, i RunIndex) mtypes.TeamInfo {
	return mtypes.TeamInfo{MaxRole: 10000}
}

func getRoomID(iGo GoroutineIndex) mtypes.RoomID {
	iRoomID, ok := goroutineMap.Load(iGo)
	if !ok {
		panic("can not get room ID")
	}
	return iRoomID.(mtypes.RoomID)
}

func getTeamID(iGo GoroutineIndex) mtypes.TeamID {
	iTeamID, ok := goroutineMap.Load(iGo)
	if !ok {
		panic("can not get team ID")
	}
	return iTeamID.(mtypes.TeamID)
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
