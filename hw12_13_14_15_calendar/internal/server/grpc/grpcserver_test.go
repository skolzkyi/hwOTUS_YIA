//go:build !integration
// +build !integration

package grpcserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/app"
	"github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/logger"
	pb "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/server/grpc/pb"
	storage "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/storage/event"
	memorystorage "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/grpc/credentials/insecure"
)

type ConfigTest struct{}

func (config *ConfigTest) Init(_ string) error {
	return nil
}

func (config *ConfigTest) GetServerURL() string {
	return "127.0.0.1:4000"
}

func (config *ConfigTest) GetAddress() string {
	return "127.0.0.1"
}

func (config *ConfigTest) GetPort() string {
	return "4000"
}

func (config *ConfigTest) GetGRPCPort() string {
	return "5000"
}

func (config *ConfigTest) GetOSFilePathSeparator() string {
	return string(os.PathSeparator)
}

func (config *ConfigTest) GetServerShutdownTimeout() time.Duration {
	return 5 * time.Second
}

func (config *ConfigTest) GetDBName() string {
	return "OTUSFinalLab"
}

func (config *ConfigTest) GetDBUser() string {
	return "imapp"
}

func (config *ConfigTest) GetDBPassword() string {
	return "LightInDark"
}

func (config *ConfigTest) GetDBConnMaxLifetime() time.Duration {
	return 5 * time.Second
}

func (config *ConfigTest) GetDBMaxOpenConns() int {
	return 20
}

func (config *ConfigTest) GetDBMaxIdleConns() int {
	return 20
}

func (config *ConfigTest) GetDBTimeOut() time.Duration {
	return 5 * time.Second
}

func (config *ConfigTest) GetDBAddress() string {
	return "127.0.0.1"
}

func (config *ConfigTest) GetDBPort() string {
	return "3306"
}

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func createGRPCserver(t *testing.T) *GRPCServer {
	t.Helper()
	config := ConfigTest{}

	fmt.Println("config: ", config)
	log, err := logger.New("debug")
	require.NoError(t, err)
	if err != nil {
		t.Fatal()
	}

	storage := memorystorage.New()
	err = storage.Init(context.Background(), log, &config)
	require.NoError(t, err)
	if err != nil {
		t.Fatal()
	}
	calendar := app.New(log, storage)

	server := GRPCServer{}
	server.logg = log
	server.app = calendar
	server.Config = &config
	server.grpcserv = grpc.NewServer()
	lis = bufconn.Listen(bufSize)
	pb.RegisterCalendarServer(server.grpcserv, &server)
	return &server
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestGetEvent(t *testing.T) {
	server := createGRPCserver(t)
	createTestEventPool(t, server)
	go func(t *testing.T) { //nolint:staticcheck,thelper
		if err := server.grpcserv.Serve(lis); err != nil {
			t.Fatal(err) //nolint: govet
		}
	}(t)
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials())) //nolint:lll,nolintlint
	require.NoError(t, err)
	defer conn.Close()
	client := pb.NewCalendarClient(conn)
	controlEvent := storage.Event{
		ID:                    0,
		Title:                 "test0 - base event",
		UserID:                "USER0",
		Description:           "",
		DateStart:             time.Date(2023, 4, 20, 12, 0, 0, 1, time.UTC),
		DateStop:              time.Date(2023, 4, 20, 16, 0, 0, 1, time.UTC),
		EventMessageTimeDelta: 4 * time.Hour,
	}

	event, err := client.GetEvent(ctx, &pb.GetEventRequest{Id: 0})
	require.NoError(t, err)
	convEvent := storage.Event{
		ID:                    int(event.GetEvent().Id),
		Title:                 event.GetEvent().Title,
		UserID:                event.GetEvent().UserID,
		Description:           event.GetEvent().Description,
		DateStart:             event.GetEvent().DateStart.AsTime(),
		DateStop:              event.GetEvent().DateStop.AsTime(),
		EventMessageTimeDelta: event.GetEvent().GetEventMessageTimeDelta().AsDuration(),
	}
	require.Equal(t, controlEvent, convEvent)
	server.grpcserv.GracefulStop()
}

