package delmark

import (
	"Cinder/Mail/mgocol"
	"Cinder/Mail/rpcproc/handler/delmark/bcmail"
	"Cinder/Mail/rpcproc/handler/delmark/mailid"
	"Cinder/Mail/rpcproc/userid"
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
)

func DeleteMail(userID userid.UserID, mailIDStr string) error {
	isBroadcast, objectID, err := mailid.ParseMailID(mailIDStr)
	if err != nil {
		return fmt.Errorf("parse mail id: %w", err)
	}
	if isBroadcast {
		return deleteBcMail(userID, objectID)
	}
	return deleteUserMail(userID, objectID)
}

// DeleteMails 删除多个邮件。
// 尽力删除。出错仍继续，只是添加错误信息。
func DeleteMails(userID userid.UserID, mailIDs []string) error {
	var build strings.Builder                           // 用于错误输出
	OIDs := make([]primitive.ObjectID, 0, len(mailIDs)) // 普通邮件OID
	errMap := make(map[string]string, len(mailIDs))     // mailID->error string
	for _, mailID := range mailIDs {
		isBroadcast, objectID, err := mailid.ParseMailID(mailID)
		if err != nil {
			errMap[mailID] = fmt.Sprintf("parse mail id: %s", err)
			continue
		}

		if isBroadcast {
			//  广播邮件总是单个删除
			if err := deleteBcMail(userID, objectID); err != nil {
				errMap[mailID] = err.Error()
			}
			continue
		}

		build.WriteString(mailID)
		build.WriteString(" ")
		OIDs = append(OIDs, objectID)
	}

	// 用户邮件可成批删除
	if err := deleteUserMails(userID, OIDs); err != nil {
		errMap[build.String()] = fmt.Sprintf("delete user mails: %s", err)
	}

	if len(errMap) == 0 {
		return nil
	}
	return fmt.Errorf("delete mails: %v", errMap)
}

// deleteBcMail 删除广播邮件
func deleteBcMail(userID userid.UserID, originalID primitive.ObjectID) error {
	return bcmail.NewBcMailUpdater(userID, originalID).SetDeleted()
}

// deleteUserMail 删除用户邮件。
// 必须要有 userID, 用于分片查找。
func deleteUserMail(userID userid.UserID, OID primitive.ObjectID) error {
	_, err := mgocol.UsersMails().DeleteOne(context.Background(), userMailSelector(userID, OID))
	return err
}

// deleteUserMails 批量删除用户邮件。
// 必须要有 userID, 用于分片查找。
func deleteUserMails(userID userid.UserID, OIDs []primitive.ObjectID) error {
	if len(OIDs) == 0 {
		return nil
	}
	_, err := mgocol.UsersMails().DeleteMany(context.Background(), bson.M{
		"mail.to": userID,
		"_id":     bson.M{"$in": OIDs},
	})
	return err
}
