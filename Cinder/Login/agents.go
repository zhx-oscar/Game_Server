package main

import (
	"Cinder/Base/Config"
	"Cinder/Base/Const"
	"errors"
	"sync"

	log "github.com/cihub/seelog"
)

var agentKeys []string = make([]string, 0, 10)
var agentAddrs []string = make([]string, 0, 10)
var currentAgent int = 0
var mux sync.Mutex

func initAgentAssign() error {

	keys, values, err := Config.Inst.GetValuesByPrefix(Const.GetNetSrvIDbySrvType(Const.Agent))
	if err != nil {
		return err
	}

	mux.Lock()
	for i := 0; i < len(values); i++ {
		agentKeys = append(agentKeys, keys[i])
		agentAddrs = append(agentAddrs, values[i])
	}
	mux.Unlock()

	_, err = Config.Inst.WatchKeys(Const.GetNetSrvIDbySrvType(Const.Agent), agentChanged)
	if err != nil {
		return err
	}

	return nil
}

func agentChanged(opType int, key string, value string) {

	mux.Lock()
	defer mux.Unlock()

	if opType == Config.KeyAdd {
		addToAgentList(key, value)
	} else if opType == Config.KeyDelete {
		removeFromAgentList(key)
	}
}

func addToAgentList(key string, addr string) {

	var i int
	for i = 0; i < len(agentKeys); i++ {
		if agentKeys[i] == key {
			break
		}
	}

	if i >= len(agentKeys) {
		agentKeys = append(agentKeys, key)
		agentAddrs = append(agentAddrs, addr)
	} else {
		agentAddrs[i] = addr
	}

}

func removeFromAgentList(key string) {

	var i int
	for i = 0; i < len(agentKeys); i++ {
		if agentKeys[i] == key {
			break
		}
	}

	if i >= len(agentKeys) {
		log.Error("remove agent list wrong ", key, agentKeys)
	} else {
		agentKeys = append(agentKeys[0:i], agentKeys[i+1:]...)
		agentAddrs = append(agentAddrs[0:i], agentAddrs[i+1:]...)
	}

}

func GetAgent() (string, error) {

	mux.Lock()
	defer mux.Unlock()

	if len(agentAddrs) == 0 {
		return "", errors.New("no agent address")
	}

	addr := agentAddrs[currentAgent%len(agentAddrs)]
	currentAgent++

	return addr, nil
}