func TestCreateEvent(t *testing.T) {
	server := createGRPCserver(t)
	createTestEventPool(t, server)
	go func(t *testing.T) { //nolint:staticcheck,thelper
		if err := server.grpcserv.Serve(lis); err != nil {
			t.Fatal(err) //nolint: govet
		}
	}(t)
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials())) //nolint:lll,nolintlint
	require.NoError(t, err)
	defer conn.Close()
	client := pb.NewCalendarClient(conn)
	pbCrREvent := pb.CreateEventRequest{}
	pbEvent := pb.Event{
		Id:                    0,
		Title:                 "from test proto",
		Description:           "",
		UserID:                "USER0",
		DateStart:             timestamppb.New(time.Date(2023, 4, 10, 12, 0, 0, 1, time.UTC)),
		DateStop:              timestamppb.New(time.Date(2023, 4, 10, 16, 0, 0, 1, time.UTC)),
		EventMessageTimeDelta: durationpb.New(4 * time.Hour),
	}
	pbCrREvent.Event = &pbEvent
	pbID, err := client.CreateEvent(ctx, &pbCrREvent)
	require.NoError(t, err)
	id := pbID.GetId()

	require.Equal(t, int(id), 6)
	newEvent, err := server.app.GetEvent(ctx, int(id))
	require.NoError(t, err)
	require.Equal(t, newEvent.Title, "from test proto")
	server.grpcserv.GracefulStop()
}

func TestUpdateEvent(t *testing.T) {
	server := createGRPCserver(t)
	createTestEventPool(t, server)
	go func(t *testing.T) { //nolint:staticcheck,thelper
		if err := server.grpcserv.Serve(lis); err != nil {
			t.Fatal(err) //nolint: govet
		}
	}(t)
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials())) //nolint:lll,nolintlint
	require.NoError(t, err)
	defer conn.Close()
	client := pb.NewCalendarClient(conn)
	pbUpEvent := pb.UpdateEventRequest{}
	pbEvent := pb.Event{
		Id:                    0,
		Title:                 "updated from test proto",
		Description:           "",
		UserID:                "USER0",
		DateStart:             timestamppb.New(time.Date(2023, 4, 10, 12, 0, 0, 1, time.UTC)),
		DateStop:              timestamppb.New(time.Date(2023, 4, 10, 16, 0, 0, 1, time.UTC)),
		EventMessageTimeDelta: durationpb.New(4 * time.Hour),
	}

	pbUpEvent.Event = &pbEvent
	pbID, err := client.UpdateEvent(ctx, &pbUpEvent)
	require.NoError(t, err)
	errstr := pbID.GetError()

	require.Equal(t, errstr, "OK!")
	UpdEvent, err := server.app.GetEvent(ctx, 0)
	require.NoError(t, err)
	require.Equal(t, UpdEvent.Title, "updated from test proto")
	server.grpcserv.GracefulStop()
}

func TestDeleteEvent(t *testing.T) {
	server := createGRPCserver(t)
	createTestEventPool(t, server)
	go func(t *testing.T) { //nolint:staticcheck,thelper
		if err := server.grpcserv.Serve(lis); err != nil {
			t.Fatal(err) //nolint: govet
		}
	}(t)
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials())) //nolint:lll,nolintlint
	require.NoError(t, err)
	defer conn.Close()
	client := pb.NewCalendarClient(conn)
	pbDelResp, err := client.DeleteEvent(ctx, &pb.DeleteEventRequest{Id: 0})
	require.NoError(t, err)
	errstr := pbDelResp.GetError()

	require.Equal(t, errstr, "OK!")
	_, err = server.app.GetEvent(ctx, 0)
	require.Truef(t, errors.Is(err, storage.ErrNoRecord), "actual error %q", err)
	server.grpcserv.GracefulStop()
}

