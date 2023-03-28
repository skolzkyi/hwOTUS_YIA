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
	stringres = stringres + "\n IN THISWEEK: \n "

	content, err = s.app.GetListEventsOnWeekByDay(context.Background(), thisDay)
	if err != nil {
		fmt.Println("handlerError: ", err.Error())
		return
	}

	for _, curContent := range content {
		stringres = stringres + " \n " + curContent.String()
	}
	stringres = stringres + "\n IN THISMOUNTH: \n "
	content, err = s.app.GetListEventsOnMonthByDay(context.Background(), thisDay)
	if err != nil {
		fmt.Println("handlerError: ", err.Error())
		return
	}

	for _, curContent := range content {
		stringres = stringres + " \n " + curContent.String()
	}
	for _, curEvent := range content {
		err = s.app.DeleteEvent(context.Background(), curEvent.ID)
		if err != nil {
			fmt.Println("handlerError: ", err.Error())
			return
		}
	}
	//w.Write([]byte("Hello world!"))
	w.Write([]byte(stringres))

}
