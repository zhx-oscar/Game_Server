package internal

import "fmt"

// Halo
type Halo struct {
	UID          uint32           // 光环uid
	caster       *Pawn            // 光环释放者
	casterEffect IEffectCallback  // 创建光环的来源效果类型
	MemberList   map[uint32]*Pawn // 成员列表
}

// _HaloFlow 光环buff流程
type _HaloFlow struct {
	scene    *Scene
	haloList []*Halo // 光环列表
}

// init 初始化
func (flow *_HaloFlow) init(scene *Scene) {
	flow.scene = scene
	flow.haloList = nil
}

// AddHalo 添加光环
func (flow *_HaloFlow) AddHalo(caster *Pawn, effect IEffectCallback) *Halo {
	scene := flow.scene

	newHalo := &Halo{
		UID:          scene.generateUID(),
		caster:       caster,
		casterEffect: effect,
		MemberList:   map[uint32]*Pawn{},
	}

	scene.PushDebugInfo(func() string {
		return fmt.Sprintf("${PawnID:%d}创建光环%d", caster.UID, newHalo.UID)
	})

	flow.haloList = append(flow.haloList, newHalo)

	pawnList := scene.GetPawnList()
	for _, pawn := range pawnList {
		if pawn.IsAlive() {
			scene.PushDebugInfo(func() string {
				return fmt.Sprintf("${PawnID:%d}进入光环%d", pawn.UID, newHalo.UID)
			})

			newHalo.MemberList[pawn.UID] = pawn
			newHalo.caster.Events.EmitHaloAddMember(newHalo, pawn)
		}
	}

	return newHalo
}

// RemoveHalo 移除光环
func (flow *_HaloFlow) RemoveHalo(haloUid uint32) {
	var halo *Halo
	for i := 0; i < len(flow.haloList); i++ {
		if flow.haloList[i].UID == haloUid {
			halo = flow.haloList[i]
			flow.haloList = append(flow.haloList[0:i], flow.haloList[i+1:]...)
			break
		}
	}

	if halo == nil {
		return
	}

	scene := flow.scene

	scene.PushDebugInfo(func() string {
		return fmt.Sprintf("${PawnID:%d}销毁光环%d", halo.caster.UID, halo.UID)
	})

	for pawnUID, member := range halo.MemberList {
		scene.PushDebugInfo(func() string {
			return fmt.Sprintf("${PawnID:%d}离开光环%d", pawnUID, halo.UID)
		})

		delete(halo.MemberList, pawnUID)
		halo.caster.Events.EmitHaloRemoveMember(halo, member)
	}
}

// AddHaloMember 增加光环成员
func (flow *_HaloFlow) AddHaloMember(pawn *Pawn) {
	if !pawn.IsAlive() {
		return
	}

	for _, halo := range flow.haloList {
		if _, ok := halo.MemberList[pawn.UID]; ok {
			continue
		}

		flow.scene.PushDebugInfo(func() string {
			return fmt.Sprintf("${PawnID:%d}进入光环%d", pawn.UID, halo.UID)
		})

		halo.MemberList[pawn.UID] = pawn
		halo.caster.Events.EmitHaloAddMember(halo, pawn)
	}
}

// RemoveHaloMember 移除光环成员
func (flow *_HaloFlow) RemoveHaloMember(pawn *Pawn) {
	for _, halo := range flow.haloList {
		if _, ok := halo.MemberList[pawn.UID]; !ok {
			continue
		}

		flow.scene.PushDebugInfo(func() string {
			return fmt.Sprintf("${PawnID:%d}离开光环%d", pawn.UID, halo.UID)
		})

		delete(halo.MemberList, pawn.UID)
		halo.caster.Events.EmitHaloRemoveMember(halo, pawn)
	}
}
