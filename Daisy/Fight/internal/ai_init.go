package internal

import (
	b3core "github.com/magicsea/behavior3go/core"
	"sync"
)

var aiWatchPath = "../res/ai/"
var behaviorTreePath = "../res/ai/chaos.b3"

func init() {
	b3core.SetSubTreeLoadFunc(func(id string) *b3core.BehaviorTree {
		t, ok := mapTreesByID.Load(id)
		if ok {
			return t.(*b3core.BehaviorTree)
		}
		return nil
	})
	mapTrees = sync.Map{}
	mapTreesByID = sync.Map{}
	extMaps = createExtStructMaps()
	LoadBehaviorTreeFile(behaviorTreePath)
	//reload()
}
