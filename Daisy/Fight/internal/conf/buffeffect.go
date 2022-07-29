package conf

import (
	"encoding/json"
	"reflect"
)

// BuffEffect buff效果
type BuffEffect struct {
	Type uint32
	Args reflect.Value
}

type _BuffEffectRaw struct {
	Type uint32
	Args json.RawMessage
}

// loadBuffEffectConfig 加载buff效果配置
func loadBuffEffectConfig(effect string) ([]*BuffEffect, error) {
	if effect == "" {
		return nil, nil
	}

	var buffEffRaws []*_BuffEffectRaw

	if err := json.Unmarshal([]byte(effect), &buffEffRaws); err != nil {
		return nil, err
	}

	var buffEffs []*BuffEffect

	for _, v := range buffEffRaws {
		var args reflect.Value

		if argsType, ok := buffEffectArgsType[int(v.Type)]; ok {
			args = reflect.New(argsType)
			if err := json.Unmarshal(v.Args, args.Interface()); err != nil {
				return nil, err
			}
		}

		buffEffs = append(buffEffs, &BuffEffect{
			Type: v.Type,
			Args: args,
		})
	}

	return buffEffs, nil
}
