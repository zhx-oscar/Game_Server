package mtypes

type RoleMap map[RoleID]*RoleInfo
type SrvIDToRoleIDs map[SrvID][]RoleID

func (r RoleMap) GetSrvIDToRoleIDs() SrvIDToRoleIDs {
	result := make(SrvIDToRoleIDs, len(r))
	for roleID, info := range r {
		srvID := info.SrvID
		result[srvID] = append(result[srvID], roleID)
	}
	return result
}

func (r RoleMap) Copy() RoleMap {
	cpy := make(RoleMap, len(r))
	for k, v := range r {
		v2 := v.Copy()
		cpy[k] = &v2
	}
	return cpy
}
