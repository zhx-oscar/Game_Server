package Cache

import (
	"fmt"
	"strings"
)

func SetUserLoginToken(userID string, token string, csk []byte, crk []byte) error {
	_, err := RedisDB.Set(getUserTokenKey(userID), combineUserTokenKey(token, csk, crk), LoginSessExpireTime).Result()
	return err
}

func ClearUserLoginToken(userID string) {
	RedisDB.Expire(getUserTokenKey(userID), LoginSessExpireTime)
}

func CancelUserLoginTokenExpire(userID string) {
	RedisDB.Persist(getUserTokenKey(userID))
}

func GetUserLoginToken(userID string) (string, []byte, []byte, error) {
	key, err := RedisDB.Get(getUserTokenKey(userID)).Result()
	if err != nil {
		return "", nil, nil, err
	}

	token, csk, crk := splitUserTokenKey(key)
	return token, csk, crk, nil
}

func getUserTokenKey(userID string) string {
	return fmt.Sprintf("UserToken:%s", userID)
}

func combineUserTokenKey(token string, csk []byte, crk []byte) string {
	return fmt.Sprintf("%s-%s-%s", token, csk, crk)
}

func splitUserTokenKey(key string) (string, []byte, []byte) {
	rets := strings.Split(key, "-")
	if len(rets) != 3 {
		panic("wrong user token format")
	}

	return rets[0], []byte(rets[1]), []byte(rets[2])
}
