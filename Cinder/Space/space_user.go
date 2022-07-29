package Space

import (
	"Cinder/Base/Message"
	BaseUser "Cinder/Base/User"
	"errors"
)

func (space *Space) EnterSpace(userID string) error {

	_, createNew, err := space.userMgr.GetOrCreateUser(userID, space.GetRealPtr())
	if err != nil {
		return err
	}

	if !createNew {
		return errors.New("the user is existed")
	}

	return nil
}

func (space *Space) LeaveSpace(userID string) error {

	_, err := space.userMgr.GetUser(userID)
	if err != nil {
		return err
	}

	err = space.userMgr.DestroyUser(userID)
	if err != nil {
		return err
	}

	return nil
}

func (space *Space) onRemoveUser(user IUser) {
	space.removeUserFromAgentMap(user.(_IUser))
	space.refreshOwnerUser()

	space.Info("User leave space", user.GetID())
}

func (space *Space) onAddUser(user IUser) {
	space.refreshOwnerUser()

	space.Info("User enter space ", user.GetID())
}

func (space *Space) refreshOwnerUser() {
	oldUser := space.ownerUser

	if space.ownerUser != nil {
		if _, err := space.GetUser(space.ownerUser.GetID()); err != nil {
			space.ownerUser = nil
		}
	}

	if space.ownerUser != nil {
		if !space.ownerUser.(_IUser).IsClientNetOK() {
			space.ownerUser = nil
		}
	}

	if space.ownerUser == nil {
		space.TraversalUser(func(user BaseUser.IUser) bool {

			if user.(_IUser).IsClientNetOK() {
				space.ownerUser = user.(IUser)
			}

			return space.ownerUser == nil
		})
	}

	if space.ownerUser != oldUser && space.ownerUser != nil {
		msg := &Message.SpaceOwnerChange{
			UserID: space.ownerUser.GetID(),
		}

		space.SendToAllClient(msg)
	}

}

func (space *Space) GetUser(userID string) (IUser, error) {
	ii, err := space.userMgr.GetUser(userID)
	if err != nil {
		return nil, errors.New("no existed")
	}

	return ii.(IUser), nil
}

func (space *Space) TraversalUser(cb func(user BaseUser.IUser) bool) {
	space.userMgr.Traversal(cb)
}

func (space *Space) removeUserFromAgentMap(user _IUser) {
	agentID := user.GetAgentID()

	agentList, ok := space.userAgentMap[agentID]
	if ok {

		for i, a := range agentList {
			if a == user.GetID() {
				agentList = append(agentList[:i], agentList[i+1:]...)
				break
			}
		}

		space.userAgentMap[agentID] = agentList
	}
}

func (space *Space) onUserAgentChanged(userID string, oldAgentID string, newAgentID string) {

	if oldAgentID != "" {

		userList, ok := space.userAgentMap[oldAgentID]
		if !ok {
			space.Warn("couldn't find agent user list AgentID", oldAgentID)
		} else {
			for i, id := range userList {
				if id == userID {
					userList = append(userList[:i], userList[i+1:]...)
					break
				}
			}

			space.userAgentMap[oldAgentID] = userList
		}
	}

	if newAgentID != "" {

		userList, ok := space.userAgentMap[newAgentID]
		if !ok {
			userList = make([]string, 0, 10)
		}

		userList = append(userList, userID)
		space.userAgentMap[newAgentID] = userList
	}
}
