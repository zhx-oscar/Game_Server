package MQNet

import "errors"

var (
	ErrorNotInit = errors.New("no init")

	ErrInvalidPostAddr = errors.New("invalid post addr")
	ErrInvalidMessage  = errors.New("invalid message")

	ErrMessageHeaderFormatInvalid = errors.New("message header format error")
	ErrMessageEndFormatInvalid    = errors.New("message end format error")
)
