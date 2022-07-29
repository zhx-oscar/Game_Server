package main

import (
	"Cinder/Base/DistributeLock"
	"Cinder/Space"
	"Daisy/Const"
	"Daisy/DB"
	"errors"
	"github.com/go-redis/redis/v7"
)

func LoadTeam(teamID string) (string, error) {
	if teamID == "" {
		return "", errors.New("teamID is empty")
	}

	createTeamLock := DistributeLock.New(Const.LoadTeamDLockPrefix + teamID)
	createTeamLock.Lock()
	defer createTeamLock.Unlock()

	srvID, err := DB.TeamUtil().GetSrvID(teamID)
	if err == redis.Nil {
		Space.Inst.CreateSpace(teamID, []byte{}, nil)
		return Space.Inst.GetServiceID(), nil
	} else if err != nil {
		return "", err
	} else {
		return srvID, nil
	}
}
