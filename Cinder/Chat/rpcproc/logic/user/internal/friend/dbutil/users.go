package dbutil

import (
	"Cinder/DB"
	"context"
	"fmt"
	"time"

	assert "github.com/arl/assertgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// _UsersUtil 操作 chat.users 集合
type _UsersUtil struct {
	userID UserID
}

func UsersUtil(userID UserID) *_UsersUtil {
	return &_UsersUtil{
		userID: userID,
	}
}

func (u *_UsersUtil) c() *mongo.Collection {
	assert.True(DB.MongoDB != nil) // 应该已初始化了
	return DB.MongoDB.Collection("chat.users")
}

func (u *_UsersUtil) DeleteFriend(friendID UserID) error {
	selector := bson.M{"userID": u.userID}
	update := bson.M{"$pull": bson.M{"friendIDs": friendID}}
	_, err := u.c().UpdateOne(context.Background(), selector, update)
	return err
}

func (u *_UsersUtil) AddFriend(friendID UserID) error {
	selector := bson.M{"userID": u.userID}
	update := bson.M{"$addToSet": bson.M{"friendIDs": friendID}}
	_, err := u.c().UpdateOne(context.Background(), selector, update)
	return err
}

func (u *_UsersUtil) GetFriendCount() (int, error) {
	matchStage := bson.D{
		{"$match", bson.M{"userID": u.userID}},
	}
	projectStage := bson.D{
		{"$project", bson.M{
			"friendCount": bson.M{
				"$size": bson.M{
					// The argument to $size must be an array
					"$cond": bson.A{
						// $isArray since MongoDB 3.2
						bson.M{"$isArray": "$friendIDs"}, // if
						"$friendIDs",                     // then
						bson.A{},                         // else
					}, // $cond
				}, // $size
			}, // friendCount
		}}, // project
	} // projectStage
	pipeline := mongo.Pipeline{matchStage, projectStage}
	opts := options.Aggregate().SetMaxTime(2 * time.Second)
	cursor, err := u.c().Aggregate(context.Background(), pipeline, opts)
	if err != nil {
		return 0, fmt.Errorf("Aggregate: %w", err)
	}
	defer cursor.Close(context.Background())

	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		return 0, fmt.Errorf("cursor All: %w", err)
	}
	if len(results) == 0 {
		return 0, nil
	}
	icount := results[0]["friendCount"]
	if icount == nil {
		return 0, nil
	}
	if count, ok := icount.(int32); !ok {
		return 0, fmt.Errorf("mongo size returns non-int32")
	} else {
		return int(count), nil
	}
}
