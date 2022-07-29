package DB

import (
	"Cinder/Cache"
	"errors"
	"fmt"
	"time"
)

// ITeamUtil 队伍信息注册接口
type ITeamUtil interface {
	GetSrvID(teamID string) (string, error)
	Register(teamID string, srvID string) error
	UnRegister(teamID string) error
	UpdateExpire(teamID string) error
}

const (
	teamPrefix = "team"

	teamExpire = 15 * time.Second
)

var ErrParamInvalid = errors.New("params invalid")

type teamUtil struct{}

var defaultTeamLoader ITeamUtil

func TeamUtil() ITeamUtil {
	return defaultTeamLoader
}

func init() {
	defaultTeamLoader = NewTeamLoader()
}

func NewTeamLoader() ITeamUtil {
	return &teamUtil{}
}

func (tl *teamUtil) GetSrvID(teamID string) (string, error) {
	if teamID == "" {
		return "", ErrParamInvalid
	}

	return Cache.RedisDB.Get(tl.genTeamkey(teamID)).Result()
}

func (tl *teamUtil) Register(teamID string, srvID string) error {
	if teamID == "" || srvID == "" {
		return ErrParamInvalid
	}

	return Cache.RedisDB.Set(tl.genTeamkey(teamID), srvID, teamExpire).Err()
}

func (tl *teamUtil) UnRegister(teamID string) error {
	if teamID == "" {
		return ErrParamInvalid
	}

	return Cache.RedisDB.Del(tl.genTeamkey(teamID)).Err()
}

func (tl *teamUtil) UpdateExpire(teamID string) error {
	if teamID == "" {
		return ErrParamInvalid
	}

	return Cache.RedisDB.Expire(tl.genTeamkey(teamID), teamExpire).Err()
}

func (tl *teamUtil) genTeamkey(teamID string) string {
	return fmt.Sprintf("%s:%s", teamPrefix, teamID)
}
