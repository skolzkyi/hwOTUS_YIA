package grpcserver

import (
	"context"
	"net"
	"time"
	"sync"

	pb "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/server/grpc/pb"
	storage "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/storage/event"
	grpc "google.golang.org/grpc"
)

type GRPCServer struct {
	pb.CalendarServer
	grpcserv *grpc.Server
	logg     Logger
	app      Application
	Config   Config
	mu        sync.RWMutex
	errMap   map[error]int32
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
	GetDBAddress() string
	GetDBPort() string
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
	CreateEvent(ctx context.Context, title string, userID string, description string, dateStart time.Time, dateStop time.Time, eventMessageTimeDelta time.Duration) (int, error) //nolint:lll
	UpdateEvent(ctx context.Context, id int, title string, userID string, description string, dateStart time.Time, dateStop time.Time, eventMessageTimeDelta time.Duration) error //nolint:lll
	DeleteEvent(ctx context.Context, id int) error
	GetListEventsonDayByDay(ctx context.Context, day time.Time) ([]storage.Event, error)
	GetListEventsOnWeekByDay(ctx context.Context, day time.Time) ([]storage.Event, error)
	GetListEventsOnMonthByDay(ctx context.Context, day time.Time) ([]storage.Event, error)
	GetListEventsNotificationByDay(ctx context.Context, dateTime time.Time) ([]storage.Event, error)
	DeleteOldEventsByDay(ctx context.Context, dateTime time.Time) (int, error)
	MarkEventNotifSended(ctx context.Context, id int) error
}

func NewServer(logger Logger, app Application, config Config) *GRPCServer {
	server := GRPCServer{}
	server.logg = logger
	server.app = app
	server.Config = config
	server.grpcserv = grpc.NewServer()
	pb.RegisterCalendarServer(server.grpcserv, &server)
	server.errMap = initErrGRPCMap()
	return &server
}

func (g *GRPCServer) Start() error {
	l, err := net.Listen("tcp", net.JoinHostPort(g.Config.GetAddress(), g.Config.GetGRPCPort()))
	if err != nil {
		return err
	}
	err = g.grpcserv.Serve(l)
	if err == nil {
		g.logg.Info("GRPCserver run on: " + g.Config.GetAddress() + ":" + g.Config.GetGRPCPort())
	}

	return err
}

func (g *GRPCServer) Stop(ctx context.Context) error {
	g.grpcserv.GracefulStop()
	err := g.app.CloseStorage(ctx)
	if err != nil {
		g.logg.Error("GRPCserver closeStorage error: " + err.Error())
		return err
	}
	g.logg.Info("GRPCserver graceful shutdown")
	return err
}

func (g *GRPCServer) getGRPCErrorCode(err error) int32 {
	g.mu.RLock()
	defer g.mu.RUnlock()
	errCode,ok:=g.errMap[err]
	if !ok {
		errCode = 500 // UNKNOWN
	}
	return errCode
}

func initErrGRPCMap() map[error]int32 {
	errMap:= make(map[error]int32)
	errMap[storage.ErrNoRecord]= 404 // NOT_FOUND
	errMap[storage.ErrStorageTimeout]= 504 // DEADLINE_EXCEEDED
	errMap[storage.ErrDateBusy]= 500 // INTERNAL
	errMap[storage.ErrIDNotUnique]= 400 // INVALID_ARGUMENT
	return errMap
}
