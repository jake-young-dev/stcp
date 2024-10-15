package stcp

import (
	"crypto/tls"
	"net"
)

/*
* A simple wrapper for TCP Client setup
 */

// client object
type Client struct {
	connection net.Conn
}

type IClient interface {
	Write(msg string) (string, error)
	Close() error
}

// creates a new tcp client to addr using timeout value as a read/write timeout
// func NewTCPClient(addr string, timeout time.Time) (*Client, error) {
func NewTCPClient(addr string) (*Client, error) {
	//setup connection
	conn, err := tls.Dial("tcp", addr, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return nil, err
	}

	return &Client{
		connection: conn,
	}, nil
}

// writes to the tcp server and waits for a response. Responses must be max 4096 bytes
func (c *Client) Write(msg string) (string, error) {
	_, err := c.connection.Write([]byte(msg))
	if err != nil {
		return "", err
	}

	buff := make([]byte, 4096)
	n, err := c.connection.Read(buff)
	if err != nil {
		return "", err
	}

	return string(buff[:n]), nil
}

// closes tcp connection to server
func (c *Client) Close() error {
	return c.connection.Close()
}
