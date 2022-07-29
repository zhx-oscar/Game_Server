package DB

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PropUtil struct {
	_id primitive.ObjectID
	col *mongo.Collection
}

func NewPropUtil(id string, typ string) (*PropUtil, error) {
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

	return &PropUtil{_id: _id, col: MongoDB.Collection(typ + "_tbl")}, nil
}

func (util *PropUtil) GetData() ([]byte, error) {
	prop := NewProp(util._id)

	if rv := util.col.FindOne(context.Background(), bson.M{"_id": util._id}, options.FindOne().SetProjection(bson.M{"data": 1})); rv.Err() != nil {
		if _, err := util.col.InsertOne(context.Background(), prop); err != nil {
			return nil, err
		}

		return []byte{}, nil

	} else {
		if err := rv.Decode(&prop); err != nil {
			return nil, err
		}
	}

	return prop.Data, nil
}

func (util *PropUtil) GetBsonData() ([]byte, error) {
	prop := NewProp(util._id)

	if rv := util.col.FindOne(context.Background(), bson.M{"_id": util._id}, options.FindOne().SetProjection(bson.M{"json_data": 1})); rv.Err() != nil {
		if _, err := util.col.InsertOne(context.Background(), prop); err != nil {
			return nil, err
		}

		return nil, nil

	} else {
		if err := rv.Decode(&prop); err != nil {
			return nil, err
		}
	}

	return bson.Marshal(prop.JsonData)
}

func (util *PropUtil) SetData(data []byte) error {
	_, err := util.col.UpdateOne(context.Background(), bson.M{"_id": util._id}, bson.M{"$set": bson.M{"data": data}}, options.Update().SetUpsert(true))
	return err
}

func (util *PropUtil) SetBsonData(data []byte) error {

	jsonData := &bson.M{}
	err := bson.Unmarshal(data, jsonData)
	if err != nil {
		return err
	}

	_, err = util.col.UpdateOne(context.Background(), bson.M{"_id": util._id}, bson.M{"$set": bson.M{"json_data": jsonData}}, options.Update().SetUpsert(true))

	return err
}

func (util *PropUtil) Remove() error {
	_, err := util.col.DeleteOne(context.Background(), bson.M{"_id": util._id})
	return err
}
