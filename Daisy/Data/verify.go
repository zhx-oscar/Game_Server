package Data

import (
	"Daisy/Const"
	"encoding/json"
	"fmt"
)

var RootPath = "../"

// VerifyExcel 检查excel数据有效性
func VerifyExcel() bool {
	var resultCode = true
	// todo 填充各个excel表格数据有效性检查

	//AI表检测
	if !ai_VerifyExcel() {
		resultCode = false
	}

	//特工表检测
	if !specialAgent_VerifyExcel() {
		resultCode = false
	}

	//技能道具表检测
	if !skillItem_VerifyExcel() {
		resultCode = false
	}

	//技能实体表检测
	if !skillEntity_VerifyExcel() {
		resultCode = false
	}

	//掉落表检测
	if !drop_VerifyExcel() {
		resultCode = false
	}

	//快速战斗表检测
	if !fastBattle_VerifyExcel() {
		resultCode = false
	}

	return resultCode
}

// ai_VerifyExcel AI配置表验证
func ai_VerifyExcel() bool {
	resultCode := true
	var data map[string]interface{}
	for id, val := range GetAIData_Config().AIData_ConfigItems {
		if len(val.BlackBoardKeys) > 0 {
			err := json.Unmarshal([]byte(val.BlackBoardKeys), &data)
			if err != nil {
				fmt.Printf("AI表中 id:%v 行为树黑板json格式错误:%v\n", id, err.Error())
				resultCode = false
			}
		}
	}

	return resultCode
}

// specialAgent_VerifyExcel 特工配置表验证
func specialAgent_VerifyExcel() bool {
	resultCode := true
	var data map[string]interface{}
	for id, val := range GetSpecialAgentConfig().SpecialAgent_ConfigItems {
		//特工必须配置AI
		if val.AIID == 0 {
			fmt.Printf("特工表中 特工id:%v AIID 为 0\n", id)
			resultCode = false
		}

		//AI信息是否存在
		_, ok := GetAIData_Config().AIData_ConfigItems[val.AIID]
		if !ok {
			fmt.Printf("特工表中 特工id:%v 配置的AIID:%v 在AI配置表中不存在\n", id, val.AIID)
			resultCode = false
		}

		//黑板json格式是否正确
		if len(val.BlackBoardKeys) > 0 {
			err := json.Unmarshal([]byte(val.BlackBoardKeys), &data)
			if err != nil {
				fmt.Printf("特工表中 特工id:%v 行为树黑板json格式错误:%v\n", id, err.Error())
				resultCode = false
			}
		}

		//普工 至少有一个技能
		if len(val.NormalAttack) == 0 {
			fmt.Printf("特工表中 特工id:%v 普攻+超能技+必杀技 至少配置一个\n", id)
			resultCode = false
		}

		//普攻技能验证
		for _, skillID := range val.NormalAttack {
			_, ok = GetSkillConfig().SkillValue_ConfigItems[skillID]
			if !ok {
				fmt.Printf("特工表中 特工id:%v 普攻技能:%v 在SkillValue表中不存在\n", val.ID, skillID)
				resultCode = false
			}
		}
	}

	//特工常量配置 需要存在
	_, ok := GetSpecialAgentConfig().SpecialAgentConst_ConfigItems[Const.SpecialAgent_buildCountMax]
	if !ok {
		fmt.Printf("特工表常量页签配置中 build列表上限 常量缺失\n")
		resultCode = false
	}

	_, ok = GetSpecialAgentConfig().SpecialAgentConst_ConfigItems[Const.SpecialAgent_buildSuperSlillCount]
	if !ok {
		fmt.Printf("特工表常量页签配置中 build超能技槽位数量 常量缺失\n")
		resultCode = false
	}

	return resultCode
}

// skillItem_VerifyExcel 技能道具配置表验证
func skillItem_VerifyExcel() bool {
	resultCode := true
	for _, val := range GetSkillConfig().SkillItem_ConfigItems {
		_, ok := GetSkillConfig().SkillMain_ConfigItems[val.SkillID]
		if !ok {
			fmt.Printf("技能道具表中 id:%v 关联的技能ID:%v 在SkillMain表中不存在\n", val.ID, val.SkillID)
			resultCode = false
		}
	}

	return resultCode
}

// skillEntity_VerifyExcel 技能实体配置表验证
func skillEntity_VerifyExcel() bool {
	resultCode := true
	for _, val := range GetSkillConfig().SkillEntity_ConfigItems {
		_, ok := GetSkillConfig().SkillMain_ConfigItems[val.ID]
		if !ok {
			fmt.Printf("技能实体表中 id:%v 关联的技能ID:%v 在SkillMain表中不存在\n", val.ID, val.ID)
			resultCode = false
		}

		//验证 等级数量 和 不同等级材料消耗 配置是否对等
		if (val.UpCoinType > 0 && val.TopSkillLevel != uint32(len(val.UpCoinNum))) || (val.UpItemID > 0 && val.TopSkillLevel != uint32(len(val.UpItemNum))) {
			fmt.Printf("技能道具表中 id:%v 技能最高等级:%v 技能升级消耗配置 UpgradeSkillCoinNum 长度:%v, UpgradeSkillmaterialsNum 长度:%v\n", val.ID, val.TopSkillLevel, len(val.UpCoinNum), len(val.UpItemNum))
			resultCode = false
		}
	}

	return resultCode
}

// drop_VerifyExcel 掉落配置表验证
func drop_VerifyExcel() bool {
	return HandleDropConfig()
}

func fastBattle_VerifyExcel() bool {
	for i := 1; i <= len(GetFastBattleConfig().FastBattle_ConfigItems); i++ {
		if _, ok := GetFastBattleConfig().FastBattle_ConfigItems[uint32(i)]; !ok {
			fmt.Print("FastBattle_Config ID必须有序递增\n")
			return false
		}
	}

	for i := 1; i <= len(GetFastBattleConfig().Energize_ConfigItems); i++ {
		if _, ok := GetFastBattleConfig().Energize_ConfigItems[uint32(i)]; !ok {
			fmt.Print("Energize_Config ID必须有序递增\n")
			return false
		}
	}

	for i := 1; i <= len(GetFastBattleConfig().SpeedUp_ConfigItems); i++ {
		if _, ok := GetFastBattleConfig().SpeedUp_ConfigItems[uint32(i)]; !ok {
			fmt.Print("SpeedUp_Config ID必须有序递增\n")
			return false
		}
	}

	totalProbability := uint32(0)
	for _, value := range GetFastBattleConfig().SpeedUp_ConfigItems {
		totalProbability += value.Probability
	}
	if totalProbability != 10000 {
		fmt.Println("SpeedUp_Config 概率之和必须为10000")
		return false
	}

	if _, ok := GetFastBattleConfig().Variable_ConfigItems[1]; !ok {
		fmt.Println("Variable_Config 必须有ID=1这一行")
		return false
	}

	return true
}
