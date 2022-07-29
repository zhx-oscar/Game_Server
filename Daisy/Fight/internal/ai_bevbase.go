package internal

import (
	"Cinder/Base/linemath"
	b3core "github.com/magicsea/behavior3go/core"
)

// getVector2Board 从黑板中读取vector2
func getVector2Board(tick *b3core.Tick, name string) (*linemath.Vector2, bool) {
	obj := tick.Blackboard.Get(name, "", "")
	if obj == nil {
		return nil, false
	}

	value := obj.(*linemath.Vector2)
	return value, true
}

// getPawnBoard 从黑板中读取pawn
func getPawnBoard(tick *b3core.Tick, name string) (*Pawn, bool) {
	obj := tick.Blackboard.Get(name, "", "")

	if obj == nil {
		return nil, false
	}
	value := obj.(*Pawn)
	return value, true
}

// getPawnBoardByBlackboard 从黑板中读取pawn
func getPawnBoardByBlackboard(blackboard *b3core.Blackboard, name string) (*Pawn, bool) {
	if blackboard == nil {
		return nil, false
	}

	obj := blackboard.Get(name, "", "")

	if obj == nil {
		return nil, false
	}
	value := obj.(*Pawn)
	return value, true
}

// getSkillBoard 从黑板中读取Skill
func getSkillBoard(tick *b3core.Tick, name string) (*_SkillItem, bool) {
	obj := tick.Blackboard.Get(name, "", "")

	if obj == nil {
		return nil, false
	}
	value := obj.(*_SkillItem)
	return value, true
}

// getSkillBoardByBlackboard 从黑板中读取Skill
func getSkillBoardByBlackboard(blackboard *b3core.Blackboard, name string) (*_SkillItem, bool) {
	if blackboard == nil {
		return nil, false
	}

	obj := blackboard.Get(name, "", "")

	if obj == nil {
		return nil, false
	}
	value := obj.(*_SkillItem)
	return value, true
}

//getEnemyListBoard 从黑板中获取敌人列表
func getEnemyListBoard(tick *b3core.Tick, name string) ([]*Pawn, bool) {
	obj := tick.Blackboard.Get(name, "", "")

	if obj == nil {
		return nil, false
	}
	value := obj.([]*Pawn)
	return value, true
}
