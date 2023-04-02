package storage

import (
	"errors"
	"strconv"
	"time"

	helpers "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/helpers"
)

var ErrNoRecord = errors.New("record not searched")
var ErrStorageTimeout = errors.New("storage timeout")
var ErrDateBusy = errors.New("this date busy by other event")
var ErrIdNotUnique = errors.New("event with this id exited")

type Config interface {
	Init(path string) error
	GetServerURL() string
	GetAddress() string
	GetPort() string
	GetOSFilePathSeparator() string
	GetServerShutdownTimeout() time.Duration
	GetDbName() string
	GetDbUser() string
	GetDbPassword() string
	GetdbConnMaxLifetime() time.Duration
	GetDbMaxOpenConns() int
	GetDbMaxIdleConns() int
	GetdbTimeOut() time.Duration
}

type Event struct {
	ID                    int
	Title                 string
	UserID                string
	Description           string
	DateStart             time.Time
	DateStop              time.Time
	EventMessageTimeDelta time.Duration
}

func (e *Event) String() string {
	message:=[]string{
		"[ID: ",
		 strconv.Itoa(e.ID),
		  ", Title: ", 
		  e.Title, 
		  ", UserID: ", 
		  e.UserID, 
		  ", Description: ", 
		  e.Description,
		  ", DateStart: ", 
		  e.DateStart.String(), 
		  ", DateStop: ", 
		  e.DateStop.String(), 
		  ", EventMessageTimeDelta: ",
		  e.EventMessageTimeDelta.String(),"]"
	}
	res := helpers.StringBuild(message)
	return res
}
