package internalhttp

import (
	"net/http"
)

func (s *Server) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", loggingMiddleware(s.helloWorld, s.logg))

	return mux
}
