package DB

import (
	"errors"
	"math"
	"sync"
)

const (
	roleStartID uint64 = 100000000
	teamStartID uint64 = 100000000

	uidGeneratorKey = "uidgenerator"
	roleUIDField    = "role"
	teamUIDField    = "team"

	segSize = 100
)

var (
	ErrRoleIDExceed = errors.New("role id exceed max")
	ErrTeamIDExceed = errors.New("team id exceed max")
)

type uidUtil struct {
	roleCurID uint64
	roleMaxID uint64
	roleMu    sync.Mutex

	teamCurID uint64
	teamMaxID uint64
	teamMu    sync.Mutex
}

var defaultUidUtil = uidUtil{}

func (util *uidUtil) FetchRoleUID() (uint64, error) {
	util.roleMu.Lock()
	defer util.roleMu.Unlock()

	if util.roleCurID < util.roleMaxID {
		util.roleCurID++
	} else {
		var err error
		util.roleMaxID, err = RedisDB.HIncrBy(uidGeneratorKey, roleUIDField, segSize).Uint64()
		if err != nil {
			return 0, err
		}

		util.roleCurID = util.roleMaxID - segSize + 1
	}

	rawid := util.roleCurID
	if rawid > math.MaxUint64-roleStartID {
		return 0, ErrRoleIDExceed
	}
	return rawid + roleStartID, nil
}

func (util *uidUtil) FetchTeamUID() (uint64, error) {
	util.teamMu.Lock()
	defer util.teamMu.Unlock()

	if util.teamCurID < util.teamMaxID {
		util.teamCurID++
	} else {
		var err error
		util.teamMaxID, err = RedisDB.HIncrBy(uidGeneratorKey, teamUIDField, segSize).Uint64()
		if err != nil {
			return 0, err
		}

		util.teamCurID = util.teamMaxID - segSize + 1
	}

	rawid := util.teamCurID
	if rawid > math.MaxUint64-teamStartID {
		return 0, ErrTeamIDExceed
	}
	return rawid + teamStartID, nil
}

func FetchRoleUID() (uint64, error) {
	return defaultUidUtil.FetchRoleUID()
}

func FetchTeamUID() (uint64, error) {
	return defaultUidUtil.FetchTeamUID()
}
