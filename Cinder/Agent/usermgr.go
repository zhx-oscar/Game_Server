package main

import (
	"Cinder/Base/Const"
	"Cinder/Base/DistributeLock"
	"Cinder/Base/Message"
	"Cinder/Base/Net"
	"Cinder/Base/User"
	"Cinder/Cache"
	"context"
	"errors"
	"sync"
	"time"

	log "github.com/cihub/seelog"
)

const (
	UserLoginLockTime        = 30 * time.Second
	UserLogoutPendingTime    = 10 * time.Second
	UserLoginTimeout         = 5 * time.Second
	UserClientMessageTimeout = 15 * time.Second
)

type UserMgr struct {
	User.IUserMgr
	pendingUsers    sync.Map
	loginInfo       sync.Map
	pendDestroyRetC sync.Map
}

func newUserMgr() *UserMgr {
	return &UserMgr{
		IUserMgr: User.NewUserMgr(&_User{}, Inst, true),
	}
}

type _LoginInfo struct {
	retC     chan *struct{}
	processC chan *struct{}
	propType string
	propData []byte
	propDef  []byte
}

func newLoginInfo() *_LoginInfo {
	return &_LoginInfo{
		retC:     make(chan *struct{}, 1),
		processC: make(chan *struct{}, 1),
	}
}

func (mgr *UserMgr) Login(sess Net.ISess, verMsg *Message.ClientValidateReq) error {

	userID := verMsg.ID
	token := verMsg.Token

	lock := DistributeLock.New(mgr.getUserLockKey(userID), DistributeLock.Expire(UserLoginLockTime))

	err := lock.Lock()
	if err != nil {
		mgr.onLoginFailed(sess, err)
		return err
	}

	defer lock.Unlock()
	mgr.clearPendLogout(userID)

	log.Debug("user start login ", userID)

	csk, crk, err := mgr.validateClientToken(userID, token)
	if err != nil {
		mgr.onLoginFailed(sess, err)
		return err
	}

	sess.SetData(userID)
	sess.SetSendSecretKey(crk)
	sess.SetRecvSecretKey(csk)
	sess.SetValidate()

	user, _ := mgr.GetAgentUser(userID)

	// is reconnect request
	if verMsg.MsgSNo > 0 {
		log.Debug("user start reconnect ", userID, " msgno ", verMsg.MsgSNo)
		if user == nil {
			err = errors.New("agent user is new")
			log.Error("reconnect fail ", err, " userID ", userID)
			mgr.onLoginFailed(sess, err)
			return err
		}

		if !user.IsMsgPoolEnough(verMsg.MsgSNo) {
			err = errors.New("message pool is not enough")
			log.Error("reconnect fail ", err, "  userID ", userID)
			mgr.onLoginFailed(sess, err)
			userMgr.kickoffUser(userID)
			return err
		}

		_ = sess.Send(&Message.ClientValidateRet{
			OK:           1,
			UserPropType: "",
			ProtoDef:     nil,
			UserData:     nil,
			ServerTime:   time.Now().UnixNano()}, 0)

		user.ReSendMsgFromPool(verMsg.MsgSNo, sess)
		user.SetClientSess(sess)

	} else {

		log.Debug("user normal login ", userID)

		mgr.removeRemotePart(userID)

		if user != nil {
			log.Debug("user exist, reset msgno ", userID)

			user.ResetSendNo()
		} else {
			var iu User.IUser
			iu, err = userMgr.CreateUser(userID, nil)
			if err != nil {
				log.Error("create agent user failed ", err, " userID ", userID)
				mgr.onLoginFailed(sess, err)
				userMgr.kickoffUser(userID)
				return err
			}

			user = iu.(*_User)
		}

		Cache.CancelUserLoginTokenExpire(userID)

		var loginGroup *_LoginInfo
		loginGroup, err = mgr.fetchLoginInfo(userID)
		if err != nil {
			log.Error("fetch login info failed ", err, " userID ", userID)
			mgr.onLoginFailed(sess, err)
			userMgr.kickoffUser(userID)
			return err
		}

		_ = sess.Send(&Message.ClientValidateRet{
			OK:           1,
			UserPropType: loginGroup.propType,
			ProtoDef:     loginGroup.propDef,
			UserData:     loginGroup.propData,
			ServerTime:   time.Now().UnixNano()}, 0)

		user.SetClientSess(sess)
	}

	return nil
}

func (mgr *UserMgr) removeRemotePart(userID string) {
	peers, err := Cache.GetUserPeersSrvID(userID)
	if err != nil {
		log.Error("removeRemotePart GetUserPeersSrvID err ", err, " userID: ", userID)
		return
	}

	srvID, ok := peers[Const.Space]
	if ok {
		mgr.removeRemotePartBySrvID(srvID, userID)
	}

	srvID, ok = peers[Const.Game]
	if ok {
		mgr.removeRemotePartBySrvID(srvID, userID)
	}

	srvID, ok = peers[Const.Agent]
	if ok {
		if srvID != mgr.GetSrvInst().GetServiceID() {
			mgr.removeRemotePartBySrvID(srvID, userID)
		}
	}
}

func (mgr *UserMgr) removeRemotePartBySrvID(srvID string, userID string) {

	destroyOtherUserC := make(chan struct{}, 1)
	mgr.pendDestroyRetC.Store(userID, destroyOtherUserC)

	msg := &Message.UserDestroyReq{UserID: userID}
	if err := mgr.GetSrvInst().Send(srvID, msg); err != nil {
		log.Error("removeRemotePartBySrvID err ", err, " srvID ", srvID)
		return
	}

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	timeout := false

	select {
	case <-destroyOtherUserC:
		break
	case <-ctx.Done():
		timeout = true
		break
	}

	if timeout {
		log.Error("remove remote part but time out ", userID, " ", srvID)
	}
}

