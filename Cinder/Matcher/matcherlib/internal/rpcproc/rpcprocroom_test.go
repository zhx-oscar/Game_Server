package rpcproc

import (
	"Cinder/Base/Core"
	"Cinder/Base/Util"
	"Cinder/Matcher/matchapi/mtypes"
	"Cinder/Matcher/matcherlib/internal/rpcproc/room"
	"Cinder/Matcher/rpcmsg"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
)

type _RoomEventHandler struct {
}

func (r *_RoomEventHandler) OnCreate(roomInfo *mtypes.RoomInfo) error {
	return nil
}

func (r *_RoomEventHandler) OnDestroy(roomInfo *mtypes.RoomInfo) {
}

func (r *_RoomEventHandler) OnAddingRoles(roomInfo *mtypes.RoomInfo, roles mtypes.RoleMap) bool {
	return len(roomInfo.Roles)+len(roles) <= 4
}

func (r *_RoomEventHandler) OnAddedRoles(roomInfo *mtypes.RoomInfo, roles mtypes.RoleMap) {
	roomInfo.IsFull = (len(roomInfo.Roles) >= 4)
	roomInfo.SetData("roleCount", len(roomInfo.Roles)) // 任意更新数据，用于后续匹配，或通知房间成员
	if roomInfo.IsFull {
		roomInfo.IsDeleting = true
	}
}

func (r *_RoomEventHandler) OnDeletingRole(roomInfo *mtypes.RoomInfo, roleID mtypes.RoleID) {
}

func InitServer() {
	if Core.Inst != nil {
		return
	}

	_ = Core.New()
	info := Core.NewDefaultInfo()
	info.ServiceType = "test"
	svcID := fmt.Sprintf("%s_%s", info.ServiceType, Util.GetGUID())
	info.ServiceID = svcID
	if err := Core.Inst.Init(info); err != nil {
		panic(err)
	}
}

func BenchmarkRPC_RoleJoinRandomRoom(b *testing.B) {
	InitServer()
	room.SetRoomEventHandler(&_RoomEventHandler{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rpcRoleJoinRandomRoom(i)
	}
}

func rpcRoleJoinRandomRoom(i int) {
	req := rpcmsg.RoleJoinRandomRoomReq{
		RoleInfo: mtypes.RoleInfo{
			RoleID: mtypes.RoleID(strconv.Itoa(i)),
			SrvID:  mtypes.SrvID(Core.Inst.GetServiceID()),
		},
		MatchMode: "2v2",
	}
	reqJson, errJson := json.Marshal(req)
	if errJson != nil {
		panic(errJson)
	}
	var proc _RPCProcRoom
	proc.RPC_RoleJoinRandomRoom(reqJson)
}
