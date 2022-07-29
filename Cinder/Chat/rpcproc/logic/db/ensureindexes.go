package db

import (
	"Cinder/DB"
	"context"
	"fmt"

	assert "github.com/arl/assertgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func EnsureIndexes() error {
	assert.True(DB.MongoDB != nil) // 应该已初始化了
	/* 有以下集合：
	chat.users
	chat.users.followers
	chat.users.friend_reqeusts
	chat.users.friend_responses
	chat.users.groups
	chat.users.offline_messages
	chat.groups.members
	chat.groups.messages
	*/
	if err := ensureIndexesOfUsers(); err != nil {
		return fmt.Errorf("ensure users index: %w", err)
	}
	if err := ensureIndexesOfUsersFollowers(); err != nil {
		return fmt.Errorf("ensure users.followers index: %w", err)
	}
	if err := ensureIndexesOfUsersFriendRequests(); err != nil {
		return fmt.Errorf("ensure users.friend_requests index: %w", err)
	}
	if err := ensureIndexesOfUsersFriendResponses(); err != nil {
		return fmt.Errorf("ensure users.friend_responses index: %w", err)
	}
	if err := ensureIndexesOfUsersGroups(); err != nil {
		return fmt.Errorf("ensuer users.groups: %w", err)
	}
	if err := ensureIndexesOfUsersOfflineMessages(); err != nil {
		return fmt.Errorf("ensuer users.offline_messages: %w", err)
	}

	if err := ensureIndexesOfGroupsMembers(); err != nil {
		return fmt.Errorf("ensure groups.members index: %w", err)
	}
	if err := ensureIndexesOfGroupsMessages(); err != nil {
		return fmt.Errorf("ensure group.messages: %w", err)
	}
	return nil
}

func ensureIndexesOfUsers() error {
	c := DB.MongoDB.Collection("chat.users")
	if _, err := c.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{"userID", "hashed"}},
	}); SkipSameIndexErr(err) != nil {
		return fmt.Errorf("ensure index $hashed:userID: %w", err)
	}
	if _, err := c.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{"userID", 1}},
		Options: options.Index().SetUnique(true),
	}); SkipSameIndexErr(err) != nil {
		return fmt.Errorf("ensure index userID: %w", err)
	}
	return nil
}

func ensureIndexesOfUsersFollowers() error {
	c := DB.MongoDB.Collection("chat.users.followers")
	if _, err := c.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{"userID", "hashed"}},
	}); SkipSameIndexErr(err) != nil {
		return fmt.Errorf("ensure index $hashed:userID: %w", err)
	}
	if _, err := c.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{"userID", 1}, {"follower", 1}},
		Options: options.Index().SetUnique(true),
	}); SkipSameIndexErr(err) != nil {
		return fmt.Errorf("ensure index (userID, follower): %w", err)
	}
	return nil
}

func ensureIndexesOfUsersFriendRequests() error {
	c := DB.MongoDB.Collection("chat.users.friend_requests")
	if _, err := c.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{"userID", "hashed"}},
	}); SkipSameIndexErr(err) != nil {
		return fmt.Errorf("ensure index $hashed:userID: %w", err)
	}
	if _, err := c.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{"userID", 1}, {"fromID", 1}},
		Options: options.Index().SetUnique(true),
	}); SkipSameIndexErr(err) != nil {
		return fmt.Errorf("ensure index (userID, fromID): %w", err)
	}
	return nil
}

func ensureIndexesOfUsersFriendResponses() error {
	c := DB.MongoDB.Collection("chat.users.friend_responses")
	if _, err := c.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{"userID", "hashed"}},
	}); SkipSameIndexErr(err) != nil {
		return fmt.Errorf("ensure index $hashed:userID: %w", err)
	}
	if _, err := c.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{"userID", 1}, {"responderID", 1}},
		Options: options.Index().SetUnique(true),
	}); SkipSameIndexErr(err) != nil {
		return fmt.Errorf("ensure index (userID, responderID): %w", err)
	}
	return nil
}

func ensureIndexesOfGroupsMembers() error {
	c := DB.MongoDB.Collection("chat.groups.members")
	if _, err := c.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{"groupID", "hashed"}},
	}); SkipSameIndexErr(err) != nil {
		return fmt.Errorf("ensure index $hashed:groupID: %w", err)
	}
	if _, err := c.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{"groupID", 1}, {"memberID", 1}},
		Options: options.Index().SetUnique(true),
	}); SkipSameIndexErr(err) != nil {
		return fmt.Errorf("ensure index (groupID, memberID): %w", err)
	}
	return nil
}

func ensureIndexesOfGroupsMessages() error {
	c := DB.MongoDB.Collection("chat.groups.messages")
	if _, err := c.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{"groupID", "hashed"}},
	}); SkipSameIndexErr(err) != nil {
		return fmt.Errorf("ensure index $hashed:groupID: %w", err)
	}
	if _, err := c.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{"groupID", 1}, {"sequenceID", 1}},
		Options: options.Index().SetUnique(true),
	}); SkipSameIndexErr(err) != nil {
		return fmt.Errorf("ensure index (groupID, sequenceID): %w", err)
	}
	return nil
}

func ensureIndexesOfUsersOfflineMessages() error {
	c := DB.MongoDB.Collection("chat.users.offline_messages")
	if _, err := c.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{"userID", "hashed"}},
	}); SkipSameIndexErr(err) != nil {
		return fmt.Errorf("ensure index $hashed:userID: %w", err)
	}
	if _, err := c.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{"userID", 1}, {"sendTime", 1}},
	}); SkipSameIndexErr(err) != nil {
		return fmt.Errorf("ensure index (userID, sendTime): %w", err)
	}
	return nil
}

func ensureIndexesOfUsersGroups() error {
	c := DB.MongoDB.Collection("chat.users.groups")
	if _, err := c.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{"userID", "hashed"}},
	}); SkipSameIndexErr(err) != nil {
		return fmt.Errorf("ensure index $hashed:userID: %w", err)
	}
	if _, err := c.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{"userID", 1}, {"groupID", 1}},
		Options: options.Index().SetUnique(true),
	}); SkipSameIndexErr(err) != nil {
		return fmt.Errorf("ensure index (userID, groupID): %w", err)
	}
	return nil
}
