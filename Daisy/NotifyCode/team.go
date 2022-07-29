package NotifyCode

// team 77-89

const (
	ApplyTeam                = 1  // 准备加入行动组 1
	ApplyTeamWithName        = 2  // 准备加入行动组  2
	BeAgreeedJoinTeam        = 3  // 已同意，等待战斗结束 3
	RefuseFormTeamInvitatoin = 4  // 拒绝了组队邀请 4
	RefuseFormTeamApply      = 5  // 拒绝了组队申请 5
	JoinTeamFailed           = 6  // 加入失败，行动组已满员 6
	JoinTeamSuccessWithName  = 7  // 谁加入行动组成功 7
	JoinTeamSuccess          = 8  // 加入行动组成功 8
	SelfTramFull             = 9  // 自己的行动组已满员 9
	TargetTeamFull           = 10 // 目标行动组已满员 10
	JoiningTeam              = 11 // 正在加入队伍中 11
	LeaveTeam                = 26 // 离开队伍 26
	KickOutTeam              = 27 // 踢出队伍 27
	NotifyJoinTeam           = 28 // 通知队友谁加入队伍 28
)
