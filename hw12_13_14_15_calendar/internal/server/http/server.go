package internalhttp

import (
	"context"
	"errors"
	"net/http"
	"time"

	// pb "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/server/grpc/pb"
	storage "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/storage/event"
	"go.uber.org/zap"
)

type Server struct {
	serv   *http.Server
	logg   Logger
	app    Application
	Config Config
}

type Config interface {
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
	GetGRPCPort() string
}

type Logger interface {
	Info(msg string)
	Warning(msg string)
	Error(msg string)
	Fatal(msg string)
	GetZapLogger() *zap.SugaredLogger
}

type Application interface {
	InitStorage(ctx context.Context, config storage.Config) error
	CloseStorage(ctx context.Context) error
	GetEvent(ctx context.Context, id int) (storage.Event, error)
	CreateEvent(ctx context.Context, title string, userID string, description string, dateStart time.Time, dateStop time.Time, eventMessageTimeDelta time.Duration) (int, error)
	UpdateEvent(ctx context.Context, id int, title string, userID string, description string, dateStart time.Time, dateStop time.Time, eventMessageTimeDelta time.Duration) error
	DeleteEvent(ctx context.Context, id int) error
	GetListEventsonDayByDay(ctx context.Context, day time.Time) ([]storage.Event, error)
	GetListEventsOnWeekByDay(ctx context.Context, day time.Time) ([]storage.Event, error)
	GetListEventsOnMonthByDay(ctx context.Context, day time.Time) ([]storage.Event, error)
	GetListEventsNotificationByDay(ctx context.Context, dateTime time.Time) ([]storage.Event, error)
	DeleteOldEventsByDay(ctx context.Context, dateTime time.Time) (int, error)
}

func NewServer(logger Logger, app Application, config Config) *Server {
	server := Server{}
	server.logg = logger
	server.app = app
	server.Config = config
	server.serv = &http.Server{
		Addr:    config.GetServerURL(),
		Handler: server.routes(),
	}

	return &server
}

func (s *Server) Start(ctx context.Context) error {
	s.logg.Info("calendar is running...")
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
	err = s.app.CloseStorage(ctx)
	if err != nil {
		s.logg.Error("server closeStorage error: " + err.Error())
		return err
	}
	s.logg.Info("server graceful shutdown")
	return err
}
