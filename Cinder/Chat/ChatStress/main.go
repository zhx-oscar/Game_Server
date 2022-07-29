package main

import (
	"Cinder/Base/Core"
	"Cinder/Base/Util"
	"Cinder/Chat/ChatStress/rpcproc"
	"Cinder/Chat/ChatStress/stress"
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	assert "github.com/arl/assertgo"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/atomic"
)

func init() {
	pflag.IntP("goroutines", "g", 1, "goroutines count")
	pflag.StringP("test", "t", "loginLogout", "test case name")
	pflag.BoolP("list", "l", false, "list all test names")
	pflag.DurationP("time", "m", 15*time.Second, "test time")
	pflag.String("area", "1", "大区ID")
}

func main() {
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	if viper.GetBool("list") {
		listAllTests()
		return
	}

	initServer()
	goroutines := viper.GetInt("goroutines")
	test := viper.GetString("test")
	tm := viper.GetDuration("time")

	fmt.Printf("chat stress test is setting up (goroutines=%d test=%s time=%s)...\n", goroutines, test, tm)
	setup()
	fmt.Printf("chat stress test is running (goroutines=%d test=%s time=%s)...\n", goroutines, test, tm)

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
	names := stress.GetStressTestNames()
	sort.Strings(names)
	for _, n := range names {
		fmt.Printf("\t%s\n", n)
	}
}

func runStressTest(ctx context.Context, wg *sync.WaitGroup, test string, iGoroutine stress.GoroutineIndex, transactions *atomic.Uint32) {
	defer wg.Done()
	runs := stress.RunStressTest(ctx, test, iGoroutine)
	transactions.Add(runs)
	return
}

func initServer() {
	_ = Core.New()
	info := Core.NewDefaultInfo()
	info.ServiceType = "chatstress"
	svcID := fmt.Sprintf("%s_%s", info.ServiceType, Util.GetGUID())
	assert.True(len(svcID) < 64) // nsq topic requires
	info.AreaID = viper.GetString("area")
	info.ServiceID = svcID
	info.RpcProc = rpcproc.NewRpcProc()
	if errInit := Core.Inst.Init(info); errInit != nil {
		panic(errInit)
	}
}

func setup() {
	var wg sync.WaitGroup
	goroutines := viper.GetInt("goroutines")
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go runStressSetup(&wg, stress.GoroutineIndex(i))
	}
	wg.Wait()
}

func runStressSetup(wg *sync.WaitGroup, iGo stress.GoroutineIndex) {
	defer wg.Done()
	stress.Setup(iGo)
}
