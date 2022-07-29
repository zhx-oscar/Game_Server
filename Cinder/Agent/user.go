package main

import (
	"Cinder/Base/Message"
	"Cinder/Base/Net"
	"Cinder/Base/User"
	"errors"
	"sync"
	"time"
)

type _User struct {
	User.User
	clientSess Net.ISess

	recvMsgNo uint32
	sendMsgNo uint32

	sendMsgPool []*_MsgInfo
	sendMsgLock sync.Mutex
}

type _MsgInfo struct {
	msg       Message.IMessage
	msgNo     uint32
	timeStamp time.Time
}

func (user *_User) SetClientSess(sess Net.ISess) {
	if user.clientSess != nil {
		user.clientSess.SetData(nil)
		user.clientSess.Close()
	}

	user.recvMsgNo = 0
	user.clientSess = sess
}

func (user *_User) GetClientSess() Net.ISess {
	return user.clientSess
}

func (user *_User) Init() {
	user.sendMsgPool = make([]*_MsgInfo, 0, 100)
	user.sendMsgNo = 0
	user.recvMsgNo = 0

	user.Info("AgentUser Init")
}

func (user *_User) Destroy() {
	if user.clientSess != nil {
		user.clientSess.SetData(nil)
		user.clientSess.Close()
		user.clientSess = nil
	}

	user.Info("AgentUser Destroy")
}

func (user *_User) Loop() {}

var errRecvNoInvalid = errors.New("recv msgno invalid")

func (user *_User) CheckRecvMsgValidate(msgNo uint32) error {
	user.recvMsgNo++

	if user.recvMsgNo != msgNo {
		return errRecvNoInvalid
	}

	return nil
}

func (user *_User) IsMsgPoolEnough(msgNo uint32) bool {
	user.sendMsgLock.Lock()
	defer user.sendMsgLock.Unlock()

	if len(user.sendMsgPool) == 0 {
		return false
	}

	return user.sendMsgPool[0].msgNo <= msgNo+1
}

func (user *_User) ReSendMsgFromPool(msgNo uint32, sess Net.ISess) {
	user.sendMsgLock.Lock()
	defer user.sendMsgLock.Unlock()

	for _, v := range user.sendMsgPool {
		if v.msgNo > msgNo {
			sess.Send(v.msg, v.msgNo)
		}
	}
}

func (user *_User) ResetSendNo() {
	user.sendMsgNo = 0
	user.sendMsgPool = user.sendMsgPool[0:0]
}

func (user *_User) SendMessageToClient(msg Message.IMessage) error {
	user.sendMsgLock.Lock()
	defer user.sendMsgLock.Unlock()

	user.sendMsgNo++
	sess := user.clientSess

	if sess != nil {
		_ = sess.Send(msg, user.sendMsgNo)
	}
	user.putMsgToSendPool(msg, user.sendMsgNo)

	return nil
}

func (user *_User) putMsgToSendPool(msg Message.IMessage, msgNo uint32) {
	msgInfo := &_MsgInfo{
		msg:       msg,
		msgNo:     msgNo,
		timeStamp: time.Now(),
	}

	user.sendMsgPool = append(user.sendMsgPool, msgInfo)

	var ek int
	for k, v := range user.sendMsgPool {
		if time.Now().Sub(v.timeStamp) < UserClientMessageTimeout {
			ek = k
			break
		}
	}

	if ek > 0 {
		user.sendMsgPool = user.sendMsgPool[ek:]
	}
}

func (user *_User) GetPropInfo() (string, string) {
	return "", ""
}
