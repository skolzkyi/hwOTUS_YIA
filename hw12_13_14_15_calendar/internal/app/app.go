package app

import (
	"context"
	"strconv"
	"time"

	helpers "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/helpers"
	storage "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/storage/event"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Info(msg string)
	Warning(msg string)
	Error(msg string)
	Fatal(msg string)
}

type Storage interface {
	Init(ctx context.Context, config storage.Config) error
	Close(ctx context.Context) error
	GetEvent(ctx context.Context, id int) (storage.Event, error)
	CreateEvent(ctx context.Context, value storage.Event) (int, error)
	UpdateEvent(ctx context.Context, value storage.Event) error
	DeleteEvent(ctx context.Context, id int) error
	GetListEventsonDayByDay(ctx context.Context, day time.Time) ([]storage.Event, error)
	GetListEventsOnWeekByDay(ctx context.Context, day time.Time) ([]storage.Event, error)
	GetListEventsOnMonthByDay(ctx context.Context, day time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	app := App{
		logger:  logger,
		storage: storage,
	}
	return &app
}

func (a *App) InitStorage(ctx context.Context, config storage.Config) error {
	return a.storage.Init(ctx, config)
}

func (a *App) CloseStorage(ctx context.Context) error {
	return a.storage.Close(ctx)
}

func (a *App) GetEvent(ctx context.Context, id int) (storage.Event, error) {
	event, err := a.storage.GetEvent(ctx, id)
	return event, err
}

func (a *App) CreateEvent(ctx context.Context, title string, userID string, description string, dateStart time.Time, dateStop time.Time, eventMessageTimeDelta time.Duration) (int, error) {
	//event := storage.Event{ID: id, Title: title, UserID: userID, Description: description, DateStart: dateStart, DateStop: dateStop, EventMessageTimeDelta: eventMessageTimeDelta}
	event, err := SimpleEventValidator(title, userID, description, dateStart, dateStop, eventMessageTimeDelta)
	if err != nil {
		message := helpers.StringBuild("event create error(title - ", title, "),error: ", err.Error())
		a.logger.Error(message)
	}
	id, err := a.storage.CreateEvent(ctx, event)
	if err != nil {
		message := helpers.StringBuild("event create error(title - ", title, "),error: ", err.Error())
		a.logger.Error(message)
	} else {
		message := helpers.StringBuild("new event created(title - ", title, ")")
		a.logger.Info(message)
	}

	return id, err
}

func (a *App) UpdateEvent(ctx context.Context, id int, title string, userID string, description string, dateStart time.Time, dateStop time.Time, eventMessageTimeDelta time.Duration) error {
	event, err := SimpleEventValidator(title, userID, description, dateStart, dateStop, eventMessageTimeDelta)
	if err != nil {
		message := helpers.StringBuild("event create error(title - ", title, "),error: ", err.Error())
		a.logger.Error(message)
	}
	err = a.storage.UpdateEvent(ctx, event)
	if err != nil {
		message := helpers.StringBuild("event update error(title - ", title, "),error: ", err.Error())
		a.logger.Error(message)
	} else {
		message := helpers.StringBuild("event updated(title - ", title, ")")
		a.logger.Info(message)
	}

	return err
}

func (a *App) DeleteEvent(ctx context.Context, id int) error {
	err := a.storage.DeleteEvent(ctx, id)
	if err != nil {
		message := helpers.StringBuild("event delete error(id - ", strconv.Itoa(id), "),error: ", err.Error())
		a.logger.Error(message)
	} else {
		message := helpers.StringBuild("event deleted(id - ", strconv.Itoa(id), ")")
		a.logger.Info(message)
	}

	return err
}

func (a *App) GetListEventsonDayByDay(ctx context.Context, day time.Time) ([]storage.Event, error) {
	eventsList, err := a.storage.GetListEventsonDayByDay(ctx, day)
	return eventsList, err
}

func (a *App) GetListEventsOnWeekByDay(ctx context.Context, day time.Time) ([]storage.Event, error) {
	eventsList, err := a.storage.GetListEventsOnWeekByDay(ctx, day)
	return eventsList, err
}

func (a *App) GetListEventsOnMonthByDay(ctx context.Context, day time.Time) ([]storage.Event, error) {
	eventsList, err := a.storage.GetListEventsOnMonthByDay(ctx, day)
	return eventsList, err
}
