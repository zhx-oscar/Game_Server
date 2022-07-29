package internal

import (
	"Daisy/Fight/internal/log"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	b3 "github.com/magicsea/behavior3go"
	b3config "github.com/magicsea/behavior3go/config"
	b3core "github.com/magicsea/behavior3go/core"
	b3loader "github.com/magicsea/behavior3go/loader"
)

//ai全局管理

//创建一个行为树
var mapTrees sync.Map
var mapTreesByID sync.Map
var extMaps *b3.RegisterStructMaps

// LoadBehaviorTreeFile 加载行为树
func LoadBehaviorTreeFile(path string) {
	createFromProject(path)
}

// reload 加载
func reload() {
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Error(err)
			return
		}
		defer watcher.Close()

		done := make(chan bool)
		go func() {
			for {
				select {
				case ev := <-watcher.Events:
					log.Debugf("#######Watch %s Op %s#######", ev.Name, ev.Op)

					if strings.HasSuffix(ev.Name, "chaos.b3") {
						time.Sleep(time.Second)
						if ev.Op&fsnotify.Write == fsnotify.Write || ev.Op&fsnotify.Create == fsnotify.Create {
							log.Info("==========>Start reload behavior!")

							LoadBehaviorTreeFile(behaviorTreePath)
						}
					}
				case err := <-watcher.Errors:
					log.Error(err)
				}
			}
		}()

		if err = watcher.Add(aiWatchPath); err != nil {
			log.Error(err)
		} else {
			log.Debug("Watch behavior path")
		}

		<-done
	}()

}

//主树

// createFromProject 创建行为树
func createFromProject(path string) {
	//log.Debug("CreateFromProject bevtree:", path)

	config, ok := b3config.LoadRawProjectCfg(path)
	if !ok {
		log.Error("CreateFromProject fail:" + path)
		return
	}

	for _, tf := range config.Data.Trees {
		//log.Debug("**********CreateFromProject create bevtree:", tf.Title)
		tree := b3loader.CreateBevTreeFromConfig(&tf, extMaps)
		//tree.Print()
		mapTrees.Store(tf.Title, tree)
		mapTreesByID.Store(tf.ID, tree)
	}
	//log.Debug("CreateFromProject bevtree success!")
}

// getBevTree 获取行为树
func getBevTree(name string) *b3core.BehaviorTree {
	t, ok := mapTrees.Load(name)
	if ok {
		return t.(*b3core.BehaviorTree)
	}
	return nil
}

// getBevTreeByID 获取行为树
func getBevTreeByID(id string) *b3core.BehaviorTree {
	t, ok := mapTreesByID.Load(id)
	if ok {
		return t.(*b3core.BehaviorTree)
	}
	return nil
}

//自定义的节点
func createExtStructMaps() *b3.RegisterStructMaps {
	st := b3.NewRegisterStructMaps()
	//actions

	st.Register("HPLowerThan", &HPLowerThan{})
	st.Register("targetDistanceLowestThan", &targetDistanceLowestThan{})
	st.Register("curSkillIsNormalAttack", &curSkillIsNormalAttack{})
	st.Register("EnableUltimateSkillIsNil", &EnableUltimateSkillIsNil{})
	st.Register("EnableSuperSkillIsNil", &EnableSuperSkillIsNil{})
	st.Register("SelfRangeHasEnemy", &SelfRangeHasEnemy{})
	st.Register("SkillAttackRangeHasEnemy", &SkillAttackRangeHasEnemy{})
	st.Register("IsRage", &IsRage{})
	st.Register("RandomSelectEnemy", &RandomSelectEnemy{})
	st.Register("selectHpLowestEnemyFromBoard", &selectHpLowestEnemyFromBoard{})
	st.Register("SelectEnemyWithMinHP", &SelectEnemyWithMinHP{})
	st.Register("DetourMove", &DetourMove{})
	st.Register("GetAttackPos", &GetAttackPos{})
	st.Register("GetCurSkill", &GetCurSkill{})
	st.Register("MoveTo", &MoveTo{})
	st.Register("CastBloadSkill", &CastBloadSkill{})
	st.Register("GetTargetBackPos", &GetTargetBackPos{})
	st.Register("FlashMoveTo", &FlashMoveTo{})
	st.Register("GetSkillByIndex", &GetSkillByIndex{})
	st.Register("ChangeForm", &ChangeForm{})
	st.Register("GetRandNormalAttack", &GetRandNormalAttack{})
	st.Register("RandSuccess", &RandSuccess{})
	st.Register("SelectSelf", &SelectSelf{})
	st.Register("AttrLowerThan", &AttrLowerThan{})
	st.Register("SelectNearestEnemy", &SelectNearestEnemy{})
	st.Register("IsOverDrive", &IsOverDrive{})
	st.Register("ResetAI", &ResetAI{})
	st.Register("ChangeMass", &ChangeMass{})

	st.Register("IsPartnerLoseHP", &IsPartnerLoseHP{})
	st.Register("IsUseableLebelsSkill", &IsUseableLebelsSkill{})
	st.Register("SelectLowestHPPartner", &SelectLowestHPPartner{})
	st.Register("IsPositiveEnemy", &IsPositiveEnemy{})
	st.Register("GetOutOfPositiveBossRangePos", &GetOutOfPositiveBossRangePos{})
	st.Register("MoveToOutOfPositivePos", &MoveToOutOfPositivePos{})

	st.Register("RandSelectEnemyReast", &RandSelectEnemyReast{})
	st.Register("SelectEnemyByMinAttr", &SelectEnemyByMinAttr{})
	st.Register("SelectEnemyByMaxAttr", &SelectEnemyByMaxAttr{})
	st.Register("SelectFurthestEnemy", &SelectFurthestEnemy{})
	st.Register("SteeringSmoothing", &SteeringSmoothing{})
	st.Register("WaitAction", &WaitAction{})
	st.Register("RandWaitAction", &RandWaitAction{})
	st.Register("GetUsableSkill", &GetUsableSkill{})
	st.Register("Retreat", &Retreat{})
	st.Register("GetSkillTargetByBlackboard", &GetSkillTargetByBlackboard{})
	st.Register("InAOERange", &InAOERange{})

	return st
}
