package server

import (
	"fmt"
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


func (hErr *HandlerError) Write(w *response.Writer) {
	w.WriteStatusLine(hErr.statusCode)
	headers := response.GetDefaultHeaders(len(hErr.errorMessage))
	w.WriteHeaders(headers)
	w.WriteBody(hErr.errorMessage)
} 

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	responseWriter := response.NewWriter(conn)
	req, err := request.RequestFromReader(conn)
	if err != nil {
		handlerError := NewHandlerError(response.StatusBadRequest, []byte(err.Error()))
		handlerError.Write(responseWriter)
		return 
	}
	s.handler(responseWriter, req)
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
