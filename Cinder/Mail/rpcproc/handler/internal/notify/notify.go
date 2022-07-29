// notify 包处理新邮件通知
package notify

import (
	"Cinder/Mail/mailapi/types"
	"Cinder/Mail/mgocol"
	"Cinder/Mail/rpcproc/handler/delmark/mailid"
	"Cinder/Mail/rpcproc/rpc"
	"Cinder/Mail/rpcproc/userid"
	"Cinder/Mail/rpcproc/usersrv"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/cihub/seelog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// NotifyMail 通知收件人新邮件。
func NotifyMail(mail types.Mail, oid primitive.ObjectID) error {
	srvID, errID := queryUsersSrvID(userid.UserID(mail.To))
	if errID != nil {
		if errors.Is(errID, mongo.ErrNoDocuments) {
			return nil // 还未上过线
		}
		return fmt.Errorf("query users srv ID: %w", errID)
	}

	if srvID == "" {
		return nil // 未在线
	}

	var err error
	mail.ID, err = mailid.GetUserMailIDStr(oid) // 需要填 mail.ID
	if err != nil {
		log.Errorf("get user mail ID str error: %v", err)
		// 忽略错误，继续通知，只是无法操作该邮件
	}
	if err := notifyMailToSrv(srvID, mail); err != nil {
		return fmt.Errorf("notify mail to server: %w", err)
	}
	return nil
}

// NotifyBroadcastMail 通知系统广播邮件
func NotifyBroadcastMail(mail types.Mail, originalID primitive.ObjectID) error {
	var errID error
	mail.ID, errID = mailid.GetBcMailIDStr(originalID) // 需要填 mail.ID
	if errID != nil {
		log.Errorf("get broadcast mail ID str error: %v", errID)
		// 但是继续通知，只是该邮件无法操作
	}

	mailJson, errJson := json.Marshal(mail)
	if errJson != nil {
		return fmt.Errorf("json marshal: %w", errJson)
	}

	userSrvIDs := usersrv.GetUserSrvIDs()
	errStrs := []string{}
	for _, id := range userSrvIDs {
		if err := notifyMailJsonToSrv(id, mailJson); err != nil {
			errStrs = append(errStrs, err.Error())
		}
	}
	if len(errStrs) == 0 {
		return nil
	}
	return fmt.Errorf("%v", errStrs)
}

// notifyMailToSrv 通知新邮件到服务器
func notifyMailToSrv(srvID string, mail types.Mail) error {
	mailJson, err := json.Marshal(mail)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}
	return notifyMailJsonToSrv(srvID, mailJson)
}

// notifyMailJsonToSrv 通知新邮件到服务器
func notifyMailJsonToSrv(srvID string, mailJson []byte) error {
	ret := rpc.RpcWithRet(srvID, "RPC_MailNotify", mailJson)
	return ret.Err
}

// queryUsersSrvID 从DB查询用户服ID, 即登录时的服务器ID
func queryUsersSrvID(userID userid.UserID) (string, error) {
	doc := struct {
		SrvID string `bson:"srvID"`
	}{}
	if rv := mgocol.Users().FindOne(context.Background(), bson.M{"userID": userID}); rv.Err() != nil {
		return "", fmt.Errorf("find userID '%s' err: %w", userID, rv.Err())
	} else {
		if err := rv.Decode(&doc); err != nil {
			return "", fmt.Errorf("find userID '%s' decode err: %w", userID, rv.Err())
		}
	}
	return doc.SrvID, nil

}
