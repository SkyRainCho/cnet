package main

import (
	"fmt"
	"os"

	cnet "github.com/SkyRainCho/cnet/net"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Invalid args")
		return
	}

	if os.Args[1] == "client" {
		runClient()
	} else if os.Args[1] == "server" {
		runServer()
	}
}

func runClient()  {
	
}

func runServer()  {
	s, err := cnet.NewServer()
	if err != nil {
		fmt.Println("NewServer::Failed")
	}
	
	err = cnet.Run(s)

	if err != nil {
		fmt.Println("RunServer::Failed")
	}
	
	s.WaitForShutdown()
}
