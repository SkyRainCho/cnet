package example

import (
	"fmt"

	cnet "github.com/SkyRainCho/cnet/net"
)

type EchoServerCodecHandler struct {
}

func (handler *EchoServerCodecHandler) Read(s *cnet.Session, buf []byte) (interface{}, int, error) {
	fmt.Println("EchoServerCodecHandler::Read::MsgLen:", len(buf))

	return buf, len(buf), nil
}

func (handler *EchoServerCodecHandler) Write(s *cnet.Session, msg interface{}) error {
	fmt.Println("EchoServerCodeHandler::Write")
	if v, ok := msg.([]byte); ok {
		s.WriteBytes(v)
	}
	return nil
}

type EchoServerEventHandler struct {
}

func (handler *EchoServerEventHandler) OnConnect(s *cnet.Session) error {
	fmt.Println("EchoServerEventHandler::OnCennect")
	return nil
}
func (handler *EchoServerEventHandler) OnDisconnect(s *cnet.Session) {
	fmt.Println("EchoServerEventHandler::OnDisconnect")

}
func (handler *EchoServerEventHandler) OnAbortConnect(s *cnet.Session, err error) {
	fmt.Println("EchoServerEventHandler::OnAbortConnect")

}
func (handler *EchoServerEventHandler) OnHeartbeat(s *cnet.Session) {
	//fmt.Println("EchoServerEventHandler::OnHeartbeat")

}
func (handler *EchoServerEventHandler) OnHandleMsg(s *cnet.Session, msg interface{}) {
	fmt.Println("EchoServerEventHandler::OnHandleMsg")
	if v, ok := msg.([]byte); ok {
		fmt.Println("EchoServerEventHandler::OnHandleMsg::", string(v))
	}
	s.WriteMsgToPending(msg)
}
func RunServer(address string) {
	s, err := cnet.NewServer(address)
	if err != nil {
		fmt.Println("NewServer::Failed")
	}

	s.Start(func(s *cnet.Session) {
		fmt.Println("Session::OnSessionCreate:", s.GetSessionID())
		s.SetIOHandler(&EchoServerCodecHandler{})
		s.SetEventHandler(&EchoServerEventHandler{})

	})
}
