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

var userQuantity = make(map[string]net.Conn, 10)

func (uh *UserHandler) HandleConnection(msgChan chan BroadPayload, mu *sync.Mutex) {
	welcome(uh.Conn)
	reader := bufio.NewReader(uh.Conn)
	if err := uh.addUserName(); err != nil {
		log.Println(err)
		return
	}
	mu.Lock()

	if len(userQuantity) > 10 {
		fmt.Fprint(uh.Conn, "\nsorry, chat is full")
		uh.Conn.Close()
		return
	}
	userQuantity[uh.Name] = uh.Conn

	mu.Unlock()
	msgChan <- BroadPayload{Msg: fmt.Sprintf("\n%s has joined our chat...", uh.Name), Name: uh.Name}
	for {
		fmt.Fprint(uh.Conn, message(uh.Name, ""))
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		if isValidMsg(msg) {
			msgChan <- BroadPayload{Msg: message(uh.Name, msg), Name: uh.Name}
		}
	}
}

func (uh *UserHandler) addUserName() error {
	name, err := bufio.NewReader(uh.Conn).ReadString('\n')
	if err != nil {
		return err
	}
	name = strings.TrimSpace(name)
	if !isValidName(uh.Conn, name) {
		return uh.addUserName()
	}
	uh.Name = name
	return nil
}
