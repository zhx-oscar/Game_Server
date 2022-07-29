package internal

import (
	"Cinder/Base/Core"
	"Cinder/Mail/mailapi/types"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type _MailService struct {
}

func NewMailService() *_MailService {
	return &_MailService{}
}

func (m *_MailService) Login(userID string) error {
	inst := Core.Inst
	return rpc("RPC_Login", userID, inst.GetServiceID())
}

func (m *_MailService) Logout(userID string) error {
	return rpc("RPC_Logout", userID)
}

func (m *_MailService) Send(mail *types.Mail) error {
	buf, err := json.Marshal(mail)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}
	return rpc("RPC_Send", buf)
}

// 弃用，请改用 Send()
func (m *_MailService) Broadcast(mail *types.Mail) error {
	mail.IsBroadcast = true
	return m.Send(mail)
}

func (m *_MailService) ListMail(userID string) ([]*types.Mail, error) {
	ret := rpcWithRet("RPC_ListMail", userID)
	if ret.Err != nil {
		return nil, ret.Err
	}

	errStr := ret.Ret[0].(string)
	if errStr != "" {
		return nil, errors.New(errStr)
	}

	var mails []*types.Mail
	if err := json.Unmarshal(ret.Ret[1].([]byte), &mails); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}
	return mails, nil
}

func (m *_MailService) Delete(userID string, mailID string) error {
	return rpc("RPC_Delete", userID, mailID)
}

func (m *_MailService) BatchDelete(userID string, mailIDs []string) error {
	mailIDsJson, err := json.Marshal(mailIDs)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}
	return rpc("RPC_BatchDelete", userID, mailIDsJson)
}

func (m *_MailService) MarkAsRead(userID string, mailID string) error {
	return rpc("RPC_MarkAsRead", userID, mailID)
}

func (m *_MailService) MarkAsUnread(userID string, mailID string) error {
	return rpc("RPC_MarkAsUnread", userID, mailID)
}

func (m *_MailService) MarkAttachmentsAsReceived(userID string, mailID string) error {
	return rpc("RPC_MarkAttachmentsAsReceived", userID, mailID)
}

func (m *_MailService) MarkAttachmentsAsUnreceived(userID string, mailID string) error {
	return rpc("RPC_MarkAttachmentsAsUnreceived", userID, mailID)
}

func (m *_MailService) SetExpireTime(userID string, mailID string, expireTime time.Time) error {
	return rpc("RPC_SetExpireTime", userID, mailID, expireTime.Unix())
}

func (m *_MailService) SetExtData(userID string, mailID string, extData []byte) error {
	return rpc("RPC_SetExtData", userID, mailID, extData)
}
