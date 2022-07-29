package Space

import (
	"Cinder/Base/Const"
	"Cinder/Base/Core"
	BaseUser "Cinder/Base/User"
	"Cinder/Base/Util"
	"errors"
	"fmt"
	"reflect"
	"sync"

	log "github.com/cihub/seelog"
)

type _Server struct {
	Core.ICore

	rpcProc interface{}

	spaces      sync.Map
	userToSpace sync.Map

	spacePT reflect.Type
	userPT  IUser
}

func newServer() *_Server {
	srv := &_Server{
		ICore: Core.New(),
	}

	return srv
}

type _InitInfo struct {
	ID           string
	Type         string
	OwnerUserID  string
	Owner        interface{}
	OwnerRealPtr interface{}
	RealPtr      interface{}
	UserPT       interface{}
	PropData     []byte
	UserData     interface{}
}

type _IInit interface {
	Init()
}

type _IDestroy interface {
	Destroy()
}

func (srv *_Server) InitSrv(areaID string, serverID string, spacePT ISpace, userPT IUser, rpcProc interface{}) error {

	info := Core.NewDefaultInfo()

	info.RpcProc = rpcProc
	info.ServiceType = Const.Space
	info.AreaID = areaID
	info.ServiceID = fmt.Sprintf("%s_%s_%s", info.ServiceType, areaID, serverID)

	srv.spacePT = reflect.ValueOf(spacePT).Elem().Type()
	srv.userPT = userPT

	if err := srv.Init(info); err != nil {
		return err
	}

	return nil
}

func (srv *_Server) DestroySrv() {

	if srv.rpcProc != nil {
		ii, ok := srv.rpcProc.(_IDestroy)
		if ok {
			ii.Destroy()
		}
	}

	srv.destroyAllSpaces()
	if srv.ICore != nil {
		srv.ICore.Destroy()
	}
}

func (srv *_Server) constructSpace() interface{} {
	s := reflect.New(srv.spacePT).Interface()
	return s
}

func (srv *_Server) CreateSpace(id string, spacePropData []byte, userData interface{}) string {
	spaceID := id
	if spaceID == "" {
		spaceID = fmt.Sprint("space_", Util.GetGUID())
	}

	space := srv.constructSpace()
	is := space.(_ISpace)

	info := &_InitInfo{
		ID:       spaceID,
		Owner:    srv,
		RealPtr:  space,
		UserPT:   srv.userPT,
		PropData: spacePropData,
		UserData: userData,
	}

	is.onInit(info)
	srv.spaces.Store(spaceID, is)
	return spaceID
}

func (srv *_Server) DestroySpace(id string) error {

	log.Info("DestroySpace ", id)

	is, ok := srv.spaces.Load(id)
	if !ok {
		return errors.New("no space exist " + id)
	}

	srv.spaces.Delete(id)
	is.(_ISpace).setDestroyFlag()

	return nil
}

func (srv *_Server) destroyAllSpaces() {
	srv.spaces = sync.Map{}

	srv.spaces.Range(func(key, value interface{}) bool {
		value.(_ISpace).setDestroyFlag()
		return true
	})
}

func (srv *_Server) GetSpace(id string) (ISpace, error) {

	is, ok := srv.spaces.Load(id)
	if !ok {
		return nil, errors.New("no space exist " + id)
	}

	return is.(ISpace), nil
}

func (srv *_Server) TraversalSpace(cb func(space ISpace)) {

	srv.spaces.Range(func(key, value interface{}) bool {
		s := value.(ISpace)
		cb(s)
		return true
	})
}

func (srv *_Server) EnterSpace(userID string, spaceID string) error {

	is, err := srv.GetSpace(spaceID)
	if err != nil {
		return err
	}

	space := is.(_ISpace)

	loadSpaceID, ok := srv.userToSpace.Load(userID)
	if ok {

		if loadSpaceID == space.GetID() {
			log.Warn("EnterSpace the user have in space, needn't enter userID ", userID, " SpaceID ", spaceID)
			return nil
		} else {
			log.Warn("EnterSpace user had in space, leave space first userID ", userID, " SpaceID ", spaceID)
			srv.LeaveSpace(userID)
		}
	}

	ret := <-space.SafeCall("EnterSpace", userID)
	if ret.Err != nil || ret.Ret[0] != nil {
		return err
	}

	srv.userToSpace.Store(userID, space.GetID())

	return nil
}

func (srv *_Server) LeaveSpace(userID string) error {

	defer srv.userToSpace.Delete(userID)

	spaceID, ok := srv.userToSpace.Load(userID)
	if !ok {
		return errors.New("couldn't find space for user " + userID)
	}

	s, err := srv.GetSpace(spaceID.(string))
	if err != nil {
		return errors.New("no space exist yet")
	}

	r := <-s.(_ISpace).SafeCall("LeaveSpace", userID)
	if r.Err != nil {
		return r.Err
	} else {
		if r.Ret[0] != nil {
			return errors.New("leave space failed")
		}
	}

	log.Info("LeaveSpace succeed userID ", userID, " SpaceID ", spaceID)
	return nil
}

func (srv *_Server) DestroyUser(id string) error {
	return srv.LeaveSpace(id)
}

func (srv *_Server) CleanUserToSpaceMap(userID string) {
	srv.userToSpace.Delete(userID)
}

func (srv *_Server) GetUser(id string) (BaseUser.IUser, error) {
	spaceID, ok := srv.userToSpace.Load(id)
	if !ok {
		return nil, errors.New("no space exist ")
	}

	space, err := srv.GetSpace(spaceID.(string))
	if err != nil {
		return nil, err
	}

	ii, err := space.GetUser(id)
	if err != nil {
		return nil, err
	}

	return ii.(BaseUser.IUser), nil
}

func (srv *_Server) Traversal(cb func(user BaseUser.IUser) bool) {
	srv.spaces.Range(func(key, value interface{}) bool {
		s := value.(_ISpace)
		s.TraversalUser(cb)
		return true
	})
}

func (srv *_Server) GetSrvInst() Core.ICore {
	return srv
}
