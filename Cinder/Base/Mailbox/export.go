package Mailbox

import (
	"Cinder/Base/CRpc"
	"Cinder/Cache"
	"errors"

	"github.com/go-redis/redis/v7"
)

type IMailbox interface {
	Rpc(methodName string, args ...interface{}) chan *CRpc.RpcRet
}

// 邮箱类型
const (
	TypeUser    uint8 = 1
	TypePropObj uint8 = 2
)

var (
	ErrMailboxNameInvalid   = errors.New("mailbox name invalid")
	ErrBadDataInCache       = errors.New("bad data in cache")
	ErrGetMailboxTypeFailed = errors.New("get mailbox type failed")
	ErrGetMailboxInfoFailed = errors.New("get mailbox info failed")
	ErrUnknownMailboxType   = errors.New("unknown mailbox type")
)

var (
	mailboxKeyPrefix = "mailbox:"
	mailboxKeyType   = "type"
	mailboxKeyValue  = "value"
)

// Set 将Mailbox保存至redis, 方便其他人获取
func Set(name string, box IMailbox) error {
	if name == "" {
		return ErrMailboxNameInvalid
	}

	im := box.(iMailBoxCtrl)
	data, err := im.marshal()
	if err != nil {
		return err
	}

	key := mailboxKeyPrefix + ":" + name
	if _, err = Cache.RedisDB.HSet(key, mailboxKeyType, im.getType(), mailboxKeyValue, data).Result(); err != nil {
		return err
	}

	defaultMgr.Store(name, box)

	return nil
}

func Get(name string) (IMailbox, error) {
	if name == "" {
		return nil, ErrMailboxNameInvalid
	}

	if v := defaultMgr.Load(name); v != nil {
		return v, nil
	}

	key := mailboxKeyPrefix + ":" + name
	result, err := Cache.RedisDB.HMGet(key, mailboxKeyType, mailboxKeyValue).Result()
	if err != nil {
		return nil, err
	}

	if len(result) != 2 {
		return nil, ErrBadDataInCache
	}

	mailboxType, err := redis.NewCmdResult(result[0], nil).Int()
	if err != nil {
		return nil, ErrGetMailboxTypeFailed
	}
	data, ok := result[1].([]byte)
	if !ok {
		return nil, ErrGetMailboxInfoFailed
	}

	switch uint8(mailboxType) {
	case TypeUser:
		userMailBox := &_UserMailBox{}
		if err = userMailBox.unmarshal(data); err != nil {
			return nil, err
		}
		defaultMgr.Store(userMailBox.ID, userMailBox)
		return userMailBox, nil

	default:
		return nil, ErrUnknownMailboxType
	}
}
