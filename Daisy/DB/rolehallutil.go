package DB

import (
	"Cinder/DB"
	"Daisy/Proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RoleProtoWrap struct {
	ID   primitive.ObjectID `bson:"_id"`
	Role *Proto.Role        `bson:"json_data"`
}

type IRoleHallUtil interface {
	Find(query interface{}, limit, skip int) ([]*RoleProtoWrap, error)
	FindID(roleID string) (*RoleProtoWrap, error)
	FindIDs(roleIDs []string) ([]*RoleProtoWrap, error)
	FindUID(roleUID uint64) (*RoleProtoWrap, error)
	Aggregate(pipeline interface{}) ([]*RoleProtoWrap, error)
}

var roleHallTabName = "RoleProp_tbl"

var defaultRoleHallUtil IRoleHallUtil

type roleHallUtil struct {
}

func init() {
	defaultRoleHallUtil = &roleHallUtil{}
}

func GetRoleHallUitl() IRoleHallUtil {
	return defaultRoleHallUtil
}

func (util *roleHallUtil) FindID(roleID string) (*RoleProtoWrap, error) {
	if DB.MongoDB == nil {
		return nil, dbInvalidErr
	}

	id, err := primitive.ObjectIDFromHex(roleID)
	if err != nil {
		return nil, err
	}

	record := &RoleProtoWrap{}
	if err = DB.MongoDB.Collection(roleHallTabName).FindOne(nil, bson.M{"_id": id}).Decode(record); err != nil {
		return nil, err
	} else {
		return record, nil
	}
}

func (util *roleHallUtil) FindIDs(roleIDs []string) ([]*RoleProtoWrap, error) {
	if DB.MongoDB == nil {
		return nil, dbInvalidErr
	}

	objectIds := make([]primitive.ObjectID, 0)
	for i := 0; i < len(roleIDs); i++ {
		if id, err := primitive.ObjectIDFromHex(roleIDs[i]); err == nil {
			objectIds = append(objectIds, id)
		}
	}

	findParam := bson.M{"_id": bson.M{"$in": objectIds}}
	if cur, err := DB.MongoDB.Collection(roleHallTabName).Find(nil, findParam); err != nil {
		return nil, err
	} else {
		records := make([]*RoleProtoWrap, 0)
		if err = cur.All(nil, &records); err != nil {
			return nil, err
		}

		return records, nil
	}
}

func (util *roleHallUtil) FindUID(roleUID uint64) (*RoleProtoWrap, error) {
	if DB.MongoDB == nil {
		return nil, dbInvalidErr
	}

	findParam := bson.M{"json_data.base.uid": roleUID}
	record := &RoleProtoWrap{}
	if err := DB.MongoDB.Collection(roleHallTabName).FindOne(nil, findParam).Decode(record); err != nil {
		return nil, err
	} else {
		return record, nil
	}
}

func (util *roleHallUtil) Aggregate(pipeline interface{}) ([]*RoleProtoWrap, error) {
	if DB.MongoDB == nil {
		return nil, dbInvalidErr
	}

	if cur, err := DB.MongoDB.Collection(roleHallTabName).Aggregate(nil, pipeline); err != nil {
		return nil, err
	} else {
		records := make([]*RoleProtoWrap, 0)
		if err = cur.All(nil, &records); err != nil {
			return nil, err
		}

		return records, nil
	}
}
func (util *roleHallUtil) Find(query interface{}, limit, skip int) ([]*RoleProtoWrap, error) {
	if DB.MongoDB == nil {
		return nil, dbInvalidErr
	}

	if cur, err := DB.MongoDB.Collection(teamHallTabName).Find(nil, query, options.Find().SetLimit(int64(limit)), options.Find().SetSkip(int64(skip))); err != nil {
		return nil, err
	} else {
		records := make([]*RoleProtoWrap, 0)
		if err = cur.All(nil, &records); err != nil {
			return nil, err
		}

		return records, nil
	}
}
