package mqnsq

import (
	"Cinder/Base/MQNet"
	"Cinder/Base/Message"
	"Cinder/Base/Util"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/cihub/seelog"
	"github.com/nsqio/go-nsq"
)

type _NSQService struct {
	options    MQNet.Options
	lookupAddr string
	adminAddr  string

	procs    []MQNet.IProc
	producer *nsq.Producer
	msgChan  chan *nsq.Message

	consumer          *nsq.Consumer
	broadcastConsumer *nsq.Consumer
	nsqLogger         *mqLogger

	sendList Util.ISafeList

	ctx        context.Context
	cancelFunc context.CancelFunc
}

func New() MQNet.IService {
	return &_NSQService{}
}

func (srv *_NSQService) Init(opts ...MQNet.Option) error {
	srv.options = MQNet.Options{}
	for _, o := range opts {
		o(&srv.options)
	}
	if addr, ok := srv.options.ExtOpts[nsqLookupKey]; ok {
		srv.lookupAddr = addr.(string)
	} else {
		return ErrLookupAddrInvalid
	}
	if addr, ok := srv.options.ExtOpts[nsqAdminKey]; ok {
		srv.adminAddr = addr.(string)
	} else {
		return ErrAdminAddrInvalid
	}

	srv.procs = make([]MQNet.IProc, 0, 10)
	srv.msgChan = make(chan *nsq.Message, 1000)
	srv.sendList = Util.NewSafeList()
	srv.nsqLogger = &mqLogger{logger: log.Current}

	srv.ctx, srv.cancelFunc = context.WithCancel(context.Background())

	if err := srv.initProducer(srv.options.Addr); err != nil {
		return fmt.Errorf("init producer: %w", err)
	}

	srv.sayHello()

	if err := srv.initConsumer(srv.lookupAddr); err != nil {
		return fmt.Errorf("init consumer: %w", err)
	}

	go srv.sendLoop()
	go srv.procloop()

	return nil
}

func (srv *_NSQService) Destroy() {
	srv.cancelFunc()

	if srv.producer != nil {
		srv.producer.Stop()
	}

	if srv.consumer != nil {
		srv.consumer.Stop()
	}

	if srv.broadcastConsumer != nil {
		srv.broadcastConsumer.Stop()
	}
}

func (srv *_NSQService) AddProc(proc MQNet.IProc) {
	srv.procs = append(srv.procs, proc)
}

func (srv *_NSQService) Post(addr string, message Message.IMessage) error {
	if srv.producer == nil {
		return MQNet.ErrorNotInit
	}

	if message == nil {
		return nil
	}

	srv.sendList.Put(&MQNet.SendInfo{
		Addr:    addr,
		Message: message,
	})
	return nil

}

func (srv *_NSQService) sendLoop() {
	for {
		select {
		case <-srv.sendList.Signal():
			for {
				info, err := srv.sendList.Pop()
				if err != nil {
					break
				}

				sendInfo := info.(*MQNet.SendInfo)
				maxLen, err := MQNet.MaxMessageSize(srv.options.ServiceAddr, sendInfo.Message)
				if err != nil {
					log.Error("sendLoop get message size err ", err)
					break
				}

				buf := Util.Get(maxLen)
				data, err := MQNet.Pack(srv.options.ServiceAddr, sendInfo.Message, buf)
				if err != nil {
					log.Error("sendLoop pack message err ", err)
					Util.Put(buf)
					break
				}

				sendBuf := make([]byte, len(data))
				copy(sendBuf, data)
				Util.Put(buf)
				err = srv.producer.PublishAsync(sendInfo.Addr+"#ephemeral", sendBuf, nil)
				if err != nil {
					log.Error("sendLoop PublishAsync err ", err, " Target: ", sendInfo.Addr, " MessageID: ", sendInfo.Message.GetID())
					break
				}
			}

		case <-srv.ctx.Done():
			log.Info("sendLoop exit ", srv.options.ServiceAddr)
			return
		}
	}
}

