package main

import (
	"Daisy/Const"
	"Daisy/DB"
	"Daisy/Proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const recruitmentsLifeTime = 10 * time.Second

type _TeamHallModel struct {
	recruitments    *Proto.Recruitments
	lastRefreshTime time.Time
	isLoading       bool
	roleIDs         map[string]bool

	RecruitmentsChan chan *Proto.Recruitments
}

func NewTeamHallModel() *_TeamHallModel {
	return &_TeamHallModel{
		RecruitmentsChan: make(chan *Proto.Recruitments, 1),
		roleIDs:          make(map[string]bool),
	}
}

func cloneRecruitments(res *Proto.Recruitments) *Proto.Recruitments {
	cloneRes := &Proto.Recruitments{
		Data: make([]*Proto.Recruitment, len(res.Data), len(res.Data)),
	}
	for i := 0; i < len(res.Data); i++ {
		cloneRes.Data[i] = res.Data[i]
	}
	return cloneRes
}

func (team *_Team) OnGetRecruitments(res *Proto.Recruitments) {
	if res == nil {
		team.hall.isLoading = false
		return
	}

	team.hall.recruitments = res
	team.hall.lastRefreshTime = time.Now()
	team.hall.isLoading = false

	for key := range team.hall.roleIDs {
		if ia, err := team.GetActor(key); err == nil {
			ia.(*_Role).OnGetRecruitments(cloneRecruitments(team.hall.recruitments))
		}
	}
	team.hall.roleIDs = make(map[string]bool)
}

func (team *_Team) GetRecruitments(r *_Role) error {
	if time.Now().Sub(team.hall.lastRefreshTime) < recruitmentsLifeTime {
		r.OnGetRecruitments(cloneRecruitments(team.hall.recruitments))
		return nil
	}

	team.hall.roleIDs[r.GetID()] = true

	if !team.hall.isLoading {
		team.hall.isLoading = true

		go func() {
			objID, err := primitive.ObjectIDFromHex(team.GetID())
			if err != nil {
				team.Errorf("ObjectIDFromHex err:%s", err)
				team.hall.RecruitmentsChan <- nil
				return
			}

			param := bson.M{"json_data.base.num": bson.M{"$gt": 0, "$lt": Const.TeamMaxMemberNum}, "_id": bson.M{"$ne": objID}, "json_data.base.published": true}

			pipeline := make([]bson.M, 0)
			pipeline = append(pipeline, bson.M{"$match": param}, bson.M{"$sample": bson.M{"size": 20}})

			teams, err := DB.GetTeamHallUitl().Aggregate(pipeline)
			if err != nil {
				team.Errorf("Aggregate err:%s", err)
				team.hall.RecruitmentsChan <- nil
				return
			}

			roleIDs := make([]string, 0)
			for _, teamProto := range teams {
				roleIDs = append(roleIDs, getTeamProtoMemberIDs(teamProto.Team)...)
			}

			roles, err := DB.GetRoleHallUitl().FindIDs(roleIDs)
			if err != nil {
				team.Errorf("RoleHallUtil FindIDs err:%s", err)
				team.hall.RecruitmentsChan <- nil
				return
			}

			res := &Proto.Recruitments{
				Data: team.assemblyRecruitments(teams, roles),
			}
			team.hall.RecruitmentsChan <- res
		}()
	}

	return nil
}

func (team *_Team) GetTeamByUID(uid uint64) (*Proto.Recruitments, error) {
	wrap, err := DB.GetTeamHallUitl().FindUID(uid)
	if err != nil {
		return nil, err
	}

	roleIDs := getTeamProtoMemberIDs(wrap.Team)
	roles, err := DB.GetRoleHallUitl().FindIDs(roleIDs)
	if err != nil {
		return nil, err
	}

	res := &Proto.Recruitments{
		Data: team.assemblyRecruitments([]*DB.TeamProtoWrap{wrap}, roles),
	}
	return res, nil
}

func (team *_Team) GetInvites(infos map[string]*Proto.RoleInviteInfo) (*Proto.Inviters, error) {
	ids := make([]string, len(infos), len(infos))
	i := 0
	for key := range infos {
		ids[i] = key
		i++
	}

	teams, err := DB.GetTeamHallUitl().FindIDs(ids)
	if err != nil {
		return nil, err
	}

	roleIDs := make([]string, 0)
	for _, teamProto := range teams {
		roleIDs = append(roleIDs, getTeamProtoMemberIDs(teamProto.Team)...)
	}

	roles, err := DB.GetRoleHallUitl().FindIDs(roleIDs)
	if err != nil {
		return nil, err
	}

	res := team.assemblyRecruitments(teams, roles)
	inviters := make([]*Proto.Inviter, len(res), len(res))
	for i := 0; i < len(res); i++ {
		inviters[i] = &Proto.Inviter{
			Re:         res[i],
			InviteTime: infos[res[i].TeamID].Time,
		}
	}

	return &Proto.Inviters{
		Data: inviters,
	}, nil
}

func (team *_Team) GetApplys(infos map[string]*Proto.TeamApplyInfo) (*Proto.Applyers, error) {
	ids := make([]string, len(infos), len(infos))
	i := 0
	for key := range infos {
		ids[i] = key
		i++
	}

	roles, err := DB.GetRoleHallUitl().FindIDs(ids)
	if err != nil {
		return nil, err
	}

	loners := team.assemblyLoners(roles)
	applyers := make([]*Proto.Applyer, len(loners), len(loners))
	for i := 0; i < len(loners); i++ {
		applyers[i] = &Proto.Applyer{
			Loner:     loners[i],
			ApplyTime: infos[loners[i].RoleID].Time,
		}
	}

	return &Proto.Applyers{
		Data: applyers,
	}, nil
}

func (team *_Team) assemblyRecruitments(teams []*DB.TeamProtoWrap, roles []*DB.RoleProtoWrap) []*Proto.Recruitment {
	roleMap := make(map[string]*Proto.Role)
	for i := 0; i < len(roles); i++ {
		roleMap[roles[i].ID.Hex()] = roles[i].Role
	}

	res := make([]*Proto.Recruitment, 0, len(teams))
	for _, value := range teams {
		if value.ID.Hex() == team.GetID() {
			continue
		}

		memberIDs := getTeamProtoMemberIDs(value.Team)
		reMembers := make([]*Proto.RecruitmentMember, 0, len(memberIDs))

		for i := 0; i < len(memberIDs); i++ {
			role, ok := roleMap[memberIDs[i]]
			if !ok {
				break
			}

			fightBuild, ok := role.BuildMap[role.FightingBuildID]
			if !ok {
				break
			}

			fightAgent, ok := role.SpecialAgentList[fightBuild.SpecialAgentID]
			if !ok {
				break
			}

			member := &Proto.RecruitmentMember{
				RoleID:            memberIDs[i],
				SpecialAgentID:    fightBuild.SpecialAgentID,
				Name:              role.Base.Name,
				SpecialAgentLevel: fightAgent.Base.Level,
				Online:            role.Base.Online,
				Status:            value.Team.Base.Members[memberIDs[i]].Status,
				JoinTime:          value.Team.Base.Members[memberIDs[i]].JoinTime,
				TotalScore:        fightBuild.FightAttr.TotalScore,
			}
			reMembers = append(reMembers, member)
		}

		if len(reMembers) < len(memberIDs) {
			continue
		}

		re := &Proto.Recruitment{
			TeamID:       value.ID.Hex(),
			TeamUID:      value.Team.Base.UID,
			Members:      reMembers,
			Topic:        value.Team.Base.Board,
			RaidProgress: value.Team.Raid.Progress,
			TeamName:     value.Team.Base.Name,
		}

		res = append(res, re)
	}

	return res
}

func (team *_Team) assemblyLoners(roles []*DB.RoleProtoWrap) []*Proto.Loner {
	loners := make([]*Proto.Loner, 0, len(roles))
	for _, roleWrap := range roles {
		role := roleWrap.Role
		fightBuild, ok := role.BuildMap[role.FightingBuildID]
		if !ok {
			continue
		}

		fightAgent, ok := role.SpecialAgentList[fightBuild.SpecialAgentID]
		if !ok {
			continue
		}

		loner := &Proto.Loner{
			RoleID:            roleWrap.ID.Hex(),
			SpecialAgentID:    fightBuild.SpecialAgentID,
			Name:              role.Base.Name,
			SpecialAgentLevel: fightAgent.Base.Level,
			Online:            role.Base.Online,
			RaidProgress:      role.Base.RaidProgress,
			TotalScore:        fightBuild.FightAttr.TotalScore,
		}

		loners = append(loners, loner)
	}

	return loners
}

func getTeamProtoMemberIDs(team *Proto.Team) []string {
	roleIDs := make([]string, len(team.Base.Members))
	i := 0
	for key := range team.Base.Members {
		roleIDs[i] = key
		i++
	}

	return roleIDs
}
