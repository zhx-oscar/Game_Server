package Mailbox

import (
	"Cinder/Base/Message"
	"errors"
	log "github.com/cihub/seelog"
	"sync"
)

type Mgr struct {
	mailboxs sync.Map
}

func GetDefaultMgr() *Mgr {
	return defaultMgr
}

var defaultMgr *Mgr

func init() {
	defaultMgr = &Mgr{}
}

func (mgr *Mgr) MessageProc(srcAddr string, message Message.IMessage) {
	if message == nil {
		return
	}

	if message.GetID() != Message.ID_Mailbox_Ret {
		return
	}

	m := message.(*Message.MailboxRet)
	v, ok := mgr.mailboxs.Load(m.MailboxID)
	if !ok {
		log.Error("Mailbox MessageProc can't find mailbox ", m.MailboxID)
		return
	}

	args, err := Message.UnPackArgs(m.Ret)
	if err != nil {
		log.Error("Mailbox MessageProc UnPackArgs err ", err)
		m.Err = m.Err + " " + err.Error()
	}

	v.(iMailBoxCtrl).onMailboxRet(m.RetID, errors.New(m.Err), args)
}

func (mgr *Mgr) Store(id string, mb IMailbox) {
	mgr.mailboxs.Store(id, mb)
}

func (mgr *Mgr) Delete(id string) {
	mgr.mailboxs.Delete(id)
}

func (mgr *Mgr) Load(id string) IMailbox {
	v, ok := mgr.mailboxs.Load(id)
	if !ok {
		return nil
	}
	return v.(IMailbox)
}
