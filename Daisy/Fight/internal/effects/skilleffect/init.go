package skilleffect

import (
	. "Daisy/Fight/internal"
)

func init() {
	Effects.Register(1, &_1_SuperSkill{})
	Effects.Register(2, &_2_NormalSkill{})
	Effects.Register(3, &_3_UltimateSkill{})
	Effects.Register(4, &_4_CombineSkill{})
}
