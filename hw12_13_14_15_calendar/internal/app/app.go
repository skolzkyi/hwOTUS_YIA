package app

import (
	"context"
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
}

type Storage interface {
	Init(ctx context.Context) error
	Close() error
	GetEvent(ctx context.Context, id string) (storage.Event, error)
	CreateEvent(ctx context.Context, value storage.Event) error
	UpdateEvent(ctx context.Context, value storage.Event) error
	DeleteEvent(ctx context.Context, id string) error
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

func (a *App) InitStorage(ctx context.Context) error {
	return a.storage.Init(ctx)
}

func (a *App) CloseStorage() error {
	return a.storage.Close()
}

func (a *App) GetEvent(ctx context.Context, id string) (storage.Event, error) {
	event, err := a.storage.GetEvent(ctx, id)
	return event, err
}

func (a *App) CreateEvent(ctx context.Context, id, title string, userID string, description string, dateStart time.Time, dateStop time.Time, eventMessageTimeDelta time.Time) error {
	err := a.storage.CreateEvent(ctx, storage.Event{ID: id, Title: title, UserID: userID, Description: description, DateStart: dateStart, DateStop: dateStop, EventMessageTimeDelta: eventMessageTimeDelta})
	if err != nil {
		message := helpers.StringBuild("создано новое событие(id - ", id, " title - ", title, ")")
		a.logger.Info(message)
	}
	return err
}

func (a *App) UpdateEvent(ctx context.Context, id, title string, userID string, description string, dateStart time.Time, dateStop time.Time, eventMessageTimeDelta time.Time) error {
	err := a.storage.UpdateEvent(ctx, storage.Event{ID: id, Title: title, UserID: userID, Description: description, DateStart: dateStart, DateStop: dateStop, EventMessageTimeDelta: eventMessageTimeDelta})
	if err != nil {
		message := helpers.StringBuild("обновлено событие(id - ", id, " title - ", title, ")")
		a.logger.Info(message)
	}
	return err
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	err := a.storage.DeleteEvent(ctx, id)
	if err != nil {
		message := helpers.StringBuild("удалено событие(id - ", id, " title - ", title, ")")
		a.logger.Info(message)
	}
	return err
}

func (a *App) GetListEventsonDayByDay(ctx context.Context, day time.Time) error {
	eventsList, err := a.storage.GetListEventsonDayByDay(ctx, day)
	return eventsList, err
}

func (a *App) GetListEventsOnWeekByDay(ctx context.Context, day time.Time) error {
	eventsList, err := a.storage.GetListEventsOnWeekByDay(ctx, day)
	return eventsList, err
}

func (a *App) GetListEventsOnMonthByDay(ctx context.Context, day time.Time) error {
	eventsList, err := a.storage.GetListEventsOnMonthByDay(ctx, day)
	return eventsList, err
}
