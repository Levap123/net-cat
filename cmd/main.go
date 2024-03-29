package main

import (
	"fmt"
	"log"
	"net-cat/internal/server"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	port := "8989"

	if len(os.Args) == 2 {
		port = os.Args[1]
	}

	address := "localhost:" + port
	server, err := server.NewServer(address)
	panic(err)
	defer server.Listener.Close()
	err = server.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
