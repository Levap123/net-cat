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
	msgBuffer []string
	mu        *sync.Mutex
}

func NewChat(mu *sync.Mutex) *Chat {
	return &Chat{
		mu: mu,
	}
}
func (c *Chat) Broadcast(msgChan chan BroadPayload, joinLeaveChan chan JoinLeave) {
	for {
		select {
		case jl := <-joinLeaveChan:
			if jl.IsJoin {
				msg := fmt.Sprintf("%s has joined our chat...", jl.Name)
				c.send(msg, jl.Name)
			} else {
				msg := fmt.Sprintf("%s has left our chat...", jl.Name)
				c.send(msg, jl.Name)
			}

		case val := <-msgChan:
			c.send(val.Msg, val.Name)
		}
	}
}

func (c *Chat) send(msg, username string) {
	c.mu.Lock()
	c.msgBuffer = append(c.msgBuffer, "\n"+msg)
	c.mu.Unlock()
	for name, conn := range userQuantity {
		if name != username {
			fmt.Fprint(conn, "\n"+msg)
			fmt.Fprint(conn, "\n"+message(name, ""))
		}
	}
}

func (chat *Chat) printAllBuffer(c net.Conn) {
	for ind, buff := range chat.msgBuffer {
		if ind == 0 {
			fmt.Fprint(c, buff[1:])
			continue
		}
		if ind == len(chat.msgBuffer)-1 {
			fmt.Fprint(c, buff+"\n")
			continue
		}
		fmt.Fprint(c, buff)
	}
}
