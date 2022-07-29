package main

import (
	"Daisy/ErrorCode"
	"Daisy/Proto"
	"fmt"
	"strings"
)

//getRedPointKey 通过红点类型 + 额外参数 获得红点key
func (r *_Role) getRedPointKey(redPointType, param string) string {
	if param == "" {
		return redPointType
	}

	return fmt.Sprint(redPointType, "_", param)
}

//RPC_RemoveRedPoint 移除红点 移除某一类红点 或者 移除精准红点
func (user *_User) RPC_RemoveRedPoint(redPointType string, redPointKey string) int32 {
	if user.role == nil {
		return ErrorCode.RoleIsNil
	}

	return user.role.RemoveRedPoint(redPointType, redPointKey)
}

//RemoveRedPoint 移除红点
func (r *_Role) RemoveRedPoint(redPointType string, redPointKey string) int32 {
	var redPointKeys Proto.StringArry

	//移除某一类红点
	if redPointType != "" {
		for key := range r.prop.Data.RedPointsData {
			if strings.HasPrefix(key, redPointType) {
				redPointKeys.Data = append(redPointKeys.Data, key)
			}
		}
	}

	//精确移除某个小红点
	redPointKeys.Data = append(redPointKeys.Data, redPointKey)

	r.prop.SyncRemoveRedPointList(&redPointKeys)
	return ErrorCode.Success
}
