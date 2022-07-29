package tosrvs

import (
	"Cinder/Matcher/matchapi/mtypes"
	"sync"
	"time"

	assert "github.com/arl/assertgo"
)

// _NotifierToSrvs 缓存通知，开启协程发送通知。
// 本身有一个协程，每个Srv有个协程。某个Srv长期为空时则删除。
type _NotifierToSrvs struct {
	notifiers  sync.Map // map[srvID]*_NotifierToSrv
	idleSrvIDs sync.Map // map[srvID]struct{}
}

var notifierToSrvs = newNotifierToSrvs()

func newNotifierToSrvs() *_NotifierToSrvs {
	n := &_NotifierToSrvs{}
	go n.Run()
	return n
}

func (n *_NotifierToSrvs) Run() {
	for {
		time.Sleep(time.Second * 60)
		// 判断已空闲：1min无消息
		n.deleteIdleNotifiers()
	}
}

func (n *_NotifierToSrvs) PostToNotify(srvID mtypes.SrvID, msg mtypes.NotifyMsgToOneSrv) {
	n.setBusy(srvID) // 非空闲

	intf, ok := n.notifiers.Load(srvID)
	if ok {
		intf.(*_NotifierToSrv).Push(msg)
		return // 多数情况
	}

	// 需要新建
	notifier := newNotifierToSrv(srvID)
	actual, loaded := n.notifiers.LoadOrStore(srvID, notifier) // 有可能已创建
	assert.True((loaded && notifier != actual) || (!loaded && notifier == actual))
	if loaded {
		notifier.Cancel() // 没用了
	}
	actual.(*_NotifierToSrv).Push(msg)
}

// deleteIdleNotifiers 删除空闲的
func (n *_NotifierToSrvs) deleteIdleNotifiers() {
	n.idleSrvIDs.Range(func(key interface{}, _ interface{}) bool {
		srvID := key.(mtypes.SrvID)
		n.idleSrvIDs.Delete(srvID)

		n.delete(srvID)
		return true
	})

	// 重新设置 idleSrvIDs
	n.notifiers.Range(func(key interface{}, _ interface{}) bool {
		n.idleSrvIDs.Store(key, struct{}{})
		return true
	})
}

// setBusy 设为非空闲的
func (n *_NotifierToSrvs) setBusy(srvID mtypes.SrvID) {
	n.idleSrvIDs.Delete(srvID)
}

func (n *_NotifierToSrvs) delete(srvID mtypes.SrvID) {
	itf, ok := n.notifiers.Load(srvID)
	if ok {
		return // 不存在
	}
	n.notifiers.Delete(srvID)
	itf.(*_NotifierToSrv).Cancel()
}
