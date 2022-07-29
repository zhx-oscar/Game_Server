package notify

import (
	"Cinder/Matcher/matchapi/mtypes"
	"Cinder/Matcher/matcherlib/internal/rpcproc/notify/tosrvs"

	"github.com/mohae/deepcopy"
)

type Notifier struct {
	srvIDToRoleIDs mtypes.SrvIDToRoleIDs
}

func NewNotifier(srvIDToRoleIDs mtypes.SrvIDToRoleIDs) *Notifier {
	return &Notifier{
		srvIDToRoleIDs: srvIDToRoleIDs,
	}
}

// PostToNotify 复制并缓存通知，然后发送
// 每个服只发一条
func (n *Notifier) PostToNotify(msg mtypes.NotifyMsg) {
	if len(n.srvIDToRoleIDs) == 0 {
		return
	}

	// 进入协程发送，需要复制消息，避免外部更改
	msgCopy := copyNotifyMsg(msg)
	for srvID, roleIDs := range n.srvIDToRoleIDs {
		toSrvMsg := mtypes.NotifyMsgToOneSrv{
			RoleIDs:   roleIDs,
			NotifyMsg: msgCopy,
		}
		tosrvs.PostToNotify(srvID, toSrvMsg)
	}
}

func copyNotifyMsg(msg mtypes.NotifyMsg) mtypes.NotifyMsg {
	return deepcopy.Copy(msg).(mtypes.NotifyMsg)
}
