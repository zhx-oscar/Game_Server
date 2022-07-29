package ErrorCode

const (
	Failure = -1
	Success = 0
)

// 通用错误定义
const (
	UserIDInvalid  = 9901 // 角色ID错误
	ArgsWrong      = 9902 // 参数错误
	ConfigNotExist = 9903 // 配置不存在
	DBOpErr        = 9904 // 数据库错误
	MarshalJsonErr = 9905 // MarshalJsonErr
	TeamPartErr    = 9906 // 获取队伍缓存失败
	Timeout        = 9907 // 请求超时
	RoleIsNil      = 9908 //角色对应的role为空
)
