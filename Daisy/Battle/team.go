package main

import (
	"Cinder/Base/User"
	"Cinder/Space"
	"Daisy/Const"
	"Daisy/DB"
	"Daisy/Prop"
	"Daisy/Proto"
	log "github.com/cihub/seelog"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const AccumulateTime = 24 // 离线收益累计积累时间（h）超过该时段不再累积收益

const requestSkillCountLimit uint32 = 1 // 乞求技能最大次数

const giveSkillCountLimit uint32 = 12 // 赠送技能最大次数

type _Team struct {
	Space.Space
	prop *Prop.TeamProp

	// loopTicker 内部逻辑帧, 每秒执行一次就够了
	loopTicker *time.Ticker

	// 队伍有效性相关
	expireTicker *time.Ticker
	offlineTimer *time.Timer

	// 定时任务
	cronJob         *cron.Cron
	cronJobChan     chan *_CronJob
	dailyResetJobID cron.EntryID

	*_RunModel
	runChest       *_RunChestModel
	raidBattle     *_RaidBattleModel
	raidBattleDrop *_RaidBattleDropModel
	fastBattle     *_FastBattleModel

	PendingLeaveUser     map[string]int64
	PlaceHolder          map[string]int64
	PendingTransferEvent []func()

	*_TeamStateMgr
	PendingStateEvent []*_PendingStateEvent

	hall *_TeamHallModel
}

type _PendingStateEvent struct {
	State uint8
	Event func()
}

func (team *_Team) Init() {
	team.RegisterActor(RoleActorType, &_Role{})

	team.prop = team.GetProp().(*Prop.TeamProp)
	team.loopTicker = time.NewTicker(time.Second)
	team.cronJob = cron.New(cron.WithSeconds())
	team.cronJobChan = make(chan *_CronJob, 10)
	team.addDailyResetJob()
	team.cronJob.Start()

	// 把队伍成员创建出来
	team.loadTeamRoles()

	if !team.prop.Data.InitedOnce {
		team.prop.SyncInitedOnce()
		team.initOnce()
	}

	team._RunModel = NewRunModel()
	team.fastBattle = NewFastBattleModel()

	// 队伍注册
	team.expireTicker = time.NewTicker(5 * time.Second)
	DB.TeamUtil().Register(team.GetID(), Space.Inst.GetServiceID())

	team.checkAndResetMemberDaily()

	team.PendingLeaveUser = make(map[string]int64)
	team.PlaceHolder = make(map[string]int64)
	team.PendingTransferEvent = make([]func(), 0)

	team._TeamStateMgr = &_TeamStateMgr{}
	team.RegState(map[uint8]ITeamState{
		TeamState_Running:       NewRunState(team),
		TeamState_Raidbattling:  NewRaidBattleState(team),
		TeamState_FastBattleing: NewFastBattleState(team),
	})
	team.PendingStateEvent = make([]*_PendingStateEvent, 0)
	team.SetState(TeamState_Running, team.RandomSpawnIdx())

	team.initApplys()
	team.calcTeamOfflineAward()

	team.hall = NewTeamHallModel()

	team.checkAndResetFastBattle()
	team.checkAndResetSupply()

	team.CheckSeasonChange()

	team.Info("Team Init")
}

func (team *_Team) Destroy() {
	team.cronJob.Stop()
	close(team.cronJobChan)

	team._RunModel.Destroy()
	team.expireTicker.Stop()
	DB.TeamUtil().UnRegister(team.GetID())
	team.loopTicker.Stop()
	team.prop.SyncSetLastLogoutTime(time.Now().Unix())
	team.teamChatDestroy()
	team.FlushToCache()
	team.Info("Team Destroy")
}

func (team *_Team) Loop() {
	team.LoopState(team.GetDeltaTime())
	team.LoopRunChest()

	select {
	case job := <-team.cronJobChan:
		team.handleCronJob(job)
	case <-team.expireTicker.C:
		DB.TeamUtil().UpdateExpire(team.GetID())
	case <-team.loopTicker.C:
		team.secondLoop()
	case res := <-team.hall.RecruitmentsChan:
		team.OnGetRecruitments(res)
	default:
	}
}

//initOnce 终身只执行一次
func (team *_Team) initOnce() {
	if uid, err := DB.FetchTeamUID(); err == nil {
		team.prop.SyncSetUID(uid)
		team.initTeamChatChannel()
	} else {
		panic(err)
	}

	team.FlushToCache()
}

func (team *_Team) secondLoop() {
	team.checkOffline()

	if team.GetState() != TeamState_Raidbattling && len(team.PendingTransferEvent) > 0 {
		for _, value := range team.PendingTransferEvent {
			value()
		}
		team.PendingTransferEvent = team.PendingTransferEvent[0:0]
	}

	//过期
	for key, value := range team.PlaceHolder {
		if time.Now().Unix()-value > 60*5 {
			team.Debugf("PlaceHolder %s 过期", key)
			team.RemovePlaceHolder(key)
		}
	}

	for key, value := range team.PendingLeaveUser {
		if time.Now().Unix()-value > 60*5 {
			team.Debugf("PendingLeaveUser %s 过期", key)
			team.ActorStopTransfer(key)
		}
	}

	for i := 0; i < len(team.PendingStateEvent); i++ {
		ev := team.PendingStateEvent[i]
		if team.CanSetState(ev.State) {
			ev.Event()

			team.PendingStateEvent = append(team.PendingStateEvent[:i], team.PendingStateEvent[i+1:]...)
			i--
		}
	}

	team.CheckSeasonChange()
}

func (team *_Team) handleCronJob(job *_CronJob) {
	switch job.Type {
	case CronJobDailyReset:
		team.doDailyResetJob()
	}
}

func (team *_Team) GetPropInfo() (string, string) {
	return Prop.TeamPropType, team.GetID()
}

func (team *_Team) loadTeamRoles() {
	for id := range team.prop.Data.Base.Members {
		team.AddActor(RoleActorType, id, id, nil, nil)
	}
}

func (team *_Team) checkOffline() {
	if team.GetState() == TeamState_Raidbattling {
		return
	}

	if len(team.PendingLeaveUser) > 0 {
		return
	}

	if len(team.PlaceHolder) > 0 {
		return
	}

	// 是否要下线的标识, 目前就看在线玩家数量是否为0
	// 后续添加更多规则
	// 当玩家数为0的时候, 启动定时器
	// 当有玩家上线的时候(EnterBattle)的时候, 关闭定时器

	userCount := 0
	team.TraversalUser(func(User.IUser) bool {
		userCount++
		return true
	})
	if (userCount == 0 || len(team.prop.Data.Base.Members) == 0) && team.offlineTimer == nil {
		team.Debug("Start offline timer")
		team.offlineTimer = time.AfterFunc(1*time.Minute, team.DestroySelf)
		go team.FlushToDB()
	}
}

func (team *_Team) GetUserCnt() uint32 {
	cnt := uint32(0)
	team.TraversalUser(func(user User.IUser) bool {
		cnt++
		return true
	})
	return cnt
}

func (team *_Team) AddPlaceHolder(userID string) {
	team.PlaceHolder[userID] = time.Now().Unix()

	if team.offlineTimer != nil {
		team.offlineTimer.Stop()
		team.offlineTimer = nil
	}
}

func (team *_Team) RemovePlaceHolder(userID string) {
	delete(team.PlaceHolder, userID)
}

func (team *_Team) ActorStartTransfer(userID string) {
	ia := team.GetActorByUserID(userID)
	if ia == nil {
		return
	}

	team.PendingLeaveUser[userID] = time.Now().Unix()

	if team.offlineTimer != nil {
		team.offlineTimer.Stop()
		team.offlineTimer = nil
	}
}

func (team *_Team) ActorStopTransfer(userID string) {
	delete(team.PendingLeaveUser, userID)
}

func (team *_Team) initApplys() {
	applys, err := DB.GetApply2InviteUtil().GetApplysInTeam(team.GetID())
	if err != nil {
		team.Error(err)
		return
	}

	if len(applys) == 0 {
		return
	}

	expireIDs := make([]primitive.ObjectID, 0)
	infos := make(map[string]*Proto.TeamApplyInfo)

	for i := 0; i < len(applys); i++ {
		apply := applys[i]
		if time.Now().Unix()-apply.Time >= Const.ApplyMaxLife {
			applys = append(applys[:i], applys[i+1:]...)
			i--

			expireIDs = append(expireIDs, apply.ID)
		} else {
			infos[apply.RoleID] = &Proto.TeamApplyInfo{
				Time:    apply.Time,
				Message: apply.ApplyParam.Message,
			}
		}
	}

	team.prop.Data.Applys = infos
	DB.GetApply2InviteUtil().RemoveByID(expireIDs)
}

func (team *_Team) calcTeamOfflineAward() {
	if team.prop.Data.Base.LastDestoryTime == 0 {
		log.Info("[calcTeamOfflineAward] 队伍离线时间为0")
		return
	}
	now := time.Now()
	progress := team.prop.Data.Raid.Progress

	team.TraversalActor(func(actor Space.IActor) {
		role := actor.(*_Role)
		if role.prop.Data.Base.LastLogoutTime == 0 {
			log.Info("[getTeamOfflineAward] 新账号，没有离线收益")
			return
		}

		teamLastOfflineTime := time.Unix(team.prop.Data.Base.LastDestoryTime, 0)
		roleLastOfflineTime := time.Unix(role.prop.Data.Base.LastLogoutTime, 0)
		maxOff := roleLastOfflineTime.Add(AccumulateTime * time.Hour)

		if teamLastOfflineTime.After(maxOff) {
			log.Info("[getTeamOfflineAward] 玩家离线时间超过24小时，不计算队伍离线收益")
			return
		}
		var offlineTime time.Duration
		if maxOff.After(now) {
			offlineTime = now.Sub(teamLastOfflineTime)
		} else {
			offlineTime = maxOff.Sub(teamLastOfflineTime)
		}

		if offlineTime > AccumulateTime*time.Hour {
			offlineTime = AccumulateTime * time.Hour
		}

		role.addOfflineAward(progress, offlineTime)
	})
}

func (team *_Team) addDailyResetJob() {
	var err error
	team.dailyResetJobID, err = team.cronJob.AddFunc("0 0 5 * * ?", func() {
		team.cronJobChan <- &_CronJob{
			Type: CronJobDailyReset,
			Arg:  nil,
		}
	})
	if err != nil {
		team.Error("addDailyResetJob err", err)
	}
}

func (team *_Team) doDailyResetJob() {
	nextResetTime := team.cronJob.Entry(team.dailyResetJobID).Next
	team.ResetFastBattle(nextResetTime)
	team.resetMemberDaily(nextResetTime)
	team.ResetSupply(nextResetTime)

}

func (team *_Team) resetMemberDaily(nextResetTime time.Time) {
	team.TraversalActor(func(actor Space.IActor) {
		role := actor.(*_Role)
		role.prop.SyncResetRoleDaily(nextResetTime.Unix(), requestSkillCountLimit, giveSkillCountLimit)
		role.prop.SyncResetBoxOpen()
		now := time.Now()
		if now.Weekday() == 1 {
			role.prop.SyncResetDiscountNum()
		}
	})

	team.prop.SyncResetMemberDaily(nextResetTime.Unix())
}

func (team *_Team) checkAndResetMemberDaily() {
	now := time.Now().Unix()
	resetTimestamp := team.prop.Data.Base.DailyResetTimestamp
	if now < resetTimestamp {
		return
	}
	nextResetTime := team.cronJob.Entry(team.dailyResetJobID).Next
	team.resetMemberDaily(nextResetTime)
}

func (team *_Team) GetTeamMemberInfo(id string) *Proto.TeamMemberInfo {
	if info, ok := team.prop.Data.Base.Members[id]; ok {
		return info
	}
	return nil
}

func (team *_Team) checkAndResetFastBattle() {
	if time.Now().Unix() > team.prop.Data.Raid.FastBattleResetTimestamp {
		team.ResetFastBattle(team.cronJob.Entry(team.dailyResetJobID).Next)
	}
}

func (team *_Team) ResetSupply(resetTime time.Time) {
	team.prop.SyncSetSupplyResetTimestamp(resetTime.Unix())
}

func (team *_Team) checkAndResetSupply() {
	now := time.Now()
	if now.Unix() > team.prop.Data.Supply.SupplyResetTimestamp {
		team.checkAndResetSupplyNum(now, team.prop.Data.Supply.SupplyResetTimestamp)
		team.ResetSupply(team.cronJob.Entry(team.dailyResetJobID).Next)
	}
}

func (team *_Team) checkAndResetSupplyNum(now time.Time, last int64) {
	t := time.Unix(last, 0)
	date := t.Weekday()
	if date == 0 {
		date = 7
	}
	diff := (8 - date) * (24 * 60 * 60)
	team.TraversalActor(func(actor Space.IActor) {
		role := actor.(*_Role)
		role.prop.SyncResetBoxOpen()
		if date == 1 || now.Unix()-last >= int64(diff) {
			role.prop.SyncResetDiscountNum()
		}
	})
}
