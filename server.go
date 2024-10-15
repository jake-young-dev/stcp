package stcp

import (
	"crypto/tls"
	"net"
)

/*
* A simple wrapper for TCP Server setup
 */

// handler function type, spawned in new goroutine for each connection accepted
type ConnectionHandler func(c net.Conn)

// server object
type Server struct {
	connection net.Listener
}

type IServer interface {
	Connect(handler ConnectionHandler)
	Close() error
}

// creates a new tcp server on addr using the supplied certificates
func NewTCPServer(addr, certPath, keyPath string) (*Server, error) {
	//load certs
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}

	//create and configure server
	cfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	listener, err := tls.Listen("tcp", addr, cfg)
	if err != nil {
		return nil, err
	}

	return &Server{
		connection: listener,
	}, nil
}

// loops indefinitely accepting and handling incoming tcp connections. This is a blocking function that is
// cancelled when the Close method is called
func (s *Server) Connect(handler ConnectionHandler) {
	for {
		//accept connection, will error when connection is closed
		conn, err := s.connection.Accept()
		if err != nil {
			return
		}

		go handler(conn)
	}
}

// closes tcp server connection
func (s *Server) Close() error {
	return s.connection.Close()
}