func (srv *_NSQService) initProducer(addr string) error {

	cfg := nsq.NewConfig()

	cfg.HeartbeatInterval = 5 * time.Second

	var err error
	srv.producer, err = nsq.NewProducer(addr, cfg)
	if err != nil {
		return err
	}

	srv.producer.SetLogger(srv.nsqLogger, nsq.LogLevelError)

	return nil
}

func (srv *_NSQService) sayHello() {
	_ = srv.Post(srv.options.ServiceAddr, &Message.MQHello{
		Greeting: "this is hello from my channel " + srv.options.ServiceAddr,
	})

	_ = srv.Post(srv.options.BoardcastAddr, &Message.MQHello{
		Greeting: "this is hello from my broad channel " + srv.options.BoardcastAddr,
	})
}

func (srv *_NSQService) initConsumer(addr string) error {

	cfg := nsq.NewConfig()

	cfg.LookupdPollInterval = 10 * time.Second
	cfg.MaxInFlight = 1000
	cfg.HeartbeatInterval = 5 * time.Second
	cfg.MaxAttempts = 0

	channelName := "channel_" + srv.options.BoardcastAddr + "_" + srv.options.ServiceAddr + "#ephemeral"

	// 清空原有channel中的消息
	var body struct {
		Action string `json:"action"`
	}
	body.Action = "empty"

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}
	url := fmt.Sprintf("http://%s/api/topics/%s/%s", srv.adminAddr, srv.options.ServiceAddr+"#ephemeral", channelName)
	bytesReader := bytes.NewReader(data)
	http.Post(url, "application/json", bytesReader)

	srv.consumer, err = nsq.NewConsumer(srv.options.ServiceAddr+"#ephemeral", channelName, cfg)
	if err != nil {
		return fmt.Errorf("mqnsq new consumer: %w", err)
	}

	srv.consumer.SetLogger(srv.nsqLogger, nsq.LogLevelError)

	srv.consumer.AddHandler(srv)
	if err = srv.consumer.ConnectToNSQLookupd(addr); err != nil {
		return fmt.Errorf("consumer connect to mqnsq lookupd: %w", err)
	}

	url = fmt.Sprintf("http://%s/api/topics/%s/%s", srv.adminAddr, srv.options.BoardcastAddr+"#ephemeral", channelName)
	http.Post(url, "application/json", bytesReader)

	srv.broadcastConsumer, err = nsq.NewConsumer(srv.options.BoardcastAddr+"#ephemeral", channelName, cfg)
	if err != nil {
		return fmt.Errorf("mqnsq new broadcast consumer: %w", err)
	}

	srv.broadcastConsumer.SetLogger(srv.nsqLogger, nsq.LogLevelError)

	srv.broadcastConsumer.AddHandler(srv)
	if err = srv.broadcastConsumer.ConnectToNSQLookupd(addr); err != nil {
		return fmt.Errorf("broadcast consumer connect to mqnsq lookupd: %w", err)
	}

	return nil
}

func (srv *_NSQService) HandleMessage(msg *nsq.Message) error {
	select {
	case srv.msgChan <- msg:
	default:
		log.Warn("HandleMessage msg chan full ", srv.options.ServiceAddr)
	}

	return nil
}

func (srv *_NSQService) procloop() {
	for {
		select {
		case <-srv.ctx.Done():
			return

		case msg := <-srv.msgChan:
			addr, imsg, err := MQNet.Unpack(msg.Body)
			if err != nil {
				log.Error("Proc Message unpark err ", err, " message: ", msg.Body)
				break
			}

			if imsg.GetID() == Message.ID_MQ_Hello {
				log.Debug("Receive mq hello message, message: ", imsg.(*Message.MQHello).Greeting)
				break
			}

			for _, p := range srv.procs {
				p.MessageProc(addr, imsg)
			}
		}
	}
}
