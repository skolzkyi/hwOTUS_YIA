package internalhttp

import (
	"net/http"
)

func (s *Server) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", loggingMiddleware(s.helloWorld, s.logg))
	mux.HandleFunc("/Event/", loggingMiddleware(s.Event_REST, s.logg))
	mux.HandleFunc("/GetEventsOnDayByDay/", loggingMiddleware(s.GetEventsOnDayByDay, s.logg))
	mux.HandleFunc("/GetEventsOnWeekByDay/", loggingMiddleware(s.GetEventsOnWeekByDay, s.logg))
	mux.HandleFunc("/GetEventsOnMonthByDay/", loggingMiddleware(s.GetEventsOnMonthByDay, s.logg))

	return mux
}
