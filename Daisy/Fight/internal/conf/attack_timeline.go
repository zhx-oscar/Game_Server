package conf

import (
	"Daisy/Proto"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
)

// attackTimelinePath 伤害体时间轴配置路径
var attackTimelinePath = "../res/Timeline/Attack"

// SetAttackTimelinePath 伤害体时间轴配置路径
func SetAttackTimelinePath(path string) {
	attackTimelinePath = path
}

// HitTimeLine 命中时间轴
type HitTimeLine struct {
	Begin         uint32             // 偏移时间（单位: ms）
	Pause         uint32             // 施法者停顿时间（单位: ms）
	HitType       Proto.HitType_Enum // 命中类型
	HitBackDis    float32            // 击退距离
	HitBackDuring uint32             // 击退时间
}

// AttackTimeLine 伤害体时间轴
type AttackTimeLine struct {
	ConfigID    uint32         // 伤害体配置ID（用于客户端查询美术特效配置）
	Begin       uint32         // 偏移时间（单位: ms）
	Delay       uint32         // 被触发后延迟时间（单位：ms）
	Speed       float32        // 移动速度（单位: 米/秒）
	LifeTime    uint32         // 生存时间（单位: ms）
	Hits        []*HitTimeLine // 命中时间轴
	FixInterval uint32         // 固定命中间隔（单位: ms）
	FixLimit    uint32         // 固定命中次数限制
	FixHit      *HitTimeLine   // 固定命中间隔hit配置
}

// loadAttackTimelineConfig 加载伤害体时间轴配置
func loadAttackTimelineConfig(fileName string) (*AttackTimeLine, error) {
	var timeLine AttackTimeLine
	filePath := path.Join(attackTimelinePath, fileName+".json")
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("TimeLine文件【%s】无法打开，%s", filePath, err.Error())
	}

	err = json.Unmarshal(file, &timeLine)
	if err != nil {
		return nil, fmt.Errorf("TimeLine文件【%s】配置错误，%s", filePath, err.Error())
	}

	return &timeLine, nil
}
