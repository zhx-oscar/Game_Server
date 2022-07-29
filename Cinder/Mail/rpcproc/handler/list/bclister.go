package list

import (
	"Cinder/Mail/mgocol"
	"Cinder/Mail/rpcproc/handler/internal/bcshard"
	"Cinder/Mail/rpcproc/handler/internal/maildoc"
	"Cinder/Mail/rpcproc/userid"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type _BcState struct {
	OriginalID primitive.ObjectID `bson:"originalID"`
	State      maildoc.MailState  `bson:"state"`
	Deleted    bool               `bson:"deleted"`
}

type _BcLister struct {
	userID   userid.UserID
	fromTime time.Time
}

func newBcLister(userID userid.UserID, fromTime time.Time) *_BcLister {
	return &_BcLister{
		userID:   userID,
		fromTime: fromTime,
	}
}

// List 列出广播邮件。
// 仅列出未删除的。状态用个人状态更新。
func (b *_BcLister) List() ([]*maildoc.BcMailDoc, error) {
	// 先获取 bc_mail_states 的记录，因为个数会比 bc_mails 的个数少。
	states, errStates := b.queryStates()
	if errStates != nil {
		return nil, fmt.Errorf("query states: %w", errStates)
	}

	// 然后读取 bc_mails, 去除已标记删除的邮件。
	stateMap, deleted := mapBcStates(states)
	bcMailDocs, errBc := b.queryBcMails(deleted)
	if errBc != nil {
		return nil, fmt.Errorf("query bc mails: %w", errBc)
	}

	// 更新状态
	for _, doc := range bcMailDocs {
		if doc == nil { // doc 必须是指针才能更新
			continue
		}
		if state, ok := stateMap[doc.OriginalID]; ok {
			doc.Mail.State = &state
		}
	}

	return bcMailDocs, nil
}

// queryBcMails 查询广播邮件，排除fromTime之前的，排除已删的。
func (b *_BcLister) queryBcMails(deleted []primitive.ObjectID) ([]*maildoc.BcMailDoc, error) {
	if deleted == nil {
		deleted = []primitive.ObjectID{} // fix mongo：$nin needs an array
	}

	var docs []*maildoc.BcMailDoc
	cursor, err := mgocol.BroadcastMails().Find(context.Background(), bson.M{
		"shard":         bcshard.GetRandBcShardID(),
		"mail.sendTime": bson.M{"$gte": b.fromTime},
		"originalID":    bson.M{"$nin": deleted},
	}, options.Find().SetLimit(1000))
	if err != nil {
		return nil, fmt.Errorf("find err: %w", err)
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &docs); err != nil {
		return nil, fmt.Errorf("cursor err: %w", err)
	}

	return docs, nil
}

func (b *_BcLister) queryStates() ([]_BcState, error) {
	docs := []_BcState{}
	cursor, err := mgocol.UsersBcMailStates().Find(context.Background(), bson.M{
		"to":       b.userID,
		"sendTime": bson.M{"$gte": b.fromTime},
	}, options.Find().SetLimit(1000))
	if err != nil {
		return nil, fmt.Errorf("find err: %w", err)
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &docs); err != nil {
		return nil, fmt.Errorf("cursor err: %w", err)
	}

	return docs, nil
}

// mapBcStates 将状态数组转成map, 并提取已删除邮件的ID
func mapBcStates(states []_BcState) (mapStates map[primitive.ObjectID]maildoc.MailState, deleted []primitive.ObjectID) {
	mapStates = make(map[primitive.ObjectID]maildoc.MailState)
	for _, st := range states {
		if st.Deleted {
			deleted = append(deleted, st.OriginalID)
			continue
		}
		mapStates[st.OriginalID] = st.State
	}
	return mapStates, deleted
}
