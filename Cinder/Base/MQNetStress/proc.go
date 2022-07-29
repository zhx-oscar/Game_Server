package main

import (
	"Cinder/Base/Message"
	"bytes"
	"encoding/gob"
	"time"

	"go.uber.org/atomic"
)

type NotifyCh chan struct{}
type Proc struct {
	// 每个协程1个 chan, 通知RPC可以开始
	RPCStartChs []NotifyCh

	// RPC 计数
	Count atomic.Uint32

	// 最大延时 ms
	MaxDelayMs atomic.Uint32

	// max delay 开始统计时间，用于忽略最初1s内的延时(可能是初始化原因，延时较大)
	startTime time.Time
}

func newProc(goroutines int) *Proc {
	chs := make([]NotifyCh, goroutines)
	for i := 0; i < goroutines; i++ {
		chs[i] = make(NotifyCh, 1)
		chs[i] <- struct{}{} // 初始化为可以开始了
	}

	return &Proc{
		RPCStartChs: chs,
		startTime:   time.Now().Add(time.Second),
	}
}

func GobEncodePingMsg(iGoroutine int) []byte {
	ping := PingMsg{
		GoroutineIndex: iGoroutine,
		Timestamp:      time.Now(),
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(ping)
	panicIfError(err)
	return buf.Bytes()
}

func GobDecodePingMsg(data []byte) PingMsg {
	ping := PingMsg{}
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	err := dec.Decode(&ping)
	panicIfError(err)
	return ping
}

func (p *Proc) MessageProc(srcAddr string, message Message.IMessage) {
	if message == nil {
		return
	}

	req := message.(*Message.RpcReq)
	ping := GobDecodePingMsg(req.Args)
	p.Count.Inc()                                    // 计数
	p.UpdateMaxDelay(ping.Timestamp)                 // 更新 max delay
	p.RPCStartChs[ping.GoroutineIndex] <- struct{}{} // 开始下一个RPC
}

func (p *Proc) UpdateMaxDelay(pingTime time.Time) {
	delayMs := uint32(time.Now().Sub(pingTime).Milliseconds())
	// fmt.Printf("delay: %dms\n", delayMs)

	if pingTime.Before(p.startTime) {
		return // 初始化时，延时不正常, 忽略
	}

	max := p.MaxDelayMs.Load()
	if delayMs > max {
		p.MaxDelayMs.Store(delayMs)
	}
}
