package rpcproc

import (
	"Cinder/Matcher/matchapi/mtypes"
	"Cinder/Matcher/matcherlib/internal/rpcproc/team"
	"Cinder/Matcher/rpcmsg"
	"encoding/json"

	log "github.com/cihub/seelog"
)

type _RPCProcTeam struct {
}

// RPC_CreateTeam 创建队伍
func (r *_RPCProcTeam) RPC_CreateTeam(reqJson []byte) (rspJson []byte, errStr string) {
	log.Debugf("RPC_CreateTeam(reqJson=%d bytes)", len(reqJson))
	req := rpcmsg.CreateTeamReq{}
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err.Error()
	}

	teamInfo, errCreate := team.GetMgr().CreateTeam(req.CreatorInfo, req.TeamInfo)
	if errCreate != nil {
		return nil, errCreate.Error()
	}
	return rpcmsg.RPCResponse{
		CreateTeamRsp: &rpcmsg.CreateTeamRsp{
			TeamInfo: teamInfo,
		},
	}.Marshal()
}

// RPC_JoinTeam 加入队伍，触发队伍广播 RPC_MatchNotifyJoinTeam
func (r *_RPCProcTeam) RPC_JoinTeam(reqJson []byte) (rspJson []byte, errStr string) {
	log.Debugf("RPC_JoinTeam(reqJson=%d bytes)", len(reqJson))
	req := rpcmsg.JoinTeamReq{}
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err.Error()
	}
	teamInfo, errJoin := team.GetMgr().JoinTeam(req.RoleInfo, req.TeamID, req.Password)
	if errJoin != nil {
		return nil, errJoin.Error()
	}
	return rpcmsg.RPCResponse{
		JoinTeamRsp: &rpcmsg.JoinTeamRsp{
			TeamInfo: teamInfo,
		},
	}.Marshal()
}

// RPC_LeaveTeam 离开队伍, 或踢除队员, 触发队伍广播 RPC_MatchNotifyLeaveTeam
func (r *_RPCProcTeam) RPC_LeaveTeam(memberID string, teamID string) (errStr string) {
	log.Debugf("RPC_LeaveTeam(memberID=%s, teamID=%s)", memberID, teamID)
	team.GetMgr().LeaveTeam(mtypes.RoleID(memberID), mtypes.TeamID(teamID))
	return ""
}

// RPC_ChangeTeamLeader 变更队长
func (r *_RPCProcTeam) RPC_ChangeTeamLeader(teamID string, newLeaderID string) (errStr string) {
	log.Debugf("RPC_ChangeTeamLeader(teamID=%s, newLeaderID=%s)", teamID, newLeaderID)
	if err := team.GetMgr().ChangeTeamLeader(mtypes.TeamID(teamID), mtypes.RoleID(newLeaderID)); err != nil {
		return err.Error()
	}
	return ""
}

// RPC_SetTeamData 设置队伍数据，触发队伍广播 RPC_MatchNotifySetTeamData
func (r *_RPCProcTeam) RPC_SetTeamData(reqJson []byte) (rspJson []byte, errStr string) {
	log.Debugf("RPC_SetTeamData(reqJson=%d bytes)", len(reqJson))
	req := rpcmsg.SetTeamDataReq{}
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err.Error()
	}
	team.GetMgr().SetTeamData(req.TeamID, req.Key, req.Data)
	return rpcmsg.RPCResponse{}.Marshal()
}

// RPC_BroadcastTeam 队伍广播
func (r *_RPCProcTeam) RPC_BroadcastTeam(reqJson []byte) (rspJson []byte, errStr string) {
	log.Debugf("RPC_BroadcastTeam(reqJson=%d bytes)", len(reqJson))
	req := rpcmsg.BroadcastTeamReq{}
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err.Error()
	}
	team.GetMgr().BroadcastTeam(req.TeamID, req.Msg)
	return rpcmsg.RPCResponse{}.Marshal()
}
