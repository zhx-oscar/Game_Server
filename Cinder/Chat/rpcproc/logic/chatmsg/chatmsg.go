package chatmsg

import (
	"Cinder/Chat/chatapi"
	"Cinder/Chat/rpcproc/logic/types"
	"time"
)

func NewChatMessage(fromID types.UserID, fromNick string, fromData []byte, msgContent []byte) *chatapi.ChatMessage {
	return &chatapi.ChatMessage{
		From:       string(fromID),
		FromNick:   fromNick,
		FromData:   fromData,
		SendTime:   time.Now().Unix(),
		MsgContent: msgContent,
	}
}
