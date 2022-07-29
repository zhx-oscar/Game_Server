package db

import (
	"Cinder/Mail/mgocol"
	"Cinder/Mail/rpcproc/handler/internal/maildoc"
	"Cinder/Mail/rpcproc/userid"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UsersBcMailStateDoc 是 mail.users.bc_mail_states 全服邮件状态集合文档
type UsersBcMailStateDoc struct {
	To         userid.UserID      `bson:"to"`
	OriginalID primitive.ObjectID `bson:"originalID"`
	State      maildoc.MailState  `bson:"state"`

	SendTime   time.Time `bson:"sendTime"` // 用于检索时分页
	ExpireTime time.Time `bson:"expireTime"`
	Deleted    bool      `bson:"deleted"` // 表示邮件已删除
}

// UsersBcMailStates 对应 mail.users.bc_mail_states 集合操作
type UsersBcMailStates struct {
	to         userid.UserID
	originalID primitive.ObjectID
}

func NewUsersBcMailStates(to userid.UserID, originalID primitive.ObjectID) *UsersBcMailStates {
	return &UsersBcMailStates{
		to:         to,
		originalID: originalID,
	}
}

func (u *UsersBcMailStates) UpdateIsRead(isRead bool) error {
	updator := getUpdatorSetIsRead(isRead)
	return u.update(updator)
}

func (u *UsersBcMailStates) UpdateIsAttachmentsReceived(isReceived bool) error {
	updator := getUpdatorSetIsAttachmentsReceived(isReceived)
	return u.update(updator)
}

func (u *UsersBcMailStates) UpdateExtData(extData []byte) error {
	updator := getUpdatorSetExtData(extData)
	return u.update(updator)
}

func (u *UsersBcMailStates) UpdateToDeleted() error {
	updator := getUpdatorDelete()
	return u.update(updator)
}

func (u *UsersBcMailStates) update(updator bson.M) error {
	rv, err := u.c().UpdateOne(context.Background(), u.getSelector(), updator)
	if err != nil {
		return err
	}
	if rv.MatchedCount <= 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (u *UsersBcMailStates) Insert(stateTime StateAndTime) error {
	doc := u.newDoc(stateTime)
	_, err := u.c().InsertOne(context.Background(), doc)
	return err
}

func (u *UsersBcMailStates) c() *mongo.Collection {
	return mgocol.UsersBcMailStates()
}

func (u *UsersBcMailStates) getSelector() bson.M {
	return bson.M{"to": u.to, "originalID": u.originalID}
}

func (u *UsersBcMailStates) newDoc(stateTime StateAndTime) *UsersBcMailStateDoc {
	return &UsersBcMailStateDoc{
		To:         u.to,
		OriginalID: u.originalID,
		State:      stateTime.State,
		SendTime:   stateTime.Mail.SendTime,
		ExpireTime: stateTime.Mail.ExpireTime,
		Deleted:    false,
	}
}

func getUpdatorSetIsRead(isRead bool) bson.M {
	return getUpdator(bson.M{"state.isRead": isRead})
}

func getUpdatorSetIsAttachmentsReceived(isReceived bool) bson.M {
	return getUpdator(bson.M{"state.isAttachmentsReceived": isReceived})
}

func getUpdatorSetExtData(extData []byte) bson.M {
	return getUpdator(bson.M{"state.extData": extData})
}

func getUpdatorDelete() bson.M {
	return getUpdator(bson.M{"deleted": true})
}

func getUpdator(setM bson.M) bson.M {
	return bson.M{"$set": setM}
}
