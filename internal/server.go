package internal

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
)

type Server struct {
	port string
	host string
}

func NewServer(port string, host string) *Server {
	return &Server{
		port: port,
		host: host,
	}
}

func (s *Server) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}

	fmt.Printf("Server started on %s:%s\n", s.host, s.port)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		<-ctx.Done()
		_ = listener.Close()
	}()

	var wg sync.WaitGroup
	for {
		conn, err := listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				fmt.Println("Server shutting down.")
				break
			}
			continue
		}

		fmt.Printf("New connection from %s \n", conn.RemoteAddr().String())
		resp := NewResp(conn)

		value, err := resp.Read(&wg, ctx)
		if err != nil {
			fmt.Printf("Error reading from connection: %v\n", err)
		}
		fmt.Println(value)
		conn.Write([]byte("+OK\r\n"))
		wg.Add(1)
	}

	wg.Wait()
	fmt.Println("Server stopped.")
}
