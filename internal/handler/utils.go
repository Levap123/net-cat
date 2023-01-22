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
