package bc

import (
	"Cinder/Chat/rpcproc/logic/types"
)

type IGroupMgr interface {
	GetOrLoadGroup(groupID types.GroupID) (types.IGroup, error)
}

var groupMgr IGroupMgr

func SetGroupMgr(mgr IGroupMgr) {
	groupMgr = mgr
}
