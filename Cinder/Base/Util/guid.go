package Util

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetGUID() string {
	return primitive.NewObjectID().Hex()
	//return xid.New().String()
}
