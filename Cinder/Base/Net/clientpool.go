package Net

import (
	"Cinder/Base/Config"
	"Cinder/Base/Const"
	"context"
	"errors"
	"net"
	"sync"

	log "github.com/cihub/seelog"
)

/*
	client pool don't handle reconnect issue \
	Client Pool depended on etcd cluster for health check
*/

type _TcpClientPool struct {
	messageProc  IProc
	watchHandles []int

	sessMap sync.Map

	ctx       context.Context
	ctxCancel context.CancelFunc
}

func NewClientPool() IClientPool {

	cp := &_TcpClientPool{
		watchHandles: make([]int, 0, 10),
	}

	cp.ctx, cp.ctxCancel = context.WithCancel(context.Background())

	return cp
}

func (cp *_TcpClientPool) Register(handler IProc) {
	cp.messageProc = handler
}

func (cp *_TcpClientPool) Init(srvTypes []string) error {

	if cp.messageProc == nil {
		return errors.New("net client pool haven't message proc")
	}

	log.Debug("net client pool init ")

	for _, srvType := range srvTypes {

		addrPrefix := Const.GetNetSrvIDbySrvType(srvType)
		var addrs []string
		var keys []string
		var err error
		if keys, addrs, err = Config.Inst.GetValuesByPrefix(addrPrefix); err != nil {
			log.Error("get remote service addr failed " + addrPrefix)
		}

		if len(addrs) != len(keys) {
			return errors.New("couldn't happen")
		}

		log.Debug("net client pool add service type ", srvType)

		for i := 0; i < len(addrs); i++ {

			err = cp.addClient(keys[i], addrs[i])
			if err != nil {
				log.Debug("add client pool failed ", err)
			} else {
				log.Debug("net client pool dial succeed  ", keys[i], addrs[i])
			}
		}

		var handle int
		if handle, err = Config.Inst.WatchKeys(addrPrefix, cp.watchClient); err != nil {
			log.Error("watch service failed")
			return err
		}

		cp.watchHandles = append(cp.watchHandles, handle)
	}

	return nil
}

func (cp *_TcpClientPool) addClient(rid, addr string) error {

	_, ok := cp.sessMap.Load(rid)
	if ok {
		log.Debug("the service had existed ", rid)
		return errors.New("the service had existed " + rid)
	}

	log.Debug("net client pool dial to  ", addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return errors.New("dail to net service failed " + addr)
	}

	go cp.handleConn(rid, conn)

	return nil
}

func (cp *_TcpClientPool) handleConn(rid string, conn net.Conn) {

	conn.RemoteAddr().String()

	sess := NewTcpSess(conn)

	cp.messageProc.OnSessConnected(sess)
	cp.sessMap.Store(rid, sess)

loop:
	for {
		msg, msgNo, err := sess.Read()
		if err != nil {
			sess.close()
			break
		}

		cp.messageProc.OnSessMessageHandle(sess, msgNo, msg)

		select {
		case <-cp.ctx.Done():
			sess.close()
			break loop
		default:

		}
	}

	cp.messageProc.OnSessClosed(sess)
	cp.sessMap.Delete(rid)
}

func (cp *_TcpClientPool) deleteClient(rid string) {

	ii, ok := cp.sessMap.Load(rid)
	if !ok {
		return
	}

	log.Debug("net client pool remove client ", rid)

	sess := ii.(*_TcpSess)
	sess.close()
}

func (cp *_TcpClientPool) watchClient(opType int, rid string, addr string) {
	if opType == Config.KeyAdd {
		err := cp.addClient(rid, addr)
		if err != nil {
			log.Debug("watch client but add client failed ", err)
		} else {
			log.Debug("watch client and add client succeed ")
		}
	} else if opType == Config.KeyDelete {
		cp.deleteClient(rid)
	}
}
func (cp *_TcpClientPool) Destroy() {

	log.Debug("net client pool closing")

	cp.ctxCancel()
	for _, handle := range cp.watchHandles {
		_ = Config.Inst.CancelWatch(handle)
	}

	cp.sessMap.Range(func(key, value interface{}) bool {

		sess := value.(*_TcpSess)
		sess.close()

		cp.sessMap.Delete(key)
		return true
	})
}
