package delmark

import (
	"Cinder/Mail/mgocol"
	"Cinder/Mail/rpcproc/handler/delmark/mailid"
	"Cinder/Mail/rpcproc/userid"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SetExpireTime(userID userid.UserID, mailIDStr string, expireTime time.Time) error {
	isBroadcast, objectID, err := mailid.ParseMailID(mailIDStr)
	if err != nil {
		return fmt.Errorf("parse mail id: %w", err)
	}
	if isBroadcast {
		return fmt.Errorf("can not set broadcasted mail's expire time")
	}
	if err := setUserMailExpireTime(userID, objectID, expireTime); err != nil {
		return fmt.Errorf("set user mail expire time: %w", err)
	}
	return nil
}

func setUserMailExpireTime(userID userid.UserID, OID primitive.ObjectID, expireTime time.Time) error {
	selector := userMailSelector(userID, OID)
	update := bson.M{"$set": bson.M{"mail.expireTime": expireTime}}
	_, err := mgocol.UsersMails().UpdateOne(context.Background(), selector, update)
	return err
}
