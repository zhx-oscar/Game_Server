package mqnsq

import (
	"Cinder/Base/MQNet"
	"errors"
)

const (
	nsqLookupKey = "nsqlookup_addr"
	nsqAdminKey  = "nsqadmin_addr"
)

var (
	ErrLookupAddrInvalid = errors.New("nsqlookup addr invalid")
	ErrAdminAddrInvalid  = errors.New("admin addr invalid")
)

func NSQLookup(addr string) MQNet.Option {
	return func(options *MQNet.Options) {
		options.ExtOpts[nsqLookupKey] = addr
	}
}

func NSQAdmin(addr string) MQNet.Option {
	return func(options *MQNet.Options) {
		options.ExtOpts[nsqAdminKey] = addr
	}
}
