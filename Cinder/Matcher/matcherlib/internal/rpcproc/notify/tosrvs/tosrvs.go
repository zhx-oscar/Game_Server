// tosrvs 开启协程通知各个服务器
package tosrvs

import (
	"Cinder/Matcher/matchapi/mtypes"
)

func PostToNotify(srvID mtypes.SrvID, msg mtypes.NotifyMsgToOneSrv) {
	notifierToSrvs.PostToNotify(srvID, msg)
}
