package handler

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type UserHandler struct {
	Name string
	Conn net.Conn
}

var userQuantity []UserHandler

func (uh *UserHandler) HandleConnection(msgChan chan string, mu *sync.Mutex) {
	mu.Lock()
	welcome(uh.Conn)
	reader := bufio.NewReader(uh.Conn)
	var err error
	uh.Name, err = reader.ReadString('\n')
	if err != nil {
		log.Println(err)
		return
	}
	uh.Name = strings.TrimSpace(uh.Name)
	userQuantity = append(userQuantity, *uh)
	mu.Unlock()
	for {
		msg, err := reader.ReadString('\n')
		fmt.Fprint(uh.Conn, message(uh.Name, ""))
		if err != nil {
			fmt.Println(err)
			return
		}
		msgChan <- msg
	}
}
