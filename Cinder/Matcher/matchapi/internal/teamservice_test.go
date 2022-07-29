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

type _TeamSuite struct {
	suite.Suite
	*require.Assertions

	svc      *TeamService
	coreInst Core.ICore

	teamID mtypes.TeamID
}

func newTeamSuite(t *testing.T) *_TeamSuite {
	assert := require.New(t)
	return &_TeamSuite{
		Assertions: assert,

		svc:      NewTeamService(),
		coreInst: Core.New(),
	}
}

func (t *_TeamSuite) SetupSuite() {
	info := Core.NewDefaultInfo()
	info.ServiceType = "test"
	svcID := fmt.Sprintf("%s_%s", info.ServiceType, Util.GetGUID())
	t.True(len(svcID) < 64) // nsq topic requires
	info.ServiceID = svcID
	info.RpcProc = test.NewRpcProc(t.Assertions)
	errInit := t.coreInst.Init(info)
	t.NoError(errInit)
}

func (t *_TeamSuite) TearDownSuite() {
	t.coreInst.Destroy()
}

func (t *_TeamSuite) SetupTest() {
	t.Empty(t.teamID)

	roleInfo := mtypes.RoleInfo{RoleID: "testRole"}
	teamInfo := mtypes.TeamInfo{
		MaxRole: 2,
	}
	teamInfo2, errCrt := t.svc.CreateTeam(roleInfo, teamInfo)
	t.NoError(errCrt)
	t.Contains(teamInfo2.Roles, mtypes.RoleID("testRole"))
	t.Equal(mtypes.RoleID("testRole"), teamInfo2.LeaderID)

	t.teamID = teamInfo2.TeamID
}

func (t *_TeamSuite) TearDownTest() {
	err := t.svc.LeaveTeam("testRole", t.teamID)
	t.NoError(err)
	t.teamID = ""
}

func (t *_TeamSuite) TestJoinTeam() {
	info, errJoin := t.svc.JoinTeam(mtypes.RoleInfo{RoleID: "role2"}, t.teamID, "" /*passwd*/)
	t.NoError(errJoin)
	t.Contains(info.Roles, mtypes.RoleID("role2"))
	t.Equal(2, len(info.Roles))

	_, errFull := t.svc.JoinTeam(mtypes.RoleInfo{RoleID: "role3"}, t.teamID, "" /*passwd*/)
	t.Error(errFull)

	err := t.svc.LeaveTeam("role2", t.teamID)
	t.NoError(err)
}

func (t *_TeamSuite) TestLeaveTeam() {
	info, errJoin := t.svc.JoinTeam(mtypes.RoleInfo{RoleID: "role2"}, t.teamID, "" /*passwd*/)
	t.NoError(errJoin)
	t.Contains(info.Roles, mtypes.RoleID("role2"))
	errLeave := t.svc.LeaveTeam("role2", t.teamID)
	t.NoError(errLeave)
	info3, err3 := t.svc.JoinTeam(mtypes.RoleInfo{RoleID: "role3"}, t.teamID, "")
	t.NoError(err3)
	t.Contains(info3.Roles, mtypes.RoleID("role3"))
	t.NotContains(info3.Roles, mtypes.RoleID("role2"))
	errLeave3 := t.svc.LeaveTeam("role3", t.teamID)
	t.NoError(errLeave3)
}

func (t *_TeamSuite) TestTeamPassword() {
	_, errJoin := t.svc.JoinTeam(mtypes.RoleInfo{RoleID: "role2"}, t.teamID, "wrongPassword")
	t.Error(errJoin)
}

func (t *_TeamSuite) TestSetTeamData() {
	err := t.svc.SetTeamData(t.teamID, "key", 123456)
	t.NoError(err)
	info2, err2 := t.svc.JoinTeam(mtypes.RoleInfo{RoleID: "role2"}, t.teamID, "")
	t.NoError(err2)
	data, ok := info2.GetData("key")
	t.True(ok)
	t.EqualValues(123456, data)
	errDel := t.svc.LeaveTeam("role2", t.teamID)
	t.NoError(errDel)
}

func TestTeamService(t *testing.T) {
	suite.Run(t, newTeamSuite(t))
}
