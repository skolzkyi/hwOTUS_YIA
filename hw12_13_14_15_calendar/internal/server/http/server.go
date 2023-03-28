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
	InitStorage(ctx context.Context, config storage.Config) error
	CloseStorage(ctx context.Context) error
	GetEvent(ctx context.Context, id int) (storage.Event, error)
	CreateEvent(ctx context.Context, title string, userID string, description string, dateStart time.Time, dateStop time.Time, eventMessageTimeDelta time.Duration) error
	UpdateEvent(ctx context.Context, id int, title string, userID string, description string, dateStart time.Time, dateStop time.Time, eventMessageTimeDelta time.Duration) error
	DeleteEvent(ctx context.Context, id int) error
	GetListEventsonDayByDay(ctx context.Context, day time.Time) ([]storage.Event, error)
	GetListEventsOnWeekByDay(ctx context.Context, day time.Time) ([]storage.Event, error)
	GetListEventsOnMonthByDay(ctx context.Context, day time.Time) ([]storage.Event, error)
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
	s.logg.Info("calendar is running...")
	//============
	std := time.Now()
	stopd := std.Add(72 * time.Hour)
	emtd := 4 * time.Hour
	s.app.CreateEvent(context.Background(), "test0 - base event", "USER0", "", std, stopd, emtd)
	s.app.CreateEvent(context.Background(), "test1 - +5days", "USER0", "", std.Add(120*time.Hour), stopd.Add(120*time.Hour), emtd)
	s.app.CreateEvent(context.Background(), "test2 - +6 days end date after week", "USER0", "", std.Add(144*time.Hour), stopd.Add(144*time.Hour), emtd)
	s.app.CreateEvent(context.Background(), "test3 - +8 days - next week", "USER0", "", std.Add(192*time.Hour), stopd.Add(192*time.Hour), emtd)
	s.app.CreateEvent(context.Background(), "test4 - start in before week and end in cur week", "USER0", "", std.Add(-48*time.Hour), std.Add(-5*time.Hour), emtd)
	s.app.CreateEvent(context.Background(), "test5 - in this day", "USER0", "", std.Add(-4*time.Hour), std.Add(-3*time.Hour), emtd)
	testEvents, err := s.app.GetListEventsOnMonthByDay(context.Background(), std)
	if len(testEvents) > 0 {
		id := testEvents[0].ID
		s.app.UpdateEvent(context.Background(), id, "test777 - updated event", "USER0", "", std, stopd, emtd)
	}

	//============
	err = s.serv.ListenAndServe()
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
