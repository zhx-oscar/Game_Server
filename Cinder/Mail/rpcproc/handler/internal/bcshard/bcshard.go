// bcshard 定义广播邮件的分片
package bcshard

import (
	"math/rand"
)

type ShardID uint16

const (
	// 全服邮件分片复制个数，shard = 0..31
	// 可以改小；但如果改大，需要先清空DB中的全服邮件，或手工补全新增副本。
	BroadcastShardCount = 32
)

func GetRandBcShardID() ShardID {
	return ShardID(rand.Intn(BroadcastShardCount)) // [0,n)
}
