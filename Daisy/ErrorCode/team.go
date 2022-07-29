package ErrorCode

// 队伍相关 1701-
const (
	JoinTeamError = 1701 //加入行动组失败
	QuitTeamError = 1702 //离开行动组失败

	TeamFull          = 1703 //行动组已满员
	AlreadInTeam      = 1704 //已加入一个行动组
	RepeatPlaceHolder = 1705 //重复占位
	NotInInvites      = 1706 //不在邀请列表中
	NotInApplys       = 1707 //不要申请列表中
	InviteOverLimit   = 1708 //邀请数目超出
	ApplyOverLimit    = 1709 //申请数目超出

	RepeatApply   = 1710 //重复申请
	FindSelf      = 1711 //查询自己
	FindNone      = 1712 //查询失败
	GetInvitesErr = 1713
	GetApplysErr  = 1714
	BoardTooLong  = 1715 //招募信息太长
	RepeatInvite  = 1716 //重复邀请

	CreateTeamOfflinePropFailed = 1717
	LoadTeamError               = 1718
	NotInTeam                   = 1719 //不在行动组中
	RoleTransfering             = 1720 //转移中
	NoPermission                = 1721 //权限不足

	TeamSrvFull     = 1722
	CreateTeamError = 1723
	EnterSpaceError = 1724
	LeaveSpaceError = 1725
)
