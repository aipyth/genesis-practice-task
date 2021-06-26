package main

import (
	"fmt"
	server "github.com/aipyth/genesis-practice-task/cmd/server"
)

func main() {
	if err := server.Run(); err != nil {
		fmt.Println(err)
	}
}
