package main

import (
	"Cinder/Base/Core"
	"Cinder/Base/Util"
	"Cinder/Matcher/MatcherStress/rpcproc"
	"Cinder/Matcher/MatcherStress/stress"
	"Cinder/Matcher/matchapi"
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"sort"
	"sync"
	"time"

	assert "github.com/arl/assertgo"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/atomic"
)

var (
	svc        = matchapi.GetRoomService()
	rolePrefix = Util.GetGUID() + "_"
	tests      = map[string]func(stress.GoroutineIndex, stress.RunIndex){
		// 房间
		"RoleJoinRandomRoom":   stress.RoleJoinRandomRoom,
		"RoleCreateDeleteRoom": stress.RoleCreateDeleteRoom,
		"RoleJoinRoom":         stress.RoleJoinRoom,
		"RoleJoinLeaveRoom":    stress.RoleJoinLeaveRoom,
		"BroadcastRoom":        stress.BroadcastRoom,
		"GetRoomList":          stress.GetRoomList,
		"GetRoomInfo":          stress.GetRoomInfo,
		"SetRoomData":          stress.SetRoomData,
		"SetRoomRoleData":      stress.SetRoomRoleData,
		// 组队
		"CreateLeaveTeam":  stress.CreateLeaveTeam,
		"JoinLeaveTeam":    stress.JoinLeaveTeam,
		"ChangeTeamLeader": stress.ChangeTeamLeader,
		"SetTeamData":      stress.SetTeamData,
		"BroadcastTeam":    stress.BroadcastTeam,
		// 组队加入房间
		"TeamCreateDeleteRoom": stress.TeamCreateDeleteRoom,
		"TeamJoinRoom":         stress.TeamJoinRoom,
		"TeamJoinRandomRoom":   stress.TeamJoinRandomRoom,
		// 其他
		"MyRPCTest": stress.MyRPCTest,
		"Ping":      stress.Ping,
	}
)

func init() {
	pflag.IntP("goroutines", "g", 1, "goroutines count")
	pflag.StringP("test", "t", "RoleJoinRandomRoom", "test case name")
	pflag.BoolP("list", "l", false, "list all test names")
	pflag.DurationP("time", "m", 15*time.Second, "test time")
	pflag.StringP("area", "a", "1", "area ID")
}

func main() {
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	if viper.GetBool("list") {
		listAllTests()
		return
	}

	initServer(viper.GetString("area"))
	goroutines := viper.GetInt("goroutines")
	test := viper.GetString("test")
	tm := viper.GetDuration("time")
	fmt.Printf("matcher stress test is running (goroutines=%d test=%s time=%s)...\n", goroutines, test, tm)

	go func() {
		fmt.Println(http.ListenAndServe("localhost:26060", nil))
	}()

	ctx, cancel := context.WithTimeout(context.Background(), tm)
	defer cancel()

	var wg sync.WaitGroup
	var transactions atomic.Uint32
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go runStressTest(ctx, &wg, test, stress.GoroutineIndex(i), &transactions)
	}
	wg.Wait()

	fmt.Printf("%s TPS: %d / %f = %f (t/s)\n", test, transactions.Load(), tm.Seconds(),
		float64(transactions.Load())/tm.Seconds())
	time.Sleep(2 * time.Second)
}

func listAllTests() {
	fmt.Printf("Test names:\n")
	names := make([]string, 0, len(tests))
	for k, _ := range tests {
		names = append(names, k)
	}
	sort.Strings(names)
	for _, n := range names {
		fmt.Printf("\t%s\n", n)
	}
}

func runStressTest(ctx context.Context, wg *sync.WaitGroup, test string, iGoroutine stress.GoroutineIndex, transactions *atomic.Uint32) {
	defer wg.Done()

	stressFun, ok := tests[test]
	if !ok {
		panic("illegal test name")
	}

	var i stress.RunIndex = 0
	for {
		select {
		case <-ctx.Done():
			transactions.Add(uint32(i))
			return
		default:
			stressFun(iGoroutine, i)
			i++
		}
	}
}

func initServer(areaID string) {
	_ = Core.New()
	info := Core.NewDefaultInfo()
	info.ServiceType = "stress"
	svcID := fmt.Sprintf("%s_%s", info.ServiceType, Util.GetGUID())
	assert.True(len(svcID) < 64) // nsq topic requires
	info.ServiceID = svcID
	info.AreaID = areaID
	info.RpcProc = rpcproc.NewRpcProc()
	if errInit := Core.Inst.Init(info); errInit != nil {
		panic(errInit)
	}
}
