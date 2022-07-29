package internal

import (
	"Cinder/Matcher/matchapi/mtypes"
	"Cinder/Matcher/rpcmsg"
)

type TeamService struct {
	rpcCaller *_RpcCaller
}

func NewTeamService(matcherAreaID string, matcherServerID string) *TeamService {
	return &TeamService{
		rpcCaller: newRpcCaller(matcherAreaID, matcherServerID),
	}
}

// CreateTeam 创建队伍
func (t *TeamService) CreateTeam(creatorInfo mtypes.RoleInfo, teamInfo mtypes.TeamInfo) (mtypes.TeamInfo, error) {
	creatorInfo.SrvID = getServiceID()
	req := rpcmsg.CreateTeamReq{
		CreatorInfo: creatorInfo,
		TeamInfo:    teamInfo,
	}
	rsp, errRsp := t.rpcCaller.rpc2("RPC_CreateTeam", req)
	return rsp.GetCreateTeamRsp().TeamInfo, errRsp
}

// JoinTeam 加入队伍，触发队伍广播 RPC_MatchNotifyJoinTeam
func (t *TeamService) JoinTeam(roleInfo mtypes.RoleInfo, teamID mtypes.TeamID, passwd string) (mtypes.TeamInfo, error) {
	roleInfo.SrvID = getServiceID()
	req := rpcmsg.JoinTeamReq{
		RoleInfo: roleInfo,
		TeamID:   teamID,
		Password: passwd,
	}
	rsp, errRsp := t.rpcCaller.rpc2("RPC_JoinTeam", req)
	return rsp.GetJoinTeamRsp().TeamInfo, errRsp
}

// LeaveTeam 离开队伍, 或踢除队员, 触发队伍广播 RPC_MatchNotifyLeaveTeam
func (t *TeamService) LeaveTeam(memberID mtypes.RoleID, teamID mtypes.TeamID) error {
	return t.rpcCaller.rpc("RPC_LeaveTeam", string(memberID), string(teamID))
}

// ChangeTeamLeader 变更队长, 触发队伍广播 RPC_MatchNotifyChangeTeamLeader
func (t *TeamService) ChangeTeamLeader(teamID mtypes.TeamID, newLeader mtypes.RoleID) error {
	return t.rpcCaller.rpc("RPC_ChangeTeamLeader", string(teamID), string(newLeader))
}

// SetTeamData 设置队伍数据, 触发队伍广播 RPC_MatchNotifySetTeamData
func (t *TeamService) SetTeamData(teamID mtypes.TeamID, key string, data interface{}) error {
	_, err := t.rpcCaller.rpc2("RPC_SetTeamData", rpcmsg.SetTeamDataReq{
		TeamID: teamID,
		Key:    key,
		Data:   data,
	})
	return err
}

// Broadcast 队伍广播. 触发队伍广播 RPC_MatchNotifyBroadcastTeam
func (t *TeamService) BroadcastTeam(teamID mtypes.TeamID, data interface{}) error {
	_, err := t.rpcCaller.rpc2("RPC_BroadcastTeam", rpcmsg.BroadcastTeamReq{
		TeamID: teamID,
		Msg:    data,
	})
	return err
}
