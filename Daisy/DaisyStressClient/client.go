package main

import (
	"Cinder/Base/Message"
	"Cinder/Base/Net"
	"Cinder/Base/Security"
	"Cinder/stats"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"time"

	log "github.com/cihub/seelog"
)

type StressClient struct {
	// 登录认证相关
	user         string
	password     string
	cspr         *big.Int
	crpr         *big.Int
	csk          []byte
	crk          []byte
	loginRespMsg LoginResp

	// 连接相关
	sendMsgNo  uint32
	clientsess Net.ISess
}

func NewClient(user, pwd string) *StressClient {
	client := &StressClient{
		user:     user,
		password: pwd,
	}
	return client
}

// LoginResp 登录消息返回格式
type LoginResp struct {
	State     int         `json:"state"`
	UserID    string      `json:"userID"`
	Token     string      `json:"token"`
	AgentAddr string      `json:"agentAddr"`
	Csk       string      `json:"csk"`
	Crk       string      `json:"crk"`
	Data      interface{} `json:"data"`
}

func (client *StressClient) String() string {
	return "Client:" + client.user
}

func (client *StressClient) RunAction(action string) {
	switch action {
	case "auth":
		client.ActionAuth()
	case "login":
		client.ActionLogin()
	default:
		log.Debug(client, " unknown action")
	}
}

func (client *StressClient) ActionAuth() {
	for {
		start := time.Now()
		if err := client.HttpAuth(); err != nil {
			log.Error(client, err)
		}
		stats.Add("Auth", time.Now().Sub(start))

		time.Sleep(5 * time.Second)
	}
}

func (client *StressClient) ActionLogin() {
	for {
		start := time.Now()

		if err := client.HttpAuth(); err != nil {
			log.Error(client, err)
			continue
		}

		if err := client.ConnectToAgent(); err != nil {
			log.Error(client, err)
			continue
		}

		stats.Add("Login", time.Now().Sub(start))

		// 保持在线30秒
		time.Sleep(5 * time.Second)

		if err := client.DisconnectToAgent(); err != nil {
			log.Error(client, err)
			continue
		}

		// 下线等待5秒后重新登录
		time.Sleep(5 * time.Second)
	}
}

// HttpAuth 尝试HTTP验证, 验证成功会保留验证成功消息
func (client *StressClient) HttpAuth() error {
	cspu, crpu := client.genSecretKey()
	resp, err := http.PostForm(LoginAddr, url.Values{"key": {"Value"}, "AccountName": {client.user}, "Password": {client.password}, "cspu": {cspu.String()}, "crpu": {crpu.String()}})
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, &client.loginRespMsg); err != nil {
		return err
	}

	cspum, _ := big.NewInt(0).SetString(client.loginRespMsg.Csk, 0)
	crpum, _ := big.NewInt(0).SetString(client.loginRespMsg.Crk, 0)

	client.csk = []byte(Security.Key(client.cspr, cspum).String())
	client.crk = []byte(Security.Key(client.crpr, crpum).String())

	return nil
}

func (client *StressClient) ConnectToAgent() error {
	if client.loginRespMsg.AgentAddr == "" {
		return errors.New("invalid Agent Address")
	}

	client.sendMsgNo = 0
	conn, err := net.Dial("tcp", client.loginRespMsg.AgentAddr)
	if err != nil {
		return errors.New("dail to net service failed " + client.loginRespMsg.AgentAddr)
	}

	client.clientsess = Net.NewTcpSess(conn)

	vm := &Message.ClientValidateReq{Version: 0, ID: client.loginRespMsg.UserID, Token: client.loginRespMsg.Token}
	client.clientsess.Send(vm, 0)
	client.clientsess.SetSendSecretKey(client.csk)
	client.clientsess.SetRecvSecretKey(client.crk)

	msg, _, err := client.clientsess.Read()
	if err != nil {
		return err
	}

	if msg.GetID() != Message.ID_Client_Validate_Ret {
		return errors.New(fmt.Sprintf("Need ClientValidateRet but recv %d", msg.GetID()))
	}

	ret := msg.(*Message.ClientValidateRet)
	if ret.OK != 1 {
		return errors.New("client validate failed")
	}

	go client.recvLoop()

	return nil
}

func (client *StressClient) DisconnectToAgent() error {
	if client.clientsess != nil {
		client.clientsess.Close()
	}
	return nil
}

func (client *StressClient) handleMessage(msg Message.IMessage) {
	//switch msg.GetID() {
	//case Message.ID_Space_Prop_Notify:
	//	fmt.Println("Client Recv SpacePropNotify", msg.(*Message.SpacePropNotify).MethodName)
	//case Message.ID_Enter_Space:
	//	fmt.Println("Client Recv EnterSpace")
	//case Message.ID_Batch_EnterAOI:
	//	fmt.Println("Client Recv BatchEnterAOI")
	//case Message.ID_User_Rpc_Req:
	//	m := msg.(*Message.UserRpcReq)
	//	Message.UnPackArgs(m.Args)
	//	fmt.Println("Client Recv RPC", m.MethodName)
	//case Message.ID_Client_Validate_Ret:
	//	client.onClientValidateRet(msg.(*Message.ClientValidateRet))
	//
	//case Message.ID_Client_Rpc_Ret:
	//	msgCt := msg.(*Message.ClientRpcRet)
	//
	//	args, err := Message.UnPackArgs(msgCt.Ret)
	//	if err != nil {
	//		log.Debug("rpc call failed , unpack arg failed ", err)
	//		return
	//	}
	//
	//	if msgCt.CBIndex == 0 {
	//		return
	//	}
	//
	//	client.onRpcRet(msgCt.CBIndex, args)
	//}
}

func (client *StressClient) recvLoop() {
	for {
		msg, _, err := client.clientsess.Read()
		if err != nil {
			client.clientsess.Close()
			return
		}
		client.handleMessage(msg)
	}
}

func (client *StressClient) genSecretKey() (*big.Int, *big.Int) {
	cspr, cspu := Security.Pair()
	crpr, crpu := Security.Pair()

	client.cspr = cspr
	client.crpr = crpr
	return cspu, crpu
}
