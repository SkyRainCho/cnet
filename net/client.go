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

type client struct {
	cn      net.Conn
	lock    sync.Mutex
	pending output
}

func (c *client) CloseConnect() {
	c.cn.Close()
}

func (c *client) FillPendingData(data []byte, length int) {
	c.pending.length += length
	c.pending.data = append(c.pending.data, data[:length]...)
	c.pending.cond.Signal()
}

func (c *client) FlushPendingData() {
	fmt.Printf("Client::FlushPendingData::%d\n", c.pending.length)
	_, err := c.cn.Write(c.pending.data[:c.pending.length])
	if err != nil {
		fmt.Println("Client::FlushPendingData::WriteError")
	}
	c.pending.length = 0
	c.pending.data = c.pending.data[0:0]
}

func (c *client) Read() {
	fmt.Println("Client::Read")

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

func (c *client) Write() {
	fmt.Println("Client::Write")

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
