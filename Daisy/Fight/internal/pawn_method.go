package internal

import (
	"Cinder/Base/linemath"
	"Daisy/Proto"
)

// Snapshot 快照pawn战斗数据
func (pawn *Pawn) Snapshot() *Pawn {
	snapshot := &Pawn{
		UID:        pawn.UID,
		Info:       pawn.Info,
		Scene:      pawn.Scene,
		isSnapshot: true,
	}

	snapshot.Attr = *pawn.Attr.copy(snapshot)
	snapshot.State = *pawn.State.copy(snapshot)

	return snapshot
}

// IsAlive 是否存活
func (pawn *Pawn) IsAlive() bool {
	return !pawn.State.Death
}

// IsRole 是否是玩家
func (pawn *Pawn) IsRole() bool {
	return pawn.Info.Type == Proto.PawnType_Role
}

// IsNpc 是否是Npc
func (pawn *Pawn) IsNpc() bool {
	return pawn.Info.Type == Proto.PawnType_Npc
}

// IsBackground 是否是背景
func (pawn *Pawn) IsBackground() bool {
	return pawn.Info.Type == Proto.PawnType_BG
}

// IsBoss 是否是Boss
func (pawn *Pawn) IsBoss() bool {
	if pawn.Info.Npc == nil {
		return false
	}

	return pawn.Info.Npc.IsBoss
}

// GetCamp 查询阵营
func (pawn *Pawn) GetCamp() Proto.Camp_Enum {
	return pawn.Info.Camp
}

// GetEmemyCamp 查询敌方阵营
func (pawn *Pawn) GetEmemyCamp() Proto.Camp_Enum {
	return GetEnemyCamp(pawn.Info.Camp)
}

//GetEnemyList 获取敌人列表
func (pawn *Pawn) GetEnemyList() []*Pawn {
	return pawn.Scene.formationList[pawn.GetEmemyCamp()].PawnList
}

//GetPartnerList 获取伙伴列表
func (pawn *Pawn) GetPartnerList() []*Pawn {
	return pawn.Scene.formationList[pawn.GetCamp()].PawnList
}

//OverlapCircleShape 自己对某点做碰撞检测
func (pawn *Pawn) OverlapCircleShape(pos linemath.Vector2) bool {
	//碰撞检测 此处位置已经被占
	targetList := pawn.Scene.overlapCircleShape(float64(pos.X), float64(pos.Y), float64(pawn.Attr.CollisionRadius))
	if len(targetList) > 0 {
		//检测如果碰撞其他对象正好相切可以站位
		for _, val := range targetList {
			//是自己
			if pawn.Equal(val) {
				continue
			}

			if Distance(val.GetPos(), pos) < val.Attr.CollisionRadius+pawn.Attr.CollisionRadius {
				return false
			}
		}
	}

	return true
}

// Equal pawn是否相同
func (pawn *Pawn) Equal(other *Pawn) bool {
	if pawn == nil || other == nil {
		return pawn == other
	}

	return pawn.UID == other.UID
}
