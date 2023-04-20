package app

import (
	"context"
	"strconv"
	"time"
	"go.uber.org/zap"

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
	GetZapLogger() *zap.SugaredLogger
}

type Storage interface {
	Init(ctx context.Context, logger storage.Logger, config storage.Config) error
	Close(ctx context.Context, logger storage.Logger) error
	GetEvent(ctx context.Context, logger storage.Logger, id int) (storage.Event, error)
	CreateEvent(ctx context.Context, logger storage.Logger, value storage.Event) (int, error)
	UpdateEvent(ctx context.Context, logger storage.Logger, value storage.Event) error
	DeleteEvent(ctx context.Context, logger storage.Logger, id int) error
	GetListEventsonDayByDay(ctx context.Context, logger storage.Logger, day time.Time) ([]storage.Event, error)
	GetListEventsOnWeekByDay(ctx context.Context, logger storage.Logger, day time.Time) ([]storage.Event, error)
	GetListEventsOnMonthByDay(ctx context.Context, logger storage.Logger, day time.Time) ([]storage.Event, error)
	GetListEventsNotificationByDay(ctx context.Context,logger storage.Logger, dateTime time.Time) ([]storage.Event, error)
	DeleteOldEventsByDay(ctx context.Context,logger storage.Logger, dateTime time.Time) (int, error)
}

func New(logger Logger, storage Storage) *App {
	app := App{
		logger:  logger,
		storage: storage,
	}
	return &app
}

func (a *App) InitStorage(ctx context.Context, config storage.Config) error {
	return a.storage.Init(ctx, a.logger, config)
}

func (a *App) CloseStorage(ctx context.Context) error {
	return a.storage.Close(ctx, a.logger)
}

func (a *App) GetEvent(ctx context.Context, id int) (storage.Event, error) {
	event, err := a.storage.GetEvent(ctx, a.logger, id)
	return event, err
}

func (a *App) CreateEvent(ctx context.Context, title string, userID string, description string, dateStart time.Time, dateStop time.Time, eventMessageTimeDelta time.Duration) (int, error) { //nolint:lll,whitespace

	event, err := SimpleEventValidator(title, userID, description, dateStart, dateStop, eventMessageTimeDelta)
	if err != nil {
		message := helpers.StringBuild("event create error(title - ", title, "),error: ", err.Error())
		a.logger.Error(message)
		return 0, err
	}
	id, err := a.storage.CreateEvent(ctx, a.logger, event)
	if err != nil {
		message := helpers.StringBuild("event create error(title - ", title, "),error: ", err.Error())
		a.logger.Error(message)
		return 0, err
	}
	message := helpers.StringBuild("new event created(title - ", title, ")")
	a.logger.Info(message)

	return id, nil
}

func (a *App) UpdateEvent(ctx context.Context, id int, title string, userID string, description string, dateStart time.Time, dateStop time.Time, eventMessageTimeDelta time.Duration) error { //nolint:lll,revive
	event, err := SimpleEventValidator(title, userID, description, dateStart, dateStop, eventMessageTimeDelta)
	if err != nil {
		message := helpers.StringBuild("event create error(title - ", title, "),error: ", err.Error())
		a.logger.Error(message)
		return err
	}
	err = a.storage.UpdateEvent(ctx, a.logger, event)
	if err != nil {
		message := helpers.StringBuild("event update error(title - ", title, "),error: ", err.Error())
		a.logger.Error(message)
		return err
	}
	message := helpers.StringBuild("event updated(title - ", title, ")")
	a.logger.Info(message)

	return nil
}

func (a *App) DeleteEvent(ctx context.Context, id int) error {
	err := a.storage.DeleteEvent(ctx, a.logger, id)
	if err != nil {
		message := helpers.StringBuild("event delete error(id - ", strconv.Itoa(id), "),error: ", err.Error())
		a.logger.Error(message)
		return err
	}
	message := helpers.StringBuild("event deleted(id - ", strconv.Itoa(id), ")")
	a.logger.Info(message)

	return nil
}

func (a *App) GetListEventsonDayByDay(ctx context.Context, day time.Time) ([]storage.Event, error) {
	eventsList, err := a.storage.GetListEventsonDayByDay(ctx, a.logger, day)
	return eventsList, err
}

func (a *App) GetListEventsOnWeekByDay(ctx context.Context, day time.Time) ([]storage.Event, error) {
	eventsList, err := a.storage.GetListEventsOnWeekByDay(ctx, a.logger, day)
	return eventsList, err
}

func (a *App) GetListEventsOnMonthByDay(ctx context.Context, day time.Time) ([]storage.Event, error) {
	eventsList, err := a.storage.GetListEventsOnMonthByDay(ctx, a.logger, day)
	return eventsList, err
}

func (a *App) GetListEventsNotificationByDay(ctx context.Context, dateTime time.Time) ([]storage.Event, error) {
	eventsList, err := a.storage.GetListEventsNotificationByDay(ctx, a.logger, dateTime)
	return eventsList, err
}

func (a *App) DeleteOldEventsByDay(ctx context.Context, dateTime time.Time) (int, error) {
	count, err := a.storage.DeleteOldEventsByDay(ctx, a.logger, dateTime)
	if err != nil {
		message := helpers.StringBuild("event delete error(count deletions - ", strconv.Itoa(count), "),error: ", err.Error())
		a.logger.Error(message)
	} else {
		message := helpers.StringBuild("event deleted(count deletions - ", strconv.Itoa(count), ")")
		a.logger.Info(message)
	}

	return count, err
}