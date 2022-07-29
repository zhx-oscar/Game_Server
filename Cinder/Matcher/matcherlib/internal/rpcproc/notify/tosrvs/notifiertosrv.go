package tosrvs

import (
	"Cinder/Matcher/matchapi/mtypes"
	"Cinder/Matcher/matcherlib/internal/rpcproc/notify/tosrvs/rpc"
	"context"
	"sync"
)

// _NotifierToSrv 通知一个Srv
type _NotifierToSrv struct {
	srvID    mtypes.SrvID
	cancel   context.CancelFunc
	hasMsgCh chan struct{} // 是否有数据，长度1

	msgsMtx sync.Mutex
	msgs    []mtypes.NotifyMsgToOneSrv
}

func newNotifierToSrv(srvID mtypes.SrvID) *_NotifierToSrv {
	ctx, cancel := context.WithCancel(context.Background())
	n := &_NotifierToSrv{
		srvID:    srvID,
		cancel:   cancel,
		hasMsgCh: make(chan struct{}),
	}
	go n.Run(ctx)
	return n
}

func (n *_NotifierToSrv) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// 最后一次发送。此时 n 已在 notifierToSrvs 中删除。可能又创建了新的 _NotifierToSrv.
			n.notify()
			return
		case <-n.hasMsgCh:
			n.notify()
			// 继续等待 hasMsgCh
		}
	}
}

func (n *_NotifierToSrv) Cancel() {
	n.cancel()
}

func (n *_NotifierToSrv) notify() {
	msgs := n.popAllMsgs()
	if len(msgs) == 0 {
		return // 最后一次发送一般为空
	}

	// json 打包。因为有包长限制，所以自适应分成几段。
	bufs := jsonMarshalMsgs(msgs)
	for _, buf := range bufs {
		rpc.RpcByID(n.srvID, "RPC_MatchNotify", buf)
	}
}

func (n *_NotifierToSrv) Push(msg mtypes.NotifyMsgToOneSrv) {
	n.msgsMtx.Lock()
	defer n.msgsMtx.Unlock()

	// 先加消息，后发信号
	n.msgs = append(n.msgs, msg)
	n.signalHasMsg()
}

func (n *_NotifierToSrv) popAllMsgs() []mtypes.NotifyMsgToOneSrv {
	n.msgsMtx.Lock()
	defer n.msgsMtx.Unlock()

	result := n.msgs
	n.msgs = nil
	return result
}

// signalHasMsg 发出有消息信号
func (n *_NotifierToSrv) signalHasMsg() {
	select {
	case n.hasMsgCh <- struct{}{}:
	default: // 已有消息
	}
}
