package main

import (
	"Cinder/Base/User"
	"Daisy/ErrorCode"
)

func (team *_Team) SendNotify(notifyID uint32, notifier, args string, isLong, showLong bool) int32 {
	notifyUsers := make([]User.IUser, 0)
	if notifier == "" {
		team.TraversalUser(func(iu User.IUser) bool {
			notifyUsers = append(notifyUsers, iu)
			return true
		})
	} else {
		iu, err := team.GetUser(notifier)
		if err == nil {
			notifyUsers = append(notifyUsers, iu)
		}
	}

	for _, iu := range notifyUsers {
		user := iu.(*_User)
		if isLong {
			if showLong {
				user.ShowCliNotify(notifyID, args)
			} else {
				user.HideCliNotify(notifyID)
			}
		} else {
			user.SendCliNotify(notifyID, args)
		}
	}

	return ErrorCode.Success
}
