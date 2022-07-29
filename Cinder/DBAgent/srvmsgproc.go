package DBAgent

import (
	"Cinder/Base/Message"

	log "github.com/cihub/seelog"
)

type _SrvMsgProc struct {
}

func (proc *_SrvMsgProc) MessageProc(srcAddr string, message Message.IMessage) {

	switch message.GetID() {
	case Message.ID_Prop_Data_Req:
		go proc.onPropDataReq(srcAddr, message.(*Message.PropDataReq))
	case Message.ID_Prop_Notify:
		proc.onPropNotify(srcAddr, message.(*Message.PropNotify))
	case Message.ID_Prop_Data_Flush_Req:
		proc.onPropDataFlush(srcAddr, message.(*Message.PropDataFlushReq))
	case Message.ID_Prop_Cache_Flush_Req:
		proc.onPropCacheFlush(srcAddr, message.(*Message.PropCacheFlushReq))
	}

}

func (proc *_SrvMsgProc) onPropDataReq(srvAddr string, msg *Message.PropDataReq) {

	prop, err := propMgr.FetchOrCreate(msg.ID, msg.PropType)

	retMsg := &Message.PropDataRet{
		TempID: msg.TempID,
		ID:     msg.ID,
		Data:   []byte{},
		Err:    "",
	}
	if err == nil {
		prop.SetCurrentOwner(msg.ServiceType, srvAddr)

		if prop.typ != msg.PropType {
			retMsg.Err = "wrong prop tye"
		} else {
			ret := <-prop.SafeCall("GetPropData")
			if ret.Err != nil {
				retMsg.Err = ret.Err.Error()
			} else {
				if ret.Ret[1] == nil {
					retMsg.Data = ret.Ret[0].([]byte)
				} else {
					retMsg.Err = ret.Ret[1].(error).Error()
				}
			}

		}
	} else {
		retMsg.Err = err.Error()
	}

	if err = Inst.Send(srvAddr, retMsg); err != nil {
		log.Error("onPropDataReq Send retMsg err ", err, " retMsg ", retMsg)
	}
}

func (proc *_SrvMsgProc) onPropNotify(srvAddr string, msg *Message.PropNotify) {

	prop, err := propMgr.Get(msg.ID, msg.Type)
	if err != nil {
		log.Error("onPropNotify Get prop err ", err, " Type: ", msg.Type, " ID: ", msg.ID)
		return
	}

	if !prop.VerifyOwner(msg.ServiceType, srvAddr) {
		log.Errorf("onPropNotify verify owner failed! Type: %s ID: %s", msg.Type, msg.ID)
		return
	}

	if prop.GetProp() != nil {
		args, _ := Message.UnPackArgs(msg.Args)
		prop.GetProp().GetCaller().SafeCall(msg.MethodName, args...)
		prop.SafeCall("Touch")
	}
}

func (proc *_SrvMsgProc) onPropDataFlush(srvAddr string, msg *Message.PropDataFlushReq) {

	retMsg := &Message.PropDataFlushRet{
		TempID: msg.TempID,
		Err:    "",
	}

	prop, err := propMgr.Get(msg.ID, msg.Type)
	if err != nil {
		retMsg.Err = "onPropDataFlush Get prop err " + err.Error() + " Type: " + msg.Type + " ID: " + msg.ID
		log.Error(retMsg.Err)
		Inst.Send(srvAddr, retMsg)
		return
	}

	if !prop.VerifyOwner(msg.ServiceType, srvAddr) {
		log.Errorf("onPropDataFlush verify owner failed! Type: %s ID: %s", msg.Type, msg.ID)
		return
	}

	retC := prop.SafeCall("WriteToDB")

	go func() {
		ret := <-retC

		if ret.Err != nil {
			retMsg.Err = ret.Err.Error()
		} else {
			if ret.Ret[0] != nil {
				retMsg.Err = ret.Ret[0].(error).Error()
			}
		}

		if err = Inst.Send(srvAddr, retMsg); err != nil {
			log.Error("onPropDataFlush Send retMsg err ", err, " retMsg ", retMsg)
		}
	}()
}

func (proc *_SrvMsgProc) onPropCacheFlush(srvAddr string, msg *Message.PropCacheFlushReq) {

	retMsg := &Message.PropCacheFlushRet{
		TempID: msg.TempID,
		Err:    "",
	}

	prop, err := propMgr.Get(msg.ID, msg.Type)
	if err != nil {
		retMsg.Err = "onPropCacheFlush Get prop err " + err.Error() + " Type: " + msg.Type + " ID: " + msg.ID
		log.Error(retMsg.Err)
		Inst.Send(srvAddr, retMsg)
		return
	}

	if !prop.VerifyOwner(msg.ServiceType, srvAddr) {
		log.Errorf("onPropCacheFlush verify owner failed! Type: %s ID: %s", msg.Type, msg.ID)
		return
	}

	retC := prop.SafeCall("WriteToCache")

	go func() {
		ret := <-retC

		if ret.Err != nil {
			retMsg.Err = ret.Err.Error()
		} else {
			if ret.Ret[0] != nil {
				retMsg.Err = ret.Ret[0].(error).Error()
			}
		}

		if err = Inst.Send(srvAddr, retMsg); err != nil {
			log.Error("onPropCacheFlush Send retMsg err ", err, " Msg ", retMsg)
		}
	}()
}
