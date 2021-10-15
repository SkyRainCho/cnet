package main

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Invalid args")
		return
	}

	if os.Args[1] == "client" {
		runClient("127.0.0.1:8090")
	} else if os.Args[1] == "server" {
		runServer("0.0.0.0:8090")
	} else if os.Args[1] == "test" {
		runTest("127.0.0.1:8090")
	}
}

func runClient(address string)  {
	fmt.Println("runClient:", address)

	conn, err := net.Dial("tcp", address)
	
	if err != nil {
		fmt.Println("runClient::Dial::Failed")
		return
	}

	rdDone := make(chan int)
	wtDone := make(chan int)

	msg_size := len(msg_list)

	go func ()  {
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

	go func () {
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

func runTest(address string) {
	done := make(chan int)

	for i := 0; i < 10000; i++ {
		go runClient(address)
	}
	<-done
}
func runServer(address string)  {
	s, err := cnet.NewServer(address)
	if err != nil {
		fmt.Println("NewServer::Failed")
	}
	
	err = cnet.Run(s)

	if err != nil {
		fmt.Println("RunServer::Failed")
	}
	
	s.WaitForShutdown()
}
