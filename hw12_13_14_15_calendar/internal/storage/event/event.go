package storage

import (
	"errors"
	"time"
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
	EventMessageTimeDelta time.Time
}
