package mtypes

import (
	"github.com/mohae/deepcopy"
)

/* M 即 map[string]interface{}, 用于自定义数据，如
mtypes.M{
	"level": 10,
	"isNewUser": true,
	"other": mtypes.M{"f1": 123, "f2": "abc"},
}
*/
type M map[string]interface{}

func (m *M) Copy() M {
	return deepcopy.Copy(*m).(M)
}
