package DB

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID   primitive.ObjectID `bson:"_id"`
	Auth *UserAuth          `bson:"auth"`
}

func NewUser(id primitive.ObjectID) *User {
	return &User{ID: id}
}

type UserAuth struct {
	AccountName string `bson:"AccountName"`
	Password    string `bson:"Password"`
	Data        string `bson:"Data"`
}
