package storage

import (
	"errors"
	"strconv"
	"time"

	helpers "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/helpers"
)

var (
	ErrNoRecord       = errors.New("record not searched")
	ErrStorageTimeout = errors.New("storage timeout")
	ErrDateBusy       = errors.New("this date busy by other event")
	ErrIDNotUnique    = errors.New("event with this id exited")
)

type Config interface {
	Init(path string) error
	GetServerURL() string
	GetAddress() string
	GetPort() string
	GetOSFilePathSeparator() string
	GetServerShutdownTimeout() time.Duration
	GetDBName() string
	GetDBUser() string
	GetDBPassword() string
	GetDBConnMaxLifetime() time.Duration
	GetDBMaxOpenConns() int
	GetDBMaxIdleConns() int
	GetDBTimeOut() time.Duration
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
	res := helpers.StringBuild("[ID: ", strconv.Itoa(e.ID), ", Title: ", e.Title, ", UserID: ", e.UserID, ", Description: ", e.Description, ", DateStart: ", e.DateStart.String(), ", DateStop: ", e.DateStop.String(), ", EventMessageTimeDelta: ", e.EventMessageTimeDelta.String(), "]") //nolint:lll
	return res
}
