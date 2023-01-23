package handler

import (
	"fmt"
	"sync"
)

type BroadPayload struct {
	Name string
	Msg  string
}

type JoinLeave struct {
	Name   string
	IsJoin bool
}

var msgBuffer = make([]string, 0)

func Broadcast(msgChan chan BroadPayload, joinLeaveChan chan JoinLeave, mu *sync.Mutex) {
	for {
		select {
		case jl := <-joinLeaveChan:
			if jl.IsJoin {
				msg := fmt.Sprintf("%s has joined our chat...", jl.Name)
				send(msg, jl.Name, mu)
			} else {
				msg := fmt.Sprintf("%s has left our chat...", jl.Name)
				send(msg, jl.Name, mu)
			}

		case val := <-msgChan:
			send(val.Msg, val.Name, mu)
		}
	}
}

func send(msg, username string, mu *sync.Mutex) {
	mu.Lock()
	msgBuffer = append(msgBuffer, "\n"+msg)
	mu.Unlock()
	for name, conn := range userQuantity {
		if name != username {
			fmt.Fprint(conn, "\n"+msg)
			fmt.Fprint(conn, "\n"+message(name, ""))
		}
	}
}
