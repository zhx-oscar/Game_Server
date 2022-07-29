package Space

/*

type _AgentSessMgr struct {
	sessMap sync.Map
}

func newAgentSessMgr() *_AgentSessMgr {
	return &_AgentSessMgr{}
}

func (mgr *_AgentSessMgr) Add(agentID string, sess Net.ISess) error {

	_, ok := mgr.sessMap.Load(agentID)
	if ok {
		return errors.New("is existed ")
	}

	mgr.sessMap.Store(agentID, sess)
	return nil
}

func (mgr *_AgentSessMgr) Get(agentID string) (Net.ISess, error) {

	i, ok := mgr.sessMap.Load(agentID)
	if !ok {
		return nil, errors.New("not existed ")
	}

	return i.(Net.ISess), nil
}

func (mgr *_AgentSessMgr) Remove(agentID string) {
	mgr.sessMap.Delete(agentID)
}

*/
