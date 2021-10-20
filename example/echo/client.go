package example

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"time"

	cnet "github.com/SkyRainCho/cnet/net"
)

var msg_list = []string{
	"Hello world!",
	"This is a test message",
	"My name is skyraincho",
	"I live in BeiJing",
	"How are you?",
	"I am fine, and you?",
}

type EchoClient struct {
	session cnet.Session
}

func RunClient(address string) {
	fmt.Println("runClient:", address)

	conn, err := net.Dial("tcp", address)

	if err != nil {
		fmt.Println("runClient::Dial::Failed")
		return
	}

	rdDone := make(chan int)
	wtDone := make(chan int)

	msg_size := len(msg_list)

	go func() {
		for {
			buf := make([]byte, 1024)
			length, err := conn.Read(buf)
			if err != nil {
				fmt.Println("runClient:read::Failed")

				rdDone <- 1

				return
			}

			fmt.Println("Received data:", string(buf[:length]))
		}
	}()

	go func() {
		tick := time.Tick(1e9)
		for {
			select {
			case <-tick:
				index := rand.Int() % msg_size
				io.WriteString(conn, msg_list[index])
			}
		}
	}()

	<-rdDone
	<-wtDone
}

func RunTest(address string) {
	done := make(chan int)

	for i := 0; i < 10000; i++ {
		go RunClient(address)
	}
	<-done
}
