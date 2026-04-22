package server

import (
	"fmt"
	"net"

	"github.com/JoStMc/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	isClosed bool
} 

func newServer(listener *net.Listener) *Server {
	return &Server{listener: *listener, isClosed: false}
}

func Serve(port uint16) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := newServer(&l)
	go s.listen()
	return s, nil
} 


func (s *Server) Close() error {
	err := s.listener.Close()
	if err != nil {
		return err
	}
	err = s.Close()
	if err != nil {
		return err
	}
	s.isClosed = true
	return nil
} 

func (s *Server) handle(conn net.Conn) error {
	err := response.WriteStatusLine(conn, response.StatusOK)
	if err != nil {
		return err
	}
	headers := response.GetDefaultHeaders(0)
	err = response.WriteHeaders(conn, headers)
	if err != nil {
		return err
	}
	return conn.Close()
} 

func (s *Server) listen() error {
	for {
		conn, err := s.listener.Accept()
		if s.isClosed {
			return nil
		} 
		if err != nil {
			return err
		}

		go s.handle(conn)
	} 
} 
