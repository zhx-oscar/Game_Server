package Prop

import (
	"Cinder/Base/Message"
	"errors"
	log "github.com/cihub/seelog"
)

func (mgr *_Mgr) MessageProc(srcAddr string, message Message.IMessage) {
	switch message.GetID() {
	case Message.ID_Prop_Data_Ret:
		msg := message.(*Message.PropDataRet)
		mgr.pushDataToPendingRetC(msg.TempID, msg.Data)

	case Message.ID_Prop_Data_Flush_Ret:
		msg := message.(*Message.PropDataFlushRet)
		mgr.pushDataToPendingRetC(msg.TempID, msg.Err)

	case Message.ID_Prop_Cache_Flush_Ret:
		msg := message.(*Message.PropCacheFlushRet)
		mgr.pushDataToPendingRetC(msg.TempID, msg.Err)

	case Message.ID_PropObject_Open_Req:
		mgr.onPropObjectOpenReq(srcAddr, message.(*Message.PropObjectOpenReq))

	case Message.ID_PropObject_Close_Req:
		mgr.onPropObjectCloseReq(srcAddr, message.(*Message.PropObjectCloseReq))

	}
}

func (mgr *_Mgr) pushDataToPendingRetC(tempID string, data interface{}) {
	ii, ok := mgr.pendingRetC.Load(tempID)
	if !ok {
		log.Debug("prop data Flush ret error , no id ", tempID)
		return
	}

	retC, ok := ii.(chan interface{})
	if ok {
		retC <- data
		close(retC)
	}
}

func (mgr *_Mgr) fetchRetC(tempID string) (<-chan interface{}, error) {

	_, ok := mgr.pendingRetC.Load(tempID)
	if ok {
		log.Error("it couldn't happen")
		return nil, errors.New("the channel had existed ")
	}
	retC := make(chan interface{}, 1)
	mgr.pendingRetC.Store(tempID, retC)
	return retC, nil
}

func (mgr *_Mgr) removeRetC(tempID string) {
	mgr.pendingRetC.Delete(tempID)
}

func (mgr *_Mgr) onPropObjectOpenReq(srvAddr string, msg *Message.PropObjectOpenReq) {

	ii, err := mgr.GetPropObject(msg.ID)
	if err != nil {
		log.Debug("prop object not exist ", msg.ID)
		return
	}

	obj := ii.(_IPropObject)

	var propData []byte

	r := <-obj.SafeCall("GetPropData")
	if r.Err == nil {
		propData = r.Ret[0].([]byte)
	}

	retMsg := &Message.PropObjectOpenRet{
		ID:       msg.ID,
		UserID:   msg.UserID,
		SrvID:    mgr.srvNode.GetID(),
		PropType: obj.GetPropType(),
		PropData: propData,
	}

	if err = mgr.srvNode.Send(srvAddr, retMsg); err != nil {
		log.Error("onPropObjectOpenReq Send retMsg err ", err, " retMsg ", retMsg)
		return
	}

	<-obj.SafeCall("AddWatcher", msg.UserID)
}

func (mgr *_Mgr) onPropObjectCloseReq(srvAddr string, msg *Message.PropObjectCloseReq) {

	ii, err := mgr.GetPropObject(msg.ID)
	if err != nil {
		log.Debug("prop object not exist ", msg.ID)
		return
	}

	obj := ii.(_IPropObject)

	<-obj.SafeCall("RemoveWatcher", msg.UserID)

	retMsg := &Message.PropObjectCloseRet{
		ID:     msg.ID,
		UserID: msg.UserID,
	}

	if err = mgr.srvNode.Send(srvAddr, retMsg); err != nil {
		log.Error("onPropObjectCloseReq Send retMsg err ", err, " retMsg ", retMsg)
	}
}
