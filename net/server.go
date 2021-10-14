package cnet

import (
	"fmt"
	"net"
)

type Server struct {
	running				bool
	shutdown 			bool
	listener 			net.Listener
	

	shutdownComplete	chan struct{}
}

func (s *Server) WaitForShutdown() {
	<- s.shutdownComplete
}

func (s *Server)isRunning() bool {
	return s.running
}

func (s *Server) Start() {
	clientListenReady := make(chan struct{})

	s.StartLoop(clientListenReady)
}

func (s *Server)StartLoop(clr chan struct{}) {
	defer func()  {
		if clr != nil {
			close(clr)
		}
	}()

	l, err := net.Listen("tcp", "127.0.0.1:8090")
	if err != nil {
		fmt.Println("ListenFailed")
	}

	s.listener = l
	s.running = true

	for s.isRunning() {
		conn, err := s.listener.Accept()
		if err != nil {

		}
	}
}

func NewServer() (*Server, error) {
	s := &Server{
		shutdown: false,
	}

	return s, nil
}

func Run(server *Server) error {
	
	server.Start()

	return nil
}