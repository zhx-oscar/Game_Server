package DB

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Prop struct {
	ID       primitive.ObjectID `bson:"_id"`
	Data     []byte             `bson:"data"`
	JsonData bson.M             `bson:"json_data"`
}

func NewProp(id primitive.ObjectID) *Prop {
	return &Prop{ID: id, Data: []byte{}}
}
