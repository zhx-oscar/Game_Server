package dbidx

import (
	"Cinder/Chat/rpcproc/logic/db"
	"Cinder/Mail/mgocol"
	"context"

	log "github.com/cihub/seelog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB字段索引设计见 README.md
const UNIQUE = true

func EnsureIndexes() error {
	log.Infof("ensure indexes...")

	if err := ensureIndexesOfBroadcastMails(); err != nil {
		return err
	}
	if err := ensuerIndexesOfUserSrvIDs(); err != nil {
		return err
	}
	if err := ensuerIndexesOfUsers(); err != nil {
		return err
	}
	if err := ensureIndexesOfUsersBcMailStates(); err != nil {
		return err
	}
	if err := ensureIndexesOfUsersMails(); err != nil {
		return err
	}

	// enableSharding 和 shardCollection 需要由管理员手动执行
	log.Infof("ensured indexes.")
	return nil
}

func ensureIndexesOfBroadcastMails() error {
	colBm := mgocol.BroadcastMails()
	if err := ensureIndex(colBm, !UNIQUE, bson.D{{"shard", "hashed"}}); err != nil {
		return err
	}
	if err := ensureIndex(colBm, UNIQUE, bson.D{{"shard", 1}, {"originalID", 1}}); err != nil {
		return err
	}
	if err := ensureIndex(colBm, UNIQUE, bson.D{{"shard", 1}, {"mail.sendTime", 1}, {"originalID", 1}}); err != nil {
		return err
	}
	if err := ensureExpireIndex(colBm, bson.D{{"mail.expireTime", 1}}); err != nil {
		return err
	}
	return nil
}

func ensuerIndexesOfUserSrvIDs() error {
	colLi := mgocol.UserSrvIDs()
	if err := ensureIndex(colLi, UNIQUE, bson.D{{"srvID", 1}}); err != nil {
		return err
	}
	return nil
}

func ensuerIndexesOfUsers() error {
	colUser := mgocol.Users()
	if err := ensureIndex(colUser, !UNIQUE, bson.D{{"userID", "hashed"}}); err != nil {
		return err
	}
	if err := ensureIndex(colUser, UNIQUE, bson.D{{"userID", 1}}); err != nil {
		return err
	}
	return nil
}

func ensureIndexesOfUsersBcMailStates() error {
	colSt := mgocol.UsersBcMailStates()
	if err := ensureIndex(colSt, !UNIQUE, bson.D{{"to", "hashed"}}); err != nil {
		return err
	}
	if err := ensureIndex(colSt, UNIQUE, bson.D{{"to", 1}, {"originalID", 1}}); err != nil {
		return err
	}
	if err := ensureIndex(colSt, !UNIQUE, bson.D{{"to", 1}, {"sendTime", 1}}); err != nil {
		return err
	}
	if err := ensureBcMailStatesExpireIndex(); err != nil {
		return err
	}
	return nil
}

func ensureIndexesOfUsersMails() error {
	colUm := mgocol.UsersMails()
	if err := ensureIndex(colUm, !UNIQUE, bson.D{{"mail.to", "hashed"}}); err != nil {
		return err
	}
	if err := ensureIndex(colUm, UNIQUE, bson.D{{"mail.to", 1}, {"_id", 1}}); err != nil {
		return err
	}
	if err := ensureIndex(colUm, UNIQUE, bson.D{{"mail.to", 1}, {"mail.sendTime", 1}, {"_id", 1}}); err != nil {
		return err
	}
	if err := ensureExpireIndex(colUm, bson.D{{"mail.expireTime", 1}}); err != nil {
		return err
	}
	return nil
}

func ensureIndex(col *mongo.Collection, isUnique bool, keys bson.D) error {
	_, err := col.Indexes().CreateOne(context.Background(), getIndex(isUnique, keys))
	return db.SkipSameIndexErr(err)
}

func ensureExpireIndex(col *mongo.Collection, keys bson.D) error {
	_, err := col.Indexes().CreateOne(context.Background(), getExpireIndex(keys, 1))
	return db.SkipSameIndexErr(err)
}

func getExpireIndex(keys bson.D, second int32) mongo.IndexModel {
	return mongo.IndexModel{
		Keys:    keys,
		Options: options.Index().SetExpireAfterSeconds(second),
	}
}

func getIndex(isUnique bool, keys bson.D) mongo.IndexModel {
	return mongo.IndexModel{
		Keys:    keys,
		Options: options.Index().SetUnique(isUnique),
	}
}

func ensureBcMailStatesExpireIndex() error {
	_, err := mgocol.UsersBcMailStates().Indexes().CreateOne(context.Background(), getExpireIndex(bson.D{{"expireTime", 1}}, 3600))
	return db.SkipSameIndexErr(err)
}
