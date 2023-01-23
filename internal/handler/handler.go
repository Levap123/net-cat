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

func Broadcast(msgChan chan BroadPayload, joinLeave chan JoinLeave) {
	for {
		select {
		case jl := <-joinLeave:
			if jl.IsJoin {
				msg := fmt.Sprintf("%s has joined our chat...", jl.Name)
				sender(msg, jl.Name)
			} else {
				msg := fmt.Sprintf("%s has left our chat...", jl.Name)
				sender(msg, jl.Name)
			}

		case val := <-msgChan:
			sender(val.Msg, val.Name)
		}
	}
}

func sender(msg, username string) {
	for name, conn := range userQuantity {
		if name != username {
			fmt.Fprint(conn, "\n"+msg)
			fmt.Fprint(conn, "\n"+message(name, ""))
		}
	}
}
