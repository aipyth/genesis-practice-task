package main

import (
	"fmt"
	server "github.com/aipyth/genesis-practice-task/internal/api"
)

func main() {
	if err := server.Run(); err != nil {
		fmt.Println(err)
	}
}
