package internalhttp

import (
	"net/http"
	"time"

	helpers "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/helpers"
)

func loggingMiddleware(next http.HandlerFunc, logg Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		next.ServeHTTP(w, r)
		message := helpers.StringBuild("[client IP: ", r.RemoteAddr, " Request DateTime: ", time.Now().String(), " Method: ", r.Method, " Request URL: ", r.RequestURI, " Request Scheme: ", r.URL.Scheme, "Request Status: ", w.Header().Get("Status"), "Time of request work: ", time.Since(t).String(), " Request User-Agent: ", r.Header.Get("User-Agent"))
		logg.Info(message)
	}
}
