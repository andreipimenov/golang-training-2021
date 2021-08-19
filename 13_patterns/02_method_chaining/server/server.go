package server

import "net/http"

type Server struct {
	addr string
	mux  *http.ServeMux
}

func New() Server {
	return Server{}
}

func (s Server) Addr(addr string) Server {
	s.addr = addr
	return s
}

func (s Server) Mux(mux *http.ServeMux) Server {
	s.mux = mux
	return s
}

func (s Server) Route(pattern string, handlerFunc http.HandlerFunc) Server {
	if s.mux == nil {
		s.mux = http.NewServeMux()
	}
	s.mux.HandleFunc(pattern, handlerFunc)
	return s
}

func (s Server) ListenAndServe() error {
	return http.ListenAndServe(s.addr, s.mux)
}
