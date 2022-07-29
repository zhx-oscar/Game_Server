package delmark

import (
	"Cinder/Mail/mgocol"
	"Cinder/Mail/rpcproc/handler/delmark/bcmail"
	"Cinder/Mail/rpcproc/handler/delmark/mailid"
	"Cinder/Mail/rpcproc/userid"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SetExtData 设置邮件额外数据
func SetExtData(userID userid.UserID, mailIDStr string, extData []byte) error {
	isBroadcast, objectID, err := mailid.ParseMailID(mailIDStr)
	if err != nil {
		return fmt.Errorf("parse mail id: %w", err)
	}
	if isBroadcast {
		return setBcMailExtData(userID, objectID, extData)
	}
	return setUserMailExtData(userID, objectID, extData)
}

// setBcMailExtData 设置广播邮件的额外数据(更改用户邮件状态)
func setBcMailExtData(userID userid.UserID, originalID primitive.ObjectID, extData []byte) error {
	return bcmail.NewBcMailUpdater(userID, originalID).SetExtData(extData)
}

// setUserMailExtData 设置用户邮件的额外数据
func setUserMailExtData(userID userid.UserID, OID primitive.ObjectID, extData []byte) error {
	selector := userMailSelector(userID, OID)
	update := bson.M{"$set": bson.M{"mail.state.extData": extData}}
	_, err := mgocol.UsersMails().UpdateOne(context.Background(), selector, update)
	return err
}
