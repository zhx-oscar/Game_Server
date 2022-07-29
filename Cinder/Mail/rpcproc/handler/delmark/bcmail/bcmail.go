package bcmail

import (
	"Cinder/Mail/rpcproc/handler/delmark/bcmail/db"
	"Cinder/Mail/rpcproc/userid"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// BcMailUpdater 是广播邮件状态更新器
type BcMailUpdater struct {
	dbStates  *db.UsersBcMailStates
	dbBcMails *db.BroadcastMails
}

func NewBcMailUpdater(to userid.UserID, originalID primitive.ObjectID) *BcMailUpdater {
	return &BcMailUpdater{
		dbStates:  db.NewUsersBcMailStates(to, originalID),
		dbBcMails: db.NewBroadcastMails(originalID),
	}
}

func (b *BcMailUpdater) SetIsRead(isRead bool) error {
	err := b.dbStates.UpdateIsRead(isRead)
	if err == nil {
		return nil // 成功
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("update isread: %w", err)
	}

	// 需要复制插入，再更新。插入时允许同时操作，忽略重复键。不可插入更新一次完成。
	if errCp := b.copyStateAndTime(); errCp != nil {
		return fmt.Errorf("copy state and time: %w", errCp)
	}
	return b.dbStates.UpdateIsRead(isRead)
}

func (b *BcMailUpdater) SetIsAttachmentsReceived(isReceived bool) error {
	err := b.dbStates.UpdateIsAttachmentsReceived(isReceived)
	if err == nil {
		return nil // 成功
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("update attachments isreceived: %w", err)
	}

	// 需要读取原状态，再插入新状态
	if errCp := b.copyStateAndTime(); errCp != nil {
		return fmt.Errorf("copy state and time: %w", errCp)
	}
	return b.dbStates.UpdateIsAttachmentsReceived(isReceived)
}

func (b *BcMailUpdater) SetDeleted() error {
	err := b.dbStates.UpdateToDeleted()
	if err == nil {
		return nil // 成功
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("udpate to set deleted: %w", err)
	}

	// 需要读取原状态，再插入新状态
	if errCp := b.copyStateAndTime(); errCp != nil {
		return fmt.Errorf("copy state and time: %w", errCp)
	}
	return b.dbStates.UpdateToDeleted()
}

// SetExtData 设置广播邮件的额外数据(更改用户邮件状态)
func (b *BcMailUpdater) SetExtData(extData []byte) error {
	err := b.dbStates.UpdateExtData(extData)
	if err == nil {
		return nil // 成功
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("update ext data: %w", err)
	}

	// 需要复制插入，再更新。插入时允许同时操作，忽略重复键。不可插入更新一次完成。
	if errCp := b.copyStateAndTime(); errCp != nil {
		return fmt.Errorf("copy state and time: %w", errCp)
	}
	return b.dbStates.UpdateExtData(extData)
}

func (b *BcMailUpdater) copyStateAndTime() error {
	stateTime, errLoad := b.dbBcMails.LoadStateAndTime()
	if errLoad != nil {
		return fmt.Errorf("load broadcast mail state and time: %w", errLoad)
	}

	return b.dbStates.Insert(stateTime) // 已忽略重复键
}
