package Space

/*

type _NetServerProc struct {
}

func (proc *_NetServerProc) OnSessConnected(sess Net.ISess) {

}

func (proc *_NetServerProc) OnSessClosed(sess Net.ISess) {

	data := sess.GetData()
	if data == nil {
		return
	}

	id, ok := data.(string)
	if !ok {
		return
	}

	sessMgr.Remove(id)

	log.Debug("agent sess remove ", id)
}

func (proc *_NetServerProc) OnSessMessageHandle(sess Net.ISess, message Message.IMessage) {
	switch message.GetID() {
	case Message.ID_Agent_Validate_Req:
		msg := message.(*Message.AgentValidateReq)
		proc.onAgentValidate(sess, msg)
	case Message.ID_Client_Forward_To_Space:
		proc.onClientForwardToSpace(sess, message.(*Message.ClientForwardToSpace))
	default:
		log.Debug("not handle message ", message.GetID())
	}
}

func (proc *_NetServerProc) onAgentValidate(sess Net.ISess, msg *Message.AgentValidateReq) {

	agentID := msg.ID
	sess.SetData(agentID)

	err := sessMgr.Add(agentID, sess)
	if err != nil {
		log.Debug("agent sess mgr add sess failed ", err)
		return
	}

	_ = sess.Send(&Message.AgentValidateRet{
		SpaceSrvID: Inst.GetServiceID(),
	})

	log.Debug("agent sess add sess ", agentID, err)
}

func (proc *_NetServerProc) onClientForwardToSpace(sess Net.ISess, msg *Message.ClientForwardToSpace) {

	space, err := Inst.GetSpace(msg.SpaceID)
	if err != nil {
		log.Debug("client forward to space , but couldn't found space ", msg.SpaceID, err)
		return
	}

	user, err := space.GetUser(msg.userID)
	if err != nil {
		log.Debug("client forward to space , but couldn't found space user ", msg.SpaceID, msg.userID, err)
		return
	}

	innerMsg, err := Message.Decode(msg.MsgID, msg.MsgFlag, msg.MsgData)
	if err != nil {
		log.Debug("client forward to space , but decode inner message failed ", err)
		return
	}

	mp, ok := user.(_IBMsgProc)
	if ok {
		mp.MsgProc(innerMsg)
	}
}
*/
