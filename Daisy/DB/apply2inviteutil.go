package DB

import (
	"Cinder/DB"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Apply2Invite struct {
	ID          primitive.ObjectID `bson:"_id"`
	Type        uint8              //0申请  1邀请
	Time        int64
	RoleID      string
	TeamID      string
	ApplyParam  *ApplyParam
	InviteParam *InviteParam
}

type ApplyParam struct {
	OriginTeamID string
	Message      string
}

type InviteParam struct {
	Instigator string
}

var apply2InviteTabName = "Apply2Invite"

func init() {
	DB.MongoDB.Collection(apply2InviteTabName).Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.M{"roleid": 1, "teamid": 1, "type": 1},
		Options: options.Index().SetUnique(true),
	})
}

type apply2InviteUtil struct {
}

func GetApply2InviteUtil() *apply2InviteUtil {
	return &apply2InviteUtil{}
}

func (util *apply2InviteUtil) RemoveByID(ids []primitive.ObjectID) error {
	filter := bson.M{"_id": bson.M{"$in": ids}}
	_, err := DB.MongoDB.Collection(apply2InviteTabName).DeleteMany(context.TODO(), filter)
	return err
}

func (util *apply2InviteUtil) ApplyIsExist(roleID, teamID string) (bool, error) {
	filter := bson.M{"roleid": roleID, "teamid": teamID, "type": 0}
	num, err := DB.MongoDB.Collection(apply2InviteTabName).CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}
	return num > 0, nil
}

func (util *apply2InviteUtil) GetApplysInRole(roleID string) ([]*Apply2Invite, error) {
	return util.find(bson.M{"type": 0, "roleid": roleID})
}

func (util *apply2InviteUtil) GetApplysInTeam(teamID string) ([]*Apply2Invite, error) {
	return util.find(bson.M{"type": 0, "teamid": teamID})
}

func (util *apply2InviteUtil) GetApply(roleID, teamID string) (*Apply2Invite, error) {
	return util.findOne(bson.M{"type": 0, "roleid": roleID, "teamid": teamID})
}

func (util *apply2InviteUtil) AddApply(roleID, teamID, originTeamID, message string) error {
	data := &Apply2Invite{
		ID:     primitive.NewObjectID(),
		Type:   0,
		Time:   time.Now().Unix(),
		RoleID: roleID,
		TeamID: teamID,
		ApplyParam: &ApplyParam{
			OriginTeamID: originTeamID,
			Message:      message,
		},
	}
	_, err := DB.MongoDB.Collection(apply2InviteTabName).InsertOne(context.TODO(), data)
	return err
}

func (util *apply2InviteUtil) RemoveApply(roleID, teamID string) error {
	filter := bson.M{"type": 0, "roleid": roleID, "teamid": teamID}
	_, err := DB.MongoDB.Collection(apply2InviteTabName).DeleteOne(context.TODO(), filter)
	return err
}

func (util *apply2InviteUtil) RemoveAllApplysInRole(roleID string) ([]*Apply2Invite, error) {
	filter := bson.M{"type": 0, "roleid": roleID}
	cur, err := DB.MongoDB.Collection(apply2InviteTabName).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var results []*Apply2Invite
	err = cur.All(nil, &results)
	if err != nil {
		return nil, err
	}

	ids := make([]primitive.ObjectID, len(results), len(results))
	for i := 0; i < len(results); i++ {
		ids[i] = results[i].ID
	}

	filter = bson.M{"_id": bson.M{"$in": ids}}
	_, err = DB.MongoDB.Collection(apply2InviteTabName).DeleteMany(context.TODO(), filter)
	return results, err
}

