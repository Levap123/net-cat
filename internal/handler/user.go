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
	Mu   *sync.Mutex
	Chat *Chat
}

var userQuantity = make(map[string]net.Conn, 10)

func NewUserHandler(conn net.Conn, mu *sync.Mutex, chat *Chat) *UserHandler {
	return &UserHandler{
		Conn: conn,
		Mu:   mu,
		Chat: chat,
	}
}

func (uh *UserHandler) HandleConnection(msgChan chan BroadPayload, joinLeaveChan chan JoinLeave) {
	welcome(uh.Conn)
	reader := bufio.NewReader(uh.Conn)
	defer uh.Conn.Close()
	if err := uh.addUserName(); err != nil {
		log.Println(err)
		return
	}

	uh.Mu.Lock()
	if len(userQuantity) >= 10 {
		fmt.Fprint(uh.Conn, "\nsorry, chat is full")
		return
	}
	userQuantity[uh.Name] = uh.Conn
	uh.Mu.Unlock()
	uh.Chat.printAllBuffer(uh.Conn)
	joinLeaveChan <- JoinLeave{IsJoin: true, Name: uh.Name}
	for {
		fmt.Fprint(uh.Conn, message(uh.Name, "\n"))
		msg, err := reader.ReadString('\n')
		if err == io.EOF {
			uh.Mu.Lock()
			delete(userQuantity, uh.Name)
			uh.Mu.Unlock()
			joinLeaveChan <- JoinLeave{IsJoin: false, Name: uh.Name}
			return
		}
		if err != nil {
			log.Println(err)
			return
		}
		if isValidMsg(msg) {
			uh.Mu.Lock()
			msgChan <- BroadPayload{Msg: message(uh.Name, msg), Name: uh.Name}
			uh.Mu.Unlock()
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
