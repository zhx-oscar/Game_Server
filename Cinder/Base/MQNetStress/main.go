package main

import (
	"Cinder/Base/MQNet"
	"Cinder/Base/MQNet/mqnats"
	"Cinder/Base/Message"
	"Cinder/Base/Util"
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type PingMsg struct {
	GoroutineIndex int
	Timestamp      time.Time
}

func init() {
	pflag.IntP("goroutines", "g", 1, "goroutines count")
	pflag.DurationP("time", "m", 15*time.Second, "test time")
	pflag.StringP("nats.addr", "a", "127.0.0.1:4222", "NATS address")
}

func main() {
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	goroutines := viper.GetInt("goroutines")
	tm := viper.GetDuration("time")

	srv := mqnats.New()
	srvID := Util.GetGUID()
	initMQNet(srv, srvID)
	proc := newProc(goroutines)
	srv.AddProc(proc)

	fmt.Printf("MQNet stress test is running (goroutines=%d time=%s)...\n", goroutines, tm)
	ctx, cancel := context.WithTimeout(context.Background(), tm)
	defer cancel()

	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go runStressTest(ctx, &wg, srv, srvID, i, proc.RPCStartChs[i])
	}
	wg.Wait()

	count := proc.Count.Load()
	fmt.Printf("RPS: %d / %f = %f (r/s), max delay: %dms\n",
		count, tm.Seconds(), float64(count)/tm.Seconds(), proc.MaxDelayMs.Load())
}

func runStressTest(ctx context.Context, wg *sync.WaitGroup, srv MQNet.IService, srvID string, iGoroutine int, rpcStartCh <-chan struct{}) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-rpcStartCh:
			request(srv, srvID, iGoroutine)
		}
	}
	return
}

func initMQNet(srv MQNet.IService, srvID string) {
	serviceAddr := srvID
	boardcastAddr := "mqnetstress"
	err := srv.Init(MQNet.InitOptions(viper.GetString("NATS.Addr"), serviceAddr, boardcastAddr))
	panicIfError(err)
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func request(srv MQNet.IService, srvID string, iGoroutine int) {
	err := srv.Post(srvID, &Message.RpcReq{
		Args: GobEncodePingMsg(iGoroutine),
	})
	panicIfError(err)
}
