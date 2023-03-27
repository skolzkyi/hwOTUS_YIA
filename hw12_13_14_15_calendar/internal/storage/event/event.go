package storage

import (
	"errors"
	"time"

	helpers "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/helpers"
)

var ErrNoRecord = errors.New("запись не найдена")
var ErrStorageTimeout = errors.New("таймаут обращения к хранилищу")

type Event struct {
	ID                    string
	Title                 string
	UserID                string
	Description           string
	DateStart             time.Time
	DateStop              time.Time
	EventMessageTimeDelta time.Duration
}

func (e *Event) String() string {
	res := helpers.StringBuild("[", "ID: ", e.ID, ", Title: ", e.Title, " UserID: ", e.UserID, " Description: ", e.Description, " DateStart: ", e.DateStart.String(), " DateStop: ", e.DateStop.String(), " EventMessageTimeDelta: ", e.EventMessageTimeDelta.String(), "]")
	return res
}
