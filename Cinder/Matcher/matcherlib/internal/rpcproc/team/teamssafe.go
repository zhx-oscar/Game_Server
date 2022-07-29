package team

import (
	"Cinder/Matcher/matchapi/mtypes"
	"sync"
)

// _TeamsSafe 协程安全的队伍map
type _TeamsSafe struct {
	mtx sync.Mutex

	teams map[mtypes.TeamID]*_Team
}

func newTeamsSafe() *_TeamsSafe {
	return &_TeamsSafe{
		teams: make(map[mtypes.TeamID]*_Team),
	}
}

func (t *_TeamsSafe) Insert(team *_Team) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	t.teams[team.ID()] = team
}

func (t *_TeamsSafe) Delete(teamID mtypes.TeamID) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	delete(t.teams, teamID)
}

func (t *_TeamsSafe) GetTeam(teamID mtypes.TeamID) *_Team {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	team, _ := t.teams[teamID]
	return team
}

func (t *_TeamsSafe) GetCount() int {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return len(t.teams)
}
