package room

import (
	"Cinder/Matcher/matchapi/mtypes"
	"errors"

	assert "github.com/arl/assertgo"
)

type _ITeamInfoGetter interface {
	GetTeamInfo(mtypes.TeamID) (mtypes.TeamInfo, bool)
}

var _teamInfoGetter _ITeamInfoGetter

func SetTeamInfoGetter(g _ITeamInfoGetter) {
	assert.True(g != nil)
	_teamInfoGetter = g
}

func getTeamRoles(teamID mtypes.TeamID) (mtypes.RoleMap, error) {
	assert.True(_teamInfoGetter != nil)
	teamInfo, ok := _teamInfoGetter.GetTeamInfo(teamID)
	if !ok {
		return nil, errors.New("can not find team")
	}
	if len(teamInfo.Roles) <= 0 {
		return nil, errors.New("team is empty")
	}
	return teamInfo.Roles, nil
}
