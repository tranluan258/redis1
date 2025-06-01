package main

import (
	"fmt"

	"github.com/tranluan258/redis1/internal"
)

func main() {
	server := internal.NewServer("6379", "localhost")
	fmt.Println("Starting server...")
	server.Run()
}
