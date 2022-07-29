package main

import (
	"Cinder/Base/Const"
	"Daisy/Data"
	"Daisy/ErrorCode"
	"Daisy/Proto"
	log "github.com/cihub/seelog"
	"math"
	"time"
)

//获得新称号
func (r *_Role) AddTitle(titleID uint32) int32 {
	title, ok := Data.GetTitleConfig().Title_ConfigItems[titleID]
	if !ok {
		return ErrorCode.Failure
	}

	startTime := time.Now().Unix()
	overTime := startTime + int64(title.TitleDuration*86400)

	//若持续时间为0则为永久称号
	if title.TitleDuration == 0 {
		overTime = 0
	}
	r.prop.SyncAddTitle(titleID, startTime, overTime)

	//更新计时器
	r.updateTitleTimer()

	//空称号直接将称号作为佩戴称号
	if r.prop.Data.Title.TitleID == 0 {
		r.prop.SyncUpdateTitle(titleID)
		user := r.GetOwnerUser()
		if user == nil {
			r.Error("[findNewTitle] role's user is nil")
		} else {
			user.Rpc(Const.Game, "RPC_UpdateChatUserActivateData")
		}
		return ErrorCode.Success
	}

	//有佩戴称号比较等级
	roleTitle, ok := Data.GetTitleConfig().Title_ConfigItems[r.prop.Data.Title.TitleID]
	if !ok {
		return ErrorCode.Failure
	}

	//比较新旧称号等级，若新称号等级较高直接替换
	if roleTitle.Level <= title.Level {
		r.prop.SyncUpdateTitle(titleID)
		user := r.GetOwnerUser()
		if user == nil {
			r.Error("[findNewTitle] role's user is nil")
		} else {
			user.Rpc(Const.Game, "RPC_UpdateChatUserActivateData")
		}
	}

	return ErrorCode.Success
}

//检测过期称号
func (r *_Role) onTitleOverTime() {
	log.Debug("StartTitleOverTime")
	//称号表为空时直接返回
	if len(r.prop.Data.Title.TitleList) == 0 {
		log.Debug("NoTitle")
		return
	}

	timeNow := time.Now().Unix()
	var deleteTitleIDList Proto.Int32Array
	//遍历表找到已超时的称号
	for ID, title := range r.prop.Data.Title.TitleList {
		if title.LostTime != 0 && title.LostTime <= timeNow {
			deleteTitleIDList.Data = append(deleteTitleIDList.Data, int32(ID))
		}
	}

	if len(deleteTitleIDList.Data) == 0 {
		r.updateTitleTimer()
		return
	}

	//删除超时称号
	r.prop.SyncClearOverTimeTitle(&deleteTitleIDList)

	//更新计时器
	r.updateTitleTimer()
	r.recheckCurTitle()
}

// 检测佩戴称号是否过期
func (r *_Role) recheckCurTitle() {
	_, ok := r.prop.Data.Title.TitleList[r.prop.Data.Title.TitleID]
	//若佩戴称号过期寻找新称号
	if !ok {
		r.findNewTitle()
	}
	return
}

//佩戴称号超时时寻找新称号

func (r *_Role) findNewTitle() {
	var maxLevel, maxTitleID uint32
	var maxStartTime int64

	for ID, title := range r.prop.Data.Title.TitleList {
		titleData := Data.GetTitleConfig().Title_ConfigItems[ID]
		//如果满足等级更高或等级相同时间更新则视为新称号候选
		if titleData.Level > maxLevel || (titleData.Level == maxLevel && title.StartTime > maxStartTime) {
			maxLevel = titleData.Level
			maxStartTime = title.StartTime
			maxTitleID = ID
		}
	}

	//约定ID为0时无称号
	r.prop.SyncUpdateTitle(maxTitleID)
	user := r.GetOwnerUser()
	if user == nil {
		r.Error("[findNewTitle] role's user is nil")
	} else {
		user.Rpc(Const.Game, "RPC_UpdateChatUserActivateData")
	}
}

//更新称号超时计时器
func (r *_Role) updateTitleTimer() {
	//遍历称号列表找到最近的超时时间
	nextLostTime := int64(math.MaxInt64)
	for _, titleInfo := range r.prop.Data.Title.TitleList {
		if titleInfo.LostTime < nextLostTime && titleInfo.LostTime != 0 {
			nextLostTime = titleInfo.LostTime
		}
	}

	//检测是否已有称号超时计时器运行
	if r.titleTimer != nil {
		r.titleTimer.Stop()
	}

	//检测是否有限时称号，若有则更新计时器
	if nextLostTime == 0 || nextLostTime == int64(math.MaxInt64) {
		return
	} else {
		timeDuartion := time.Second * time.Duration(nextLostTime-time.Now().Unix())
		r.titleTimer = time.AfterFunc(timeDuartion, func() {
			r.titleChan <- true
		})
	}
}
