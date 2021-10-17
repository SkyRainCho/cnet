package cnet

import (
	"fmt"
	"net"
	"sync"
)

type output struct {
	data   []byte
	length int

	cond *sync.Cond
}

//每个Session都应该具有一个SessionID用于唯一标记
//同时具有一个输出缓冲区
type Session struct {
	id      int64
	cn      net.Conn
	lock    sync.Mutex
	pending output
}

func (c *Session) CloseConnect() {
	c.cn.Close()
}

func (c *Session) FillPendingData(data []byte, length int) {
	c.pending.length += length
	c.pending.data = append(c.pending.data, data[:length]...)
	c.pending.cond.Signal()
}

func (c *Session) FlushPendingData() {
	fmt.Printf("Session::FlushPendingData::%d\n", c.pending.length)
	_, err := c.cn.Write(c.pending.data[:c.pending.length])
	if err != nil {
		fmt.Println("Session::FlushPendingData::WriteError")
	}
	c.pending.length = 0
	c.pending.data = c.pending.data[0:0]
}

func (c *Session) Read() {
	fmt.Println("Session::Read")

	//创建一个读取缓冲区
	buf := make([]byte, 2048)

	for {
		n, err := c.cn.Read(buf)
		if n == 0 && err != nil {
			c.CloseConnect()
			return
		}
		fmt.Println("ReadMsg:", string(buf[:n]))
		c.FillPendingData(buf, n)
	}
}

func (c *Session) Write() {
	fmt.Println("Session::Write")

	//输出缓冲区之中没有数据的话，协程等待
	for {
		c.lock.Lock()
		if c.pending.length == 0 {
			c.pending.cond.Wait()
		}
		c.FlushPendingData()
		c.lock.Unlock()
	}
}
