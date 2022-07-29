package Net

import (
	"Cinder/Base/Message"
	"Cinder/Base/Security"
	"Cinder/Base/Util"
	"context"
	"errors"
	"io"
	"net"

	log "github.com/cihub/seelog"
)

type _TcpSess struct {
	conn net.Conn

	sendList Util.ISafeList

	data       interface{}
	isValidate bool

	ctx       context.Context
	ctxCancel context.CancelFunc

	sendCrypt Security.ICrypt
	recvCrypt Security.ICrypt

	trafficCount int //每秒流量计数
}

func NewTcpSess(conn net.Conn) *_TcpSess {
	sess := &_TcpSess{
		conn:     conn,
		sendList: Util.NewSafeList(),
	}

	tcpConn, ok := conn.(*net.TCPConn)
	if ok {
		_ = tcpConn.SetNoDelay(true)
	}

	sess.ctx, sess.ctxCancel = context.WithCancel(context.Background())

	go sess.loop()

	return sess
}

func (sess *_TcpSess) SetSendSecretKey(key []byte) {
	sess.sendCrypt = Security.NewCrypt(key)
}

func (sess *_TcpSess) SetRecvSecretKey(key []byte) {
	sess.recvCrypt = Security.NewCrypt(key)
}

func (sess *_TcpSess) Read() (Message.IMessage, uint32, error) {
	header := Util.Get(Message.HeadLen)
	defer func() {
		Util.Put(header)
	}()
	if _, err := io.ReadFull(sess.conn, header); err != nil {
		return nil, 0, err
	}
	l := int(header[0]) | int(header[1])<<8 | int(header[2])<<16

	body := Util.Get(l + Message.HeadLen)
	defer func() {
		Util.Put(body)
	}()
	if _, err := io.ReadFull(sess.conn, body[Message.HeadLen:Message.HeadLen+l]); err != nil {
		return nil, 0, errors.New("read from conn failed ")
	}
	copy(body[:Message.HeadLen], header)

	msg, msgNo, err := Message.UnpackWithSKey(body[:Message.HeadLen+l], sess.recvCrypt)
	if err != nil {
		return nil, 0, err
	}

	return msg, msgNo, nil
}

func (sess *_TcpSess) send(buf []byte) error {
	if sess.conn == nil {
		return errors.New("no connection")
	}

	if _, err := sess.conn.Write(buf); err != nil {
		return err
	}
	return nil
}

func (sess *_TcpSess) loop() {

loop:
	for {
		select {
		case <-sess.sendList.Signal():
			for {
				ii, err := sess.sendList.Pop()
				if err != nil {
					break
				}
				msgInfo := ii.([]byte)
				err = sess.send(msgInfo)
				if err != nil {
					log.Error("sess ", sess.GetData(), " send message err ", err.Error(), " RemoteAddr ", sess.conn.RemoteAddr().String())
					sess.close()
					break
				}
			}
		case <-sess.ctx.Done():
			log.Info("sess loop exit ", sess.GetData())
			break loop
		}
	}
}

func (sess *_TcpSess) Send(message Message.IMessage, msgNo uint32) error {

	if sess.conn == nil {
		return errors.New("no connection")
	}
	buf, err := Message.PackWithSKey(message, sess.sendCrypt, msgNo)
	if err != nil {
		return err
	}

	sess.sendList.Put(buf)
	return nil
}

func (sess *_TcpSess) Close() {
	sess.close()
}

func (sess *_TcpSess) close() {
	sess.ctxCancel()
	if sess.conn != nil {
		_ = sess.conn.Close()
	}
}

func (sess *_TcpSess) SetData(data interface{}) {
	sess.data = data
}

func (sess *_TcpSess) GetData() interface{} {
	return sess.data
}

func (sess *_TcpSess) SetValidate() {
	sess.isValidate = true
}

func (sess *_TcpSess) IsValidate() bool {
	return sess.isValidate
}
