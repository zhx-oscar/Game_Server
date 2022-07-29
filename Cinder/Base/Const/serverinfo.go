package Const

import (
	"errors"
	"fmt"
	"strings"
)

// service Info
// Key
// Rpc/GameService/1/10101001
// Net/GameService/1/10101001
// format GameService/ServiceType/ServiceID
//

const (
	GameServicePrefix = "Service"
	//	GameServiceRpcPrefix = "Rpc"
	GameServiceNetPrefix = "Net"
	GameServiceMQPrefix  = "Mq"
)

func GetNetSrvID(serviceType string, serviceID string) string {
	return fmt.Sprintf("%s/%s/%s/%s", GameServiceNetPrefix, GameServicePrefix, serviceType, serviceID)
}

func GetNetSrvIDbySrvType(serviceType string) string {
	return fmt.Sprintf("%s/%s/%s/", GameServiceNetPrefix, GameServicePrefix, serviceType)
}

func GetMQSrvID(serviceID string) string {
	return fmt.Sprintf("%s/%s/%s", GameServiceMQPrefix, GameServicePrefix, serviceID)
}

func GetMQSrvPrefix() string {
	return fmt.Sprintf("%s/%s", GameServiceMQPrefix, GameServicePrefix)
}

func GetMQInfoByRID(rid string) (string, error) {
	addrs := strings.Split(rid, "/")

	if len(addrs) != 3 {
		return "", errors.New("wrong format RID " + rid)
	}
	var serviceID string

	serviceID = addrs[2]

	return serviceID, nil
}
