package internal

import (
	"Daisy/Fight/internal/conf"
	"errors"
	"fmt"
)

// BuffKey Buff key
type BuffKey struct {
	MainID uint32 // buff配置ID
	UID    uint32 // buff uid
}

// ToUint64 转换为uint64类型
func (buffKey *BuffKey) ToUint64() uint64 {
	return uint64(buffKey.MainID)<<32 + uint64(buffKey.UID)
}

// Buff Buff
type Buff struct {
	Config       *conf.BuffConfig  // buff配置
	BuffKey      BuffKey           // Buff key
	Caster, Pawn *Pawn             // 释放buff者与被施加buff者
	ExtValue     int64             // 外部数值
	effectTab    []IEffectCallback // buff效果表
	IsDestroy    bool              // 已删除
	destroyTime  uint32            // 销毁时间
}

// init 初始化
func (buff *Buff) init(buffKey BuffKey, pawn, caster *Pawn, extValue int64) error {
	if pawn == nil || caster == nil {
		return errors.New("args invalid")
	}

	buff.BuffKey = buffKey
	buff.Pawn = pawn
	buff.Caster = caster
	buff.ExtValue = extValue

	var ok bool
	if buff.Config, ok = pawn.Scene.GetBuffConfig(buffKey.MainID); !ok {
		return fmt.Errorf("buff %d not found", buffKey.MainID)
	}

	for _, v := range buff.Config.EffectConfs {
		effect, err := createEffect(v.Type, v.Args)
		if err != nil {
			return fmt.Errorf("buff %d create effect failed, %s", buffKey.MainID, err.Error())
		}

		if buffEff, ok := effect.(IBuffEffect); ok {
			if err := buffEff.Init(buff); err != nil {
				return fmt.Errorf("buff %d effect %d init error, %s", buffKey.MainID, v.Type, err.Error())
			}
		}

		buff.effectTab = append(buff.effectTab, effect)
	}

	return nil
}
