package Cache

import (
	"errors"
	"time"
)

// Expire time
const (
	UserPeerSrvIDExpireTime   = 10 * time.Second
	PropObjectSrvIDExpireTime = 10 * time.Second
	PropDBSrvIDExpireTime     = 10 * time.Second
	LoginSessExpireTime       = 30 * time.Second
)

var (
	ErrInvalidParam = errors.New("invalid param")
)
