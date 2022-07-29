package main

import (
	_ "Cinder/Base/Log"
	"Cinder/Base/Security"
	_ "Cinder/Base/ServerConfig"
	"Cinder/Base/Util"
	"Cinder/Cache"
	"Cinder/DB"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"strings"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"

	log "github.com/cihub/seelog"
)

/*
http://127.0.0.1:1010/Login?AccountName="wxj";Password="xxwwwgggxx"
*/

type LoginResp struct {
	State     int         `json:"state"`
	UserID    string      `json:"userID"`
	Token     string      `json:"token"`
	AgentAddr string      `json:"agentAddr"`
	Csk       string      `json:"csk"`
	Crk       string      `json:"crk"`
	Data      interface{} `json:"data"`
}

const (
	StatePasswordWrong     = 1
	StateInvalidAccoutName = 2
	StateInvalidKey        = 3
	StateInternalErr       = 4
)

var (
	errPasswordIncorrect = errors.New("password incorrect")
)

func main() {
	defer log.Flush()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// 启动服务器
	if err := initAgentAssign(); err != nil {
		log.Error("start login failed ", err)
		return
	}

	listenAddr := viper.GetString("Login.ListenAddr")
	http.HandleFunc("/Login", HandleLogin)

	log.Debug("login server is starting , listen ", listenAddr)

	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Error("login server start failed ", err)
		return
	}

	log.Debug("login server is shutdown")
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Error("HandleLogin ParseForm err ", err)
		return
	}

	accountName := r.Form.Get("AccountName")
	password := r.Form.Get("Password")
	cspums := r.Form.Get("cspu")
	crpums := r.Form.Get("crpu")

	resp := LoginResp{}
	defer func() {
		data, _ := json.Marshal(resp)
		log.Debug("HandleLogin Account: ", accountName, " Resp: ", string(data))
		w.Write(data)
	}()

	if accountName == "" {
		resp.State = StateInvalidAccoutName
		return
	}

	if cspums == "" || crpums == "" {
		resp.State = StateInvalidKey
		return
	}

	agentAddr, err := GetAgent()
	if err != nil {
		resp.State = StateInternalErr
		return
	}

	cspum, b := big.NewInt(0).SetString(cspums, 0)
	if !b {
		resp.State = StateInternalErr
		return
	}
	crpum, b := big.NewInt(0).SetString(crpums, 0)
	if !b {
		resp.State = StateInternalErr
		return
	}

	cspr, crpr, cspu, crpu := genSecretKey()

	csk := Security.Key(cspr, cspum).String()
	crk := Security.Key(crpr, crpum).String()

	id, token, err := login(accountName, password, "", []byte(csk), []byte(crk))
	if err != nil {
		resp.State = StatePasswordWrong
		return
	}

	resp.UserID = id
	resp.Token = token
	resp.AgentAddr = agentAddr
	resp.Csk = cspu.String()
	resp.Crk = crpu.String()

	log.Info("HandleLogin Success Account: ", accountName)
}

func genSecretKey() (*big.Int, *big.Int, *big.Int, *big.Int) {
	cspr, cspu := Security.Pair()
	crpr, crpu := Security.Pair()
	return cspr, crpr, cspu, crpu
}

func checkAccount(accountName, pwd, data string, authCreate bool) (string, error) {
	util, err := DB.NewUserUtil("")
	if err != nil {
		return "", err
	}

	id, auth, err := util.GetAuthByAccName(accountName)

	if err != nil {
		if authCreate {
			user := DB.NewUser(primitive.NilObjectID)

			user.Auth = &DB.UserAuth{
				AccountName: accountName,
				Password:    pwd,
				Data:        data,
			}

			err = util.Insert(user)

			if err != nil {
				return "", err
			}

			return util.GetID(), nil
		}

		return "", err
	}

	if auth.Password != pwd {
		return "", errPasswordIncorrect
	}

	return id, nil
}

func login(accountName, pwd, data string, csk, crk []byte) (string, string, error) {
	id, err := checkAccount(accountName, pwd, data, viper.GetBool("Login.AutoCreate"))
	if err != nil {
		return "", "", err
	}

	token := "token_" + Util.GetGUID()
	err = Cache.SetUserLoginToken(id, token, csk, crk)
	if err != nil {
		return "", "", err
	}

	return id, token, nil
}
