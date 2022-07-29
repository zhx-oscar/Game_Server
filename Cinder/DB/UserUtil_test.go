package DB

//import (
//	"go.mongodb.org/mongo-driver/bson/primitive"
//	"reflect"
//	"testing"
//)
//
//func TestUserUtil_Base(t *testing.T) {
//	userObjID := primitive.NewObjectID()
//	user := NewUser(userObjID)
//	//user.Data = []byte("hello world")
//	user.Auth = &UserAuth{
//		AccountName: userObjID.Hex(),
//		Password:    "unknown",
//		Data:        "bmw",
//	}
//
//	util, err := NewUserUtil(userObjID.Hex())
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if err = util.Insert(user); err != nil {
//		t.Fatal(err)
//	}
//
//	var newUser *User
//	if newUser, err = util.GetUser(); err != nil {
//		t.Fatal(err)
//	}
//	if !reflect.DeepEqual(user, newUser) {
//		t.Fatal("user mismatch")
//	}
//
//	var newAuth *UserAuth
//	if newAuth, err = util.GetAuth(); err != nil {
//		t.Fatal(err)
//	}
//	if !reflect.DeepEqual(user.Auth, newAuth) {
//		t.Fatal("auth mismatch")
//	}
//
//	var newID string
//	if newID, newAuth, err = util.GetAuthByAccName(userObjID.Hex()); err != nil {
//		t.Fatal(err)
//	}
//	if !reflect.DeepEqual(user.Auth, newAuth) {
//		t.Fatal("auth mismatch")
//	}
//	if newID != user.ID.Hex() {
//		t.Fatal("user id mismatch")
//	}
//}
