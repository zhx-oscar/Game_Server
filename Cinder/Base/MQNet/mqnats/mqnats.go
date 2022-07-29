package mqnats

import (
	"Cinder/Base/MQNet"
	"Cinder/Base/Message"
	"Cinder/Base/Util"
	"context"
	"errors"
	log "github.com/cihub/seelog"
	"github.com/nats-io/nats.go"
)

var ErrMessageChanFull = errors.New("message chan full")

type _NATSService struct {
	options MQNet.Options

	procs   []MQNet.IProc
	nc      *nats.Conn
	msgChan chan []byte

	ctx        context.Context
	cancelFunc context.CancelFunc
}

func New() MQNet.IService {
	return &_NATSService{}
}

func (srv *_NATSService) Init(opts ...MQNet.Option) error {
	srv.options = MQNet.Options{}
	for _, o := range opts {
		o(&srv.options)
	}

	srv.procs = make([]MQNet.IProc, 0, 10)
	srv.msgChan = make(chan []byte, 1000)

	var err error
	srv.nc, err = nats.Connect(srv.options.Addr)
	if err != nil {
		return err
	}
	if _, err = srv.nc.Subscribe(srv.options.ServiceAddr, srv.onRecvMsg); err != nil {
		return err
	}
	if _, err = srv.nc.Subscribe(srv.options.BoardcastAddr, srv.onRecvMsg); err != nil {
		return err
	}

	srv.ctx, srv.cancelFunc = context.WithCancel(context.Background())

	go srv.procloop()

	return nil
}

func (srv *_NATSService) Destroy() {
	srv.cancelFunc()

	if srv.nc != nil {
		srv.nc.Close()
	}
}

func (srv *_NATSService) AddProc(proc MQNet.IProc) {
	srv.procs = append(srv.procs, proc)
}

func (srv *_NATSService) Post(addr string, msg Message.IMessage) error {
	if addr == "" {
		return MQNet.ErrInvalidPostAddr
	}
	if msg == nil {
		return MQNet.ErrInvalidMessage
	}
	if srv.nc == nil {
		return MQNet.ErrorNotInit
	}
	if !srv.nc.IsConnected() {
		return MQNet.ErrorNotInit
	}

	maxLen, err := MQNet.MaxMessageSize(srv.options.ServiceAddr, msg)
	if err != nil {
		return err
	}

	buf := Util.Get(maxLen)
	data, err := MQNet.Pack(srv.options.ServiceAddr, msg, buf)
	if err != nil {
		Util.Put(buf)
		return err
	}

	err = srv.nc.Publish(addr, data)
	Util.Put(buf)
	if err != nil {
		return err
	}

	return nil
}

func (srv *_NATSService) onRecvMsg(msg *nats.Msg) {
	select {
	case srv.msgChan <- msg.Data:
	default:
		log.Warn("msg chan full ", srv.options.ServiceAddr)
	}
}

func (srv *_NATSService) procloop() {
	for {
		select {
		case <-srv.ctx.Done():
			return

		case data := <-srv.msgChan:
			addr, imsg, err := MQNet.Unpack(data)
			if err != nil {
				log.Error("Proc Message unpark err ", err, " message: ", data)
				break
			}

			for _, p := range srv.procs {
				p.MessageProc(addr, imsg)
			}
		}
	}
}
