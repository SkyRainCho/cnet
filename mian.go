package main

import (
	"fmt"
	"os"

	example "github.com/SkyRainCho/cnet/example/echo"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Invalid args")
		return
	}

	if os.Args[1] == "client" {
		example.RunClient("127.0.0.1:8090")
	} else if os.Args[1] == "server" {
		example.RunServer("0.0.0.0:8090")
	} else if os.Args[1] == "test" {
		example.RunClient("127.0.0.1:8090")
	}
}
