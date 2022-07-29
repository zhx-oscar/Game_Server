package main

import (
	"Daisy/DHDB"
	"Daisy/ErrorCode"
	"Daisy/Proto"
)

//RPC_GetRoleInfoList 批量获取玩家信息
func (user *_User) RPC_GetRoleInfoList(idList *Proto.StringArry) (int32, *Proto.RoleArry) {
	if user.role == nil || idList == nil {
		return ErrorCode.RoleIsNil, nil
	}

	return user.role.GetRoleInfoList(idList)
}

//GetRoleInfoList 批量获取玩家信息
func (r *_Role) GetRoleInfoList(idList *Proto.StringArry) (int32, *Proto.RoleArry) {
	result := &Proto.RoleArry{RoleList: []*Proto.Role{}}
	for _, id := range idList.Data {
		roleCache, err := DHDB.GetRoleCache(id)
		if err != nil {
			r.Error("DHDB.GetRolecaChe error ", id, err)
			continue
		}

		role := &Proto.Role{
			Base: roleCache.Base,
		}

		result.RoleList = append(result.RoleList, role)
	}

	return ErrorCode.Success, result
}

//RPC_GetTeamPartInfo 获取队伍信息
func (user *_User) RPC_GetTeamPartInfo(roleID string) (int32, *Proto.TeamPart) {
	if user.role == nil {
		return ErrorCode.RoleIsNil, nil
	}

	return user.role.GetTeamPartInfo(roleID)
}

//GetTeamPartInfo 获取队伍信息
func (r *_Role) GetTeamPartInfo(roleID string) (int32, *Proto.TeamPart) {
	data, err := DHDB.GetTeamPart(roleID)
	if err != nil {
		r.Error("DHDB.GetTeamPart error ", roleID, err)
		return ErrorCode.TeamPartErr, nil
	}

	return ErrorCode.Success, data
}
