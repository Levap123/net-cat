package server

import (
	"fmt"
	"log"
	"net"
	"net-cat/internal/handler"
	"strings"
	"sync"
)

var userQuantity = 0

type Server struct {
	listener net.Listener
}

func NewServer(address string) (*Server, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		if strings.Contains(err.Error(), "Servname not supported") {
			return nil, fmt.Errorf("port should be a number, %w", err)
		}

		return nil, err
	}

	return &Server{
		listener: listener,
	}, nil
}

func (s *Server) Run() error {
	log.Printf("Listening on the %v\n", s.listener.Addr())
	ch := make(chan string)
	var mu sync.Mutex
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}
		user := &handler.UserHandler{
			Conn: conn,
		}
		go user.HandleConnection(ch, &mu)
	}
}
