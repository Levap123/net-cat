package handler

import (
	"fmt"
)

type BroadPayload struct {
	Name string
	Msg  string
}

func Broadcast(msgChan chan BroadPayload) {
	for val := range msgChan {
		for name, conn := range userQuantity {
			if name != val.Name {
				fmt.Fprint(conn, "\n"+val.Msg)
			}
		}
	}
}
