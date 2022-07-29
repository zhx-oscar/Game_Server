package Cache

import (
	"errors"
	"fmt"
)

func getPropObjectKey(typ, objID string) string {
	return fmt.Sprintf("PropObject:%s:%s", typ, objID)
}

func SetPropObjectSrvID(typ, objID string, srvID string) error {
	b, err := RedisDB.SetNX(getPropObjectKey(typ, objID), srvID, PropObjectSrvIDExpireTime).Result()

	if err != nil {
		return err
	}

	if !b {
		return errors.New("object have existed ")
	}

	return nil
}

func ClearPropObjectSrvID(typ, objID string) {
	_ = RedisDB.Del(getPropObjectKey(typ, objID))
}

func GetPropObjectSrvID(typ, objID string) (string, error) {

	ret, err := RedisDB.Get(getPropObjectKey(typ, objID)).Result()

	if err != nil {
		return "", err
	}

	return ret, nil
}

func KeepAlivePropObjectSrvID(typ, objID string) {
	RedisDB.Expire(getPropObjectKey(typ, objID), PropObjectSrvIDExpireTime)
}
