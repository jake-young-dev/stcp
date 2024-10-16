package stcp

import (
	"context"
	"crypto/tls"
	"net"
	"time"
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

// creates a new tcp client to addr using the certificates supplied. Connection is made with
// a defaulted timeout of 5 seconds
func NewTCPClient(addr, certPath, keyPath string) (*Client, error) {
	//load certs
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}
	cfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	//setup connection
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	dlr := tls.Dialer{
		Config: cfg,
	}
	conn, err := dlr.DialContext(ctx, "tcp", addr)
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
