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
	Listener net.Listener
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
		Listener: listener,
	}, nil
}

func (s *Server) Run() error {
	log.Printf("Listening on the %v\n", s.Listener.Addr())
	ch := make(chan handler.BroadPayload)
	var mu sync.Mutex
	go handler.Broadcast(ch)
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			return err
		}
		user := &handler.UserHandler{
			Conn: conn,
		}
		go user.HandleConnection(ch, &mu)
	}
}
