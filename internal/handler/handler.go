package handler

import (
	"fmt"
	"net"
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

type Chat struct {
	MsgBuffer []string
	Mu        *sync.Mutex
}

func NewChat(mu *sync.Mutex) *Chat {
	return &Chat{
		Mu: mu,
	}
}
func (c *Chat) Broadcast(msgChan chan BroadPayload, joinLeaveChan chan JoinLeave, mu *sync.Mutex) {
	for {
		select {
		case jl := <-joinLeaveChan:
			if jl.IsJoin {
				msg := fmt.Sprintf("%s has joined our chat...", jl.Name)
				c.send(msg, jl.Name, mu)
			} else {
				msg := fmt.Sprintf("%s has left our chat...", jl.Name)
				c.send(msg, jl.Name, mu)
			}

		case val := <-msgChan:
			c.send(val.Msg, val.Name, mu)
		}
	}
}

func (c *Chat) send(msg, username string, mu *sync.Mutex) {
	c.Mu.Lock()
	c.MsgBuffer = append(c.MsgBuffer, "\n"+msg)
	c.Mu.Unlock()
	for name, conn := range userQuantity {
		if name != username {
			fmt.Fprint(conn, "\n"+msg)
			fmt.Fprint(conn, "\n"+message(name, ""))
		}
	}
}

func (chat *Chat) printAllBuffer(c net.Conn) {
	// fmt.Println(msgBuffer)
	for ind, buff := range chat.MsgBuffer {
		if ind == 0 {
			fmt.Fprint(c, buff[1:])
			continue
		}
		if ind == len(chat.MsgBuffer)-1 {
			fmt.Fprint(c, buff+"\n")
			continue
		}
		fmt.Fprint(c, buff)
	}
}
