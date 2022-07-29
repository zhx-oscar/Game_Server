package stats

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// Action 单个统计项的统计结果
type Action struct {
	Name string // 动作名字

	// 统计项
	AvgCost  time.Duration // 平均每次消耗
	Min      time.Duration // 最低消耗
	Max      time.Duration // 最高消耗
	Cost50   time.Duration // 排序后50%处耗时
	Cost70   time.Duration // 排序后70%处耗时
	Cost90   time.Duration // 排序后90%处耗时
	QPS      float32       // 每秒处理数
	TotalQPS float32       // 总QPS数

	results        resultList
	startTime      time.Time
	lastStatsTime  time.Time
	lastStatsCount int

	mux sync.Mutex
}

func NewAction(name string) *Action {
	now := time.Now()
	act := &Action{
		Name:          name,
		results:       make([]time.Duration, 0, 100),
		startTime:     now,
		lastStatsTime: now,
	}
	return act
}

func (act *Action) String() string {
	return fmt.Sprintf("Action:%s QPS:%.2f TotalQPS:%.2f Avg:%s Min:%s Max:%s 50:%s 70:%s 90:%s",
		act.Name, act.QPS, act.TotalQPS, act.AvgCost, act.Min, act.Max, act.Cost50, act.Cost70, act.Cost90)
}

func (act *Action) Add(cost time.Duration) {
	act.mux.Lock()
	defer act.mux.Unlock()

	act.results = append(act.results, cost)
}

func (act *Action) Calc() {
	act.mux.Lock()
	defer act.mux.Unlock()

	total := len(act.results)
	if total == 0 {
		return
	}

	sort.Sort(act.results)

	act.Min = act.results[0]
	act.Max = act.results[total-1]
	act.Cost50 = act.results[total/2]
	act.Cost70 = act.results[total*7/10]
	act.Cost90 = act.results[total*9/10]
	act.AvgCost = act.results.Sum() / time.Duration(act.results.Len())
	act.QPS = float32(total-act.lastStatsCount) / float32(time.Now().Sub(act.lastStatsTime).Seconds())
	act.TotalQPS = float32(total) / float32(time.Now().Sub(act.startTime).Seconds())

	act.lastStatsTime = time.Now()
	act.lastStatsCount = total
}

func (act *Action) Reset() {
	act.mux.Lock()
	defer act.mux.Unlock()

	act.results = make([]time.Duration, 0, 10)
	act.AvgCost = 0
	act.Cost50 = 0
	act.Cost70 = 0
	act.Cost90 = 0
	act.QPS = 0
}

type resultList []time.Duration

func (l resultList) Len() int           { return len(l) }
func (l resultList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l resultList) Less(i, j int) bool { return l[i] < l[j] }

func (l resultList) Sum() time.Duration {
	var sum time.Duration
	for i := range l {
		sum += l[i]
	}
	return sum
}
