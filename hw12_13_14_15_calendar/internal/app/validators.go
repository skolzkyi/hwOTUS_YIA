package app

import (
	"errors"
	"time"

	storage "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/storage/event"
)

var ErrVoidTitle = errors.New("title is void")
var ErrVoidUserID = errors.New("userID is void")
var ErrVoidDateStart = errors.New("dateStart is void")
var ErrVoidDateStop = errors.New("dateStop is void")
var ErrTitleTooLong = errors.New("title too long")
var ErrUserIDTooLong = errors.New("userID too long")
var ErrDescTooLong = errors.New("description too long")
var ErrEndDateBefstartDate = errors.New("endDate before startDate or equal")

func SimpleEventValidator(title string, userID string, description string, dateStart time.Time, dateStop time.Time, eventMessageTimeDelta time.Duration) (storage.Event, error) {
	event := storage.Event{ID: 0, Title: title, UserID: userID, Description: description, DateStart: dateStart, DateStop: dateStop, EventMessageTimeDelta: eventMessageTimeDelta}
	switch {
	case event.Title == "":
		return storage.Event{}, ErrVoidTitle
	case len(event.Title) > 254:
		return storage.Event{}, ErrTitleTooLong
	case event.UserID == "":
		return storage.Event{}, ErrVoidUserID
	case len(event.UserID) > 49:
		return storage.Event{}, ErrUserIDTooLong
	case len(event.Description) > 1499:
		return storage.Event{}, ErrDescTooLong
	case event.DateStart.IsZero():
		return storage.Event{}, ErrVoidDateStart
	case event.DateStop.IsZero():
		return storage.Event{}, ErrVoidDateStop
	case event.DateStop.Before(event.DateStart) || event.DateStop.Equal(event.DateStart):
		return storage.Event{}, ErrEndDateBefstartDate
	default:
	}

	return event, nil
}
