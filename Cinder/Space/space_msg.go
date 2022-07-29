package Space

import (
	"Cinder/Base/Message"
	log "github.com/cihub/seelog"
)

func (space *Space) SendToAllClient(msg Message.IMessage) {
	if msg == nil {
		return
	}

	msgBuf, err := Message.Pack(msg)
	if err != nil {
		space.Error("SendToAllClient Pack err", err, "MessageID: ", msg.GetID())
		return
	}

	for agentID, userList := range space.userAgentMap {

		for i := 0; i < len(userList)/100+1; i++ {

			startIndex := i * 100
			endIndex := i*100 + 100
			if endIndex > len(userList) {
				endIndex = len(userList)
			}

			sendList := userList[startIndex:endIndex]

			if len(sendList) > 0 {

				broadcastMsg := &Message.SpaceBroadcastToClient{}

				broadcastMsg.UserList = sendList
				broadcastMsg.MsgData = msgBuf

				if err = Inst.Send(agentID, broadcastMsg); err != nil {
					log.Error("SendToAllClient Send broadcastMsg err ", err, " Target ", agentID)
				}
			}
		}
	}
}

func (space *Space) SendToAllClientExceptOne(msg Message.IMessage, exceptUserID string) {
	if msg == nil {
		return
	}

	msgBuf, err := Message.Pack(msg)
	if err != nil {
		space.Error("SendToAllClientExceptOne Pack err", err, "MessageID: ", msg.GetID())
		return
	}

	for agentID, userList := range space.userAgentMap {

		for i := 0; i < len(userList)/100+1; i++ {

			startIndex := i * 100
			endIndex := i*100 + 100
			if endIndex > len(userList) {
				endIndex = len(userList)
			}

			sendList := userList[startIndex:endIndex]

			if len(sendList) > 0 {

				broadcastMsg := &Message.SpaceBroadcastToClient{
					ExceptUserID: exceptUserID,
				}

				broadcastMsg.UserList = sendList
				broadcastMsg.MsgData = msgBuf

				if err = Inst.Send(agentID, broadcastMsg); err != nil {
					log.Error("SendToAllClientExceptOne Send broadcastMsg err ", err, " Target ", agentID)
				}
			}
		}
	}

}
