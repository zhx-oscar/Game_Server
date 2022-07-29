package internal

import "Daisy/Fight/internal/conf"

type BuffList []*Buff

// _PawnBuff pawn buff模块
type _PawnBuff struct {
	pawn         *Pawn
	buffTab      map[BuffKey]*Buff    // buff表
	sameBuffList map[uint32]*BuffList // 相同ID buff列表
	buffList     BuffList             // buff列表
}

// init 初始化
func (pawnBuff *_PawnBuff) init(pawn *Pawn) {
	pawnBuff.pawn = pawn

	// 创建buff表
	pawn.buffTab = make(map[BuffKey]*Buff)
	pawn.sameBuffList = make(map[uint32]*BuffList)
}

// BuffExist buff是否存在
func (pawn *Pawn) BuffExist(buffKey BuffKey) bool {
	_, ok := pawn.Scene.getBuff(pawn, buffKey)
	return ok
}

// GetBuff 查询buff
func (pawn *Pawn) GetBuff(buffKey BuffKey) (*Buff, bool) {
	return pawn.Scene.getBuff(pawn, buffKey)
}

// GetSameBuffs 查询相同id的buff列表
func (pawn *Pawn) GetSameBuffs(mainID uint32) ([]*Buff, bool) {
	return pawn.Scene.getSameBuffs(pawn, mainID)
}

// AddBuff 添加buff
func (pawn *Pawn) AddBuff(caster *Pawn, mainID uint32, extValue int64) (*Buff, bool) {
	return pawn.Scene.addBuff(pawn, caster, mainID, extValue)
}

// RemoveBuff 删除buff
func (pawn *Pawn) RemoveBuff(buffKey BuffKey) bool {
	buff, ok := pawn.Scene.getBuff(pawn, buffKey)
	if !ok {
		return false
	}

	return pawn.Scene.removeBuff(pawn, buff, false)
}

// ClearBuffKind 清除指定类型buff
func (pawn *Pawn) ClearBuffKind(kind conf.BuffKind, num int) int {
	return pawn.Scene.clearBuffKind(pawn, kind, num)
}

// ClearBuff 驱散buff
func (pawn *Pawn) ClearBuff(clearGroup uint32) int {
	return pawn.Scene.clearBuff(pawn, clearGroup)
}

// RefreshBuffDuration 刷新buff时长
func (pawn *Pawn) RefreshBuffDuration(buffKey BuffKey) bool {
	return pawn.Scene.refreshBuffDuration(pawn, buffKey)
}

// ExtendBuffDuration 延长buff时长
func (pawn *Pawn) ExtendBuffDuration(buffKey BuffKey, duration uint32) bool {
	return pawn.Scene.extendBuffDuration(pawn, buffKey, duration)
}

// BatchAddBuffs 批量添加buff
func (pawn *Pawn) BatchAddBuffs(caster *Pawn, mainIDs []uint32, extValue int64) []*Buff {
	return pawn.Scene.batchAddBuffs(pawn, caster, mainIDs, extValue)
}

// BatchRemoveBuffs 批量删除buff
func (pawn *Pawn) BatchRemoveBuffs(mainIDs []uint32, casterID uint32) int {
	return pawn.Scene.batchRemoveBuffs(pawn, mainIDs, casterID)
}
