package conf

import (
	"Daisy/Data"
	"Daisy/Proto"
	"os"
)

// HotUpdateConfs 热更新配置
func HotUpdateConfs(simulatorMode, errPanic bool) {
	if !simulatorMode {
		// 非模拟器需要检测 ../res目录是否存在
		if _, err := os.Stat("../res"); err != nil {
			if os.IsNotExist(err) {
				return
			} else {
				panic(err)
			}
		}
	}

	tConfigs := &Configs{}
	var err error

	// 加载伤害体配置表
	tConfigs.attackExcelConf = Data.GetAttackConfig()

	// 加载技能配置表
	tConfigs.skillExcelConf = Data.GetSkillConfig()

	// 加载技能目标策略配置表
	tConfigs.targetStrategy = Data.GetTargetStrategyConfig()

	// 加载技能配置
	tConfigs.skillConfs, err = loadSkillConfigs(tConfigs.skillExcelConf, tConfigs.attackExcelConf)
	if err != nil {
		if errPanic {
			panic(err)
		} else {
			return
		}
	}

	// 加载buff配置表
	tConfigs.buffExcelConf = Data.GetBuffConfig()

	// 加载buff配置
	tConfigs.buffConfs, err = loadBuffConfigs(tConfigs.buffExcelConf)
	if err != nil {
		if errPanic {
			panic(err)
		} else {
			return
		}
	}

	// 加载数值配置表
	tConfigs.propValueExcelConf = Data.GetPropConfig()

	// 加载role配置表
	tConfigs.specialAgentExcelConf = Data.GetSpecialAgentConfig()

	// 加载role受击配置
	tConfigs.playerBeHitConfs, err = loadRoleBeHitConfigs(tConfigs.specialAgentExcelConf)
	if err != nil {
		if errPanic {
			panic(err)
		} else {
			return
		}
	}

	tConfigs.pawnConfs = make(map[Proto.PawnType_Enum]map[uint32]*PawnConfig)

	// 加载role配置
	rolePawnConfig, err := loadRolePawnConfig(tConfigs.specialAgentExcelConf, tConfigs.propValueExcelConf, tConfigs.playerBeHitConfs)
	if err != nil {
		if errPanic {
			panic(err)
		} else {
			return
		}
	}
	tConfigs.pawnConfs[Proto.PawnType_Role] = rolePawnConfig

	// 加载怪物配置表
	tConfigs.monsterExcelConf = Data.GetMonsterConfig()

	// 加载怪物受击配置
	tConfigs.monsterBeHitConfs, err = loadMonsterBeHitConfigs(tConfigs.monsterExcelConf)
	if err != nil {
		if errPanic {
			panic(err)
		} else {
			return
		}
	}

	// 加载怪物表演配置
	tConfigs.npcActConfs, err = loadNpcActConfigs(tConfigs.monsterExcelConf)
	if err != nil {
		if errPanic {
			panic(err)
		} else {
			return
		}
	}

	// 加载pawnconfig配置
	monsterPawnConfig, err := loadMonsterPawnConfig(tConfigs.monsterExcelConf, tConfigs.propValueExcelConf, tConfigs.monsterBeHitConfs, tConfigs.npcActConfs)
	if err != nil {
		if errPanic {
			panic(err)
		} else {
			return
		}
	}
	tConfigs.pawnConfs[Proto.PawnType_Npc] = monsterPawnConfig

	// 加载background配置
	tConfigs.pawnConfs[Proto.PawnType_BG] = loadBGPawnConfig()

	// 加载ai配置表
	tConfigs.aiDataExcelConf = Data.GetAIData_Config()

	// 加载ai配置
	tConfigs.aiDataConfs = loadAIInfo(tConfigs.aiDataExcelConf)

	// 加载质量表
	tConfigs.massExcelConf = Data.GetMassConfig()

	// 加载常量表
	tConfigs.constExcelConf = Data.GetFightConstConfig()

	// 加载内置伤害体模板
	tConfigs.innerAttackTemplate, err = loadInnerAttackTemplate(tConfigs.constExcelConf, tConfigs.attackExcelConf)
	if err != nil {
		if errPanic {
			panic(err)
		} else {
			return
		}
	}

	// 切换全局配置
	configs = tConfigs
}
