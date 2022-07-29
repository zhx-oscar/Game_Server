package Core

import (
	"Cinder/Base/CRpc"
	"Cinder/Base/Config"
	"Cinder/Base/Const"
	"Cinder/Base/Net"
	"Cinder/Base/Prop"
	"Cinder/Base/SrvNet"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"math/rand"
	"net"
	"time"

	log "github.com/cihub/seelog"
)

type _Core struct {
	isInit  bool
	srvInfo *Info

	mqService  SrvNet.INode
	rpcService CRpc.IService
	rpcClient  CRpc.IClient
	netService Net.IService

	SrvNet.INodeCaller
	SrvNet.INodeQuery
	CRpc.IRpcCaller
	Prop.IMgr
}

func New() ICore {
	Inst = &_Core{isInit: false}
	return Inst
}

func (srv *_Core) Init(info *Info) error {

	if srv.isInit {
		return errors.New("server had been initialized")
	}

	if info.ServiceType == "" {
		return errors.New("ServerType empty")
	}
	if info.AreaID == "" {
		return errors.New("AreaID empty")
	}
	if info.ServiceID == "" {
		return errors.New("ServerID empty")
	}

	srv.srvInfo = info

	if err := srv.initMQService(); err != nil {
		return fmt.Errorf("init MQ service: %w", err)
	}

	if err := srv.initRpcService(); err != nil {
		return fmt.Errorf("init RPC service: %w", err)
	}

	if err := srv.initRpcClient(); err != nil {
		return fmt.Errorf("init RPC client: %w", err)
	}

	if err := srv.initNetService(); err != nil {
		return fmt.Errorf("init net service: %w", err)
	}

	if err := srv.initPropMgr(); err != nil {
		return fmt.Errorf("init prop mgr: %w", err)
	}

	srv.isInit = true

	go srv.keepAlive()

	return nil
}

func (srv *_Core) keepAlive() {

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Debug("i'm alive ")
		}
	}
}

func (srv *_Core) GetServiceID() string {
	return srv.srvInfo.ServiceID
}

func (srv *_Core) GetServiceType() string {
	return srv.srvInfo.ServiceType
}

func (srv *_Core) GetNetNode() SrvNet.INode {
	return srv.mqService
}

func (srv *_Core) isValidListenAddr(addr string) bool {

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}

	_ = listener.Close()

	return true
}

func (srv *_Core) getValidListenPort() (int, error) {

	var addr string

	for i := 0; i < 15; i++ {
		port := srv.srvInfo.PortMax
		if srv.srvInfo.PortMax > srv.srvInfo.PortMin {
			port = rand.Intn(srv.srvInfo.PortMax-srv.srvInfo.PortMin) + srv.srvInfo.PortMin
		}
		addr = fmt.Sprintf("%s:%d", srv.srvInfo.ListenAddr, port)

		if srv.isValidListenAddr(addr) {
			return port, nil
		}
	}

	return 0, errors.New("couldn't find a valid listen port")
}

func (srv *_Core) initRpcService() error {

	if srv.srvInfo.RpcProc == nil {
		return nil
	}

	srv.rpcService = CRpc.NewService()

	if err := srv.rpcService.Init(srv.mqService, srv.srvInfo.RpcProc); err != nil {
		return err
	}

	return nil
}

func (srv *_Core) initRpcClient() error {

	srv.rpcClient = CRpc.NewClient()

	if err := srv.rpcClient.Init(srv.mqService); err != nil {
		return err
	}

	srv.IRpcCaller = srv.rpcClient

	return nil
}

func (srv *_Core) initNetService() error {

	if srv.srvInfo.NetServerMessageProc == nil {
		return nil
	}

	srv.netService = Net.NewService(viper.GetInt("TrafficControl.MaxTrafficCount"), viper.GetBool("TrafficControl.On_Off"))
	srv.netService.Register(srv.srvInfo.NetServerMessageProc)

	var port int
	var err error
	if port, err = srv.getValidListenPort(); err != nil {
		return err
	}

	listenAddr := fmt.Sprintf("%s:%d", srv.srvInfo.ListenAddr, port)
	if err = srv.netService.Init(listenAddr); err != nil {
		return err
	}

	// 向ETCD注册信息
	regAddr := fmt.Sprintf("%s:%d", srv.srvInfo.OuterAddr, port)
	if err = Config.Inst.SetValueAndKeepAlive(Const.GetNetSrvID(srv.GetServiceType(), srv.GetServiceID()), regAddr); err != nil {
		log.Error("register service failed", regAddr)
		return err
	} else {
		log.Info("register service success", regAddr)
	}

	return nil
}

func (srv *_Core) initMQService() error {
	srv.mqService = SrvNet.NewMQNode()

	if err := srv.mqService.Init(srv.srvInfo.AreaID, srv.srvInfo.ServiceID, srv.srvInfo.ServiceType); err != nil {
		return err
	}

	srv.INodeCaller = srv.mqService
	srv.INodeQuery = srv.mqService

	return nil
}

func (srv *_Core) initPropMgr() error {
	srv.IMgr = Prop.NewMgr(srv.mqService, srv.rpcClient)
	return nil
}

func (srv *_Core) Destroy() {

	if !srv.isInit {
		return
	}

	srv.isInit = false

	if srv.netService != nil {
		srv.netService.Destroy()
		srv.netService = nil
	}
	log.Info("net service destroy")

	if srv.rpcService != nil {
		srv.rpcService.Destroy()
		srv.rpcService = nil
	}
	log.Info("rpc service destroy")

	if srv.rpcClient != nil {
		srv.rpcClient.Destroy()
		srv.rpcClient = nil
	}
	log.Info("rpc client destroy")

	if srv.mqService != nil {
		srv.mqService.Destroy()
		srv.mqService = nil
	}
	log.Info("mq service destroy")
}