func TestGetEventsOnDayByDay(t *testing.T) {
	server := createGRPCserver(t)
	createTestEventPool(t, server)
	go func(t *testing.T) { //nolint:staticcheck,thelper
		if err := server.grpcserv.Serve(lis); err != nil {
			t.Fatal(err) //nolint: govet
		}
	}(t)
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials())) //nolint:lll,nolintlint
	require.NoError(t, err)
	defer conn.Close()
	client := pb.NewCalendarClient(conn)

	pbResp, err := client.GetEventsOnDayByDay(ctx, &pb.GetEventsOnDayRequest{Date: timestamppb.New(time.Date(2023, 4, 20, 12, 0, 0, 1, time.UTC))}) //nolint:lll,nolintlint
	require.NoError(t, err)
	resID := make(map[int32]struct{})
	for _, curEvent := range pbResp.GetEvents() {
		resID[curEvent.GetId()] = struct{}{}
	}
	_, ok := resID[0]
	require.Equal(t, ok, true)
	_, ok = resID[4]
	require.Equal(t, ok, true)
	_, ok = resID[5]
	require.Equal(t, ok, true)
	require.Equal(t, len(resID), 3)
	server.grpcserv.GracefulStop()
}

func TestGetListEventsOnWeekByDay(t *testing.T) {
	server := createGRPCserver(t)
	createTestEventPool(t, server)
	go func(t *testing.T) { //nolint:staticcheck,thelper
		if err := server.grpcserv.Serve(lis); err != nil {
			t.Fatal(err) //nolint: govet
		}
	}(t)
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials())) //nolint:lll,nolintlint
	require.NoError(t, err)
	defer conn.Close()
	client := pb.NewCalendarClient(conn)

	pbResp, err := client.GetEventsOnWeekByDay(ctx, &pb.GetEventsOnDayRequest{Date: timestamppb.New(time.Date(2023, 4, 20, 12, 0, 0, 1, time.UTC))}) //nolint:lll,nolintlint
	require.NoError(t, err)
	resID := make(map[int32]struct{})
	for _, curEvent := range pbResp.GetEvents() {
		resID[curEvent.GetId()] = struct{}{}
	}
	_, ok := resID[0]
	require.Equal(t, ok, true)
	_, ok = resID[1]
	require.Equal(t, ok, true)
	_, ok = resID[2]
	require.Equal(t, ok, true)
	_, ok = resID[4]
	require.Equal(t, ok, true)
	_, ok = resID[5]
	require.Equal(t, ok, true)
	require.Equal(t, len(resID), 5)
	server.grpcserv.GracefulStop()
}

func TestGetListEventsOnMonthByDay(t *testing.T) {
	server := createGRPCserver(t)
	createTestEventPool(t, server)
	go func(t *testing.T) { //nolint:staticcheck,thelper
		if err := server.grpcserv.Serve(lis); err != nil {
			t.Fatal(err) //nolint: govet
		}
	}(t)
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials())) //nolint:lll,nolintlint
	require.NoError(t, err)
	defer conn.Close()
	client := pb.NewCalendarClient(conn)

	pbResp, err := client.GetEventsOnMonthByDay(ctx, &pb.GetEventsOnDayRequest{Date: timestamppb.New(time.Date(2023, 4, 20, 12, 0, 0, 1, time.UTC))}) //nolint:lll,nolintlint
	require.NoError(t, err)
	resID := make(map[int32]struct{})
	for _, curEvent := range pbResp.GetEvents() {
		resID[curEvent.GetId()] = struct{}{}
	}
	_, ok := resID[0]
	require.Equal(t, ok, true)
	_, ok = resID[1]
	require.Equal(t, ok, true)
	_, ok = resID[2]
	require.Equal(t, ok, true)
	_, ok = resID[3]
	require.Equal(t, ok, true)
	_, ok = resID[4]
	require.Equal(t, ok, true)
	_, ok = resID[5]
	require.Equal(t, ok, true)
	require.Equal(t, len(resID), 6)
	server.grpcserv.GracefulStop()
}

