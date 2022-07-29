// mgocol 包返回 mongodb collection
package mgocol

import (
	"Cinder/DB"

	assert "github.com/arl/assertgo"
	"go.mongodb.org/mongo-driver/mongo"
)

func Users() *mongo.Collection {
	assert.True(DB.MongoDB != nil) // 应该已初始化了
	return DB.MongoDB.Collection("mail.users")
}

func UsersMails() *mongo.Collection {
	assert.True(DB.MongoDB != nil) // 应该已初始化了
	return DB.MongoDB.Collection("mail.users.mails")
}

func UsersBcMailStates() *mongo.Collection {
	assert.True(DB.MongoDB != nil) // 应该已初始化了
	return DB.MongoDB.Collection("mail.users.bc_mail_states")
}

func UserSrvIDs() *mongo.Collection {
	assert.True(DB.MongoDB != nil) // 应该已初始化了
	return DB.MongoDB.Collection("mail.user_srv_ids")
}

func BroadcastMails() *mongo.Collection {
	assert.True(DB.MongoDB != nil) // 应该已初始化了
	return DB.MongoDB.Collection("mail.broadcast_mails")
}
