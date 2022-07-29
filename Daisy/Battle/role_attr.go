package main

import (
	"Daisy/Const"
	"Daisy/Data"
	"Daisy/Fight"
	"Daisy/Fight/attraffix"
	"Daisy/Proto"
)

// refreshBuildFightAttr 刷新build战斗属性
func (r *_Role) refreshBuildFightAttr(buildID string) {
	build, ok := r.prop.Data.BuildMap[buildID]
	if !ok {
		return
	}

	specialAgent, ok := r.prop.Data.SpecialAgentList[build.SpecialAgentID]
	if !ok {
		return
	}

	// 战斗属性计算器
	calcAttr := &Fight.CalcAttr{}
	if err := calcAttr.Init(Proto.PawnType_Role, specialAgent.Base.ConfigID, int32(specialAgent.Base.Level), nil); err != nil {
		return
	}

	var AttackScore, DefenceScore, PsychokinesisScore, extraScore float32

	// 装备词条
	for _, itemID := range build.EquipmentMap {
		// 装备词条
		equipAffix := r.GetEquipAttrAffix(itemID)

		// 生效属性词条
		calcAttr.AddAttrAffix(equipAffix.AttrList)

		// buff评分
		for _, buffAfx := range equipAffix.BuffList {
			switch buffAfx.ScoreType {
			case Const.AffixScoreType_Extra:
				extraScore += float32(buffAfx.Score)
			case Const.AffixScoreType_Attack:
				AttackScore += float32(buffAfx.Score)
			case Const.AffixScoreType_Defence:
				DefenceScore += float32(buffAfx.Score)
			case Const.AffixScoreType_Psychokinesis:
				PsychokinesisScore += float32(buffAfx.Score)
			}
		}
	}

	// 统计属性评分
	statsFun := func(field attraffix.Field) {
		attrScore, ok := Data.GetEquipConfig().AttEnumEration_ConfigItems[uint32(field)]
		if !ok || attrScore.ScoreParam <= 0 {
			return
		}

		score := float32(calcAttr.GetAttr(field) / float64(attrScore.ScoreParam))

		switch attrScore.ScoreType {
		case Const.AffixScoreType_Extra:
			extraScore += score
		case Const.AffixScoreType_Attack:
			AttackScore += score
		case Const.AffixScoreType_Defence:
			DefenceScore += score
		case Const.AffixScoreType_Psychokinesis:
			PsychokinesisScore += score
		}
	}

	for i := attraffix.Field_Prop_Begin; i < attraffix.Field_Prop_End; i++ {
		statsFun(i)
	}

	for i := attraffix.Field_Equip_Begin; i < attraffix.Field_Equip_End; i++ {
		statsFun(i)
	}

	for i := attraffix.Field_Logic_Begin; i < attraffix.Field_Logic_End; i++ {
		statsFun(i)
	}

	// 填充当前属性
	fightAttr := &Proto.FightAttr{
		MaxHP:                         calcAttr.MaxHP,
		RecoverHP:                     calcAttr.RecoverHP,
		Attack:                        int64(calcAttr.Attack),
		HitRate:                       calcAttr.HitRate,
		AttackLucky:                   calcAttr.AttackLucky,
		NormalAttackSpeed:             calcAttr.NormalAttackSpeed,
		Defence:                       int64(calcAttr.Armor),
		DodgeRate:                     calcAttr.DodgeRate,
		Strength:                      calcAttr.Strength,
		Agility:                       calcAttr.Agility,
		Intelligence:                  calcAttr.Intelligence,
		Vitality:                      calcAttr.Vitality,
		SkillPowerLimit:               calcAttr.SkillPowerLimit,
		RecoverUltimateSkillPowerRate: calcAttr.RecoverUltimateSkillPowerRate,
		Lucky:                         int32(calcAttr.Lucky),
		CritRate:                      calcAttr.CritRate,
		CritDamageRate:                calcAttr.CritDamageRate,
		BlockRate:                     calcAttr.BlockRate,
		BlockValue:                    int64(calcAttr.BlockValue),
		BeDamageNormalDeduct:          int64(calcAttr.BeDamageNormalDeduct),
		ResistancePoison:              calcAttr.ResistancePoison,
		ResistanceFire:                calcAttr.ResistanceFire,
		ResistanceCold:                calcAttr.ResistanceCold,
		ResistanceLightning:           calcAttr.ResistanceLightning,
		AttackScore:                   int32(AttackScore),
		DefenceScore:                  int32(DefenceScore),
		PsychokinesisScore:            int32(PsychokinesisScore),
		TotalScore:                    int32(extraScore + AttackScore + DefenceScore + PsychokinesisScore),
	}

	r.prop.SyncSetBuildFightAttr(buildID, fightAttr)
	r.FlushToCache()
}
