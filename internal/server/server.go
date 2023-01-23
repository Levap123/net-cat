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
	ch1 := make(chan handler.JoinLeave)
	var mu sync.Mutex
	var mu1 sync.Mutex
	chat := handler.NewChat(&mu1)
	go chat.Broadcast(ch, ch1, &mu1)
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			return err
		}
		user := handler.NewUserHandler(conn, &mu, chat)
		go user.HandleConnection(ch, ch1)
	}
}
