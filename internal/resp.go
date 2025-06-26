package internal

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strconv"
	"sync"
)

// Example client data sending "$6\r\nfoobar\r\n" to server

const (
	STRING  = '+'
	ARRAY   = '*'
	INTEGER = ':'
	ERROR   = '-'
	BULK    = '$'
)

type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(conn net.Conn) *Resp {
	return &Resp{reader: bufio.NewReader(conn)}
}

func (r *Resp) Read(wg *sync.WaitGroup, ctx context.Context) (value Value, err error) {
	defer wg.Done()

	_type, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch _type {
	case BULK:
		return r.readBulk()
	case ARRAY:
		return r.readArray()
	default:
		fmt.Printf("Unknown type: %v", string(_type))
		return Value{}, nil
	}
}

func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		// Read a line from the reader until we hit \r\n
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	// Ensure the last two characters are \r\n
	return line[:len(line)-2], n, nil
}

func (r *Resp) readBulk() (Value, error) {
	v := Value{typ: "bulk"}

	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, len)

	r.reader.Read(bulk)

	v.bulk = string(bulk)

	r.readLine()

	return v, nil
}

func (r *Resp) readInteger() (x int, n int, err error) {
	// Read the line containing the integer: // e.g., ":123\r\n"
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}

func (r *Resp) readArray() (Value, error) {
	v := Value{typ: "array"}

	// Read the length of the array
	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	v.array = make([]Value, len)

	for i := 0; i < len; i++ {
		// Read each value in the array
		value, err := r.Read(nil, context.Background())
		if err != nil {
			return v, err
		}

		v.array[i] = value
	}
	return v, nil
}
