package MQNet

import (
	"Cinder/Base/Message"
	"reflect"
	"testing"
)

func TestPack(t *testing.T) {
	msg := &Message.ClientValidateReq{
		Version: 2,
		ID:      "Hello",
		Token:   "World",
		MsgSNo:  33,
	}

	maxLen, err := MaxMessageSize("agent", msg)
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, maxLen)

	data, err := Pack("agent", msg, buf)
	if err != nil {
		t.Fatal(err)
	}

	addr, newMessage, err := Unpack(data)
	if err != nil {
		t.Fatal(err)
	}
	if addr != "agent" {
		t.Fatal("Addr mismatch")
	}
	if !reflect.DeepEqual(msg, newMessage) {
		t.Fatal("mismatch")
	}
}