func (mgr *UserMgr) onUserDestroyRet(msg *Message.UserDestroyRet) {

	ii, ok := mgr.pendDestroyRetC.Load(msg.UserID)
	if !ok {
		return
	}

	retC := ii.(chan struct{})
	retC <- struct{}{}
	close(retC)

	mgr.pendDestroyRetC.Delete(msg.UserID)
}

func (mgr *UserMgr) onLoginFailed(sess Net.ISess, err error) {
	_ = sess.Send(&Message.ClientValidateRet{OK: 0, ERR: err.Error(), UserPropType: "", UserData: []byte{}, ServerTime: time.Now().UnixNano()}, 0)
}

func (mgr *UserMgr) fetchLoginInfo(userID string) (*_LoginInfo, error) {

	loginGroup := newLoginInfo()
	mgr.loginInfo.Store(userID, loginGroup)
	defer func() {
		mgr.loginInfo.Delete(userID)
		loginGroup.processC <- &struct{}{}
		close(loginGroup.processC)
	}()

	gameSrvID, err := Inst.GetSrvIDByType(Const.Game)
	if err != nil {
		return nil, err
	}

	if err = Inst.Send(gameSrvID, &Message.UserLoginReq{UserID: userID}); err != nil {
		return nil, err
	}

	ctx, ctxCancel := context.WithTimeout(context.Background(), UserLoginTimeout)
	defer ctxCancel()

	timeout := false
forLoop:

	select {
	case <-loginGroup.retC:
		break forLoop
	case <-ctx.Done():
		timeout = true
		break forLoop
	}

	if timeout {
		err = errors.New("login failed , time out")
		return nil, err
	}

	return loginGroup, nil
}

func (mgr *UserMgr) onLoginRet(userID string, propType string, propDef []byte, userData []byte) <-chan *struct{} {

	info, ok := mgr.loginInfo.Load(userID)
	if ok {
		info.(*_LoginInfo).propData = userData
		info.(*_LoginInfo).propType = propType
		info.(*_LoginInfo).propDef = propDef
		info.(*_LoginInfo).retC <- &struct{}{}
		close(info.(*_LoginInfo).retC)
		return info.(*_LoginInfo).processC
	} else {
		log.Error("couldn't find login group when user login ret from game " + userID)
		ch := make(chan *struct{}, 1)
		ch <- &struct{}{}
		close(ch)
		return ch
	}
}

func (mgr *UserMgr) getUserLockKey(userID string) string {
	return "UserLogin:" + userID
}

func (mgr *UserMgr) Logout(userID string) {
	lock := DistributeLock.New(mgr.getUserLockKey(userID), DistributeLock.Expire(UserLoginLockTime))
	err := lock.Lock()

	if err != nil {
		log.Error("Logout err", err, " userID: ", userID)
		return
	}
	defer lock.Unlock()

	mgr.kickoffUser(userID)

	log.Info("Logout userID ", userID)
}

func (mgr *UserMgr) kickoffUser(userID string) {
	_, err := mgr.GetAgentUser(userID)
	if err != nil {
		log.Error("kickoffUser GetAgentUser err ", err, " userID: ", userID)
		return
	}

	msg := &Message.UserDestroyReq{
		UserID: userID,
	}

	Inst.Broadcast(Const.Game, msg)
	Inst.Broadcast(Const.Space, msg)

	_ = mgr.DestroyUser(userID)

	Cache.ClearUserLoginToken(userID)
}

func (mgr *UserMgr) PendLogout(userID string) {
	log.Debug("add user to pending logout list ", userID)

	i, ok := mgr.pendingUsers.Load(userID)
	if ok {
		i.(*time.Timer).Stop()
	}

	t := time.AfterFunc(UserLogoutPendingTime, func() {
		if _, ok = mgr.pendingUsers.Load(userID); ok {
			log.Debug("user have already pended timeout , so logout ", userID)
			mgr.pendingUsers.Delete(userID)
			mgr.Logout(userID)
		}
	})

	mgr.pendingUsers.Store(userID, t)

	user, err := userMgr.GetAgentUser(userID)
	if err != nil {
		log.Error("PendLogout GetAgentUser err ", err, "userID: ", userID)
		return
	}
	user.SetClientSess(nil)
}

func (mgr *UserMgr) clearPendLogout(userID string) {
	i, ok := mgr.pendingUsers.Load(userID)
	if ok {
		i.(*time.Timer).Stop()
		mgr.pendingUsers.Delete(userID)
	}
}

func (mgr *UserMgr) GetAgentUser(userID string) (*_User, error) {
	iu, err := mgr.GetUser(userID)
	if err != nil {
		return nil, err
	}

	return iu.(*_User), nil
}

func (mgr *UserMgr) CheckRecvMsgValidate(userID string, msgNo uint32) error {

	user, err := mgr.GetAgentUser(userID)
	if err != nil {
		return err
	}

	return user.CheckRecvMsgValidate(msgNo)
}

func (mgr *UserMgr) validateClientToken(id string, token string) ([]byte, []byte, error) {

	t, csk, crk, err := Cache.GetUserLoginToken(id)
	if err != nil {
		return nil, nil, err
	}

	if t != token {
		return nil, nil, errors.New("token invalid")
	}
	return csk, crk, nil
}
