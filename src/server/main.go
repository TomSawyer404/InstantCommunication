package main

import (
	"fmt"
)

func main() {
	server := NewServer("127.0.0.1", 8888)
	if server == nil {
		fmt.Println("Unable to create a server object!!!")
		return
	}

	server.Start()
}
