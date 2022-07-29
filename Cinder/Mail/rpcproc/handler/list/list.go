// list 从DB列出邮件
package list

import (
	"Cinder/Mail/mailapi/types"
	"Cinder/Mail/mgocol"
	"Cinder/Mail/rpcproc/handler/internal/maildoc"
	"Cinder/Mail/rpcproc/userid"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// ListMails 从DB列邮件
func ListMails(userID userid.UserID) ([]types.Mail, error) {
	createTime, err := queryUserCreateTime(userID)
	if err != nil {
		return nil, fmt.Errorf("query user create time: %w", err)
	}
	userMails, errUsr := queryUserMails(userID, createTime)
	if errUsr != nil {
		return nil, fmt.Errorf("query user mails: %w", errUsr)
	}
	cbMails, errBc := queryBcMails(userID, createTime)
	if errBc != nil {
		return nil, fmt.Errorf("query broadcast mails: %w", errBc)
	}

	mails := append(userMails, cbMails...)
	return mails, nil
}

// queryUserCreateTime 查询用户创建时间
func queryUserCreateTime(userID userid.UserID) (time.Time, error) {
	doc := struct {
		ID primitive.ObjectID `bson:"_id"`
	}{}
	if rv := mgocol.Users().FindOne(context.Background(), bson.M{"userID": userID}); rv.Err() != nil {
		return time.Time{}, rv.Err()
	} else {
		if err := rv.Decode(&doc); err != nil {
			return time.Time{}, err
		}
	}
	return doc.ID.Timestamp(), nil
}

// queryUserMails 查询用户邮件
func queryUserMails(userID userid.UserID, fromTime time.Time) ([]types.Mail, error) {
	var docs []*maildoc.UserMailDoc
	cursor, err := mgocol.UsersMails().Find(context.Background(), bson.M{
		"mail.to":       userID,
		"mail.sendTime": bson.M{"$gte": fromTime},
	}, options.Find().SetLimit(1000))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &docs); err != nil {
		return nil, err
	}

	result := make([]types.Mail, 0, len(docs))
	for _, doc := range docs {
		mail := maildoc.DocUserToMail(*doc.Mail, doc.OID) // DocUserToMail()中会设置ID
		result = append(result, mail)
	}
	return result, nil
}

func queryBcMails(userID userid.UserID, fromTime time.Time) ([]types.Mail, error) {
	mailDocs, err := newBcLister(userID, fromTime).List()
	if err != nil {
		return nil, fmt.Errorf("broadcast lister: %w", err)
	}

	result := make([]types.Mail, 0, len(mailDocs))
	for _, doc := range mailDocs {
		mail := maildoc.DocBcToMail(*doc.Mail, doc.OriginalID) // DocBcToMail()中会设置ID
		result = append(result, mail)
	}
	return result, nil
}
