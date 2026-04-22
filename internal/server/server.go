package server

import (
	"fmt"
	"net"
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
	out := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!\n"
	_, err := conn.Write([]byte(out))
	if err != nil {
		return err
	}
	return s.Close()
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
