package handler

import (
	"fmt"
	"net"
	"os"
	"time"
)

func welcome(conn net.Conn) {
	file, _ := os.ReadFile("hello.txt")
	fmt.Fprint(conn, string(file))
}

func message(name, data string) string {
	realTtime := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("\n[%s][%s]:%s", realTtime, name, data)
}

func isValidMsg(msg string) bool {
	return !(msg == "" || (msg[0] >= 0 && msg[0] <= 31))
}

func isValidName(c net.Conn, name string) bool {
	for _, ch := range name {
		if ch >= 0 && ch <= 31 {
			fmt.Fprint(c, "Not valid name\nPlease try again...\n[ENTER YOUR NAME]:")
			return false
		}
	}

	if name == "" {
		fmt.Fprint(c, "Not valid name\nPlease try again...\n[ENTER YOUR NAME]:")
		return false
	}

	if _, ok := userQuantity[name]; ok {
		fmt.Fprint(c, "Username has already taken\nPlease try again...\n[ENTER YOUR NAME]:")
		return false
	}

	return true
}
