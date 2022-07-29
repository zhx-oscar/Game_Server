package ErrorCode

// 上下线相关错误定义 100~199
const (
	GetTeamError           = 9801
	EnterTeamSpaceError    = 9802
	LeaveTeamSpaceError    = 9803
	TeamMemberOnlineError  = 9804
	TeamMemberOfflineError = 9805
	GetUserError           = 9806 // 进入场景获取账号错误
	GetRoleError           = 9807 // 进入场景获取角色错误
)
