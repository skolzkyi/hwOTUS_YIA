package internalhttp

import "net/http"

func (s *Server) helloWorld(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 Not Found"))
		return
	}

	w.Write([]byte("Hello world!"))

}
