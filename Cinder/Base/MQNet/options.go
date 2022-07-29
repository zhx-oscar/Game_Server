package MQNet

type Options struct {
	Addr string // MQ服务地址

	ServiceAddr   string // 即 srvID
	BoardcastAddr string // 即 srvType_areaID

	ExtOpts map[string]interface{}
}

type Option func(*Options)

func InitOptions(addr, service, boardcast string) Option {
	return func(options *Options) {
		options.Addr = addr
		options.ServiceAddr = service
		options.BoardcastAddr = boardcast
	}
}
