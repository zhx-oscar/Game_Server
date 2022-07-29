package main

import (
	"Daisy/Prop"
	"Daisy/Proto"
	"time"
)

func (team *_Team) AddRobot(name string, specialAgentID uint32) (string, error) {
	specialAgent := &Proto.SpecialAgent{
		Base: &Proto.SpecialAgentBase{
			ConfigID: specialAgentID,
			Level:    1,
			Exp:      0,
			GainTime: time.Now().Unix(),
		},
	}

	build := &Proto.BuildData{
		BuildID:        name + "1",
		Name:           name + "预设 1",
		SpecialAgentID: specialAgentID,
		Skill: &Proto.BuildSkillData{
			UltimateSkillID: 0,
			SuperSkill:      map[uint32]uint32{},
		},
		EquipmentMap: map[int32]string{},
		CreateTime:   time.Now().Unix(),
	}

	roleProp := &Prop.RoleProp{
		Data: &Proto.Role{
			Base: &Proto.RoleBase{
				Name:       name,
				CreateTime: time.Now().Unix(),
				TeamID:     team.GetID(),
			},
			SpecialAgentList: map[uint32]*Proto.SpecialAgent{specialAgentID: specialAgent},
			BuildMap:         map[string]*Proto.BuildData{build.BuildID: build},
			FightingBuildID:  build.BuildID,
		},
	}

	propData, err := roleProp.Marshal()
	if err != nil {
		return "", err
	}

	return team.AddActor(RoleActorType, "", "", propData, nil)
}

func (team *_Team) RemoveRobot(id string) error {
	return team.RemoveActor(id)
}
