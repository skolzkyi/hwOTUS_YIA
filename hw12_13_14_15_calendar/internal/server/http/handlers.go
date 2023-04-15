package internalhttp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	helpers "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/helpers"
	storage "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/storage/event"
)

type outputJSON struct {
	Text string
	Code int
}

type EventRawData struct {
	EventMessageTimeDelta int64
	Title                 string
	UserID                string
	Description           string
	DateStart             string
	DateStop              string
	ID                    int
}
type EventAnswer struct {
	Events  []storage.Event
	Message outputJSON
}

type InputDate struct {
	Date string
}

var (
	ErrInJSONBadParse     = errors.New("error parsing input json")
	ErrOutJSONBadParse    = errors.New("error parsing output json")
	ErrUnsupportedMethod  = errors.New("unsupported method")
	ErrNoIDInEventHandler = errors.New("no ID in event handler")
)

func apiErrHandler(err error, w *http.ResponseWriter) {
	W := *w
	newMessage := outputJSON{}
	newMessage.Text = err.Error()
	newMessage.Code = 1
	jsonstring, err := json.Marshal(newMessage)
	if err != nil {
		errMessage := helpers.StringBuild(http.StatusText(http.StatusInternalServerError), " (", err.Error(), ")")
		http.Error(W, errMessage, http.StatusInternalServerError)
	}

	_, err = W.Write(jsonstring)
	if err != nil {
		errMessage := helpers.StringBuild(http.StatusText(http.StatusInternalServerError), " (", err.Error(), ")")
		http.Error(W, errMessage, http.StatusInternalServerError)
	}
}

func (s *Server) helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}

func (s *Server) Event_REST(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, _ := context.WithTimeout(context.Background(), s.Config.GetDBTimeOut())

	switch r.Method {

	case http.MethodGet:

		fmt.Println("Get")
		newMessage := outputJSON{}
		EvAnswer := EventAnswer{}
		EvAnswer.Events = make([]storage.Event, 1, 1)
		path := strings.Trim(r.URL.Path, "/")
		pathParts := strings.Split(path, "/")
		if len(pathParts) < 2 {
			apiErrHandler(ErrNoIDInEventHandler, &w)
			return
		}
		id, err := strconv.Atoi(pathParts[1])
		if err != nil {
			apiErrHandler(err, &w)
			return
		}
		fmt.Println("ID: ", id)
		event, errInner := s.app.GetEvent(ctx, id)
		if errInner != nil {
			newMessage.Text = errInner.Error()
			newMessage.Code = 1
		} else {
			newMessage.Text = "OK!"
			newMessage.Code = 0
		}
		EvAnswer.Events[0] = event
		EvAnswer.Message = newMessage
		jsonstring, err := json.Marshal(EvAnswer)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		_, err = w.Write(jsonstring)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		return

	case http.MethodPost:

		fmt.Println("Post")

		newEvent := EventRawData{}
		newMessage := outputJSON{}
		tflayout := "2006-01-02 15:04:05"

		body, err := io.ReadAll(r.Body)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		err = json.Unmarshal(body, &newEvent)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		dateStart, err := time.Parse(tflayout, newEvent.DateStart)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		dateStop, err := time.Parse(tflayout, newEvent.DateStop)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		eventMessageTimeDelta := time.Duration(newEvent.EventMessageTimeDelta) * time.Millisecond

		fmt.Println("PostItem: ", newEvent)

		id, errInner := s.app.CreateEvent(ctx, newEvent.Title, newEvent.UserID, newEvent.Description, dateStart, dateStop, eventMessageTimeDelta)
		if errInner != nil {
			newMessage.Text = errInner.Error()
			newMessage.Code = 1
		} else {
			newMessage.Text = "OK!"
			newMessage.Code = id
		}

		jsonstring, err := json.Marshal(newMessage)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		_, err = w.Write(jsonstring)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		return

	case http.MethodPut:

		fmt.Println("Put")

		newEvent := EventRawData{}
		newMessage := outputJSON{}
		tflayout := "2006-01-02 15:04:05"

		body, err := io.ReadAll(r.Body)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		err = json.Unmarshal(body, &newEvent)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		dateStart, err := time.Parse(tflayout, newEvent.DateStart)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		dateStop, err := time.Parse(tflayout, newEvent.DateStop)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		eventMessageTimeDelta := time.Duration(newEvent.EventMessageTimeDelta) * time.Millisecond

		fmt.Println("PutItem: ", newEvent)

		errInner := s.app.UpdateEvent(ctx, newEvent.ID, newEvent.Title, newEvent.UserID, newEvent.Description, dateStart, dateStop, eventMessageTimeDelta)
		if errInner != nil {
			newMessage.Text = errInner.Error()
			newMessage.Code = 1
		} else {
			newMessage.Text = "OK!"
			newMessage.Code = 0
		}

		jsonstring, err := json.Marshal(newMessage)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		_, err = w.Write(jsonstring)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		return

	case http.MethodDelete:

		fmt.Println("Delete")
		newMessage := outputJSON{}

		path := strings.Trim(r.URL.Path, "/")
		pathParts := strings.Split(path, "/")
		if len(pathParts) < 2 {
			apiErrHandler(ErrNoIDInEventHandler, &w)
			return
		}
		id, err := strconv.Atoi(pathParts[1])
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		errInner := s.app.DeleteEvent(ctx, id)
		if errInner != nil {
			newMessage.Text = errInner.Error()
			newMessage.Code = 1
		} else {
			newMessage.Text = "OK!"
			newMessage.Code = 0
		}

		jsonstring, err := json.Marshal(newMessage)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		_, err = w.Write(jsonstring)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		return

	default:
		apiErrHandler(ErrUnsupportedMethod, &w)
		return
	}
}

