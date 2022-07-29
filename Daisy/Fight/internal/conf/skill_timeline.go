package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
)

// skillTimelinePath 技能时间轴配置路径
var skillTimelinePath = "../res/Timeline/Skill"

// SetSkillTimeLinePath 技能时间轴配置路径
func SetSkillTimeLinePath(path string) {
	skillTimelinePath = path
}

// Turn 转向 采用最小角度
type Turn struct {
	Begin    uint32  //转向修正开始时间 毫秒
	Duration uint32  //转向修正持续时间
	Speed    float32 // 转向速度（单位: 度/秒）
}

// SkillTimeLine 技能流程配置
type SkillTimeLine struct {
	ShowTime   uint32            // 特写时间（单位: ms）
	BeforeTime uint32            // 前摇时间（单位: ms）
	AttackTime uint32            // 伤害时间（单位: ms）
	Attacks    []*AttackTimeLine // 伤害时间轴
	LaterTime  uint32            // 后摇时间（单位: ms）
	Turn       *Turn
}

// loadSkillTimelineConfig 加载技能时间轴配置
func loadSkillTimelineConfig(fileName string) (*SkillTimeLine, error) {
	var timeLine SkillTimeLine
	filePath := path.Join(skillTimelinePath, fileName+".json")
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("TimeLine文件【%s】无法打开，%s", filePath, err.Error())
	}

	err = json.Unmarshal(file, &timeLine)
	if err != nil {
		return nil, fmt.Errorf("TimeLine文件【%s】配置错误，%s", filePath, err.Error())
	}

	for _, attack := range timeLine.Attacks {
		for _, hit := range attack.Hits {
			if !floatEqual(float64(hit.HitBackDis), 0) && hit.HitBackDuring == 0 {
				return nil, fmt.Errorf("TimeLine文件【%s】击退配置错误，击退距离:%v, 击退时间%v", filePath, hit.HitBackDis, hit.HitBackDuring)
			}
		}
	}

	//沟通策划既定规则 turn开始时间不能在前摇之前
	if timeLine.Turn != nil && timeLine.Turn.Begin < timeLine.BeforeTime {
		return nil, fmt.Errorf("TimeLine文件: %v, Turn开始时间不能比前摇早. turn.begin:%v, BeforeTime:%v", fileName, timeLine.Turn.Begin, timeLine.BeforeTime)
	}

	return &timeLine, nil
}
