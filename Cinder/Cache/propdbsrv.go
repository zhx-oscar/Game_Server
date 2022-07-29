package Cache

import "fmt"

func GetOrSetPropDBSrvID(propType string, propID string, dbSrvID string) (string, error) {
	isSet, err := RedisDB.SetNX(getPropDBSrvIDKey(propType, propID), dbSrvID, PropDBSrvIDExpireTime).Result()
	if err != nil {
		return "", err
	}
	if isSet {
		return dbSrvID, nil
	}

	var existedSrvID string
	existedSrvID, err = RedisDB.Get(getPropDBSrvIDKey(propType, propID)).Result()
	if err != nil {
		return "", err
	}
	return existedSrvID, nil
}

func KeepAlivePropDBSrvID(propType string, propID string) {
	RedisDB.Expire(getPropDBSrvIDKey(propType, propID), PropDBSrvIDExpireTime)
}

func getPropDBSrvIDKey(propType string, propID string) string {
	return fmt.Sprintf("PropDB:%s:%s", propType, propID)
}
