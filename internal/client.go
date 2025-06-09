package internal

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type Conn struct {
	conn net.Conn
}

func NewCon(conn net.Conn) *Conn {
	return &Conn{conn: conn}
}

func (c *Conn) HandleRequest(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	defer c.conn.Close()

	reader := bufio.NewReader(c.conn)

	_ = c.conn.SetReadDeadline(time.Now().Add(5 * time.Minute))

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			c.conn.Close()
			return
		}

		select {
		case <-ctx.Done():
			fmt.Println("Server shutting down, closing connection.")
			return
		default:
		}

		fmt.Printf("Incoming message: %s", string(msg))

		splitMsg := strings.Split(strings.TrimSpace(msg), " ")

		if len(splitMsg) < 2 {
			c.conn.Write([]byte("Invalid message format. Use 'command argument'.\n"))
		}

		action := splitMsg[0]
		key := splitMsg[1]
		value := splitMsg[2]

		command := NewCommand(action, key, value)
		if !command.IsValid() {
			c.conn.Write([]byte("Invalid command.\n"))
		}
		response := command.Execute()
		c.conn.Write([]byte(response + "\n"))
	}
}
