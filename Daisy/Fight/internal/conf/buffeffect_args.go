package conf

import "reflect"

var buffEffectArgsType = map[int]reflect.Type{
	2000: reflect.TypeOf(ChangeAttrArgs{}),
	2001: reflect.TypeOf(ChangeStateArgs{}),
	2002: reflect.TypeOf(HPShieldArgs{}),
	2003: reflect.TypeOf(DelayEffectArgs{}),
	2004: reflect.TypeOf(SummonArgs{}),
}

type SummonArgs struct {
	Delay      uint32
	SummonType SummonType
	Npcs       []SummonNpc
}

type SummonNpc struct {
	NpcID    uint32
	PosIndex uint32
	LifeTime uint32
}

type AttrArgs struct {
	Rate float32
	Fix  float64
}

type ChangeAttrArgs struct {
	ChangeAttr map[uint32]*AttrArgs
}

type ChangeStateArgs struct {
	ChangeState []uint32
}

type HPShieldArgs struct {
	HPShieldSource  ShieldValueSrc // 护盾数值源
	HPShieldValue   int64          // 护盾数值
	HPShieldAddRate float64        // 护盾增加比例
	HPShieldAddFix  int64          // 护盾增加固定值
}

type DelayEffectArgs struct {
	Time   uint32   // 生效时间
	BuffId []uint32 // 添加buff列表
}
