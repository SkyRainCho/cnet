package example

import (
	"fmt"

	cnet "github.com/SkyRainCho/cnet/net"
)

func RunServer(address string) {
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