func TestGetListEventsNotificationByDay(t *testing.T) {
	server := createGRPCserver(t)
	createTestEventPool(t, server)
	go func(t *testing.T) { //nolint:staticcheck,thelper
		if err := server.grpcserv.Serve(lis); err != nil {
			t.Fatal(err) //nolint: govet
		}
	}(t)
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials())) //nolint:lll,nolintlint
	require.NoError(t, err)
	defer conn.Close()
	client := pb.NewCalendarClient(conn)

	pbResp, err := client.GetListEventsNotificationByDay(ctx, &pb.GetEventsOnDayRequest{Date: timestamppb.New(time.Date(2023, 4, 20, 9, 0, 0, 1, time.UTC))}) //nolint:lll,nolintlint
	require.NoError(t, err)
	resID := make(map[int32]struct{})
	for _, curEvent := range pbResp.GetEvents() {
		resID[curEvent.GetId()] = struct{}{}
	}
	_, ok := resID[0]
	require.Equal(t, ok, true)
	require.Equal(t, len(resID), 1)
	server.grpcserv.GracefulStop()
}

func TestDeleteOldEventsByDay(t *testing.T) {
	server := createGRPCserver(t)
	createTestEventPool(t, server)
	go func(t *testing.T) { //nolint:staticcheck,thelper
		if err := server.grpcserv.Serve(lis); err != nil {
			t.Fatal(err) //nolint: govet
		}
	}(t)
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials())) //nolint:lll,nolintlint
	require.NoError(t, err)
	defer conn.Close()
	client := pb.NewCalendarClient(conn)

	pbResp, err := client.DeleteOldEvents(ctx, &pb.DeleteOldEventsRequest{Date: timestamppb.New(time.Date(2024, 4, 20, 9, 0, 0, 1, time.UTC))}) //nolint:lll,nolintlint
	require.NoError(t, err)
	errstr := pbResp.GetError()
	count := pbResp.GetCount()
	require.Equal(t, errstr, "OK!")
	require.Equal(t, count, int32(3))

	std := time.Date(2023, 4, 20, 0, 0, 0, 1, time.Local)
	events, err := server.app.GetListEventsOnMonthByDay(ctx, std)
	require.NoError(t, err)
	resID := make(map[int]struct{})
	for _, curEvent := range events {
		resID[curEvent.ID] = struct{}{}
	}
	_, ok := resID[1]
	require.Equal(t, ok, true)
	_, ok = resID[2]
	require.Equal(t, ok, true)
	_, ok = resID[3]
	require.Equal(t, ok, true)

	require.Equal(t, len(resID), 3)
	server.grpcserv.GracefulStop()
}

func createTestEventPool(t *testing.T, server *GRPCServer) {
	t.Helper()
	ctx := context.Background()
	std := time.Date(2023, 4, 20, 12, 0, 0, 1, time.UTC)
	emtd := 4 * time.Hour
	_, err := server.app.CreateEvent(ctx, "test0 - base event", "USER0", "", std, std.Add(4*time.Hour), emtd) //nolint:lll,nolintlint
	if err != nil {
		t.Fatal()
	}
	_, err = server.app.CreateEvent(ctx, "test1 - +5days", "USER0", "", std.Add(120*time.Hour), std.Add(124*time.Hour), emtd) //nolint:lll,nolintlint
	if err != nil {
		t.Fatal()
	}
	_, err = server.app.CreateEvent(ctx, "test2 - +6 days end date after week", "USER0", "", std.Add(144*time.Hour), std.Add(150*time.Hour), emtd) //nolint:lll,nolintlint
	if err != nil {
		t.Fatal()
	}
	_, err = server.app.CreateEvent(ctx, "test3 - +8 days - next week", "USER0", "", std.Add(192*time.Hour), std.Add(200*time.Hour), emtd) //nolint:lll,nolintlint
	if err != nil {
		t.Fatal()
	}
	_, err = server.app.CreateEvent(ctx, "test4 - start in before week and end in cur week", "USER0", "", std.Add(-48*time.Hour), std.Add(-5*time.Hour), emtd) //nolint:lll
	if err != nil {
		t.Fatal()
	}
	_, err = server.app.CreateEvent(ctx, "test5 - in this day", "USER0", "", std.Add(-4*time.Hour), std.Add(-3*time.Hour), emtd) //nolint:lll,nolintlint
	if err != nil {
		t.Fatal()
	}
}
