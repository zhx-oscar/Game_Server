package mtypes

// TeamInfo 是队伍信息
type TeamInfo struct {
	TeamID  TeamID // 队伍ID, 创建时输入为空
	Passwd  string // 队伍密码
	MaxRole int    // 容量，必须大于0

	Data M // map[string]interface{} 自定义数据

	Roles    RoleMap // 成员列表, 创建时输入为空
	LeaderID RoleID  // 队长ID, 创建时输入为空
}

// 队伍ID
type TeamID string

func (t *TeamInfo) GetData(key string) (interface{}, bool) {
	if t.Data == nil {
		return nil, false
	}
	result, ok := t.Data[key]
	return result, ok
}

func (t *TeamInfo) SetData(key string, val interface{}) {
	if t.Data != nil {
		t.Data[key] = val
		return
	}
	t.Data = map[string]interface{}{key: val}
}

func (t *TeamInfo) GetSrvIDToRoleIDs() SrvIDToRoleIDs {
	return t.Roles.GetSrvIDToRoleIDs()
}

func (t *TeamInfo) Copy() TeamInfo {
	cpy := *t
	cpy.Data = t.Data.Copy()
	cpy.Roles = t.Roles.Copy()
	return cpy
}
