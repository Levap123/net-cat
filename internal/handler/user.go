package handler

import (
	"bufio"
	"fmt"
	"io"
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

func (uh *UserHandler) HandleConnection(msgChan chan BroadPayload, jlChan chan JoinLeave, mu *sync.Mutex) {
	welcome(uh.Conn)
	reader := bufio.NewReader(uh.Conn)
	defer uh.Conn.Close()
	if err := uh.addUserName(); err != nil {
		log.Println(err)
		return
	}
	mu.Lock()

	if len(userQuantity) > 10 {
		fmt.Fprint(uh.Conn, "\nsorry, chat is full")
		return
	}
	userQuantity[uh.Name] = uh.Conn

	mu.Unlock()
	jlChan <- JoinLeave{IsJoin: true, Name: uh.Name}
	for {
		fmt.Fprint(uh.Conn, message(uh.Name, "\n"))
		msg, err := reader.ReadString('\n')
		if err == io.EOF {
			mu.Lock()
			delete(userQuantity, uh.Name)
			mu.Unlock()
			jlChan <- JoinLeave{IsJoin: false, Name: uh.Name}
		}
		if err != nil {
			log.Println(err)
			return
		}
		if isValidMsg(msg) {
			// fmt.Fprintf(uh.Conn, "\n%s", message(uh.Name, "\n"))
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
