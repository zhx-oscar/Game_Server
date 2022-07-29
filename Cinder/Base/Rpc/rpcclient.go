package Rpc

import (
	"Cinder/Base/Config"
	"Cinder/Base/Const"
	"errors"
	log "github.com/cihub/seelog"
	"math/rand"
	"net/rpc"
	"sync"
)

type _RpcClient struct {
	cli *rpc.Client

	srvType string
	srvID   string
}

func newRpcClient(srvType, srvID string) *_RpcClient {
	return &_RpcClient{
		srvType: srvType,
		srvID:   srvID,
	}
}

func (c *_RpcClient) init(addr string) error {

	var err error
	c.cli, err = rpc.Dial("tcp", addr)
	if err != nil {
		return err
	}

	return nil
}

func (c *_RpcClient) close() {
	if c.cli != nil {
		_ = c.cli.Close()
		c.cli = nil
	}
}

func (c *_RpcClient) call(serviceMethod string, args interface{}, reply interface{}) error {
	return c.cli.Call(serviceMethod, args, reply)
}

type _RpcClientPool struct {
	mtx           sync.Mutex
	clientsByType sync.Map
	clientsByID   sync.Map
}

func NewRpcClientPool() *_RpcClientPool {
	return &_RpcClientPool{}
}

func (p *_RpcClientPool) Init() error {

	log.Debug("rpc client pool init ")
	return nil
}

func (p *_RpcClientPool) Destroy() {

	log.Debug("rpc client pool closing")

	p.clientsByID.Range(func(key, value interface{}) bool {

		cli := value.(*_RpcClient)
		cli.close()

		p.clientsByID.Delete(key)
		return true
	})

	p.clientsByType.Range(func(key, value interface{}) bool {
		cli := value.(*_RpcClient)
		cli.close()

		p.clientsByID.Delete(key)
		return true
	})

}

var maxTryTime int = 2

func (p *_RpcClientPool) CallBySrvType(srvType string, serviceMethod string, args interface{}, reply interface{}) error {

	for i := 0; i < maxTryTime; i++ {

		c, err := p.getClientByType(srvType)
		if err != nil {
			log.Debug("get rpc client by type ", srvType, " failed , try again ", err.Error())
			continue
		}

		// 此处无法区分是网络原因导致Rpc调用错误还是Rpc函数返回错误
		// 所以统一采用重新连接的方式处理
		err = c.call(serviceMethod, args, reply)
		if err == nil {
			return nil
		} else {
			log.Debug("call rpc error ", err.Error())
			p.deleteClient(c)
		}

	}

	return errors.New("call rpc failed ")
}

func (p *_RpcClientPool) CallBySrvID(srvType, srvID string, serviceMethod string, args interface{}, reply interface{}) error {

	c, err := p.getClientByID(srvType, srvID)
	if err != nil {
		return err
	}

	return c.call(serviceMethod, args, reply)
}

func (p *_RpcClientPool) getClientByType(srvType string) (*_RpcClient, error) {

	p.mtx.Lock()
	defer p.mtx.Unlock()

	key := Const.GetRpcSrvIDbySrvType(srvType)

	v, ok := p.clientsByType.Load(srvType)
	if ok {
		return v.(*_RpcClient), nil
	}

	c := newRpcClient(srvType, "")
	_, addrs, err := Config.Inst.GetValuesByPrefix(key)
	if err != nil {
		return nil, err
	}

	if len(addrs) == 0 {
		return nil, errors.New("no rpc service found")
	}

	addr := addrs[rand.Intn(len(addrs))]

	log.Debug("random get rpc service addr ", addr)

	if err = c.init(addr); err != nil {
		return nil, err
	}

	log.Debug("dial to rpc service succeed ")

	ac, ok := p.clientsByType.LoadOrStore(srvType, c)
	if ok {
		c.close()
	}

	return ac.(*_RpcClient), nil
}

func (p *_RpcClientPool) getClientByID(srvType, srvID string) (*_RpcClient, error) {

	p.mtx.Lock()
	defer p.mtx.Unlock()

	key := Const.GetRpcSrvID(srvType, srvID)

	v, ok := p.clientsByID.Load(srvID)
	if ok {
		return v.(*_RpcClient), nil
	}

	c := newRpcClient("", srvID)
	addr, err := Config.Inst.GetValue(key)
	if err != nil {
		return nil, err
	}

	if err = c.init(addr); err != nil {
		return nil, err
	}

	ac, ok := p.clientsByID.LoadOrStore(srvID, c)
	if ok {
		c.close()
	}

	return ac.(*_RpcClient), nil
}

func (p *_RpcClientPool) deleteClient(c *_RpcClient) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	c.close()

	if c.srvType != "" {
		p.clientsByType.Delete(c.srvType)
	}

	if c.srvID != "" {
		p.clientsByID.Delete(c.srvID)
	}
}
