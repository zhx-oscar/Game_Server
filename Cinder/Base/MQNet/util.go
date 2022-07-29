package MQNet

import (
	"Cinder/Base/Message"
)

type SendInfo struct {
	Addr    string
	Message Message.IMessage
}

func MaxMessageSize(addr string, msg Message.IMessage) (int, error) {
	size, err := Message.GetMaxMessageEncodedSize(msg)
	if err != nil {
		return 0, err
	}

	return size + len(addr) + 2, nil
}

func Unpack(data []byte) (string, Message.IMessage, error) {
	if len(data) < 2 {
		return "", nil, ErrInvalidMessage
	}

	if data[0] != '|' {
		return "", nil, ErrMessageHeaderFormatInvalid
	}

	var offset int
	for offset = 1; offset < len(data); offset++ {
		if data[offset] == '|' {
			break
		}
	}
	addr := string(data[1:offset])

	msg, err := Message.Unpack(data[offset+1:])
	if err != nil {
		return "", nil, err
	}

	return addr, msg, nil
}

func Pack(addr string, message Message.IMessage, buf []byte) ([]byte, error) {
	headLen := len(addr) + 2
	data, err := Message.PackWithBuf(message, buf[headLen:])
	if err != nil {
		return nil, err
	}

	buf[0] = '|'
	copy(buf[1:], addr)
	buf[len(addr)+1] = '|'
	return buf[:len(data)+headLen], nil
}
