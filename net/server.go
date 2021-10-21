package cnet

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"
)

type OnSessionCreate func(s *Session)

type Server struct {
	nextSID  int64
	running  bool
	shutdown bool
	address  string
	listener net.Listener

	connMap map[int64]*Session

	waitgroup sync.WaitGroup
}

func (s *Server) isRunning() bool {
	return s.running
}

func (s *Server) Start(callback OnSessionCreate) {

	s.StartEventLoop(callback)

	s.waitgroup.Wait()

	s.onServerStop()
}

func (s *Server) Stop() {
	fmt.Println("Server::Stop")
	s.shutdown = true
	s.waitgroup.Done()
}

func (s *Server) onServerStop() {
	fmt.Println("Server::onServerStop")
}

func (s *Server) startRoutine(f func()) bool {
	started := false
	if s.isRunning() {
		go f()
		started = true
	}
	return started
}

func (s *Server) createSession(conn net.Conn) *Session {

	atomic.AddInt64(&(s.nextSID), 1)

	c := &Session{
		id:         s.nextSID,
		cn:         conn,
		readQueue:  make(chan interface{}),
		writeQueue: make(chan interface{}),
		done:       make(chan interface{}),
	}
	fmt.Println("Server::createSession")

	s.connMap[c.id] = c

	s.startRoutine(func() { c.StartEventLoop() })
	return c
}

func (s *Server) acceptConnect(createFunc func(conn net.Conn)) {

	defer func() {
		s.waitgroup.Done()
	}()

	for {
		if s.shutdown {
			break
		}
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("Accept Failed")
			continue
		}

		s.startRoutine(func() {
			fmt.Println("Prepare create session in routine")
			createFunc(conn)
		})
	}
}

func (s *Server) StartEventLoop(callback OnSessionCreate) {
	s.waitgroup.Add(1)

	l, err := net.Listen("tcp", s.address)
	if err != nil {
		fmt.Println("ListenFailed")
	}

	s.listener = l
	s.running = true

	//启动一个监听协程用于接收网络连接，创建客户端
	go s.acceptConnect(func(conn net.Conn) {
		session := s.createSession(conn)
		callback(session)
	})
}

func NewServer(addr string) (*Server, error) {
	s := &Server{
		nextSID:  1,
		shutdown: false,
		address:  addr,
		connMap:  make(map[int64]*Session),
	}

	return s, nil
}
