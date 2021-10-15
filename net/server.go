package cnet

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	running  bool
	shutdown bool
	listener net.Listener

	shutdownComplete chan struct{}
}

func (s *Server) WaitForShutdown() {
	<-s.shutdownComplete
}

func (s *Server) isRunning() bool {
	return s.running
}

func (s *Server) Start() {
	clientListenReady := make(chan struct{})

	s.StartLoop(clientListenReady)
}

func (s *Server) startRoutine(f func()) bool {
	started := false
	if s.isRunning() {
		go f()
		started = true
	}
	return started
}

func (s *Server) createClient(conn net.Conn) *client {
	c := &client{
		cn: conn,
	}
	fmt.Println("Server::createClient")

	c.lock.Lock()
	c.pending.cond = sync.NewCond(&(c.lock))
	c.lock.Unlock()

	s.startRoutine(func() { c.Read() })

	s.startRoutine(func() { c.Write() })
	return c
}

func (s *Server) acceptConnect(createFunc func(conn net.Conn)) {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("Accept Failed")
			continue
		}

		s.startRoutine(func() {
			fmt.Println("Prepare create client in routine")
			createFunc(conn)
		})
	}
}

func (s *Server) StartLoop(clr chan struct{}) {
	defer func() {
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

	//启动一个监听协程用于接收网络连接，创建客户端
	go s.acceptConnect(func(conn net.Conn) { s.createClient(conn) })
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
