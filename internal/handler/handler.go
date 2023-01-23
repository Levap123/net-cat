package handler

import (
	"fmt"
)

type BroadPayload struct {
	Name string
	Msg  string
}

type JoinLeave struct {
	Name   string
	IsJoin bool
}

func Broadcast(msgChan chan BroadPayload, joinLeaveChan chan JoinLeave) {
	for {
		select {
		case jl := <-joinLeaveChan:
			if jl.IsJoin {
				msg := fmt.Sprintf("%s has joined our chat...", jl.Name)
				send(msg, jl.Name)
			} else {
				msg := fmt.Sprintf("%s has left our chat...", jl.Name)
				send(msg, jl.Name)
			}

		case val := <-msgChan:
			send(val.Msg, val.Name)
		}
	}
}

func send(msg, username string) {
	for name, conn := range userQuantity {
		if name != username {
			fmt.Fprint(conn, "\n"+msg)
			fmt.Fprint(conn, "\n"+message(name, ""))
		}
	}
}
