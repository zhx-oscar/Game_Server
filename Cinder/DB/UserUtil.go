package DB

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserUtil struct {
	_id primitive.ObjectID
	col *mongo.Collection
}

func NewUserUtil(id string) (*UserUtil, error) {
	if MongoDB == nil {
		return nil, errors.New("no database connect")
	}

	var _id primitive.ObjectID

	if id == "" {
		_id = primitive.NewObjectID()
	} else {
		var err error
		if _id, err = primitive.ObjectIDFromHex(id); err != nil {
			return nil, err
		}
	}

	return &UserUtil{_id: _id, col: MongoDB.Collection("User")}, nil
}

func (util *UserUtil) GetID() string {
	return util._id.Hex()
}

func (util *UserUtil) Insert(user *User) error {
	user.ID = util._id
	_, err := util.col.InsertOne(context.Background(), user)
	return err
}

func (util *UserUtil) GetUser() (*User, error) {
	rv := util.col.FindOne(context.Background(), bson.M{"_id": util._id})
	if rv != nil {
		return nil, rv.Err()
	}
	user := NewUser(util._id)
	if err := rv.Decode(&user); err != nil {
		return nil, err
	}
	return user, nil
}

func (util *UserUtil) GetAuth() (*UserAuth, error) {
	rv := util.col.FindOne(context.Background(), bson.M{"_id": util._id}, options.FindOne().SetProjection(bson.M{"auth": 1}))
	if rv.Err() != nil {
		return nil, rv.Err()
	}
	user := NewUser(util._id)
	if err := rv.Decode(&user); err != nil {
		return nil, err
	}
	return user.Auth, nil
}

func (util *UserUtil) UpdateAuthData(authData string) error {
	_, err := util.col.UpdateOne(context.Background(), bson.M{"_id": util._id}, bson.M{"$set": bson.M{"auth.Data": authData}})
	return err
}

func (util *UserUtil) GetAuthByAccName(accName string) (string, *UserAuth, error) {
	rv := util.col.FindOne(context.Background(), bson.M{"auth.AccountName": accName}, options.FindOne().SetProjection(bson.M{"auth": 1, "_id": 1}))
	if rv.Err() != nil {
		return "", nil, rv.Err()
	}
	user := NewUser(primitive.NilObjectID)
	if err := rv.Decode(&user); err != nil {
		return "", nil, err
	}
	return user.ID.Hex(), user.Auth, nil
}
