package main

import (
	"Cinder/Base/MQNet"
	"Cinder/Base/Message"
	"Cinder/Base/Util"
	"fmt"
	log "github.com/cihub/seelog"
)

func main() {
	srv := MQNet.NewService()

	err := srv.Init("agent_1", "agent")
	if err != nil {
		log.Debug(err)
		return
	}

	srv.AddProc(&_TProc{})

	go func() {

		index := 0

		for {
			msg := &Message.UserLoginReq{
				UserID:   Util.GetGUID(),
				UserData: []byte("hi , this is go language "),
			}

			if err := srv.Post("agent_1", msg); err != nil {
				fmt.Println(err)
			} else {
				index++

				if index%100 == 0 {
					fmt.Println("index", index)
				}
			}

		}

	}()

	select {}

}

type _TProc struct {
}

func (p *_TProc) MessageProc(srcAddr string, message Message.IMessage) {
	//fmt.Println("MessageProc")
	//time.Sleep(5 * time.Millisecond)
}
