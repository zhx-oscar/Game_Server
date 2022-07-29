package buffeffect

import (
	. "Daisy/Fight/internal"
)

func init() {
	// 程序内置buff效果
	Effects.Register(1000, &_1000_OverDrive{})
	Effects.Register(1001, &_1001_EnergyShield{})
	Effects.Register(1002, &_1002_Unbalance{})
	Effects.Register(1003, &_1003_BlockBreak{})
	Effects.Register(1004, &_1004_BornAct{})
	Effects.Register(1005, &_1005_RecoverHP{})
	Effects.Register(1006, &_1006_RecoverUltimateSkillPower{})
	Effects.Register(1007, &_1007_Thorns{})
	Effects.Register(1008, &_1008_StealUltimateSkillPower{})

	// 策划可配置效果
	Effects.Register(2000, &_2000_ChangeAttr{})
	Effects.Register(2001, &_2001_ChangeState{})
	Effects.Register(2002, &_2002_HPShield{})
	Effects.Register(2003, &_2003_DelayEffect{})
	Effects.Register(2004, &_2004_Summon{})
}
