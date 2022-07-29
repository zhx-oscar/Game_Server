package SrvNet

import (
	"Cinder/Base/Config"
	"Cinder/Base/Const"
	"Cinder/Base/MQNet"
	"Cinder/Base/MQNet/mqnats"
	"Cinder/Base/Message"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"math/rand"
	"sync"
	"time"

	log "github.com/cihub/seelog"
)

const (
	keepAliveDuration = 5 //暂定5秒
)

type loadBlanceInfo struct {
	SrvType      string
	LoadCurValue float32
	LoadMaxValue float32
}

type _Node struct {
	mqSrv MQNet.IService

	NodesByType sync.Map
	NodesByID   sync.Map

	areaID  string
	srvID   string
	srvType string

	watcherHandle int

	loadGetFunc     LoadGetFunc
	defaultLoadData sync.Map
}

func NewMQNode() INode {
	return &_Node{mqSrv: mqnats.New(), watcherHandle: -1}
}

func (n *_Node) Init(areaID string, srvID string, srvType string) error {

	n.areaID = areaID
	n.srvID = srvID
	n.srvType = srvType

	n.loadGetFunc = n.getDefaultLoadData

	if err := n.updateNodeInfo(); err != nil {
		return fmt.Errorf("register node: %w", err)
	}

	if err := n.collectNodeInfo(); err != nil {
		return fmt.Errorf("collect node info: %w", err)
	}

	serviceAddr := srvID
	boardcastAddr := srvType + "_" + areaID

	err := n.mqSrv.Init(MQNet.InitOptions(viper.GetString("NATS.Addr"), serviceAddr, boardcastAddr))
	if err != nil {
		return err
	}

	go n.keepalive()
	go n.serverLoadloop()

	return nil
}

func (n *_Node) SetLoadBlanceGetter(f LoadGetFunc) {
	n.loadGetFunc = f
}

//keepalive
func (n *_Node) keepalive() {
	//暂定 keepAliveDuration 秒处理一次
	ticker := time.NewTicker(keepAliveDuration * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			n.updateNodeInfo()
		}
	}
}

func (n *_Node) Destroy() {
	n.mqSrv.Destroy()
	Config.Inst.CancelWatch(n.watcherHandle)
}

func (n *_Node) GetID() string {
	return n.srvID
}

func (n *_Node) GetType() string {
	return n.srvType
}

func (n *_Node) Send(srvID string, msg Message.IMessage) error {
	return n.mqSrv.Post(srvID, msg)
}

func (n *_Node) Broadcast(srvType string, msg Message.IMessage) error {
	return n.mqSrv.Post(srvType+"_"+n.areaID, msg)
}

func (n *_Node) AddMessageProc(proc MQNet.IProc) {
	n.mqSrv.AddProc(proc)
}

func (n *_Node) updateNodeInfo() error {
	info := loadBlanceInfo{
		SrvType:      n.srvType,
		LoadCurValue: n.loadGetFunc(),
		LoadMaxValue: float32(viper.GetFloat64("LoadBlance.LimitValue")),
	}

	data, err := json.Marshal(info)
	if err != nil {
		return err
	}

	return Config.Inst.SetValueAndOvertime(Const.GetMQSrvID(n.srvID), string(data), keepAliveDuration*2)
}

func (n *_Node) collectNodeInfo() error {

	addrPrefix := Const.GetMQSrvPrefix()
	var addrs []string
	var keys []string
	var err error
	if keys, addrs, err = Config.Inst.GetValuesByPrefix(addrPrefix); err != nil {

		log.Error("get mq service addr failed " + addrPrefix)
		//continue
	}

	if len(addrs) != len(keys) {
		return errors.New("couldn't happen")
	}

	for i := 0; i < len(addrs); i++ {

		err = n.addNode(keys[i], addrs[i])
		if err != nil {
			log.Debug("add mq node failed ", err)
		}
	}

	if n.watcherHandle, err = Config.Inst.WatchKeys(addrPrefix, n.watchNode); err != nil {
		log.Error("watch mq service failed")
		return err
	}

	return nil
}

