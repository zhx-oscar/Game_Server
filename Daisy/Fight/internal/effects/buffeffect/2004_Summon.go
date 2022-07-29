package buffeffect

import (
	. "Daisy/Fight/internal"
	"Daisy/Fight/internal/conf"
	"Daisy/Fight/internal/effects"
	"Daisy/Fight/internal/log"
)

type _2004_Summon struct {
	effects.Blank
	buff    *Buff
	Args    *conf.SummonArgs
	AddTime uint32
}

func (effect *_2004_Summon) Init(buff *Buff) error {
	effect.buff = buff
	effect.AddTime = buff.Pawn.Scene.NowTime

	return nil
}

// OnBuffUpdate buff状态更新
func (effect *_2004_Summon) OnBuffUpdate(buff *Buff) {
	pawn := buff.Caster

	//非延迟召唤判断
	if effect.Args.SummonType == conf.SummonTypeDead {
		return
	}

	if pawn.Scene.NowTime < effect.AddTime+effect.Args.Delay {
		return
	}

	for _, npc := range effect.Args.Npcs {
		info, err := pawn.Scene.BuildNPCInfo(npc.NpcID)
		if err != nil {
			log.Errorf("_2004_Summon buffID: %v, NPCID: %v", buff.Config.ID, npc.NpcID)
			continue
		}

		bornPoints := pawn.Scene.Info.Formation[pawn.GetCamp()].BornPoints
		value := bornPoints[int(npc.PosIndex)%len(bornPoints)]

		info.BornPos = &value.Point
		info.BornAngle = value.Angle
		info.LifeTime = npc.LifeTime

		summonPawn, err := pawn.Scene.Summon(info)
		if err != nil {
			log.Errorf("_2004_Summon buffID: %v, NPCID: %v Summon fail err: %v", effect.buff.Config.ID, npc.NpcID, err)
			continue
		}

		summonPawn.BatchAddBuffs(summonPawn, summonPawn.Info.InnerBornBuffs, 0)
		summonPawn.BatchAddBuffs(summonPawn, summonPawn.Info.BornBuffs, 0)
	}

	buff.Destroy()
}

// OnDead 死亡后（所有buff能收到，attack为nil表示非伤害造成的死亡）
func (effect *_2004_Summon) OnDead(attack *Attack) {
	pawn := effect.buff.Pawn

	//非死亡召唤判断
	if effect.Args.SummonType == conf.SummonTypeDelay {
		return
	}

	for _, npc := range effect.Args.Npcs {
		info, err := pawn.Scene.BuildNPCInfo(npc.NpcID)
		if err != nil {
			log.Errorf("_2004_Summon buffID: %v, NPCID: %v", effect.buff.Config.ID, npc.NpcID)
			continue
		}

		bornPoints := pawn.Scene.Info.Formation[pawn.GetCamp()].BornPoints
		value := bornPoints[int(npc.PosIndex)%len(bornPoints)]

		info.BornPos = &value.Point
		info.BornAngle = value.Angle
		info.LifeTime = npc.LifeTime

		summonPawn, err := pawn.Scene.Summon(info)
		if err != nil {
			log.Errorf("_2004_Summon buffID: %v, NPCID: %v Summon fail err: %v", effect.buff.Config.ID, npc.NpcID, err)
			continue
		}

		summonPawn.BatchAddBuffs(summonPawn, summonPawn.Info.InnerBornBuffs, 0)
		summonPawn.BatchAddBuffs(summonPawn, summonPawn.Info.BornBuffs, 0)
	}

	effect.buff.Destroy()
}
