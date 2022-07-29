package Cache

import (
	"Cinder/Base/Const"
	"errors"
	"fmt"
	"strings"

	log "github.com/cihub/seelog"
)

func getUserPeerRedisKey(userID string, srvType string) string {
	return fmt.Sprintf("UserPeer:%s:%s", userID, srvType)
}

func SetUserPeerRedis(userID string, srvType string, srvID string) error {
	_, err := RedisDB.Set(getUserPeerRedisKey(userID, srvType), srvID, UserPeerSrvIDExpireTime).Result()
	return err
}

func KeepAliveUserPeerSrvID(userID string, srvType string) {
	RedisDB.Expire(getUserPeerRedisKey(userID, srvType), UserPeerSrvIDExpireTime)
}

func GetUserPeerSrvID(userID string, srvType string) (string, error) {
	srvID, err := RedisDB.Get(getUserPeerRedisKey(userID, srvType)).Result()
	return srvID, err
}

func ClearUserPeerSrvID(userID string, srvType string, srvID string) {
	oSrvID, err := GetUserPeerSrvID(userID, srvType)
	if err != nil {
		return
	}
	if oSrvID == srvID {
		RedisDB.Del(getUserPeerRedisKey(userID, srvType))
	} else {
		log.Error("clear user peer srvID , but the data is wrong ", userID, srvType, oSrvID, srvID)
	}
}

func GetUserPeersSrvID(userID string) (map[string]string, error) {

	if RedisDB == nil {
		return nil, errors.New("no redis")
	}

	keys := []string{
		getUserPeerRedisKey(userID, Const.Game),
		getUserPeerRedisKey(userID, Const.Agent),
		getUserPeerRedisKey(userID, Const.Space),
	}

	vals, err := RedisDB.MGet(keys...).Result()
	if err != nil {
		return nil, err
	}

	ret := make(map[string]string)

	for i := 0; i < len(keys); i++ {

		a := strings.Split(keys[i], ":")
		if len(a) != 3 {
			return nil, errors.New("wrong key format")
		}

		srvType := a[2]
		if vals[i] != nil {
			ret[srvType] = vals[i].(string)
		}
	}

	return ret, nil
}
