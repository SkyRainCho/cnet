package cnet

import (
	"bytes"
	"fmt"
	"net"
	"time"
)

//每个Session都应该具有一个SessionID用于唯一标记
//同时具有一个输出缓冲区
type Session struct {
	id           int64
	cn           net.Conn
	ioHandler    IOHandler
	eventHandler EventHandler

	readQueue  chan interface{}
	writeQueue chan interface{}
}

func (s *Session) SetIOHandler(io IOHandler) {
	fmt.Println("Session::SetIOHandler")
	s.ioHandler = io
}

func (s *Session) SetEventHandler(event EventHandler) {
	fmt.Println("Session::SetEventHandler")
	s.eventHandler = event
}

func (s *Session) WriteMsgToPending(msg interface{}) {
	fmt.Println("Session::WriteMsgToPending")

	if v, ok := msg.([]byte); ok {
		fmt.Println("EchoServerEventHandler::WriteMsgToPending::", string(v))
	}
	s.writeQueue <- msg

	fmt.Println("Session::WriteMsgToPending::SendSuccess")
}

func (s *Session) WriteBytes(pkg []byte) error {
	_, err := s.cn.Write(pkg)
	return err
}

func (s *Session) StartEventLoop() {

	//先初始化两个通道
	s.readQueue = make(chan interface{})
	s.writeQueue = make(chan interface{})

	//初次启动事件循环，需要调动EventHandler的OnConnect的接口来处理建立连接的回调函数
	s.eventHandler.OnConnect(s)

	//启动协程用于处理io输入输出
	go s.runDecodeLoop()
	go s.runProcessLoop()
	go s.runEncodeLoop()
}

func (s *Session) runDecodeLoop() {
	fmt.Println("Seesion::runDecodeLoop")
	//尝试从缓冲区之中读取数据，并将其进行解码反序列化
	//如果反序列化成功，则将反序序列话的数据放入readQueue

	buf := make([]byte, 2048)
	bufStream := new(bytes.Buffer)
	var err error
	var inputMsg interface{}
	var msgLen int
	var bufLen int
	exit := false
	for {
		bufLen = 0
		for {
			fmt.Println("Seesion::runDecodeLoop::ReadRawDataFromSocket")
			bufLen, err = s.cn.Read(buf)

			if err != nil {
				fmt.Println("Seesion::runDecodeLoop::ReadFaile:", err.Error())
				exit = true
				break
			}
			break
		}

		if exit {
			break
		}

		if bufLen == 0 {
			continue
		}

		bufStream.Write(buf[:bufLen])
		//尽可能多的从bufStream之中解析出结构化的数据来
		for {
			if bufStream.Len() == 0 {
				break
			}
			inputMsg, msgLen, err = s.ioHandler.Read(s, bufStream.Bytes())
			if err != nil {
				exit = true
				break
			}
			if inputMsg == nil {
				break
			}
			s.readQueue <- inputMsg
			bufStream.Next(msgLen)
		}

		if exit {
			break
		}
	}
}

func (s *Session) runProcessLoop() {
	fmt.Println("Seesion::runProcessLoop")
	var (
		inputMsg interface{}
	)
	ticker := time.NewTicker(1e9)
	for {
		select {
		case inputMsg = <-s.readQueue:
			fmt.Println("Seesion::runProcessLoop::readQueueReady")
			s.eventHandler.OnHandleMsg(s, inputMsg)
		case <-ticker.C:
			s.eventHandler.OnHeartbeat(s)
		}
	}
}

func (s *Session) runEncodeLoop() {
	fmt.Println("Seesion::runEncodeLoop")
	var (
		outputMsg interface{}
	)

	for {
		select {
		case outputMsg = <-s.writeQueue:
			fmt.Println("Seesion::runEncodeLoop::writeQueueReady")
			s.ioHandler.Write(s, outputMsg)
		}
	}
}

func (s *Session) GetSessionID() int64 {
	return s.id
}
func (c *Session) CloseConnect() {
	c.cn.Close()
}
