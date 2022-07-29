package internal

import (
	"fmt"
	"reflect"
)

// Effects 效果注册器
var Effects = &EffectRegister{
	effectTypeMap: map[uint32]reflect.Type{},
}

// EffectRegister 效果注册器
type EffectRegister struct {
	effectTypeMap map[uint32]reflect.Type
}

// Register 注册效果
func (effects *EffectRegister) Register(effectType uint32, effect IEffectCallback) {
	if effect == nil {
		panic("nil effect")
	}

	effects.effectTypeMap[effectType] = reflect.Indirect(reflect.ValueOf(effect).Elem()).Type()
}

// createEffect 创建效果
func createEffect(effectType uint32, args reflect.Value) (IEffectCallback, error) {
	effectClass, ok := Effects.effectTypeMap[effectType]
	if !ok {
		return nil, fmt.Errorf("effect type %d not found", effectType)
	}

	effect := reflect.New(effectClass)
	if !effect.IsValid() {
		return nil, fmt.Errorf("new effect type %d failed", effectType)
	}

	if field := effect.Elem().FieldByName("Args"); field.IsValid() {
		field.Set(args)
	}

	return effect.Interface().(IEffectCallback), nil
}
