package tcp

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/sadovnikoff/GoConcurrencyCourse/homework_2/internal/common"
)

// Client - TCP client
type Client struct {
	conn        net.Conn
	bufferSize  int
	idleTimeout time.Duration
}

// NewClient - returns *Client
func NewClient(addr, maxBufSize string, idleTimeout time.Duration) (*Client, error) {

	if idleTimeout == 0 {
		idleTimeout = defaultIdleTimeout
	}

	bufSize, err := common.ParseBufSize(maxBufSize)
	if err != nil {
		bufSize = defaultBufSize
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to server: %s", err)

	}

	client := &Client{
		conn:        conn,
		idleTimeout: idleTimeout,
		bufferSize:  bufSize,
	}

	return client, nil
}

func (c *Client) Communicate(request string) ([]byte, error) {
	if c.idleTimeout != 0 {
		if err := c.conn.SetDeadline(time.Now().Add(c.idleTimeout)); err != nil {
			return nil, fmt.Errorf("failed to set deadline for connection: %s", err)
		}
	}

	_, err := c.conn.Write([]byte(request))
	if err != nil {
		return nil, fmt.Errorf("error sending request: %s", err)
	}

	response := make([]byte, c.bufferSize)
	count, err := c.conn.Read(response)
	if err != nil {
		return nil, fmt.Errorf("error reading server response: %s", err)
	}

	if count == c.bufferSize {
		return nil, errors.New("error reading server response: too small buffer size")
	}

	return response[:count], nil
}

func (c *Client) Close() {
	if c.conn != nil {
		_ = c.conn.Close()
	}
}
