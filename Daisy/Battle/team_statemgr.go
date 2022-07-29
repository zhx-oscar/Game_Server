package main

import (
	"errors"
	"fmt"
	"time"
)

const (
	TeamState_Init uint8 = iota
	TeamState_Running
	TeamState_Raidbattling
	TeamState_FastBattleing

	TeamState_Max
)

type ITeamState interface {
	OnLeave()
	OnEnter(args ...interface{})
	OnLoop(delta time.Duration)
	CanChangeTo(state uint8) bool
}

type _TeamStateMgr struct {
	state    uint8
	stateObj ITeamState
	stateMap map[uint8]ITeamState
}

func (mgr *_TeamStateMgr) RegState(states map[uint8]ITeamState) {
	mgr.stateMap = states
}

func (mgr *_TeamStateMgr) GetState() uint8 {
	return mgr.state
}

func (mgr *_TeamStateMgr) SetState(state uint8, args ...interface{}) error {
	if !mgr.CanSetState(state) {
		return errors.New(fmt.Sprintf("can't change state %d -> %d", mgr.state, state))
	}

	if _, ok := mgr.stateMap[state]; !ok {
		return errors.New("dest state not register")
	}

	if mgr.stateObj != nil {
		mgr.stateObj.OnLeave()
	}
	mgr.state = state
	mgr.stateObj = mgr.stateMap[state]
	mgr.stateObj.OnEnter(args...)
	return nil
}

func (mgr *_TeamStateMgr) CanSetState(state uint8) bool {
	if state >= TeamState_Max || state == mgr.state {
		return false
	}

	if mgr.stateObj != nil && !mgr.stateObj.CanChangeTo(state) {
		return false
	}

	if _, ok := mgr.stateMap[state]; !ok {
		return false
	}

	return true
}

func (mgr *_TeamStateMgr) LoopState(delta time.Duration) {
	if state, ok := mgr.stateMap[mgr.state]; ok {
		state.OnLoop(delta)
	}
}
