package internal

import (
	"Cinder/Base/Core"
	"Cinder/Base/Util"
	"Cinder/Matcher/matchapi/internal/test"
	"Cinder/Matcher/matchapi/mtypes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type _RoomSuite struct {
	suite.Suite
	*require.Assertions

	svc      *RoomService
	coreInst Core.ICore

	roomInfo mtypes.RoomInfo
}

func newRoomSuite(t *testing.T) *_RoomSuite {
	return &_RoomSuite{
		Assertions: require.New(t),

		svc:      NewRoomService(),
		coreInst: Core.New(),
	}
}

func (r *_RoomSuite) SetupSuite() {
	info := Core.NewDefaultInfo()
	info.ServiceType = "test"
	svcID := fmt.Sprintf("%s_%s", info.ServiceType, Util.GetGUID())
	r.True(len(svcID) < 64) // nsq topic requires
	info.ServiceID = svcID
	info.RpcProc = test.NewRpcProc(r.Assertions)
	errInit := r.coreInst.Init(info)
	r.NoError(errInit)
}

func (r *_RoomSuite) TearDownSuite() {
	r.coreInst.Destroy()
}

func (r *_RoomSuite) SetupTest() {
	r.Empty(r.roomInfo)

	var err error
	r.roomInfo, err = r.svc.RoleJoinRandomRoom(mtypes.RoleInfo{RoleID: "testRole"}, "2v2")
	r.NoError(err)
	r.Equal(1, len(r.roomInfo.Roles))
	r.Contains(r.roomInfo.Roles, mtypes.RoleID("testRole"))
}

func (r *_RoomSuite) TearDownTest() {
	errDel := r.svc.DeleteRoom(r.roomInfo.RoomID)
	r.NoError(errDel)
	r.roomInfo = mtypes.RoomInfo{}
}

func (r *_RoomSuite) TestRoleJoinRandomRoom() {
	roomInfo2, errJoin := r.svc.RoleJoinRandomRoom(mtypes.RoleInfo{RoleID: "testRole2"}, "1v1")
	r.NoError(errJoin)
	r.NotEqual(r.roomInfo.RoomID, roomInfo2.RoomID)
	errDel := r.svc.DeleteRoom(roomInfo2.RoomID)
	r.NoError(errDel)
}

func (r *_RoomSuite) TestRoleJoinRoom() {
	roomInfo := r.roomInfo
	roomInfo2, errJoin := r.svc.RoleJoinRoom(mtypes.RoleInfo{RoleID: "testRole2"}, roomInfo.RoomID)
	r.NoError(errJoin)
	r.Equal(roomInfo.RoomID, roomInfo2.RoomID)
	r.Equal(2, len(roomInfo2.Roles))
	r.Contains(roomInfo2.Roles, mtypes.RoleID("testRole"))
	r.Contains(roomInfo2.Roles, mtypes.RoleID("testRole2"))
}

func (r *_RoomSuite) TestRoleLeaveRoom() {
	roomInfo := r.roomInfo
	_, errJoin := r.svc.RoleJoinRoom(mtypes.RoleInfo{RoleID: "testRole2"}, roomInfo.RoomID)
	r.NoError(errJoin)
	roomInfo2, errLeave := r.svc.RoleLeaveRoom("testRole2", roomInfo.RoomID)
	r.NoError(errLeave)
	r.NotContains(roomInfo2.Roles, "testRole2")
}

func (r *_RoomSuite) TestBroadcastRoom() {
	roomInfo := r.roomInfo
	_, errJoin := r.svc.RoleJoinRoom(mtypes.RoleInfo{RoleID: "testRole2"}, roomInfo.RoomID)
	r.NoError(errJoin)
	errBc := r.svc.BroadcastRoom(roomInfo.RoomID, "test room broadcast")
	r.NoError(errBc)
}

func (r *_RoomSuite) TestGetRoomList() {
	roomInfo2, err2 := r.svc.RoleCreateRoom(mtypes.RoleInfo{RoleID: "testRole2"}, "1v1")
	r.NoError(err2)

	roomIDs1, errList1 := r.svc.GetRoomList("2v2")
	r.NoError(errList1)
	r.Contains(roomIDs1, r.roomInfo.RoomID)
	r.NotContains(roomIDs1, roomInfo2.RoomID)

	roomIDs2, errList2 := r.svc.GetRoomList("1v1")
	r.NoError(errList2)
	r.Contains(roomIDs2, roomInfo2.RoomID)
	r.NotContains(roomIDs2, r.roomInfo.RoomID)

	errDel2 := r.svc.DeleteRoom(roomInfo2.RoomID)
	r.NoError(errDel2)
}

func (r *_RoomSuite) TestGetRoomInfo() {
	_, errGet := r.svc.GetRoomInfo("NoSuchRoom")
	r.Error(errGet)
	_, errGet2 := r.svc.GetRoomInfo(r.roomInfo.RoomID)
	r.NoError(errGet2)
}

func (r *_RoomSuite) TestSetRoomData() {
	roomInfo := r.roomInfo
	errSet1 := r.svc.SetRoomData(roomInfo.RoomID, "testKey1", "TestData")
	r.NoError(errSet1)
	errSet2 := r.svc.SetRoomData(roomInfo.RoomID, "testKey2", 12345)
	r.NoError(errSet2)

	roomInfo2, errGet := r.svc.GetRoomInfo(roomInfo.RoomID)
	r.NoError(errGet)
	val1, ok1 := roomInfo2.GetData("testKey1")
	r.True(ok1)
	r.EqualValues("TestData", val1)
	val2, ok2 := roomInfo2.GetData("testKey2")
	r.True(ok2)
	r.EqualValues(12345, val2)
}

func (r *_RoomSuite) TestUpdateRoomRoleData() {
	roomInfo := r.roomInfo
	errFloat := r.svc.SetRoomRoleFloatData(roomInfo.RoomID, "testRole", "testFloatKey", 123)
	r.NoError(errFloat)
	errString := r.svc.SetRoomRoleStringData(roomInfo.RoomID, "testRole", "testStringKey", "aaaa")
	r.NoError(errString)
	errAddTag := r.svc.AddRoomRoleTag(roomInfo.RoomID, "testRole", "addTag")
	r.NoError(errAddTag)

	roomInfo2, errGet := r.svc.GetRoomInfo(roomInfo.RoomID)
	r.NoError(errGet)
	roleInfo := roomInfo2.Roles["testRole"]
	f, okF := roleInfo.GetFloatData("testFloatKey")
	r.True(okF)
	r.Equal(float64(123), f)
	s, okS := roleInfo.GetStringData("testStringKey")
	r.True(okS)
	r.Equal("aaaa", s)
	r.True(roleInfo.HasTag("addTag"))

	errDelTag := r.svc.DelRoomRoleTag(roomInfo.RoomID, "testRole", "addTag")
	r.NoError(errDelTag)
	roomInfo3, errGet3 := r.svc.GetRoomInfo(roomInfo.RoomID)
	r.NoError(errGet3)
	roleInfo3, okRole := roomInfo3.Roles["testRole"]
	r.True(okRole)
	r.False(roleInfo3.HasTag("addTag"))
}

// 让 go test 执行测试
func TestRoomService(t *testing.T) {
	suite.Run(t, newRoomSuite(t))
}
