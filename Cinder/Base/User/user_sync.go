package User

import (
	"Cinder/Base/Const"
	"Cinder/Base/Message"
	"Cinder/Base/Prop"
	"Cinder/Cache"
	"errors"
)

type _IUserPeerChange interface {
	OnPeerChange(srvType string, srvID string)
}

func (u *User) notifyCreate() {
	u.notifyCreateToType(Const.Agent)
	u.notifyCreateToType(Const.Game)
	u.notifyCreateToType(Const.Space)
}

func (u *User) notifyDestroy() {
	u.notifyDestroyToType(Const.Agent)
	u.notifyDestroyToType(Const.Game)
	u.notifyDestroyToType(Const.Space)
}

func (u *User) notifyCreateToType(srvType string) {

	if srvType == u.GetType() {
		return
	}

	srvID, err := u.GetPeerServerID(srvType)
	if err != nil {
		return
	}

	msg := &Message.UserBroadcastCreate{
		UserID:  u.GetID(),
		SrvType: u.GetSrvInst().GetServiceType(),
		SrvID:   u.GetSrvInst().GetServiceID(),
	}

	if err = u.GetSrvInst().Send(srvID, msg); err != nil {
		u.Error("notifyCreateToType Send msg err", err)
	}
}

func (u *User) notifyDestroyToType(srvType string) {

	if srvType == u.GetType() {
		return
	}

	srvID, err := u.GetPeerServerID(srvType)
	if err != nil {
		return
	}

	msg := &Message.UserBroadcastDestroy{
		UserID:  u.GetID(),
		SrvType: u.GetSrvInst().GetServiceType(),
		SrvID:   u.GetSrvInst().GetServiceID(),
	}

	if err = u.GetSrvInst().Send(srvID, msg); err != nil {
		u.Error("notifyDestroyToType Send msg err", err)
	}
}

func (u *User) initPeers() {
	peers, err := Cache.GetUserPeersSrvID(u.GetID())
	if err != nil {
		u.Error("initPeers GetUserPeersSrvID err", err)
		return
	}

	for k, v := range peers {
		u.peers.Store(k, v)
	}
}

func (u *User) GetPeerServerID(srvType string) (string, error) {
	srvID, ok := u.peers.Load(srvType)
	if !ok {
		return "", errors.New("no peer existed")
	}
	return srvID.(string), nil
}

func (u *User) syncPeerCreate(srvID string, srvType string) {
	if srvType == u.GetType() {
		if srvID != u.GetSrvInst().GetServiceID() {
			u.Error("you should have been destroyed  ....")
		}
	} else {
		u.peers.Store(srvType, srvID)

		ii, ok := u.GetRealPtr().(_IUserPeerChange)
		if ok {
			ii.OnPeerChange(srvType, srvID)
		}
	}
}

func (u *User) syncPeerDestroy(srvID string, srvType string) {
	id, ok := u.peers.Load(srvType)
	if ok {
		if id == srvID {
			u.peers.Delete(srvType)

			ii, ok := u.GetRealPtr().(_IUserPeerChange)
			if ok {
				ii.OnPeerChange(srvType, "")
			}
		}
	}
}

func (u *User) SendToPeerServer(srvType string, msg Message.IMessage) error {

	id, ok := u.peers.Load(srvType)
	if !ok {
		return errors.New("no peer exist " + srvType)
	}

	err := u.GetSrvInst().Send(id.(string), msg)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) SendToPeerUser(srvType string, msg Message.IMessage) error {

	buf, err := Message.Pack(msg)
	if err != nil {
		return err
	}

	forwardmsg := &Message.ForwardUserMessage{
		TargetSrv: srvType,
		UserID:    u.GetID(),
		MsgData:   buf,
	}

	return u.SendToPeerServer(srvType, forwardmsg)
}

func (u *User) SendToClient(msg Message.IMessage) error {

	if !u.IsReady() {
		return errNotInited
	}

	if u.GetType() == Const.Agent {

		sender, ok := u.GetRealPtr().(IClientMessageSender)
		if ok {
			return sender.SendMessageToClient(msg)
		}

	} else {

		buf, err := Message.Pack(msg)
		if err != nil {
			return err
		}

		forwardmsg := &Message.ForwardUserMessage{
			TargetSrv: Const.Agent,
			UserID:    u.GetID(),
			MsgData:   buf,
		}

		return u.SendToPeerServer(Const.Agent, forwardmsg)
	}

	return nil
}

type _ISpaceUser interface {
	SendToAllClient(msg Message.IMessage)
	SendToAllClientExceptMe(msg Message.IMessage)
}

func (u *User) SyncProp(methodName string, args []byte, targets ...int) {

	as, _ := Message.UnPackArgs(args)
	_, _ = u.GetProp().GetCaller().CallMethod(methodName, as...)

	for _, target := range targets {

		msg := &Message.UserPropNotify{
			UserID:     u.GetID(),
			MethodName: methodName,
			Args:       args,
		}

		switch target {
		case Prop.Target_Game:
			if u.GetType() != Const.Game {
				_ = u.SendToPeerServer(Const.Game, msg)
			}
		case Prop.Target_Space:
			if u.GetType() != Const.Space {
				msg.Target = Prop.Target_Space
				_ = u.SendToPeerServer(Const.Space, msg)
			}
		case Prop.Target_Client:
			_ = u.SendToPeerServer(Const.Agent, msg)

		case Prop.Target_All_Clients:
			if u.GetType() == Const.Space {
				su, ok := u.GetRealPtr().(_ISpaceUser)
				if ok {
					su.SendToAllClient(msg)
				}

			} else {
				msg.Target = Prop.Target_All_Clients
				_ = u.SendToPeerServer(Const.Space, msg)

			}
		case Prop.Target_Other_Clients:
			if u.GetType() == Const.Space {
				su, ok := u.GetRealPtr().(_ISpaceUser)
				if ok {
					su.SendToAllClientExceptMe(msg)
				}
			} else {
				msg.Target = Prop.Target_Other_Clients
				_ = u.SendToPeerServer(Const.Space, msg)
			}

		default:
			u.Debug("no support target type", target)
		}

	}
}
