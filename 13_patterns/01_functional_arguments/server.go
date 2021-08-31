package main

import (
	"fmt"
	"net/http"
)

type Server struct {
	addr string
	mux  *http.ServeMux
}

func NewServer(ops ...serverOption) *Server {
	s := &Server{}
	for _, op := range ops {
		op(s)
	}
	return s
}

func (s *Server) ListenAndServe() error {
	if s.addr == "" {
		return fmt.Errorf("addr must not be empty")
	}
	return http.ListenAndServe(s.addr, s.mux)
}

type serverOption func(*Server)

func WithAddr(addr string) serverOption {
	return func(s *Server) {
		s.addr = addr
	}
}

func WithMux(mux *http.ServeMux) serverOption {
	return func(s *Server) {
		s.mux = mux
	}
}

func WithRoute(pattern string, handlerFunc http.HandlerFunc) serverOption {
	return func(s *Server) {
		if s.mux == nil {
			s.mux = http.NewServeMux()
		}
		s.mux.HandleFunc(pattern, handlerFunc)
	}
}
