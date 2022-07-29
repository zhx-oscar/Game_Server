package DHDB

import (
	Prop2 "Cinder/Base/Prop"
	"Daisy/Prop"
	"Daisy/Proto"
	"errors"
	"fmt"
)

var ErrRoleIDInvalid = errors.New("role ids invalid")
var ErrTeamIDInvalid = errors.New("team ids invalid")

/**
获取缓存
*/

func GetRoleCache(roleID string) (*Proto.RoleCache, error) {
	if roleID == "" {
		return nil, ErrRoleIDInvalid
	}

	cache, err := Prop2.GetCacheProp(Prop.RolePropType, roleID)
	if err != nil {
		return nil, err
	}

	return cache.(*Proto.RoleCache), nil
}

func GetTeamBase(teamID string) (*Proto.TeamBase, error) {
	if teamID == "" {
		return nil, ErrTeamIDInvalid
	}

	cache, err := Prop2.GetCacheProp(Prop.TeamPropType, teamID)
	if err != nil {
		return nil, err
	}

	return cache.(*Proto.TeamCache).Base, nil
}

//GetTeamPart 目前只是封装 角色base  + teambase
func GetTeamPart(roleID string) (*Proto.TeamPart, error) {
	if roleID == "" {
		return nil, ErrRoleIDInvalid
	}

	roleCache, err := Prop2.GetCacheProp(Prop.RolePropType, roleID)
	if err != nil {
		return nil, err
	}

	teamCache, err := Prop2.GetCacheProp(Prop.TeamPropType, roleCache.(*Proto.RoleCache).Base.TeamID)
	if err != nil {
		return nil, err
	}

	result := &Proto.TeamPart{
		RoleList: []*Proto.Role{},
		TeamInfo: &Proto.Team{Base: teamCache.(*Proto.TeamCache).Base},
	}

	for id := range teamCache.(*Proto.TeamCache).Base.Members {
		var tempRole Proto.Role

		if id == roleID {
			tempRole.Base = roleCache.(*Proto.RoleCache).Base
		} else {
			temp, err := Prop2.GetCacheProp(Prop.RolePropType, id)
			if err != nil {
				return nil, fmt.Errorf("RPC_GetRoleBase id:%s err:%s", roleID, err)
			}

			tempRole.Base = temp.(*Proto.RoleCache).Base
		}

		result.RoleList = append(result.RoleList, &tempRole)
	}

	return result, nil
}