func (s *Server) GetEventsOnDayByDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		apiErrHandler(ErrUnsupportedMethod, &w)
		return
	} else {
		ctx, _ := context.WithTimeout(context.Background(), s.Config.GetDBTimeOut())
		newMessage := outputJSON{}
		tflayout := "2006-01-02 15:04:05"
		inpDate := InputDate{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		err = json.Unmarshal(body, &inpDate)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		dateStart, err := time.Parse(tflayout, inpDate.Date)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}
		fmt.Println("controlDate: ", dateStart)
		List, errInner := s.app.GetListEventsonDayByDay(ctx, dateStart)
		if errInner != nil {
			newMessage.Text = errInner.Error()
			newMessage.Code = 1
		} else {
			newMessage.Text = "OK!"
			newMessage.Code = 0
		}

		EvAnswer := EventAnswer{}
		EvAnswer.Events = List
		EvAnswer.Message = newMessage
		jsonstring, err := json.Marshal(EvAnswer)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		_, err = w.Write(jsonstring)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		return
	}
}

func (s *Server) GetEventsOnWeekByDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		apiErrHandler(ErrUnsupportedMethod, &w)
		return
	} else {
		ctx, _ := context.WithTimeout(context.Background(), s.Config.GetDBTimeOut())
		newMessage := outputJSON{}
		tflayout := "2006-01-02 15:04:05"
		inpDate := InputDate{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		err = json.Unmarshal(body, &inpDate)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		dateStart, err := time.Parse(tflayout, inpDate.Date)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		List, errInner := s.app.GetListEventsOnWeekByDay(ctx, dateStart)
		if errInner != nil {
			newMessage.Text = errInner.Error()
			newMessage.Code = 1
		} else {
			newMessage.Text = "OK!"
			newMessage.Code = 0
		}

		EvAnswer := EventAnswer{}
		EvAnswer.Events = List
		EvAnswer.Message = newMessage
		jsonstring, err := json.Marshal(EvAnswer)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		_, err = w.Write(jsonstring)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		return
	}
}

func (s *Server) GetEventsOnMonthByDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		apiErrHandler(ErrUnsupportedMethod, &w)
		return
	} else {
		ctx, _ := context.WithTimeout(context.Background(), s.Config.GetDBTimeOut())
		newMessage := outputJSON{}
		tflayout := "2006-01-02 15:04:05"
		inpDate := InputDate{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		err = json.Unmarshal(body, &inpDate)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		dateStart, err := time.Parse(tflayout, inpDate.Date)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		List, errInner := s.app.GetListEventsOnMonthByDay(ctx, dateStart)
		if errInner != nil {
			newMessage.Text = errInner.Error()
			newMessage.Code = 1
		} else {
			newMessage.Text = "OK!"
			newMessage.Code = 0
		}

		EvAnswer := EventAnswer{}
		EvAnswer.Events = List
		EvAnswer.Message = newMessage
		jsonstring, err := json.Marshal(EvAnswer)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		_, err = w.Write(jsonstring)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		return
	}
}

func (s *Server) GetListEventsNotificationByDay(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		apiErrHandler(ErrUnsupportedMethod, &w)
		return
	} else {
		ctx, cancel := context.WithTimeout(r.Context(), s.Config.GetDBTimeOut())
		defer cancel()
		newMessage := outputJSON{}
		tflayout := "2006-01-02 15:04:05"
		inpDate := InputDate{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		err = json.Unmarshal(body, &inpDate)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		date, err := time.Parse(tflayout, inpDate.Date)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		List, errInner := s.app.GetListEventsNotificationByDay(ctx, date)
		if errInner != nil {
			newMessage.Text = errInner.Error()
			newMessage.Code = 1
		} else {
			newMessage.Text = "OK!"
			newMessage.Code = 0
		}

		EvAnswer := EventAnswer{}
		EvAnswer.Events = List
		EvAnswer.Message = newMessage
		jsonstring, err := json.Marshal(EvAnswer)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		_, err = w.Write(jsonstring)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		return
	}
}

func (s *Server) DeleteOldEventsByDay(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		apiErrHandler(ErrUnsupportedMethod, &w)
		return
	} else {
		ctx, cancel := context.WithTimeout(r.Context(), s.Config.GetDBTimeOut())
		defer cancel()
		newMessage := outputJSON{}
		tflayout := "2006-01-02 15:04:05"
		inpDate := InputDate{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		err = json.Unmarshal(body, &inpDate)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		date, err := time.Parse(tflayout, inpDate.Date)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		count, errInner := s.app.DeleteOldEventsByDay(ctx, date)
		if errInner != nil {
			newMessage.Text = errInner.Error()
			newMessage.Code = 1
		} else {
			newMessage.Text = "OK!"
			newMessage.Code = count
		}

		jsonstring, err := json.Marshal(newMessage)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		_, err = w.Write(jsonstring)
		if err != nil {
			apiErrHandler(err, &w)
			return
		}

		return
	}
}