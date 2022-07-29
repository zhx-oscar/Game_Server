package StressClient

import (
	"Cinder/Base/Const"
	"Cinder/Base/Message"
	"Cinder/Base/Net"
	"Cinder/Base/Security"
	"Cinder/Base/Util"
	"context"
	"fmt"
	log "github.com/cihub/seelog"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type _IClient interface {
	IClient
	init(realPtr interface{}, addr string, userName string, password string, userData []byte)
	setDelegateToLoginSucceed(del func())
	setDelegateToLoginFailed(del func(string))
	Login()
	Logout()
}

type Client struct {
	userName  string
	password  string
	loginAddr string

	userID    string
	agentAddr string
	token     string

	realPtr interface{}

	loginSucceedDel func()
	LoginFailedDel  func(string)

	netMessageEvt []func(string, ...interface{})

	sess       Net.ISess
	ctx        context.Context
	cancelFunc context.CancelFunc
	msgList    Util.ISafeList
	closeC     chan struct{}

	iInit    IInit
	iDestroy IDestroy
	iLoop    ILoop

	lastHeartbeatTime time.Time
	ttl               time.Duration
	serverTime        time.Time

	sendMsgNo uint32
	recvMsgNo uint32

	cspr *big.Int
	crpr *big.Int

	csk []byte
	crk []byte
}

func (c *Client) init(realPtr interface{}, addr string, userName string, password string, userData []byte) {
	c.realPtr = realPtr
	c.loginAddr = addr
	c.userName = userName
	c.password = password
}

func (c *Client) setDelegateToLoginSucceed(del func()) {
	c.loginSucceedDel = del
}

func (c *Client) setDelegateToLoginFailed(del func(string)) {
	c.LoginFailedDel = del
}

func (c *Client) AddDelegateToNetMessage(del func(string, ...interface{})) {
	if c.netMessageEvt == nil {
		c.netMessageEvt = make([]func(string, ...interface{}), 0, 5)
	}

	c.netMessageEvt = append(c.netMessageEvt, del)
}

func (c *Client) onLoginSucceed() {
	log.Debug("login succeed ", c.GetUserName())

	if c.loginSucceedDel != nil {
		c.loginSucceedDel()
	}
}

func (c *Client) onLoginFailed(err string) {

	log.Debug("login failed ", c.GetUserName(), "  ", err)

	if c.LoginFailedDel != nil {
		c.LoginFailedDel(err)
	}
}

func (c *Client) onNetMessage(methodName string, args []interface{}) {
	if c.netMessageEvt != nil {
		for _, v := range c.netMessageEvt {
			v(methodName, args...)
		}
	}
}

func (c *Client) setRealPtr(realPtr interface{}) {
	c.realPtr = realPtr
}

func (c *Client) onInit() {
	c.ctx, c.cancelFunc = context.WithCancel(context.Background())
	c.msgList = Util.NewSafeList()
	c.closeC = make(chan struct{}, 1)

	ii, ok := c.realPtr.(IInit)
	if ok {
		c.iInit = ii
	}

	id, ok := c.realPtr.(IDestroy)
	if ok {
		c.iDestroy = id
	}

	ll, ok := c.realPtr.(ILoop)
	if ok {
		c.iLoop = ll
	}

	if c.iInit != nil {
		c.iInit.Init()
	}
}

func (c *Client) onDestroy() {

	if c.iDestroy != nil {
		c.iDestroy.Destroy()
	}

	c.netMessageEvt = nil
}

func (c *Client) onLoop() {
	if c.iLoop != nil {
		c.iLoop.Loop()
	}
}

func (c *Client) GetID() string {
	return c.userID
}

func (c *Client) GetUserName() string {
	return c.userName
}

func (c *Client) Login() {

	cspu, crpu := c.genSecretKey()
	resp, err := http.PostForm(c.loginAddr+"/Login", url.Values{"AccountName": {c.userName}, "Password": {c.password}, "cspu": {cspu.String()}, "crpu": {crpu.String()}})
	if err != nil {
		c.onLoginFailed(err.Error())
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.onLoginFailed(err.Error())
		return
	}
	var result string
	result = string(body)

	index := strings.Index(result, "ok")
	if index < 0 {
		c.onLoginFailed(fmt.Sprint("Http auth fail", result))
		return
	}

	r := strings.Split(result, "-")

	c.userID = r[1]
	c.agentAddr = r[3]
	c.token = r[2]

	cspum, _ := big.NewInt(0).SetString(r[4], 0)
	crpum, _ := big.NewInt(0).SetString(r[5], 0)

	c.csk = []byte(Security.Key(c.cspr, cspum).String())
	c.crk = []byte(Security.Key(c.crpr, crpum).String())

	//log.Debug("csk: ", c.csk)
	//log.Debug("crk: ", c.crk)

	log.Debug("user ", c.userName, "  connect to login succeed , agent addr ", c.agentAddr)

	c.connectToAgent()
}

func (c *Client) genSecretKey() (*big.Int, *big.Int) {
	cspr, cspu := Security.Pair()
	crpr, crpu := Security.Pair()

	c.cspr = cspr
	c.crpr = crpr
	return cspu, crpu
}

func (c *Client) connectToAgent() {
	c.sendMsgNo = 0

	conn, err := net.Dial("tcp", c.agentAddr)
	if err != nil {
		c.onLoginFailed(fmt.Sprint("User ", c.userName, "connect to agent failed ", c.agentAddr))
		return
	}

	c.sess = Net.NewTcpSess(conn)

	vm := &Message.ClientValidateReq{Version: 0, ID: c.userID, Token: c.token}
	c.sess.Send(vm, 0)

	c.sess.SetSendSecretKey(c.csk)
	c.sess.SetRecvSecretKey(c.crk)

	msg, _, err := c.sess.Read()
	if err != nil {
		c.onLoginFailed(fmt.Sprint("User ", c.userName, " validate failed ", err))
		return
	}

	if msg.GetID() != Message.ID_Client_Validate_Ret {
		c.onLoginFailed(fmt.Sprint("User ", c.userName, " validate failed , read wrong message ", msg.GetID()))
		return
	}

	ret := msg.(*Message.ClientValidateRet)

	if ret.OK != 1 {
		c.onLoginFailed(fmt.Sprint("User ", c.userName, " validate failed ", ret.ERR))
		return
	}
	c.sess.SetValidate()
	c.serverTime = time.Unix(0, ret.ServerTime)
	c.lastHeartbeatTime = time.Now()

	c.createLocalUser(ret.UserPropType, ret.UserData)
	c.onLoginSucceed()
	go c.recvMessage(c.sess)
	go c.mainLoop()
}

func (c *Client) createLocalUser(propType string, propData []byte) {
	// add later
}

func (c *Client) Logout() {
	if c.sess == nil {
		return
	}

	c.sess.Close()
	c.sess = nil

	c.cancelFunc()
	<-c.closeC
	close(c.closeC)
}

func (c *Client) Rpc(methodName string, args ...interface{}) {
	m := &Message.ClientRpcReq{
		UserID:     c.userID,
		SrvType:    Const.Game,
		MethodName: methodName,
		Args:       Message.PackArgs(args),
		CBIndex:    0,
	}

	c.sendToServer(m)
}

func (c *Client) SpaceRpc(methodName string, args ...interface{}) {
	m := &Message.ClientRpcReq{
		UserID:     c.userID,
		SrvType:    Const.Space,
		MethodName: methodName,
		Args:       Message.PackArgs(args),
		CBIndex:    0,
	}

	c.sendToServer(m)
}

func (c *Client) sendToServer(msg Message.IMessage) {
	if c.sess != nil && c.sess.IsValidate() {
		c.sendMsgNo++
		_ = c.sess.Send(msg, c.sendMsgNo)
	}
}

func (c *Client) mainLoop() {

	ticker := time.NewTicker(100 * time.Millisecond)
	hbTicker := time.NewTicker(5000 * time.Millisecond)

	defer func() {
		ticker.Stop()
		hbTicker.Stop()
	}()

	c.onInit()
	for {
		select {
		case <-ticker.C:
			c.onLoop()
		case <-hbTicker.C:
			c.onHeartBeat()
		case <-c.ctx.Done():
			c.onDestroy()
			c.closeC <- struct{}{}
			return
		case <-c.msgList.Signal():
			for {
				msg, err := c.msgList.Pop()
				if err != nil {
					break
				}
				c.onMessageProc(msg.(Message.IMessage))
			}
		}
	}

}

func (c *Client) onHeartBeat() {
	c.sendHeartBeat()
}

func (c *Client) isHeartBeatTimeout() bool {
	if c.lastHeartbeatTime.IsZero() {
		return false
	}

	if time.Now().Sub(c.lastHeartbeatTime) < 11*time.Second {
		return false
	}

	return true
}

func (c *Client) CloseSess() {
	c.sess.Close()
}

/*
func (c *Client) TryReconnect() error {

	c.sendMsgNo = 0

	var retErr error
	for i := 0; i < 3; i++ {

		log.Debug("try reconnect at ", i, " times", "   ", c.userName, "   msgno ", c.recvMsgNo)

		if c.sess != nil {
			c.sess.Close()
			c.sess = nil
		}

		conn, err := net.Dial("tcp", c.agentAddr)
		if err != nil {
			retErr = errors.New("connect failed")
			continue
		}

		c.sess = Net.NewTcpSess(conn)

		vm := &Message.ClientValidateReq{Version: 0, ID: c.userID, Token: c.token, MsgSNo: c.recvMsgNo}
		c.sess.Send(vm, 0)

		c.sess.SetSendSecretKey(c.csk)
		c.sess.SetRecvSecretKey(c.crk)

		msg, _, err := c.sess.Read()
		if err != nil {
			retErr = errors.New("validate error")
			continue
		}

		if msg.GetID() != Message.ID_Client_Validate_Ret {
			retErr = errors.New("wrong message")
			break
		}

		ret := msg.(*Message.ClientValidateRet)

		if ret.OK == 0 {
			retErr = errors.New("validate failed " + ret.ERR)
			break
		}
		c.sess.SetValidate()

		c.serverTime = time.Unix(0, ret.ServerTime)
		c.lastHeartbeatTime = time.Now()
		go c.recvMessage(c.sess)
		break
	}

	return retErr
}
*/

func (c *Client) sendHeartBeat() {
	msg := &Message.HeartBeat{
		ClientSendTime: time.Now().UnixNano(),
		ServerTime:     0,
	}

	c.sendToServer(msg)
}

func (c *Client) onHeartbeatRecv(msg *Message.HeartBeat) {
	c.ttl = time.Now().Sub(time.Unix(0, msg.ClientSendTime))
	c.lastHeartbeatTime = time.Now()
	c.serverTime = time.Unix(0, msg.ServerTime)
}

func (c *Client) GetSrvTime() time.Time {
	return c.serverTime.Add((c.ttl / 2) + time.Now().Sub(c.lastHeartbeatTime))
}

func (c *Client) recvMessage(sess Net.ISess) {
	for {

		msg, msgNo, err := sess.Read()
		if err != nil {
			sess.Close()
			/*
				if sess.IsValidate() {
					log.Debug("sess read error , so reconnect ", c.userName)
					_ = c.tryReconnect()
				}
			*/
			return
		}

		log.Debug("recv message ", c.userName, "   ", msgNo)

		c.recvMsgNo = msgNo
		c.msgList.Put(msg)
	}
}

func (c *Client) onMessageProc(msg Message.IMessage) {

	switch msg.GetID() {
	case Message.ID_User_Rpc_Req:
		m := msg.(*Message.UserRpcReq)
		args, _ := Message.UnPackArgs(m.Args)
		c.onNetMessage(m.MethodName, args)
	case Message.ID_Heart_Beat:
		c.onHeartbeatRecv(msg.(*Message.HeartBeat))
	}

}
