package DB

import (
	"Cinder/DB"
	"Daisy/Proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TeamProtoWrap struct {
	ID   primitive.ObjectID `bson:"_id"`
	Team *Proto.Team        `bson:"json_data"`
}

type ITeamHallUtil interface {
	Find(query interface{}, limit, skip int) ([]*TeamProtoWrap, error)
	FindID(teamID string) (*TeamProtoWrap, error)
	FindIDs(teamIDs []string) ([]*TeamProtoWrap, error)
	FindUID(teamUID uint64) (*TeamProtoWrap, error)
	Aggregate(pipeline interface{}) ([]*TeamProtoWrap, error)
}

var teamHallTabName = "TeamProp_tbl"

var defaultTeamHallUtil ITeamHallUtil

type teamHallUtil struct {
}

func init() {
	defaultTeamHallUtil = &teamHallUtil{}
}

func GetTeamHallUitl() ITeamHallUtil {
	return defaultTeamHallUtil
}

func (util *teamHallUtil) Find(query interface{}, limit, skip int) ([]*TeamProtoWrap, error) {
	if DB.MongoDB == nil {
		return nil, dbInvalidErr
	}

	if cur, err := DB.MongoDB.Collection(teamHallTabName).Find(nil, query, options.Find().SetLimit(int64(limit)), options.Find().SetSkip(int64(skip))); err != nil {
		return nil, err
	} else {
		records := make([]*TeamProtoWrap, 0)
		if err = cur.All(nil, &records); err != nil {
			return nil, err
		}

		return records, nil
	}
}

func (util *teamHallUtil) FindID(teamID string) (*TeamProtoWrap, error) {
	if DB.MongoDB == nil {
		return nil, dbInvalidErr
	}

	id, err := primitive.ObjectIDFromHex(teamID)
	if err != nil {
		return nil, err
	}

	record := &TeamProtoWrap{}
	if err = DB.MongoDB.Collection(teamHallTabName).FindOne(nil, bson.M{"_id": id}).Decode(record); err != nil {
		return nil, err
	} else {
		return record, nil
	}
}

func (util *teamHallUtil) FindUID(teamUID uint64) (*TeamProtoWrap, error) {
	if DB.MongoDB == nil {
		return nil, dbInvalidErr
	}

	findParam := bson.M{"json_data.base.uid": teamUID}
	record := &TeamProtoWrap{}
	if err := DB.MongoDB.Collection(teamHallTabName).FindOne(nil, findParam).Decode(record); err != nil {
		return nil, err
	} else {
		return record, nil
	}
}

func (util *teamHallUtil) FindIDs(teamIDs []string) ([]*TeamProtoWrap, error) {
	if DB.MongoDB == nil {
		return nil, dbInvalidErr
	}

	objectIds := make([]primitive.ObjectID, 0)
	for i := 0; i < len(teamIDs); i++ {
		if id, err := primitive.ObjectIDFromHex(teamIDs[i]); err == nil {
			objectIds = append(objectIds, id)
		}
	}

	findParam := bson.M{"_id": bson.M{"$in": objectIds}}
	if cur, err := DB.MongoDB.Collection(teamHallTabName).Find(nil, findParam); err != nil {
		return nil, err
	} else {
		records := make([]*TeamProtoWrap, 0)
		if err = cur.All(nil, &records); err != nil {
			return nil, err
		}

		return records, nil
	}
}

func (util *teamHallUtil) Aggregate(pipeline interface{}) ([]*TeamProtoWrap, error) {
	if DB.MongoDB == nil {
		return nil, dbInvalidErr
	}

	if cur, err := DB.MongoDB.Collection(teamHallTabName).Aggregate(nil, pipeline); err != nil {
		return nil, err
	} else {
		records := make([]*TeamProtoWrap, 0)
		if err = cur.All(nil, &records); err != nil {
			return nil, err
		}

		return records, nil
	}
}
