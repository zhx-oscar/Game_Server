package delmark

import (
	"Cinder/Mail/rpcproc/userid"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
)

func userMailSelector(userID userid.UserID, OID primitive.ObjectID) bson.M {
	return bson.M{
		"mail.to": userID,
		"_id":     OID,
	}
}
