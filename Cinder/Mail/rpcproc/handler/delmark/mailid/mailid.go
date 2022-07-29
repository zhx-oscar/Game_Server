package mailid

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// _MailID 合并邮件ID信息，用于打包后填入 types.Mail.ID.
// _MailID 不会存入 DB, 这样可以任意更改打包算法。
type _MailID struct {
	IsBroadcast bool
	ObjectID    primitive.ObjectID // 如果是广播邮件，该字段为原始ID
}

// GetUserMailIDStr 生成普通邮件 Mail ID 字符串
func GetUserMailIDStr(objectID primitive.ObjectID) (string, error) {
	return GetMailIDStr(false /*isBroadcast*/, objectID)
}

// GetBcMailIDStr 生成广播邮件 Mail ID 字符串
func GetBcMailIDStr(objectID primitive.ObjectID) (string, error) {
	return GetMailIDStr(true /*isBroadcast*/, objectID)
}

// GetMailIDStr 生成 Mail ID 字符串
func GetMailIDStr(isBroadcast bool, oid primitive.ObjectID) (string, error) {
	mailID := _MailID{
		IsBroadcast: isBroadcast,
		ObjectID:    oid,
	}
	buf, err := json.Marshal(mailID)
	if err != nil {
		return "", fmt.Errorf("json marshal: %w", err)
	}
	return base64.URLEncoding.EncodeToString(buf), nil
}

// ParseMailID 将字符串解包为 _MailID 信息
func ParseMailID(mailIDStr string) (isBroadcast bool, oid primitive.ObjectID, err error) {
	buf, err := base64.URLEncoding.DecodeString(mailIDStr)
	if err != nil {
		return isBroadcast, oid, fmt.Errorf("base64 decode: %w", err)
	}
	mailID := _MailID{}
	if err := json.Unmarshal(buf, &mailID); err != nil {
		return isBroadcast, oid, fmt.Errorf("json unmarshal: %w", err)
	}
	return mailID.IsBroadcast, mailID.ObjectID, nil
}
