package internalhttp

import "net/http"

func (s *Server) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/404", app.page404)

	return mux
}