func (n *_Node) addNode(rid string, value string) error {
	srvID, err := Const.GetMQInfoByRID(rid)
	if err != nil {
		return err
	}

	info := &loadBlanceInfo{}
	err = json.Unmarshal([]byte(value), info)
	if err != nil {
		return err
	}

	//node 如果已经存在,则只更新NodeInfo
	if _, ok := n.NodesByID.Load(srvID); ok {
		n.NodesByID.Store(srvID, info)
		return nil
	}

	n.NodesByID.Store(srvID, info)
	nl, _ := n.NodesByType.LoadOrStore(info.SrvType, newNodeList())

	inl := nl.(*_NodeList)
	inl.addNode(srvID)

	log.Debug("add mq node succeed  ", rid, value)
	return nil
}

func (n *_Node) deleteNode(rid string) {
	srvID, err := Const.GetMQInfoByRID(rid)
	if err != nil {
		return
	}

	v, ok := n.NodesByID.Load(srvID)
	if !ok {
		return
	}

	n.NodesByID.Delete(srvID)

	info := v.(*loadBlanceInfo)
	v, ok = n.NodesByType.Load(info.SrvType)
	if !ok {
		return
	}

	nl := v.(*_NodeList)
	nl.removeNode(srvID)
}

func (n *_Node) watchNode(opType int, rid string, addr string) {
	if opType == Config.KeyAdd {
		err := n.addNode(rid, addr)
		if err != nil {
			log.Debug("watch mq client but add client failed ", err, rid)
		}
	} else if opType == Config.KeyDelete {
		n.deleteNode(rid)
		log.Debug("watch mq client and remove client ", rid)
	}
}

func (n *_Node) GetSrvIDSByType(srvType string) ([]string, error) {

	v, ok := n.NodesByType.Load(srvType)
	if !ok {
		return nil, errors.New("no srvType node " + srvType)
	}

	nl := v.(*_NodeList)

	return nl.nodeList, nil
}

func (n *_Node) GetSrvIDByType(srvType string) (string, error) {

	v, ok := n.NodesByType.Load(srvType)
	if !ok {
		return "", errors.New("no srvType node " + srvType)
	}

	nl := v.(*_NodeList)

	return nl.getNode(n)
}

func (n *_Node) GetSrvTypeByID(srvID string) (string, error) {

	v, ok := n.NodesByID.Load(srvID)

	if !ok {
		return "", errors.New("no exist ")
	}

	info := v.(*loadBlanceInfo)

	return info.SrvType, nil
}

type _NodeList struct {
	nodeList []string
	pos      int
	mutex    sync.RWMutex
}

func newNodeList() *_NodeList {
	return &_NodeList{
		nodeList: make([]string, 0, 10),
		pos:      -1,
	}
}

func (nl *_NodeList) removeNode(srvID string) {
	nl.mutex.Lock()
	defer nl.mutex.Unlock()

	var i int
	for i = 0; i < len(nl.nodeList); i++ {
		if nl.nodeList[i] == srvID {
			break
		}
	}

	if i < len(nl.nodeList) && len(nl.nodeList) > 0 {
		nl.nodeList = append(nl.nodeList[:i], nl.nodeList[i+1:]...)
	}

	if len(nl.nodeList) == 0 {
		nl.pos = -1
	} else {
		nl.pos = rand.Intn(len(nl.nodeList))
	}

}

func (nl *_NodeList) addNode(srvID string) {
	nl.mutex.Lock()
	defer nl.mutex.Unlock()

	nl.nodeList = append(nl.nodeList, srvID)
	if nl.pos == -1 {
		nl.pos = 0
	}
}

func (nl *_NodeList) getNode(n *_Node) (string, error) {
	if nl.pos == -1 {
		return "", errors.New("no node")
	}

	nl.mutex.RLock()
	defer nl.mutex.RUnlock()

	var tempNodeList []string
	for _, id := range nl.nodeList {
		v, ok := n.NodesByID.Load(id)
		if !ok {
			continue
		}

		//负载检测
		if n.isOverloadLimit(v.(*loadBlanceInfo)) {
			continue
		}

		tempNodeList = append(tempNodeList, id)
	}

	if len(tempNodeList) == 0 {
		return "", errors.New("no useable node")
	}

	o := nl.pos
	nl.pos = (nl.pos + 1) % len(tempNodeList)

	return tempNodeList[o], nil
}

//isOverloadLimit 负载检测
func (n *_Node) isOverloadLimit(info *loadBlanceInfo) bool {
	if info == nil {
		return false
	}

	//负载开关
	onoff := viper.GetBool("LoadBlance.OnOff")
	if !onoff {
		return false
	}

	if info.LoadCurValue > info.LoadMaxValue {
		return true
	}

	return false
}