func (util *apply2InviteUtil) RemoveAllApplysInTeam(teamID string) ([]*Apply2Invite, error) {
	filter := bson.M{"type": 0, "teamid": teamID}
	cur, err := DB.MongoDB.Collection(apply2InviteTabName).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var results []*Apply2Invite
	err = cur.All(nil, &results)
	if err != nil {
		return nil, err
	}

	ids := make([]primitive.ObjectID, len(results), len(results))
	for i := 0; i < len(results); i++ {
		ids[i] = results[i].ID
	}

	filter = bson.M{"_id": bson.M{"$in": ids}}
	_, err = DB.MongoDB.Collection(apply2InviteTabName).DeleteMany(context.TODO(), filter)
	return results, err
}

func (util *apply2InviteUtil) InviteIsExist(roleID, teamID string) (bool, error) {
	filter := bson.M{"roleid": roleID, "teamid": teamID, "type": 1}
	num, err := DB.MongoDB.Collection(apply2InviteTabName).CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}
	return num > 0, nil
}

func (util *apply2InviteUtil) GetInvitesInRole(roleID string) ([]*Apply2Invite, error) {
	return util.find(bson.M{"type": 1, "roleid": roleID})
}

func (util *apply2InviteUtil) GetInvitesInTeam(teamID string) ([]*Apply2Invite, error) {
	return util.find(bson.M{"type": 1, "teamid": teamID})
}

func (util *apply2InviteUtil) GetInvite(roleID, teamID string) (*Apply2Invite, error) {
	return util.findOne(bson.M{"type": 1, "roleid": roleID, "teamid": teamID})
}

func (util *apply2InviteUtil) AddInvite(roleID, teamID, instigator string) error {
	data := &Apply2Invite{
		ID:     primitive.NewObjectID(),
		Type:   1,
		Time:   time.Now().Unix(),
		RoleID: roleID,
		TeamID: teamID,
		InviteParam: &InviteParam{
			Instigator: instigator,
		},
	}
	_, err := DB.MongoDB.Collection(apply2InviteTabName).InsertOne(context.TODO(), data)
	return err
}

func (util *apply2InviteUtil) RemoveInvite(roleID, teamID string) error {
	filter := bson.M{"type": 1, "roleid": roleID, "teamid": teamID}
	_, err := DB.MongoDB.Collection(apply2InviteTabName).DeleteOne(context.TODO(), filter)
	return err
}

func (util *apply2InviteUtil) RemoveAllInvitesInRole(roleID string) ([]*Apply2Invite, error) {
	filter := bson.M{"type": 1, "roleid": roleID}
	cur, err := DB.MongoDB.Collection(apply2InviteTabName).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var results []*Apply2Invite
	err = cur.All(nil, &results)
	if err != nil {
		return nil, err
	}

	ids := make([]primitive.ObjectID, len(results), len(results))
	for i := 0; i < len(results); i++ {
		ids[i] = results[i].ID
	}

	filter = bson.M{"_id": bson.M{"$in": ids}}
	_, err = DB.MongoDB.Collection(apply2InviteTabName).DeleteMany(context.TODO(), filter)
	return results, err
}

func (util *apply2InviteUtil) RemoveAllInvitesInTeam(teamID string) ([]*Apply2Invite, error) {
	filter := bson.M{"type": 1, "teamid": teamID}
	cur, err := DB.MongoDB.Collection(apply2InviteTabName).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var results []*Apply2Invite
	err = cur.All(nil, &results)
	if err != nil {
		return nil, err
	}

	ids := make([]primitive.ObjectID, len(results), len(results))
	for i := 0; i < len(results); i++ {
		ids[i] = results[i].ID
	}

	filter = bson.M{"_id": bson.M{"$in": ids}}
	_, err = DB.MongoDB.Collection(apply2InviteTabName).DeleteMany(context.TODO(), filter)
	return results, err
}

func (util *apply2InviteUtil) find(filter interface{}) ([]*Apply2Invite, error) {
	cur, err := DB.MongoDB.Collection(apply2InviteTabName).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var results []*Apply2Invite
	err = cur.All(nil, &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (util *apply2InviteUtil) findOne(filter interface{}) (*Apply2Invite, error) {
	result := &Apply2Invite{}
	err := DB.MongoDB.Collection(apply2InviteTabName).FindOne(context.TODO(), filter).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
