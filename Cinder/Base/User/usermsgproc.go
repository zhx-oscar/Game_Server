package User

import (
	"Cinder/Base/Const"
	"Cinder/Base/Core"
	"Cinder/Base/MQNet"
	"Cinder/Base/Mailbox"
	"Cinder/Base/Message"
	"Cinder/Base/Prop"
	"Cinder/Base/Util"
	"errors"
	"reflect"
	"sync"

	log "github.com/cihub/seelog"
)

type _IUserQuery interface {
	GetUser(id string) (IUser, error)
	DestroyUser(id string) error
	Traversal(cb func(user IUser) bool)
	GetSrvInst() Core.ICore
}

func NewUserMessageProc(query _IUserQuery) MQNet.IProc {
	return &_UserMsgProc{
		mgr: query,
	}
}

type _UserMsgProc struct {
	mgr _IUserQuery

	methodRegTypeMap sync.Map
}

type _IBMsgProc interface {
	MsgProc(msg Message.IMessage)
}

type _IClientSender interface {
	SendToAllClient(msg Message.IMessage)
	SendToAllClientExceptMe(msg Message.IMessage)
}

func (p *_UserMsgProc) MessageProc(srvAddr string, msg Message.IMessage) {

	switch msg.GetID() {
	case Message.ID_User_Broadcast_Create:
		m := msg.(*Message.UserBroadcastCreate)

		i, err := p.mgr.GetUser(m.UserID)
		if err != nil {
			return
		}

		u, ok := i.(_IUser)
		if ok {
			u.syncPeerCreate(m.SrvID, m.SrvType)
		}

	case Message.ID_User_Broadcast_Destroy:
		m := msg.(*Message.UserBroadcastDestroy)

		i, err := p.mgr.GetUser(m.UserID)
		if err != nil {
			return
		}

		u, ok := i.(_IUser)
		if ok {
			u.syncPeerDestroy(m.SrvID, m.SrvType)
		}

	case Message.ID_User_Prop_Notify:
		m := msg.(*Message.UserPropNotify)

		i, err := p.mgr.GetUser(m.UserID)
		if err != nil {
			return
		}

		u, ok := i.(_IUser)
		if ok {
			if u.GetType() == Const.Agent {
				_ = u.SendToClient(m)
			} else if u.GetType() == Const.Space {
				if m.Target == Prop.Target_Space {
					var args []interface{}
					args, err = Message.UnPackArgs(m.Args)
					if err == nil {
						if u.GetProp() != nil {
							u.GetProp().GetCaller().SafeCall(m.MethodName, args...)
						}
					} else {
						log.Error("MessageProc UserPropNotify UnPackArgs err ", err, " Method: ", m.MethodName, " userID: ", m.UserID)
					}
				} else if m.Target == Prop.Target_All_Clients {
					u.GetRealPtr().(_IClientSender).SendToAllClient(m)
				} else if m.Target == Prop.Target_Other_Clients {
					u.GetRealPtr().(_IClientSender).SendToAllClientExceptMe(m)
				}
			} else {
				var args []interface{}
				args, err = Message.UnPackArgs(m.Args)
				if err == nil {
					if u.GetProp() != nil {
						u.GetProp().GetCaller().SafeCall(m.MethodName, args...)
					}
				} else {
					log.Error("MessageProc UserPropNotify UnPackArgs err ", err, " Method: ", m.MethodName, " userID: ", m.UserID)
				}
			}
		}

	case Message.ID_User_Rpc_Req:

		go func() {

			m := msg.(*Message.UserRpcReq)
			retMsg := &Message.RpcRet{
				RetID: m.RetID,
				Ret:   nil,
				Err:   "",
			}

			i, err := p.mgr.GetUser(m.UserID)
			if err == nil {
				u := i.(_IUser)
				if u.GetType() == Const.Agent {
					u.SendToClient(m)
					return
				}

				var args []interface{}
				args, err = Message.UnPackArgs(m.Args)
				if err == nil {
					r := <-u.SafeCall(m.MethodName, args...)
					if r.Err == nil {
						retMsg.Ret = Message.PackArgs(r.Ret...)
					} else {
						retMsg.Err = "MessageProc UserRpcReq SafeCall err " + r.Err.Error() + " Method " + m.MethodName + " userID " + m.UserID
						log.Error(retMsg.Err)
					}

				} else {
					retMsg.Err = "MessageProc UserRpcReq UnPackArgs err " + err.Error() + " Method " + m.MethodName + " userID " + m.UserID
					log.Error(retMsg.Err)
				}

			} else {
				retMsg.Err = "MessageProc UserRpcReq GetUser err " + err.Error() + " Method " + m.MethodName + " userID " + m.UserID
			}

			p.mgr.GetSrvInst().Send(srvAddr, retMsg)
		}()

	case Message.ID_Users_Rpc_Req:

		m := msg.(*Message.UsersRpcReq)

		for _, userID := range m.UserIDS {
			i, err := p.mgr.GetUser(userID)
			if err != nil {
				continue
			}

			u, ok := i.(_IUser)
			if ok {
				if u.GetType() == Const.Agent {
					u.SendToClient(&Message.UserRpcReq{
						UserID:     userID,
						MethodName: m.MethodName,
						Args:       m.Args,
					})
				} else {
					go func() {
						var args []interface{}
						args, err = Message.UnPackArgs(m.Args)
						if err != nil {
							log.Error("MessageProc UsersRpcReq UnPackArgs err ", err, " Method ", m.MethodName, " userID ", u.GetID())
							return
						}

						r := <-u.SafeCall(m.MethodName, args...)
						if r.Err != nil {
							log.Error("MessageProc UsersRpcReq SafeCall err ", err, " Method ", m.MethodName, " userID ", u.GetID())
						}
					}()
				}
			}
		}

	case Message.ID_All_Users_Rpc_Req:

		m := msg.(*Message.AllUsersRpcReq)

		p.mgr.Traversal(func(user IUser) bool {

			u, ok := user.(_IUser)
			if ok {
				if u.GetType() == Const.Agent {
					u.SendToClient(&Message.UserRpcReq{
						UserID:     u.GetID(),
						MethodName: m.MethodName,
						Args:       m.Args,
					})
				} else {
					go func() {
						args, err := Message.UnPackArgs(m.Args)
						if err != nil {
							log.Error("MessageProc AllUsersRpcReq UnPackArgs err ", err, " Method ", m.MethodName, " userID ", u.GetID())
							return
						}

						r := <-u.SafeCall(m.MethodName, args...)
						if r.Err != nil {
							log.Error("MessageProc AllUsersRpcReq SafeCall err ", err, " Method ", m.MethodName, " userID ", u.GetID())
						}
					}()
				}
			}

			return true
		})

	case Message.ID_User_Rpc_Ret:

		m := msg.(*Message.UserRpcRet)

		i, err := p.mgr.GetUser(m.UserID)
		if err != nil {
			return
		}

		args, err := Message.UnPackArgs(m.Ret)
		if err != nil {
			m.Err = "MessageProc UserRpcRet UnPackArgs err " + err.Error()
			log.Error(m.Err)
			return
		}

		u, ok := i.(_IUser)
		if ok {
			u.OnRpcRet(m.RetID, m.Err, args)
		}

	case Message.ID_Client_Rpc_Req:
		m := msg.(*Message.ClientRpcReq)

		i, err := p.mgr.GetUser(m.UserID)
		if err != nil {
			return
		}

		u, ok := i.(_IUser)
		if ok {
			if _InstMethodForbidden.isMethodForbidden(m.MethodName) {
				u.SendToPeerServer(Const.Agent, &Message.RpcForbiddenRet{
					UserID:  u.GetID(),
					CBIndex: m.CBIndex,
					Source:  m.Source,
				})
				return
			}

			var args []interface{}
			args, err = Message.UnPackArgs(m.Args)
			if err != nil {
				log.Error("MessageProc ClientRpcReq UnPackArgs err ", err, " Method ", m.MethodName, " userID ", m.UserID)
				return
			}

			args, err = p.conventArgs(u, m.MethodName, args)
			if err != nil {
				log.Error("MessageProc ClientRpcReq conventArgs err ", err, " Method ", m.MethodName, " userID ", m.UserID)
				return
			}

			go func() {
				var r *Util.SafeCallRet
				r = <-u.SafeCall(m.MethodName, args...)

				if r.Err == nil {
					retData := Message.PackArgs(r.Ret...)
					u.SendToPeerServer(Const.Agent, &Message.ClientRpcRet{
						UserID:  u.GetID(),
						CBIndex: m.CBIndex,
						Ret:     retData,
						Source:  m.Source,
					})
				} else {
					log.Error("MessageProc ClientRpcReq SafeCall err ", r.Err, " Method ", m.MethodName, " userID ", m.UserID)
				}
			}()
		}

	case Message.ID_Forward_User_Message:

		m := msg.(*Message.ForwardUserMessage)

		user, err := p.mgr.GetUser(m.UserID)
		if err != nil {
			return
		}

		innerMsg, err := Message.Unpack(m.MsgData)
		if err != nil {
			log.Error("MessageProc ForwardUserMessage Unpack err ", err, " userID ", m.UserID)
			return
		}

		if user.GetType() == Const.Agent {
			user.SendToClient(innerMsg)
		} else {
			mp, ok := user.(_IBMsgProc)
			if ok {
				mp.MsgProc(innerMsg)
			}
		}

	case Message.ID_Forward_Users_Message:

		m := msg.(*Message.ForwardUsersMessage)

		for _, userID := range m.UserIDS {
			user, err := p.mgr.GetUser(userID)
			if err != nil {
				continue
			}

			innerMsg, err := Message.Unpack(m.MsgData)
			if err != nil {
				log.Error("MessageProc ForwardUsersMessage Unpack err ", err)
				continue
			}

			if user.GetType() == Const.Agent {
				user.SendToClient(innerMsg)
			} else {
				mp, ok := user.(_IBMsgProc)
				if ok {
					mp.MsgProc(innerMsg)
				}
			}
		}

	case Message.ID_Forward_All_Users_Message:

		m := msg.(*Message.ForwardAllUsersMessage)

		innerMsg, err := Message.Unpack(m.MsgData)
		if err != nil {
			log.Error("MessageProc ForwardAllUsersMessage Unpack err ", err)
			return
		}

		p.mgr.Traversal(func(user IUser) bool {
			if user.GetType() == Const.Agent {
				user.SendToClient(innerMsg)
			} else {
				mp, ok := user.(_IBMsgProc)
				if ok {
					mp.MsgProc(innerMsg)
				}
			}
			return true
		})

	case Message.ID_Mailbox_Req:
		go func() {
			m := msg.(*Message.MailboxReq)
			if m.MailBoxType != Mailbox.TypeUser {
				return
			}

			retMsg := &Message.MailboxRet{
				MailboxID: m.MailBoxID,
				RetID:     m.RetID,
				Ret:       nil,
				Err:       "",
			}

			i, err := p.mgr.GetUser(m.TargetID)
			if err == nil {
				u := i.(_IUser)

				var args []interface{}
				args, err = Message.UnPackArgs(m.Args)
				if err == nil {
					r := <-u.SafeCall(m.MethodName, args...)
					if r.Err == nil {
						retMsg.Ret = Message.PackArgs(r.Ret...)
					} else {
						retMsg.Err = "MessageProc MailboxReq SafeCall err " + r.Err.Error() + " Method " + m.MethodName + " userID " + m.TargetID
						log.Error(retMsg.Err)
					}

				} else {
					retMsg.Err = "MessageProc MailboxReq UnPackArgs err " + err.Error() + " Method " + m.MethodName + " userID " + m.TargetID
					log.Error(retMsg.Err)
				}

			} else {
				retMsg.Err = "MessageProc MailboxReq GetUser err " + err.Error() + " Method " + m.MethodName + " userID " + m.TargetID
			}

			p.mgr.GetSrvInst().Send(srvAddr, retMsg)
		}()
	}
}

