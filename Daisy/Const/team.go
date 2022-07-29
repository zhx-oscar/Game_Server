package Const

const TeamMaxMemberNum = 4 //队伍人数上限

const (
	TeamStatus_NORMAL uint32 = iota //普通队员
	TeamStatus_SECOND               //副队长
	TeamStatus_LEADER               //队长
)

const MaxInvites = 999              //邀请列表上限
const MaxApplys = 999               //申请列表上限
const InviteMaxLife = 7 * 24 * 3600 //邀请记录时限
const ApplyMaxLife = 7 * 24 * 3600  //申请记录时限

const LoadTeamDLockPrefix = "createTeam:"

const (
	TransferReason_Invite   uint8 = iota //邀请
	TransferReason_Apply                 //申请
	TransferReason_Kick                  //踢人
	TransferReason_Quit                  //离队
	TransferReason_AutoJoin              //自动加入
)
