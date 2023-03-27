package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func (s *Server) helloWorld(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 Not Found"))
		return
	}
	thisDay := time.Now()
	var stringres string
	stringres = "IN THISDAY: \n "
	content, err := s.app.GetListEventsonDayByDay(context.Background(), thisDay)
	if err != nil {
		fmt.Println("handlerError: ", err.Error())
		return
	}

	for _, curContent := range content {
		stringres = stringres + " \n " + curContent.String()
	}
	stringres = " \n " + stringres + "IN THISWEEK: \n "

	content, err = s.app.GetListEventsOnWeekByDay(context.Background(), thisDay)
	if err != nil {
		fmt.Println("handlerError: ", err.Error())
		return
	}

	for _, curContent := range content {
		stringres = stringres + " \n " + curContent.String()
	}
	stringres = " \n " + stringres + "IN THISMOUNTH: \n "
	content, err = s.app.GetListEventsOnMonthByDay(context.Background(), thisDay)
	if err != nil {
		fmt.Println("handlerError: ", err.Error())
		return
	}

	for _, curContent := range content {
		stringres = stringres + " \n " + curContent.String()
	}
	//w.Write([]byte("Hello world!"))
	w.Write([]byte(stringres))

}
