package delmark

import (
	"Cinder/Mail/mgocol"
	"Cinder/Mail/rpcproc/handler/delmark/bcmail"
	"Cinder/Mail/rpcproc/handler/delmark/mailid"
	"Cinder/Mail/rpcproc/handler/internal/maildoc"
	"Cinder/Mail/rpcproc/userid"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MarkAsRead(userID userid.UserID, mailIDStr string) error {
	return markRead(userID, mailIDStr, true)
}

func MarkAsUnread(userID userid.UserID, mailIDStr string) error {
	return markRead(userID, mailIDStr, false)
}

// MarkAttachmentsAsReceived 设置附件已收。
// 与其他设置不同，附件只能收取一次，如果附件已收，则不能再次设置。
func MarkAttachmentsAsReceived(userID userid.UserID, mailIDStr string) error {
	return markAttach(userID, mailIDStr, true)
}

func MarkAttachmentsAsUnreceived(userID userid.UserID, mailIDStr string) error {
	return markAttach(userID, mailIDStr, false)
}

func markRead(userID userid.UserID, mailIDStr string, read bool) error {
	isBroadcast, objectID, err := mailid.ParseMailID(mailIDStr)
	if err != nil {
		return fmt.Errorf("parse mail id: %w", err)
	}
	if isBroadcast {
		if err := markBcMailRead(userID, objectID, read); err != nil {
			return fmt.Errorf("mark broad mail read(%v): %w", read, err)
		}
		return nil
	}
	if err := markUserMailRead(userID, objectID, read); err != nil {
		return fmt.Errorf("mark user mail read(%v): %w", read, err)
	}
	return nil
}

func markBcMailRead(userID userid.UserID, originalID primitive.ObjectID, read bool) error {
	return bcmail.NewBcMailUpdater(userID, originalID).SetIsRead(read)
}

func markUserMailRead(userID userid.UserID, OID primitive.ObjectID, read bool) error {
	selector := userMailSelector(userID, OID)
	update := bson.M{"$set": bson.M{"mail.state.isRead": read}}
	_, err := mgocol.UsersMails().UpdateOne(context.Background(), selector, update)
	return err
}

func markAttach(userID userid.UserID, mailIDStr string, received bool) error {
	isBroadcast, objectID, err := mailid.ParseMailID(mailIDStr)
	if err != nil {
		return fmt.Errorf("parse mail id: %w", err)
	}

	// 如果已领取，则不能再次领取
	if received {
		if err := checkAlreadyReceived(userID, isBroadcast, objectID); err != nil {
			return fmt.Errorf("check already received: %w", err)
		}
	}

	if isBroadcast {
		if err := markBcMailAttach(userID, objectID, received); err != nil {
			return fmt.Errorf("mark broadcast mail attachments received(%v): %w", received, err)
		}
		return nil
	}
	if err := markUserMailAttach(userID, objectID, received); err != nil {
		return fmt.Errorf("mark user mail attachments received(%v): %w", received, err)
	}
	return nil
}

func markBcMailAttach(userID userid.UserID, originalID primitive.ObjectID, received bool) error {
	return bcmail.NewBcMailUpdater(userID, originalID).SetIsAttachmentsReceived(received)
}

func markUserMailAttach(userID userid.UserID, OID primitive.ObjectID, received bool) error {
	selector := userMailSelector(userID, OID)
	update := bson.M{"$set": bson.M{"mail.state.isAttachmentsReceived": received}}
	_, err := mgocol.UsersMails().UpdateOne(context.Background(), selector, update)
	return err
}

func checkAlreadyReceived(userID userid.UserID, isBroadcast bool, OID primitive.ObjectID) error {
	isReceivedFun := isUserMailAttachmentsReceived
	if isBroadcast {
		isReceivedFun = isBcMailAttachmentsReceived
	}

	isReceived, err := isReceivedFun(userID, OID)
	if err != nil {
		return fmt.Errorf("check is received: %w", err)
	}
	if isReceived {
		return errors.New("attachments are already received")
	}
	return nil
}

func isUserMailAttachmentsReceived(userID userid.UserID, OID primitive.ObjectID) (bool, error) {
	query := userMailSelector(userID, OID)
	selector := bson.M{"mail.state.isAttachmentsReceived": 1}
	doc := maildoc.UserMailDoc{}
	if rv := mgocol.UsersMails().FindOne(context.Background(), query, options.FindOne().SetProjection(selector)); rv.Err() != nil {
		return false, rv.Err()
	} else {
		if err := rv.Decode(&doc); err != nil {
			return false, err
		}
	}
	if doc.Mail == nil || doc.Mail.State == nil {
		return false, nil
	}
	return doc.Mail.State.IsAttachmentsReceived, nil
}

func isBcMailAttachmentsReceived(userID userid.UserID, OID primitive.ObjectID) (bool, error) {
	query := bson.M{
		"to":         userID,
		"originalID": OID,
	}
	selector := bson.M{"state.isAttachmentsReceived": 1}
	doc := struct {
		State maildoc.MailState `bson:"state"`
	}{}
	if rv := mgocol.UsersBcMailStates().FindOne(context.Background(), query, options.FindOne().SetProjection(selector)); rv.Err() != nil {
		if errors.Is(rv.Err(), mongo.ErrNoDocuments) {
			return false, nil // 还未设置状态
		}
		return false, rv.Err()
	} else {
		if err := rv.Decode(&doc); err != nil {
			return false, err
		}
	}
	return doc.State.IsAttachmentsReceived, nil
}
