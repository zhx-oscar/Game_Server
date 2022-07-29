package DB

import (
	"Cinder/DB"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 角色ID工具类
// 当创建角色的时候, 先用LockName方法锁定该角色名, 如果锁定成功, 则说明该名字没有被占用
// 角色创建成功后, 使用SetRoleID保存角色ID
// 角色创建失败, 需要调用Del释放该角色名

var RoleNameTabName = "RoleNames"
var dbInvalidErr = errors.New("MongoDB is invalid")

type RoleDBProp struct {
	RoleName string `bson:"RoleName"`
	RoleID   string `bson:"RoleID"`
}

type NameUtilImpl struct {
	name string
}

func NameUtil(name string) *NameUtilImpl {
	return &NameUtilImpl{
		name: name,
	}
}

func (util *NameUtilImpl) SetRoleID(id string) error {
	if DB.MongoDB == nil {
		return dbInvalidErr
	}
	_, err := DB.MongoDB.Collection(RoleNameTabName).UpdateOne(context.Background(), bson.M{"RoleName": util.name}, bson.M{"$set": bson.M{"RoleID": id}})
	return err
}

func (util *NameUtilImpl) GetID() (string, error) {
	if DB.MongoDB == nil {
		return "", dbInvalidErr
	}
	roleInfo := &RoleDBProp{}
	err := DB.MongoDB.Collection(RoleNameTabName).FindOne(context.Background(), bson.M{"RoleName": util.name}).Decode(roleInfo)
	if err != nil {
		return "", err
	}
	return roleInfo.RoleID, nil
}

func (util *NameUtilImpl) LockName() (bool, error) {

	if DB.MongoDB == nil {
		return false, dbInvalidErr
	}
	changeInfo, err := DB.MongoDB.Collection(RoleNameTabName).UpdateOne(context.Background(),
		bson.M{"RoleName": util.name}, bson.M{"$set": bson.M{"RoleName": util.name}}, options.Update().SetUpsert(true))
	if err != nil {
		return false, err
	}
	if changeInfo.MatchedCount == 1 {
		return false, nil
	}
	return true, nil
}

func (util *NameUtilImpl) Del() error {
	_, err := DB.MongoDB.Collection(RoleNameTabName).DeleteMany(context.Background(), bson.M{"RoleName": util.name})
	return err
}
