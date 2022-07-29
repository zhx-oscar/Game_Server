package rpcproc

import (
	"Cinder/Mail/mailapi"
	"Cinder/Mail/mailapi/types"
	"Cinder/Mail/rpcproc/handler/bcmail"
	"Cinder/Mail/rpcproc/handler/delmark"
	"Cinder/Mail/rpcproc/handler/list"
	"Cinder/Mail/rpcproc/handler/loginout"
	"Cinder/Mail/rpcproc/handler/usermail"
	"Cinder/Mail/rpcproc/userid"
	"Cinder/Mail/rpcproc/usersrv"
	"encoding/json"
	"fmt"
	"time"

	log "github.com/cihub/seelog"
)

type UserID = userid.UserID

type _RPCProc struct{}

func NewRPCProc() *_RPCProc {
	return &_RPCProc{}
}

func (r *_RPCProc) RPC_Login(userID string, peerSrvID string) string {
	log.Infof("RPC_Login(userID=%s, peerSrvID=%s)", userID, peerSrvID)
	if err := loginout.Login(UserID(userID), peerSrvID); err != nil {
		return err.Error()
	}
	return ""
}

func (r *_RPCProc) RPC_Logout(userID string) string {
	log.Infof("RPC_Logout(userID=%s)", userID)
	if err := loginout.Logout(UserID(userID)); err != nil {
		return err.Error()
	}
	return ""
}

func (r *_RPCProc) RPC_Send(mailJsonData []byte) string {
	log.Infof("RPC_Send(mailJsonData(%d bytes))", len(mailJsonData))
	mail, err := mailapi.UnmarshalMail(mailJsonData)
	if err != nil {
		return err.Error()
	}

	if mail.IsBroadcast {
		// 存DB, 然后通知
		if err := bcmail.Broadcast(mail); err != nil {
			return err.Error()
		}
		return ""
	}

	// 存DB, 然后查询用户，如果在线就通知
	if err := usermail.Send(mail); err != nil {
		return err.Error()
	}
	return ""
}

// RPC_Broadcast 已弃用。改用 RPC_Send
func (r *_RPCProc) RPC_Broadcast(mailJsonData []byte) string {
	log.Infof("RPC_Broadcast(mailJsonData(%d bytes))", len(mailJsonData))
	mail, err := mailapi.UnmarshalMail(mailJsonData)
	if err != nil {
		return err.Error()
	}

	// 存DB, 然后通知
	if err := bcmail.Broadcast(mail); err != nil {
		return err.Error()
	}
	return ""
}

// RPC_ListMail 列举邮件。
// mailsJsonData 是 []types.Mail 的 json 打包
func (r *_RPCProc) RPC_ListMail(userID string) (errStr string, mailsJsonData []byte) {
	log.Infof("RPC_ListMail(userID=%s)", userID)
	var mails []types.Mail
	var err error
	mails, err = list.ListMails(UserID(userID))
	if err != nil {
		return err.Error(), nil
	}
	buf, errJson := json.Marshal(mails)
	if errJson != nil {
		return errJson.Error(), nil
	}

	return "", buf
}

func (r *_RPCProc) RPC_Delete(userID string, mailID string) string {
	log.Infof("RPC_Delete(userID=%s, mailID=%s)", userID, mailID)
	if err := delmark.DeleteMail(UserID(userID), mailID); err != nil {
		return err.Error()
	}
	return ""
}

func (r *_RPCProc) RPC_BatchDelete(userID string, mailIDsJson []byte) string {
	log.Infof("RPC_BatchDelete(userID=%s, mailIDsJson=%d bytes)", userID, len(mailIDsJson))
	mailIDs := []string{}
	if err := json.Unmarshal(mailIDsJson, &mailIDs); err != nil {
		return fmt.Sprintf("json unmarshal error: %s", err)
	}
	if err := delmark.DeleteMails(UserID(userID), mailIDs); err != nil {
		return err.Error()
	}
	return ""
}

func (r *_RPCProc) RPC_MarkAsRead(userID string, mailID string) string {
	log.Infof("RPC_MarkAsRead(userID=%s, mailID=%s)", userID, mailID)
	if err := delmark.MarkAsRead(UserID(userID), mailID); err != nil {
		return err.Error()
	}
	return ""
}

func (r *_RPCProc) RPC_MarkAsUnread(userID string, mailID string) string {
	log.Infof("RPC_MarkAsUnread(userID=%s, mailID=%s)", userID, mailID)
	if err := delmark.MarkAsUnread(UserID(userID), mailID); err != nil {
		return err.Error()
	}
	return ""
}

func (r *_RPCProc) RPC_MarkAttachmentsAsReceived(userID string, mailID string) string {
	log.Infof("RPC_MarkAttachmentsAsReceived(userID=%s, mailID=%s)", userID, mailID)
	if err := delmark.MarkAttachmentsAsReceived(UserID(userID), mailID); err != nil {
		return err.Error()
	}
	return ""
}

func (r *_RPCProc) RPC_MarkAttachmentsAsUnreceived(userID string, mailID string) string {
	log.Infof("RPC_MarkAttachmentsAsUnreceived(userID=%s, mailID=%s)", userID, mailID)
	if err := delmark.MarkAttachmentsAsUnreceived(UserID(userID), mailID); err != nil {
		return err.Error()
	}
	return ""
}

func (r *_RPCProc) RPC_SetExpireTime(userID string, mailID string, expireUnixSec int64) string {
	expireTime := time.Unix(expireUnixSec, 0)
	log.Infof("RPC_SetExpireTime(userID=%s, mailID=%s, expireTime=%s)", userID, mailID, expireTime)
	if err := delmark.SetExpireTime(UserID(userID), mailID, expireTime); err != nil {
		return err.Error()
	}
	return ""
}

func (r *_RPCProc) RPC_SetExtData(userID string, mailID string, extData []byte) string {
	log.Infof("RPC_SetExtData(userID=%s, mailID=%s, extDataLen=%d)", userID, mailID, len(extData))
	if err := delmark.SetExtData(UserID(userID), mailID, extData); err != nil {
		return err.Error()
	}
	return ""
}

// RPC_SyncUserSrvID 所有Mail服同步增加用户服ID
func (r *_RPCProc) RPC_SyncUserSrvID(srvID string) {
	log.Infof("RPC_SyncUserSrvID(srvID=%s)", srvID)
	usersrv.Sync(srvID)
}
