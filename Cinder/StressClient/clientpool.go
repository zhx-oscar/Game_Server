package StressClient

import (
	_ "Cinder/Base/Log"
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"
)

func NewClientPool(loginAddr string, accountPrefix string, loginPeriod time.Duration, loginNumPrePeriod int, onlineNumMax int, loginLeastTime time.Duration, client IClient) IClientPool {
	cp := &_ClientPool{
		loginAddr:         loginAddr,
		accountPrefix:     accountPrefix,
		loginPeriod:       loginPeriod,
		loginNumPrePeriod: loginNumPrePeriod,
		onlineNumMax:      onlineNumMax,
		loginLeastTime:    loginLeastTime,
		clientType:        reflect.TypeOf(client).Elem(),
		closeC:            make(chan struct{}, 1),
	}

	cp.ctx, cp.ctxCancel = context.WithCancel(context.Background())

	return cp
}

type _ClientPool struct {
	loginAddr         string
	accountPrefix     string
	loginPeriod       time.Duration
	loginNumPrePeriod int
	onlineNumMax      int
	loginLeastTime    time.Duration

	clientType reflect.Type

	initList    sync.Map
	pendingList sync.Map
	succeedList sync.Map
	failedList  sync.Map

	ctx       context.Context
	ctxCancel context.CancelFunc

	closeC chan struct{}
}

func (cp *_ClientPool) Start() {
	cp.initSimClients()
	go cp.mainLoop()
}

func (cp *_ClientPool) Close() {
	cp.ctxCancel()
	<-cp.closeC
	close(cp.closeC)
}

func (cp *_ClientPool) mainLoop() {

	ticker := time.NewTicker(cp.loginPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cp.oneSecLoop()
		case <-cp.ctx.Done():
			cp.clearOnlineSimClient()
			cp.closeC <- struct{}{}
			return
		}
	}

}

func (cp *_ClientPool) oneSecLoop() {
	cp.simLogin()
}

func (cp *_ClientPool) simLogin() {

	var i int
	for i = 0; i < cp.loginNumPrePeriod; i++ {

		c, err := cp.fetchSimClient()
		if err != nil {
			break
		}

		if !c.CouldLogin() {
			break
		}

		cp.onLoginStart(c)
		go c.Login()
	}

	leftNum := cp.loginNumPrePeriod - i

	if leftNum > 0 {
		cp.returnSimClient(leftNum)
	}

}

func (cp *_ClientPool) fetchSimClient() (*_SimClient, error) {

	var sim *_SimClient
	cp.initList.Range(func(key, value interface{}) bool {
		sim = value.(*_SimClient)
		return false
	})

	if sim != nil {
		return sim, nil
	}

	cp.failedList.Range(func(key, value interface{}) bool {
		sim = value.(*_SimClient)
		return false
	})

	if sim != nil {
		return sim, nil
	}

	return nil, errors.New("no sim client ")
}

func (cp *_ClientPool) returnSimClient(num int) {

	retNum := 0
	retList := make([]*_SimClient, 0, num)

	cp.succeedList.Range(func(key, value interface{}) bool {

		sim := value.(*_SimClient)
		if sim.CouldRecycle(cp.loginLeastTime) {
			retNum++
			retList = append(retList, sim)
		}

		if retNum >= num {
			return false
		} else {
			return true
		}
	})

	for _, v := range retList {
		cp.onLogout(v)
		v.Logout()
	}

}

func (cp *_ClientPool) clearOnlineSimClient() {

	cp.succeedList.Range(func(key, value interface{}) bool {
		c := value.(*_SimClient)
		c.Logout()
		return true
	})

}

func (cp *_ClientPool) initSimClients() {
	num := cp.onlineNumMax + cp.loginNumPrePeriod

	for i := 0; i < num; i++ {
		cp.initList.Store(i, NewSimClient(i, cp))
	}
}

func (cp *_ClientPool) createClient() _IClient {
	return reflect.New(cp.clientType).Interface().(_IClient)
}

func (cp *_ClientPool) getAccount(id int) string {
	return fmt.Sprint(cp.accountPrefix, id)
}

func (cp *_ClientPool) onLoginStart(sim *_SimClient) {
	cp.clearSimClientFromList(sim)
	cp.pendingList.Store(sim.GetID(), sim)
}

func (cp *_ClientPool) onLoginSucceed(sim *_SimClient) {
	cp.clearSimClientFromList(sim)
	cp.succeedList.Store(sim.GetID(), sim)
}

func (cp *_ClientPool) onLoginFailed(sim *_SimClient) {
	cp.clearSimClientFromList(sim)
	cp.failedList.Store(sim.GetID(), sim)
}

func (cp *_ClientPool) onLogout(sim *_SimClient) {
	cp.clearSimClientFromList(sim)
	cp.initList.Store(sim.GetID(), sim)
}

func (cp *_ClientPool) clearSimClientFromList(sim *_SimClient) {
	_, ok := cp.initList.Load(sim.GetID())
	if ok {
		cp.initList.Delete(sim.GetID())
	}

	_, ok = cp.pendingList.Load(sim.GetID())
	if ok {
		cp.pendingList.Delete(sim.GetID())
	}

	_, ok = cp.failedList.Load(sim.GetID())
	if ok {
		cp.failedList.Delete(sim.GetID())
	}

	_, ok = cp.succeedList.Load(sim.GetID())
	if ok {
		cp.succeedList.Delete(sim.GetID())
	}
}
