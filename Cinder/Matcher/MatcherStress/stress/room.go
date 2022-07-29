package stress

import (
	"Cinder/Base/Const"
	"Cinder/Base/Core"
	"Cinder/Matcher/matchapi/mtypes"
)

func RoleJoinRandomRoom(iGo GoroutineIndex, i RunIndex) {
	rs.RoleJoinRandomRoom(getRoleInfo(iGo, i), "2v2")
}

func MyRPCTest(iGo GoroutineIndex, i RunIndex) {
	ret := <-Core.Inst.RpcByType(Const.Matcher, "RPC_MyRPCTest", "test")
	panicIfError(ret.Err)
}

func RoleCreateDeleteRoom(iGo GoroutineIndex, i RunIndex) {
	roomInfo, errCrt := rs.RoleCreateRoom(getRoleInfo(iGo, i), "1v1")
	panicIfError(errCrt)
	rs.DeleteRoom(roomInfo.RoomID)
}

func GetRoomList(iGo GoroutineIndex, i RunIndex) {
	rs.GetRoomList("no_such_match_mode")
}

func Ping(iGo GoroutineIndex, i RunIndex) {
	ret := <-Core.Inst.RpcByType(Const.Matcher, "RPC_Ping")
	panicIfError(ret.Err)
}

func RoleJoinRoom(iGo GoroutineIndex, i RunIndex) {
	roleInfo := getRoleInfo(iGo, i)
	if 0 == i%4 {
		roomInfo, errCrt := rs.RoleCreateRoom(roleInfo, "2v2")
		panicIfError(errCrt)
		goroutineMap.Store(iGo, roomInfo.RoomID)
		return
	}

	roomID := getRoomID(iGo)
	if _, err := rs.RoleJoinRoom(roleInfo, roomID); err != nil {
		panic(err)
	}
}

func RoleJoinLeaveRoom(iGo GoroutineIndex, i RunIndex) {
	firstCreateRoom(iGo, i)
	roomID := getRoomID(iGo)
	if i == 0 {
		// 多加一个成员
		roleID := mtypes.RoleID(rolePrefix + getRoleSuffix(iGo, i) + "_joined")
		_, errJoin := rs.RoleJoinRoom(mtypes.RoleInfo{RoleID: roleID}, roomID)
		panicIfError(errJoin)
	}

	roleInfo := getRoleInfo(iGo, i)
	if _, errJoin := rs.RoleJoinRoom(roleInfo, roomID); errJoin != nil {
		panic(errJoin)
	}
	if _, errLeave := rs.RoleLeaveRoom(roleInfo.RoleID, roomID); errLeave != nil {
		panic(errLeave)
	}
}

func BroadcastRoom(iGo GoroutineIndex, i RunIndex) {
	firstCreateRoom(iGo, i)
	roomID := getRoomID(iGo)
	if i == 0 {
		sRoleID := rolePrefix + getRoleSuffix(iGo, i)
		if _, err1 := rs.RoleJoinRoom(mtypes.RoleInfo{RoleID: mtypes.RoleID(sRoleID + "_join1")}, roomID); err1 != nil {
			panic(err1)
		}
		if _, err2 := rs.RoleJoinRoom(mtypes.RoleInfo{RoleID: mtypes.RoleID(sRoleID + "_join2")}, roomID); err2 != nil {
			panic(err2)
		}
	}

	err := rs.BroadcastRoom(roomID, "test msg")
	panicIfError(err)
}

func GetRoomInfo(iGo GoroutineIndex, i RunIndex) {
	firstCreateRoom(iGo, i)

	roomID := getRoomID(iGo)
	_, err := rs.GetRoomInfo(roomID)
	panicIfError(err)
}

func SetRoomData(iGo GoroutineIndex, i RunIndex) {
	firstCreateRoom(iGo, i)
	roomID := getRoomID(iGo)
	err := rs.SetRoomData(roomID, "testKey", "testData")
	panicIfError(err)
}

func SetRoomRoleData(iGo GoroutineIndex, i RunIndex) {
	firstCreateRoom(iGo, i)
	roomID := getRoomID(iGo)
	roleID := getRoleID(iGo, 0)

	var err error
	switch i % 4 {
	case 0:
		err = rs.SetRoomRoleFloatData(roomID, roleID, "testFloatKey", 1.234)
	case 1:
		err = rs.SetRoomRoleStringData(roomID, roleID, "testStringKey", "test")
	case 2:
		err = rs.AddRoomRoleTag(roomID, roleID, "testTag")
	case 3:
		err = rs.DelRoomRoleTag(roomID, roleID, "testTag")
	}
	panicIfError(err)
}

func firstCreateRoom(iGo GoroutineIndex, i RunIndex) {
	if i != 0 {
		return
	}

	roomInfo, errCrt := rs.RoleCreateRoom(getRoleInfo(iGo, i), "2v2")
	panicIfError(errCrt)
	goroutineMap.Store(iGo, roomInfo.RoomID)
}
