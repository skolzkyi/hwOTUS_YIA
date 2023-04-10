package internalhttp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/app"
	"github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/logger"
	storage "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/storage/event"
	memorystorage "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
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

func (config *ConfigTest) GetGRPSPort() string {
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

func TestCreateEvent(t *testing.T) {
	data := bytes.NewBufferString(getTestEventData(t))
	server := createServer(t)

	r := httptest.NewRequest("POST", "/Event/", data)
	w := httptest.NewRecorder()
	server.Event_REST(w, r)

	res := w.Result()
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	respExp := `{"Text":"OK!","Code":0}`
	require.Equal(t, respExp, string(respBody))
}

func TestCreateEventBadTimeBusy(t *testing.T) {
	data := bytes.NewBufferString(getTestEventData(t))
	server := createServer(t)

	r := httptest.NewRequest("POST", "/Event/", data)
	w := httptest.NewRecorder()
	server.Event_REST(w, r)

	res := w.Result()

	respBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	respExp := `{"Text":"OK!","Code":0}`
	require.Equal(t, respExp, string(respBody))
	res.Body.Close()

	data = bytes.NewBufferString(getTestEventData(t))
	r = httptest.NewRequest("POST", "/Event/", data)
	w = httptest.NewRecorder()
	server.Event_REST(w, r)

	res = w.Result()
	defer res.Body.Close()
	respBody, err = io.ReadAll(res.Body)
	require.NoError(t, err)

	respExp = `{"Text":"this date busy by other event","Code":1}`
	require.Equal(t, respExp, string(respBody))
}

func TestGetEvent(t *testing.T) {
	server := createServer(t)
	ctx := context.Background()
	std := time.Date(2023, 4, 20, 0, 0, 0, 1, time.Local)
	emtd := 4 * time.Hour
	id, err := server.app.CreateEvent(ctx, "testData", "USER0", "", std, std.Add(4*time.Hour), emtd)
	if err != nil {
		t.Fatal()
	}
	r := httptest.NewRequest("GET", "/Event/"+strconv.Itoa(id), nil)
	w := httptest.NewRecorder()
	server.Event_REST(w, r)

	res := w.Result()
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	respExp := `{"Events":[{"ID":0,"Title":"testData","UserID":"USER0","Description":"","DateStart":"2023-04-20T00:00:00.000000001Z","DateStop":"2023-04-20T04:00:00.000000001Z","EventMessageTimeDelta":14400000000000}],"Message":{"Text":"OK!","Code":0}}` // nolint:lll
	require.Equal(t, respExp, string(respBody))
}

func TestGetEventBadID(t *testing.T) {
	server := createServer(t)

	r := httptest.NewRequest("GET", "/Event/99", nil)
	w := httptest.NewRecorder()
	server.Event_REST(w, r)

	res := w.Result()
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	respExp := `{"Events":[{"ID":0,"Title":"","UserID":"","Description":"","DateStart":"0001-01-01T00:00:00Z","DateStop":"0001-01-01T00:00:00Z","EventMessageTimeDelta":0}],"Message":{"Text":"record not searched","Code":1}}` // nolint:lll
	require.Equal(t, respExp, string(respBody))
}

func TestUpdateEvent(t *testing.T) {
	server := createServer(t)
	ctx := context.Background()
	std := time.Date(2023, 4, 20, 0, 0, 0, 1, time.Local)
	emtd := 4 * time.Hour
	id, err := server.app.CreateEvent(ctx, "testData - not Updated", "USER0", "", std, std.Add(4*time.Hour), emtd)
	if err != nil {
		t.Fatal()
	}
	data := bytes.NewBufferString(getTestEventData(t))
	r := httptest.NewRequest("PUT", "/Event/"+strconv.Itoa(id), data)
	w := httptest.NewRecorder()
	server.Event_REST(w, r)

	res := w.Result()
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	respExp := `{"Text":"OK!","Code":0}`
	require.Equal(t, respExp, string(respBody))
}

func TestUpdateEventBadID(t *testing.T) {
	server := createServer(t)

	data := bytes.NewBufferString(getTestEventData(t))
	r := httptest.NewRequest("PUT", "/Event/99", data)
	w := httptest.NewRecorder()
	server.Event_REST(w, r)

	res := w.Result()
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	respExp := `{"Text":"record not searched","Code":1}`
	require.Equal(t, respExp, string(respBody))
}

func TestDeleteEvent(t *testing.T) {
	server := createServer(t)
	ctx := context.Background()
	std := time.Date(2023, 4, 20, 0, 0, 0, 1, time.Local)
	emtd := 4 * time.Hour
	id, err := server.app.CreateEvent(ctx, "testData", "USER0", "", std, std.Add(4*time.Hour), emtd)
	if err != nil {
		t.Fatal()
	}

	r := httptest.NewRequest("DELETE", "/Event/"+strconv.Itoa(id), nil)
	w := httptest.NewRecorder()
	server.Event_REST(w, r)

	res := w.Result()
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	respExp := `{"Text":"OK!","Code":0}`
	require.Equal(t, respExp, string(respBody))

	_, err = server.app.GetEvent(ctx, id)
	require.Truef(t, errors.Is(err, storage.ErrNoRecord), "actual error %q", err)
}

func TestDeleteEventBadID(t *testing.T) {
	server := createServer(t)

	r := httptest.NewRequest("DELETE", "/Event/99", nil)
	w := httptest.NewRecorder()
	server.Event_REST(w, r)

	res := w.Result()
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	respExp := `{"Text":"record not searched","Code":1}`
	require.Equal(t, respExp, string(respBody))
}

func TestGetEventsOnDayByDay(t *testing.T) {
	server := createServer(t)
	createTestEventPool(t, server)
	r := httptest.NewRequest("GET", "/GetEventsOnDayByDay/", bytes.NewBufferString(`{"Date":"2023-04-20 17:51:00"}`))
	w := httptest.NewRecorder()
	server.GetEventsOnDayByDay(w, r)

	res := w.Result()
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	respExp := `{"Events":[{"ID":4,"Title":"test4 - start in before week and end in cur week","UserID":"USER0","Description":"","DateStart":"2023-04-18T12:00:00.000000001Z","DateStop":"2023-04-20T07:00:00.000000001Z","EventMessageTimeDelta":14400000000000},{"ID":5,"Title":"test5 - in this day","UserID":"USER0","Description":"","DateStart":"2023-04-20T08:00:00.000000001Z","DateStop":"2023-04-20T09:00:00.000000001Z","EventMessageTimeDelta":14400000000000},{"ID":0,"Title":"test0 - base event","UserID":"USER0","Description":"","DateStart":"2023-04-20T12:00:00.000000001Z","DateStop":"2023-04-20T16:00:00.000000001Z","EventMessageTimeDelta":14400000000000}],"Message":{"Text":"OK!","Code":0}}` // nolint:lll
	require.Equal(t, respExp, string(respBody))
}

func TestGetEventsOnWeekByDay(t *testing.T) {
	server := createServer(t)
	createTestEventPool(t, server)
	r := httptest.NewRequest("GET", "/GetEventsOnWeekByDay/", bytes.NewBufferString(`{"Date":"2023-04-20 17:51:00"}`))
	w := httptest.NewRecorder()
	server.GetEventsOnWeekByDay(w, r)

	res := w.Result()
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	respExp := `{"Events":[{"ID":4,"Title":"test4 - start in before week and end in cur week","UserID":"USER0","Description":"","DateStart":"2023-04-18T12:00:00.000000001Z","DateStop":"2023-04-20T07:00:00.000000001Z","EventMessageTimeDelta":14400000000000},{"ID":5,"Title":"test5 - in this day","UserID":"USER0","Description":"","DateStart":"2023-04-20T08:00:00.000000001Z","DateStop":"2023-04-20T09:00:00.000000001Z","EventMessageTimeDelta":14400000000000},{"ID":0,"Title":"test0 - base event","UserID":"USER0","Description":"","DateStart":"2023-04-20T12:00:00.000000001Z","DateStop":"2023-04-20T16:00:00.000000001Z","EventMessageTimeDelta":14400000000000},{"ID":1,"Title":"test1 - +5days","UserID":"USER0","Description":"","DateStart":"2023-04-25T12:00:00.000000001Z","DateStop":"2023-04-25T16:00:00.000000001Z","EventMessageTimeDelta":14400000000000},{"ID":2,"Title":"test2 - +6 days end date after week","UserID":"USER0","Description":"","DateStart":"2023-04-26T12:00:00.000000001Z","DateStop":"2023-04-26T18:00:00.000000001Z","EventMessageTimeDelta":14400000000000}],"Message":{"Text":"OK!","Code":0}}` // nolint:lll
	require.Equal(t, respExp, string(respBody))
}

func TestGetEventsOnMonthByDay(t *testing.T) {
	server := createServer(t)
	createTestEventPool(t, server)
	r := httptest.NewRequest("GET", "/GetEventsOnMonthByDay/", bytes.NewBufferString(`{"Date":"2023-04-20 17:51:00"}`))
	w := httptest.NewRecorder()
	server.GetEventsOnMonthByDay(w, r)

	res := w.Result()
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	respExp := `{"Events":[{"ID":4,"Title":"test4 - start in before week and end in cur week","UserID":"USER0","Description":"","DateStart":"2023-04-18T12:00:00.000000001Z","DateStop":"2023-04-20T07:00:00.000000001Z","EventMessageTimeDelta":14400000000000},{"ID":5,"Title":"test5 - in this day","UserID":"USER0","Description":"","DateStart":"2023-04-20T08:00:00.000000001Z","DateStop":"2023-04-20T09:00:00.000000001Z","EventMessageTimeDelta":14400000000000},{"ID":0,"Title":"test0 - base event","UserID":"USER0","Description":"","DateStart":"2023-04-20T12:00:00.000000001Z","DateStop":"2023-04-20T16:00:00.000000001Z","EventMessageTimeDelta":14400000000000},{"ID":1,"Title":"test1 - +5days","UserID":"USER0","Description":"","DateStart":"2023-04-25T12:00:00.000000001Z","DateStop":"2023-04-25T16:00:00.000000001Z","EventMessageTimeDelta":14400000000000},{"ID":2,"Title":"test2 - +6 days end date after week","UserID":"USER0","Description":"","DateStart":"2023-04-26T12:00:00.000000001Z","DateStop":"2023-04-26T18:00:00.000000001Z","EventMessageTimeDelta":14400000000000},{"ID":3,"Title":"test3 - +8 days - next week","UserID":"USER0","Description":"","DateStart":"2023-04-28T12:00:00.000000001Z","DateStop":"2023-04-28T20:00:00.000000001Z","EventMessageTimeDelta":14400000000000}],"Message":{"Text":"OK!","Code":0}}` // nolint:lll
	require.Equal(t, respExp, string(respBody))
}

func createTestEventPool(t *testing.T, server *Server) {
	t.Helper()
	ctx := context.Background()
	std := time.Date(2023, 4, 20, 12, 0, 0, 1, time.Local)
	emtd := 4 * time.Hour
	_, err := server.app.CreateEvent(ctx, "test0 - base event", "USER0", "", std, std.Add(4*time.Hour), emtd)
	if err != nil {
		t.Fatal()
	}
	_, err = server.app.CreateEvent(ctx, "test1 - +5days", "USER0", "", std.Add(120*time.Hour), std.Add(124*time.Hour), emtd)
	if err != nil {
		t.Fatal()
	}
	_, err = server.app.CreateEvent(ctx, "test2 - +6 days end date after week", "USER0", "", std.Add(144*time.Hour), std.Add(150*time.Hour), emtd)
	if err != nil {
		t.Fatal()
	}
	_, err = server.app.CreateEvent(ctx, "test3 - +8 days - next week", "USER0", "", std.Add(192*time.Hour), std.Add(200*time.Hour), emtd)
	if err != nil {
		t.Fatal()
	}
	_, err = server.app.CreateEvent(ctx, "test4 - start in before week and end in cur week", "USER0", "", std.Add(-48*time.Hour), std.Add(-5*time.Hour), emtd)
	if err != nil {
		t.Fatal()
	}
	_, err = server.app.CreateEvent(ctx, "test5 - in this day", "USER0", "", std.Add(-4*time.Hour), std.Add(-3*time.Hour), emtd)
	if err != nil {
		t.Fatal()
	}
}

func getTestEventData(t *testing.T) string {
	t.Helper()
	return `{
		"ID":                    0,
		"Title":                 "testData",
		"UserID":                "USER0",
		"Description":           "",
		"DateStart":             "2023-04-20 15:04:05",
		"DateStop":              "2023-04-22 15:04:05",
		"EventMessageTimeDelta": 10800000
	}`
}

func createServer(t *testing.T) *Server {
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

	server := NewServer(log, calendar, &config)

	return server
}
