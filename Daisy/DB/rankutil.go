package DB

import (
	"fmt"
	"github.com/go-redis/redis/v7"
)

var DefaultIRankUtil *rankUtil

type IRankUtil interface {
	GetList(typ string, rankID uint32, num int64) []string
	GetPlace(typ string, rankID uint32, member string) uint32
	Remove(typ string, rankID uint32, member string) (int64,error)
	UpdateScore(typ string, rankID uint32, score float64, member string) (int64,error)
}

type IRankItem interface {
	GetRankData() interface{}
	GetRankScore() float64
	GetRankType() string
	GetRankID() uint32
	GetID() string
}

const (
	SeasonRank = "seasonrank"
	SeasonLegendRank = "seasonlegendrank"
)

func init() {
	DefaultIRankUtil = CreateIRankUtil()
}

func CreateIRankUtil() *rankUtil {
	return  &rankUtil{}
}
type rankUtil struct {

}

func (r *rankUtil) GetList(typ string, rankID uint32, num int64) []string{
	list := make([]string, 0)
	ret := RedisDB.ZRevRange(r.getRankKey(typ, rankID), 0, num)
	//ret := RedisDB.ZRange(r.getRankKey(typ, rankID), 0, num)
	members,err := ret.Result()
	if err != nil {
		return nil
	}
	for k,v := range members{
		fmt.Println(k,v)
		_data,_err := r.GetMemberData(typ, rankID, v)
		if _err != nil {
			fmt.Println(_err)
			continue
		}
		list = append(list, _data)
	}
	return list
}

func (r *rankUtil) GetPlace(typ string, rankID uint32, member string) uint32{
	ret := RedisDB.ZRevRank(r.getRankKey(typ, rankID), member)
	if ret == nil || ret.Err() != nil{
		return 0
	}
	return uint32(ret.Val()+1)
}

func (r *rankUtil) GetPlacePercentage(typ string, rankID uint32, member string) uint32{
	ret := RedisDB.ZRevRank(r.getRankKey(typ, rankID), member)
	if ret == nil || ret.Err() != nil{
		return 0
	}
	all := RedisDB.ZCard(r.getRankKey(typ, rankID))
	if all == nil || all.Err()!=nil || all.Val() == 0{
		return 0
	}
	return uint32(all.Val() - ret.Val())*100/uint32(all.Val())
}

func (r *rankUtil) UpdateScore(typ string, rankID uint32, score float64, member string) (int64,error) {
	ret := RedisDB.ZAdd(r.getRankKey(typ, rankID), &redis.Z{
		Score:score,Member:member,
	})
	return  ret.Result()
}

func (r *rankUtil) Remove(typ string, rankID uint32, member string) (int64,error){
	ret := RedisDB.ZRem(r.getRankKey(typ, rankID), member)
	return ret.Result()
}

func (r *rankUtil) getRankKey(typ string, rankID uint32) string {
	return fmt.Sprintf("rank:%s:%d",  typ, rankID)
}

func (r *rankUtil) getRankDataKey(typ string, rankID uint32) string{
	return fmt.Sprintf("rankdata:%s:%d", typ, rankID)
}

func (r *rankUtil) UpdateMemberData(typ string, rankID uint32, member string, data[]byte) (int64,error){
	ret := RedisDB.HSet(r.getRankDataKey(typ, rankID), member, data)
	return ret.Result()
}

func (r *rankUtil) RemoveMemberData(typ string, rankID uint32, member string) (int64,error){
	ret := RedisDB.HDel(r.getRankDataKey(typ, rankID), member)
	return ret.Result()
}

func (r *rankUtil) GetMemberData(typ string, rankID uint32, member string) (string, error){
	ret := RedisDB.HGet(r.getRankDataKey(typ, rankID), member)
	return ret.Result()
}