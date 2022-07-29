package rpcproc

import (
	"Cinder/Matcher/matcherlib/internal/rpcproc/room"
	"Cinder/Matcher/matcherlib/internal/rpcproc/team"
	"time"

	log "github.com/cihub/seelog"
)

type RPCProc struct {
	_RPCProcTeam
	_RPCProcRoom
}

func init() {
	room.SetTeamInfoGetter(team.GetMgr())

	go runLogStatus()
}

func NewRPCProc() *RPCProc {
	return &RPCProc{}
}

func init() {
	go runLogStatus()
}

func runLogStatus() {
	for {
		time.Sleep(5 * time.Second)

		rooms := room.GetMgr().GetRoomCount()
		roomModes := room.GetMgr().GetRoomModeCount()
		teams := team.GetMgr().GetTeamCount()
		log.Infof("rooms=%d(modes=%d), teams=%d", rooms, roomModes, teams)
	}
}
