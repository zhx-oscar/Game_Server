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

type _TeamRoomSuite struct {
	suite.Suite
	*require.Assertions

	svc      *RoomService
	coreInst Core.ICore

	teamID mtypes.TeamID
}

func newTeamRoomSuite(t *testing.T) *_TeamRoomSuite {
	return &_TeamRoomSuite{
		Assertions: require.New(t),

		svc:      NewRoomService(),
		coreInst: Core.New(),
	}
}

func (t *_TeamRoomSuite) SetupSuite() {
	info := Core.NewDefaultInfo()
	info.ServiceType = "test"
	svcID := fmt.Sprintf("%s_%s", info.ServiceType, Util.GetGUID())
	t.True(len(svcID) < 64) // nsq topic requires
	info.ServiceID = svcID
	info.RpcProc = test.NewRpcProc(t.Assertions)
	errInit := t.coreInst.Init(info)
	t.NoError(errInit)
}

func (t *_TeamRoomSuite) TearDownSuite() {
	t.coreInst.Destroy()
}

func (t *_TeamRoomSuite) SetupTest() {
	t.Empty(t.teamID)

	teamInfo, errCrt := NewTeamService().CreateTeam(mtypes.RoleInfo{RoleID: "role1"}, mtypes.TeamInfo{MaxRole: 2})
	t.NoError(errCrt)
	t.Equal(1, len(teamInfo.Roles))
	t.Contains(teamInfo.Roles, mtypes.RoleID("role1"))

	t.teamID = teamInfo.TeamID
	t.NotEmpty(t.teamID)
}

func (t *_TeamRoomSuite) TearDownTest() {
	errLeave := NewTeamService().LeaveTeam("role1", t.teamID)
	t.NoError(errLeave)

	t.teamID = ""
}

func (t *_TeamRoomSuite) TestTeamJoinRandomRoom_CreatNew() {
	roomInfo, errJoin := t.svc.TeamJoinRandomRoom(t.teamID, "2v2")
	t.NoError(errJoin)
	t.Equal(1, len(roomInfo.Roles))
	t.Equal(t.teamID, roomInfo.Roles["role1"].TeamID)

	errDel := t.svc.DeleteRoom(roomInfo.RoomID)
	t.NoError(errDel)
}

func (t *_TeamRoomSuite) TestTeamJoinRandomRoom_ExistRoom() {
	roomInfo0, errCrt := t.svc.RoleCreateRoom(mtypes.RoleInfo{RoleID: "role0"}, "2v2")
	t.NoError(errCrt)

	roomInfo, errJoin := t.svc.TeamJoinRandomRoom(t.teamID, "2v2")
	t.NoError(errJoin)

	t.Equal(roomInfo0.RoomID, roomInfo.RoomID)
	t.Equal(2, len(roomInfo.Roles))
	t.Equal(t.teamID, roomInfo.Roles["role1"].TeamID)
	t.Empty(roomInfo.Roles["role0"].TeamID)

	errDel := t.svc.DeleteRoom(roomInfo.RoomID)
	t.NoError(errDel)
}

func (t *_TeamRoomSuite) TestTeamCreateRoom() {
	roomInfo, errJoin := t.svc.TeamCreateRoom(t.teamID, "2v2")
	t.NoError(errJoin)
	t.Equal(1, len(roomInfo.Roles))
	t.Equal(t.teamID, roomInfo.Roles["role1"].TeamID)

	errDel := t.svc.DeleteRoom(roomInfo.RoomID)
	t.NoError(errDel)
}

func (t *_TeamRoomSuite) TestTeamJoinRoom() {
	roomInfo0, errCrt := t.svc.RoleCreateRoom(mtypes.RoleInfo{RoleID: "role0"}, "2v2")
	t.NoError(errCrt)

	roomInfo, errJoin := t.svc.TeamJoinRoom(t.teamID, roomInfo0.RoomID)
	t.NoError(errJoin)

	t.Equal(roomInfo0.RoomID, roomInfo.RoomID)
	t.Equal(2, len(roomInfo.Roles))
	t.Equal(t.teamID, roomInfo.Roles["role1"].TeamID)
	t.Empty(roomInfo.Roles["role0"].TeamID)

	errDel := t.svc.DeleteRoom(roomInfo.RoomID)
	t.NoError(errDel)
}

func (t *_TeamRoomSuite) TestLeaveRoom() {
	roomInfo, errCrt := t.svc.TeamCreateRoom(t.teamID, "2v2")
	t.NoError(errCrt)
	roomInfo2, errLeave := t.svc.RoleLeaveRoom("role1", roomInfo.RoomID)
	t.NoError(errLeave)
	t.NotContains(roomInfo2.Roles, "role1")

	errDel := t.svc.DeleteRoom(roomInfo.RoomID)
	t.NoError(errDel)
}

// 让 go test 执行测试
func TestTeamRoomService(t *testing.T) {
	suite.Run(t, newTeamRoomSuite(t))
}