var errMethodNotExisted = errors.New("method not existed")

func (p *_UserMsgProc) conventArgs(user _IUser, methodName string, in []interface{}) ([]interface{}, error) {
	if len(in) <= 0 {
		return in, nil
	}

	var regTypeList []reflect.Type
	v, ok := p.methodRegTypeMap.Load(methodName)
	if ok {
		regTypeList = v.([]reflect.Type)
	} else {
		m, mok := reflect.TypeOf(user.GetRealPtr()).MethodByName(methodName)
		if !mok {
			return nil, errMethodNotExisted
		}

		regTypeList = make([]reflect.Type, 0, 1)
		for i := 1; i < m.Type.NumIn(); i++ {
			regTypeList = append(regTypeList, m.Type.In(i))
		}
		p.methodRegTypeMap.Store(methodName, regTypeList)
	}

	for i := range in {
		in[i] = conventNum(regTypeList[i], in[i])
	}

	return in, nil
}

// 转换数字类型, 从客户端上来的int64和float64类型转换成目标数字类型
func conventNum(targetType reflect.Type, source interface{}) interface{} {
	sourceType := reflect.TypeOf(source)
	if targetType.Kind() == sourceType.Kind() {
		return source
	}
	if sourceType.Kind() != reflect.Int64 && sourceType.Kind() != reflect.Float64 {
		return source
	}

	target := reflect.New(targetType)
	switch targetType.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if sourceType.Kind() == reflect.Int64 {
			target.Elem().SetUint(uint64(source.(int64)))
		} else {
			target.Elem().SetUint(uint64(source.(float64)))
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if sourceType.Kind() == reflect.Int64 {
			target.Elem().SetInt(source.(int64))
		} else {
			target.Elem().SetInt(int64(source.(float64)))
		}
	case reflect.Float32, reflect.Float64:
		if sourceType.Kind() == reflect.Int64 {
			target.Elem().SetFloat(float64(source.(int64)))
		} else {
			target.Elem().SetFloat(source.(float64))
		}
	default:
		return source
	}

	return target.Elem().Interface()
}
