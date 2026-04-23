package server

import (
	"bytes"
	"fmt"
	"io"
	"net"

	"github.com/JoStMc/httpfromtcp/internal/request"
	"github.com/JoStMc/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	isClosed bool
	handler Handler
} 

func newServer(listener *net.Listener, handler Handler) *Server {
	return &Server{listener: *listener, isClosed: false, handler: handler}
}

func Serve(port uint16, handler Handler) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := newServer(&l, handler)
	go s.listen()
	return s, nil
} 


func (s *Server) Close() error {
	err := s.listener.Close()
	if err != nil {
		return err
	}
	s.isClosed = true
	return nil
} 

func writeHandlerError(w io.Writer, handlerError *HandlerError) error {
	err := response.WriteStatusLine(w, handlerError.statusCode)
	if err != nil {
		return err
	}
	_, err = w.Write(handlerError.errorMessage)
	return err
} 

func (s *Server) handle(conn net.Conn) error {
	req, err := request.RequestFromReader(conn)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer([]byte{})

	handlerError := s.handler(conn, req)
	if handlerError != nil {
		err = writeHandlerError(buf, handlerError)
		if err != nil {
			return err
		}
		return nil
	} 

	err = response.WriteStatusLine(conn, response.StatusOK)
	if err != nil {
		return err
	}
	headers := response.GetDefaultHeaders(buf.Len())
	err = response.WriteHeaders(conn, headers)
	if err != nil {
		return err
	}
	conn.Write(buf.Bytes())
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
