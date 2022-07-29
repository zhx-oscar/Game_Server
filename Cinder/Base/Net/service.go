package Net

import (
	"context"
	"errors"
	"net"
	"time"

	log "github.com/cihub/seelog"
)

type _Service struct {
	messageProc IProc
	listener    net.Listener
	ctx         context.Context
	cancelCtx   context.CancelFunc

	maxTrafficCount     int  //每秒流量计数上限
	trafficControlOnOff bool //流量控制开关
}

func NewService(maxTrafficCount int, trafficControlOnOff bool) IService {
	srv := &_Service{
		maxTrafficCount:     maxTrafficCount,
		trafficControlOnOff: trafficControlOnOff,
	}
	srv.ctx, srv.cancelCtx = context.WithCancel(context.Background())
	return srv
}

func (srv *_Service) Register(handler IProc) {
	srv.messageProc = handler
}

func (srv *_Service) Init(addr string) error {

	if srv.messageProc == nil {
		return errors.New("no message proc")
	}

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	srv.listener = l

	go srv.acceptConn()

	log.Info("net service listen success ", addr)

	return nil
}

func (srv *_Service) acceptConn() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("accept conn coroutine quit panic ", err)
		}
	}()

loop:
	for {
		conn, err := srv.listener.Accept()
		if err != nil {
			log.Error("acceptConn err ", err)
			continue
		}

		select {
		case <-srv.ctx.Done():
			break loop
		default:
		}

		go srv.handleConn(conn)
	}
}

func (srv *_Service) handleConn(conn net.Conn) {
	log.Info("handleConn ", conn.RemoteAddr().String())

	sess := NewTcpSess(conn)

	srv.messageProc.OnSessConnected(sess)

	//2秒检测空链接，如果2秒之后sess还没有验证通过。则关闭sess
	time.AfterFunc(2*time.Second, func() {
		if !sess.IsValidate() {
			sess.Close()
		}
	})

	//流量控制ticker
	trafficCountTicker := time.NewTicker(1 * time.Second)
	defer trafficCountTicker.Stop()

loop:
	for {
		select {
		case <-srv.ctx.Done():
			sess.close()
			break loop
		case <-trafficCountTicker.C:
			if srv.trafficControlOnOff {
				if sess.trafficCount > srv.maxTrafficCount {
					log.Error("traffic control upper limit exceeded Sess: ", sess.GetData(), "RemoteAddr: ", conn.RemoteAddr().String())
					sess.close()
					break
				}

				sess.trafficCount = 0
			}
		default:
		}

		msg, msgNo, err := sess.Read()
		if err != nil {
			log.Error("handleConn read data err ", err, "Sess: ", sess.GetData(), "RemoteAddr: ", conn.RemoteAddr().String())
			sess.close()
			break
		}

		//流量计数统计
		if srv.trafficControlOnOff {
			sess.trafficCount++
		}

		srv.messageProc.OnSessMessageHandle(sess, msgNo, msg)
	}

	srv.messageProc.OnSessClosed(sess)

	log.Info("handleConn goroutina exit ", sess.GetData())

}

func (srv *_Service) Destroy() {
	srv.cancelCtx()
	if srv.listener != nil {
		_ = srv.listener.Close()
		srv.listener = nil
	}

	log.Info("net service closing")
}
