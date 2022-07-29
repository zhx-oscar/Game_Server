package mtypes

//

// RoleInfo 角色信息
// 对应 open-match 中的 Ticket.
// Role 角色是匹配的基本单位。
type RoleInfo struct {
	RoleID RoleID // 角色ID
	SrvID  SrvID  // 服务器ID, 用于RPC回调, matchapi 会自动赋值
	TeamID TeamID // 未组队则为空

	// 固定类型数据，用于匹配时判断。
	// 可以包含等级，积分等数据用于自定义的匹配算法。
	// 也可以包含昵称，头像等数据用于客户端显示。
	// 复杂的数据还可以json打包后存为字符串。
	FloatData  map[string]float64  // 数值
	StringData map[string]string   // 字符串
	Tags       map[string]struct{} // 标签
}

type RoleID string // 角色ID
type SrvID string  // 玩家所在服ID, 用于回调

func (r *RoleInfo) GetFloatData(key string) (float64, bool) {
	if r.FloatData == nil {
		return 0, false
	}
	result, ok := r.FloatData[key]
	return result, ok
}

func (r *RoleInfo) SetFloatData(key string, val float64) {
	if r.FloatData != nil {
		r.FloatData[key] = val
		return
	}
	r.FloatData = map[string]float64{key: val}
}

func (r *RoleInfo) GetStringData(key string) (string, bool) {
	if r.StringData == nil {
		return "", false
	}
	result, ok := r.StringData[key]
	return result, ok
}

func (r *RoleInfo) SetStringData(key string, val string) {
	if r.StringData != nil {
		r.StringData[key] = val
		return
	}
	r.StringData = map[string]string{key: val}
}

func (r *RoleInfo) HasTag(tag string) bool {
	if r.Tags == nil {
		return false
	}
	_, ok := r.Tags[tag]
	return ok
}

func (r *RoleInfo) AddTag(tag string) {
	if r.Tags != nil {
		r.Tags[tag] = struct{}{}
		return
	}
	r.Tags = map[string]struct{}{
		tag: struct{}{},
	}
}

func (r *RoleInfo) DelTag(tag string) {
	if r.Tags != nil {
		delete(r.Tags, tag)
		return
	}
}

func (r *RoleInfo) Copy() RoleInfo {
	cpy := *r
	cpy.FloatData = r.CopyFloatData()
	cpy.StringData = r.CopyStringData()
	cpy.Tags = r.CopyTags()
	return cpy
}

func (r *RoleInfo) CopyFloatData() map[string]float64 {
	cpy := make(map[string]float64, len(r.FloatData))
	for k, v := range r.FloatData {
		cpy[k] = v
	}
	return cpy
}

func (r *RoleInfo) CopyStringData() map[string]string {
	cpy := make(map[string]string, len(r.StringData))
	for k, v := range r.StringData {
		cpy[k] = v
	}
	return cpy
}

func (r *RoleInfo) CopyTags() map[string]struct{} {
	cpy := make(map[string]struct{}, len(r.Tags))
	for k, _ := range r.Tags {
		cpy[k] = struct{}{}
	}
	return cpy
}
