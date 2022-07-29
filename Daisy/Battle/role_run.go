package main

import (
	"Daisy/ErrorCode"
	"Daisy/Proto"
)

func (user *_User) RPC_LootRunChest() (int32, *Proto.Items) {
	return user.role.LootRunChest()
}

func (r *_Role) LootRunChest() (int32, *Proto.Items) {
	team := r.GetSpace().(*_Team)
	if team.runChest != nil {
		return team.OnLootRunChest(r.GetOwnerUserID())
	}

	r.Error("LootRunChest err:not in Runchesting")
	return ErrorCode.Failure, nil
}
