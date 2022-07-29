package Data

import (
	"Cinder/LDB"
	"Daisy/DataTables"
	"encoding/json"
	log "github.com/cihub/seelog"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

// 配置文件初始化，增加一张表时，需要手动添加注册表结构的代码

var inst LDB.IDataTables
var doWatch *sync.Once
var resPath = "../res/DataTables"

func init() {
	inst = LDB.NewDataTables()
	doWatch = &sync.Once{}

	// 注册表结构
	inst.Register("Skill_Config", &DataTables.Skill_Config_Data{})
	inst.Register("SpecialAgent_Config", &DataTables.SpecialAgent_Config_Data{})
	inst.Register("Mass_Config", &DataTables.Mass_Config_Data{})
	inst.Register("AIData_Config", &DataTables.AIData_Config_Data{})
	inst.Register("Monster_Config", &DataTables.Monster_Config_Data{})
	inst.Register("Prop_Config", &DataTables.Prop_Config_Data{})
	inst.Register("Scene_Config", &DataTables.Scene_Config_Data{})
	inst.Register("FightConst_Config", &DataTables.FightConst_Config_Data{})
	inst.Register("Guide_Config", &DataTables.Guide_Config_Data{})
	inst.Register("Drop_Config", &DataTables.Drop_Config_Data{})
	inst.Register("Equip_Config", &DataTables.Equip_Config_Data{})
	inst.Register("Buff_Config", &DataTables.Buff_Config_Data{})
	inst.Register("Attack_Config", &DataTables.Attack_Config_Data{})
	inst.Register("ItemType_Config", &DataTables.ItemType_Config_Data{})
	inst.Register("Talent_Config", &DataTables.Talent_Config_Data{})
	inst.Register("FastBattle_Config", &DataTables.FastBattle_Config_Data{})
	inst.Register("Mail_Config", &DataTables.Mail_Config_Data{})
	inst.Register("Begging_Config", &DataTables.Begging_Config_Data{})
	inst.Register("Supply_Config", &DataTables.Supply_Config_Data{})
	inst.Register("TargetStrategy_Config", &DataTables.TargetStrategy_Config_Data{})
	inst.Register("Title_Config", &DataTables.Title_Config_Data{})
	inst.Register("Season_Config", &DataTables.Season_Config_Data{})
	inst.Register("PlayerUpgrade_Config", &DataTables.PlayerUpgrade_Config_Data{})

	if _, err := os.Stat(resPath); err != nil {
		if os.IsNotExist(err) {
			return
		} else {
			panic(err)
		}
	}

	if err := LoadDataTables(resPath); err != nil {
		panic(err)
	}

	//服务器启动 excel检测
	if !VerifyExcel() {
		panic("excel verify fail")
	}

	//二次处理数据 初始化
	initHandleData()

	// 注册HandleData热更
	AttachHotUpdate(hotUpdateHandleData)
}

// LoadDataTables 加载数据表
func LoadDataTables(resPath string) error {
	if err := inst.Load(resPath); err != nil {
		return err
	}

	execHotUpdate(false)

	go doWatch.Do(func() {
		hotUpdateWatch()
	})

	return nil
}

//hotUpdateWatch 热更配置监听
func hotUpdateWatch() {
	//fmt.Println("hotUpdateWatch ===============   ", resPath)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error(err)
		return
	}
	defer watcher.Close()

	reloadPath := "../res/ExcelReload/"
	done := make(chan bool)
	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				log.Debugf("Watch %s Op %s", ev.Name, ev.Op)

				if strings.Contains(ev.Name, "reload.json") {
					if ev.Op&fsnotify.Write == fsnotify.Write || ev.Op&fsnotify.Create == fsnotify.Create {
						log.Info("Start Reload!")

						if err = Reload(reloadPath + "reload.json"); err != nil {
							log.Error("Reload failed ", err)
						} else {
							execHotUpdate(true)
							log.Info("Reload excel success")
						}
					}
				}
			case err = <-watcher.Errors:
				log.Error(err)
			}
		}
	}()

	if err = watcher.Add(reloadPath); err != nil {
		log.Error(err)
	} else {
		log.Debug("Watch res path:", reloadPath)
	}
	<-done
}

func Reload(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	config := make(map[string]bool)
	err = json.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	for k, v := range config {
		if v {
			err = inst.HotUpdateLoadFile(resPath + "/" + k + ".bytes")
			if err != nil {
				log.Error("Reload failed ", k, err)
			}
		}
	}

	return nil
}
