package Rpc

import (
	"Cinder/Base/Config"
	"Cinder/Base/Const"
	"context"
	"errors"
	log "github.com/cihub/seelog"
	"net"
	"net/rpc"
	"strconv"
	"strings"
)

type _RpcService struct {
	srvID   string
	srvType string

	listener net.Listener

	ctx       context.Context
	ctxCancel context.CancelFunc

	isInit bool
}

func NewRpcService(srvID, srvType string) *_RpcService {

	s := &_RpcService{
		srvID:   srvID,
		srvType: srvType,
		isInit:  false,
	}

	s.ctx, s.ctxCancel = context.WithCancel(context.Background())
	return s
}

func (srv *_RpcService) Init(addr string, rpcProc interface{}) error {

	if srv.isInit {
		return nil
	}

	if rpcProc == nil {
		return errors.New("rpc proc object is nil")
	}

	var err error
	if err = rpc.Register(rpcProc); err != nil {
		return err
	}

	var port int

	if port, err = srv.getListenPort(addr); err != nil {
		log.Error("get listen port failed ", addr)
		return err
	}

	log.Debug("rpc service listen port ", port)

	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error("net service listen failed ", addr)
		return err
	}

	log.Debug("rpc service listen succeed ", addr)

	if err := srv.registerService(port); err != nil {
		log.Error("register service failed")
		return err
	}

	srv.listener = l
	go srv.handleRpc(l)

	srv.isInit = true

	return nil
}

func (srv *_RpcService) handleRpc(l net.Listener) {

loop:
	for {
		conn, err := l.Accept()

		if err != nil {
			continue
		}

		log.Debug("rpc service accept a rpc connection , remote addr = ", conn.RemoteAddr().String())

		go rpc.ServeConn(conn)

		select {
		case <-srv.ctx.Done():
			break loop
		default:
		}
	}

	log.Debug("rpc service handle rpc connect coroutine exit ")
}

func (srv *_RpcService) getListenPort(addr string) (int, error) {
	rets := strings.Split(addr, ":")

	if len(rets) != 2 {
		return 0, errors.New("invalid addr " + addr)
	}

	return strconv.Atoi(rets[1])
}

func (srv *_RpcService) registerService(port int) error {

	var addr string
	var err error
	if addr, err = Const.GetSrvAddr(port, true); err != nil {
		return err
	}

	log.Debug("register rpc service to service center key = ", Const.GetRpcSrvID(srv.srvType, srv.srvID), "  value = ", addr)

	return Config.Inst.SetValueAndKeepAlive(Const.GetRpcSrvID(srv.srvType, srv.srvID), addr)
}

func (srv *_RpcService) Destroy() {

	log.Debug("rpc service closing")

	if srv.isInit {
		_ = srv.listener.Close()
		srv.ctxCancel()
		srv.isInit = false
	}

}
