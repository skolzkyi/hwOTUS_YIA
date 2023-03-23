package internalhttp

import (
	"context"
	"errors"
	"net/http"
	"time"

	storage "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/storage/event"
)

type Server struct {
	serv *http.Server
	logg Logger
	app  Application
}

type Config interface {
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
}

type Logger interface {
	Info(msg string)
	Warning(msg string)
	Error(msg string)
}

type Application interface {
	InitStorage(ctx context.Context) error
	CloseStorage() error
	GetEvent(ctx context.Context, id string) (storage.Event, error)
	CreateEvent(ctx context.Context, id, title string, userID string, description string, dateStart time.Time, dateStop time.Time, eventMessageTimeDelta time.Time) error
	UpdateEvent(ctx context.Context, id, title string, userID string, description string, dateStart time.Time, dateStop time.Time, eventMessageTimeDelta time.Time) error
	DeleteEvent(ctx context.Context, id string) error
	GetListEventsonDayByDay(ctx context.Context, day time.Time) error
	GetListEventsOnWeekByDay(ctx context.Context, day time.Time) error
	GetListEventsOnMonthByDay(ctx context.Context, day time.Time) error
}

func NewServer(logger Logger, app Application, config Config) *Server {
	server := Server{}
	server.logg = logger
	server.app = app
	server.serv = &http.Server{
		Addr:    config.GetServerURL(),
		Handler: server.routes(),
	}

	return &server
}

func (s *Server) Start(ctx context.Context) error {
	err := s.serv.ListenAndServe()
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			s.logg.Error("server start error: " + err.Error())
			return err
		}
	}
	<-ctx.Done()
	return err
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.serv.Shutdown(ctx)
	if err != nil {
		s.logg.Error("server shutdown error: " + err.Error())
		return err
	}
	err = s.app.CloseStorage()
	if err != nil {
		s.logg.Error("server closeStorage error: " + err.Error())
		return err
	}
	return err
}

// TODO
