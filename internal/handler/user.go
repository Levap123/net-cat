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
	conn net.Conn
	mu   *sync.Mutex
	chat *Chat
}

var userQuantity = make(map[string]net.Conn, 10)

func NewUserHandler(conn net.Conn, mu *sync.Mutex, chat *Chat) *UserHandler {
	return &UserHandler{
		conn: conn,
		mu:   mu,
		chat: chat,
	}
}

func (uh *UserHandler) HandleConnection(msgChan chan BroadPayload, joinLeaveChan chan JoinLeave) {
	welcome(uh.conn)
	defer uh.conn.Close()
	if err := uh.addUserName(); err != nil {
		log.Println(err)
		return
	}

	uh.mu.Lock()
	reader := bufio.NewReader(uh.conn)
	if len(userQuantity) >= 10 {
		fmt.Fprint(uh.conn, "\nsorry, chat is full")
		return
	}
	userQuantity[uh.Name] = uh.conn
	uh.chat.printAllBuffer(uh.conn)
	joinLeaveChan <- JoinLeave{IsJoin: true, Name: uh.Name}
	uh.mu.Unlock()
	for {
		fmt.Fprint(uh.conn, message(uh.Name, "\n"))
		msg, err := reader.ReadString('\n')
		if err == io.EOF {
			uh.mu.Lock()
			delete(userQuantity, uh.Name)
			joinLeaveChan <- JoinLeave{IsJoin: false, Name: uh.Name}
			uh.mu.Unlock()
			return
		}
		if err != nil {
			log.Println(err)
			return
		}
		if isValidMsg(msg) {
			uh.mu.Lock()
			msgChan <- BroadPayload{Msg: message(uh.Name, msg), Name: uh.Name}
			uh.mu.Unlock()
		}
	}
}

func (uh *UserHandler) addUserName() error {
	name, err := bufio.NewReader(uh.conn).ReadString('\n')
	if err != nil {
		return err
	}
	name = strings.TrimSpace(name)
	if !isValidName(uh.conn, name) {
		return uh.addUserName()
	}
	uh.Name = name
	return nil
}
