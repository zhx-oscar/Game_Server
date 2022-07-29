package stats

import (
	"context"
	"sync"
	"time"

	log "github.com/cihub/seelog"
)

type IManager interface {
	Add(name string, cost time.Duration)
	Start()
	Stop()
}

var defaultMgr IManager

func init() {
	defaultMgr = &fakeMgr{}
}

// Enable 启用真正的统计功能
func Enable() {
	_, ok := defaultMgr.(*fakeMgr)
	if ok {
		defaultMgr = NewManager()
		defaultMgr.Start()
	}
}

// Disable 停用统计功能
func Disable() {
	defaultMgr.Stop()

	_, ok := defaultMgr.(*fakeMgr)
	if !ok {
		defaultMgr = &fakeMgr{}
	}
}

type fakeMgr struct{}

func (mgr *fakeMgr) Add(name string, cost time.Duration) {}
func (mgr *fakeMgr) Start()                              {}
func (mgr *fakeMgr) Stop()                               {}

type Manager struct {
	actionMap sync.Map

	isRunning bool
	ctx       context.Context
	ctxFunc   context.CancelFunc
}

func NewManager() *Manager {
	mgr := &Manager{
		actionMap: sync.Map{},
	}
	mgr.ctx, mgr.ctxFunc = context.WithCancel(context.Background())
	return mgr
}

func Add(name string, cost time.Duration) { defaultMgr.Add(name, cost) }
func (mgr *Manager) Add(name string, cost time.Duration) {
	v, _ := mgr.actionMap.LoadOrStore(name, NewAction(name))
	v.(*Action).Add(cost)
}

func Start() { defaultMgr.Start() }
func (mgr *Manager) Start() {
	if mgr.isRunning {
		return
	}
	mgr.isRunning = true

	go mgr.loop()
}

func Stop() { defaultMgr.Stop() }
func (mgr *Manager) Stop() {
	if !mgr.isRunning {
		return
	}
	mgr.isRunning = false

	mgr.ctxFunc()
}

func (mgr *Manager) loop() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-mgr.ctx.Done():
			return

		case <-ticker.C:
			mgr.actionMap.Range(func(key, value interface{}) bool {
				act := value.(*Action)
				//act.Reset()
				act.Calc()
				log.Debug(act)

				return true
			})
		}
	}
}
