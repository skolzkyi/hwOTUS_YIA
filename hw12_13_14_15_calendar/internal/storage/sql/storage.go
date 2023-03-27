package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	helpers "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/helpers"
	storage "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/storage/event"
)

type Storage struct {
	DB *sql.DB
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Init(ctx context.Context, config storage.Config) error {
	err := s.Connect(ctx, config)
	if err != nil {
		return err
	}
	if err = s.DB.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Storage) Connect(ctx context.Context, config storage.Config) error {
	dsn := helpers.StringBuild(config.GetDbName(), ":", config.GetDbPassword(), "@/", config.GetDbName(), "?parseTime=true")
	var err error
	s.DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	s.DB.SetConnMaxLifetime(config.GetdbConnMaxLifetime())
	s.DB.SetMaxOpenConns(config.GetDbMaxOpenConns())
	s.DB.SetMaxIdleConns(config.GetDbMaxIdleConns())

	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return storage.ErrStorageTimeout
	default:
		err := s.DB.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) GetEvent(ctx context.Context, id string) (storage.Event, error) {
	stmt := "SELECT id, title , userID, description , dateStart, dateStop, eventMessageTimeDelta FROM snippets WHERE id = ?"

	row := s.DB.QueryRow(stmt, id)

	event := &storage.Event{}

	err := row.Scan(&event.ID, &event.Title, &event.UserID, &event.Description, &event.DateStart, &event.DateStop, &event.EventMessageTimeDelta)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return storage.Event{}, storage.ErrNoRecord
		} else {
			return storage.Event{}, err
		}
	}

	return *event, nil

}

func (s *Storage) CreateEvent(ctx context.Context, value storage.Event) error {
	stmt := "INSERT INTO Events(id, title , userID, description , dateStart, dateStop, eventMessageTimeDelta) VALUES (?,?,?,?,?,?)"

	_, err := s.DB.ExecContext(ctx, stmt, value.ID, value.Title, value.UserID, value.Description, value.DateStart, value.DateStop, value.EventMessageTimeDelta)
	if err != nil {
		fmt.Println("CreateEvent error: ", err.Error())
		return err
	}

	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, value storage.Event) error {
	stmt := "UPDATE Events SET title =?,userID=?,description=?,dateStart=?, dateStop=?, eventMessageTimeDelta=? WHERE id=?"

	_, err := s.DB.ExecContext(ctx, stmt, value.Title, value.UserID, value.Description, value.DateStart, value.DateStop, value.EventMessageTimeDelta, value.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrNoRecord
		} else {
			return err
		}
	}
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	stmt := "DELETE from Events WHERE id=?"

	_, err := s.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) getListEventsBetweenTwoDateInclude(ctx context.Context, StartTime time.Time, EndTime time.Time) ([]storage.Event, error) {
	resEvents := make([]storage.Event, 0)

	return resEvents, nil

}

func (s *Storage) GetListEventsonDayByDay(ctx context.Context, day time.Time) ([]storage.Event, error) {
	dayStart := helpers.DateStartTime(day)
	dayEnd := helpers.DateEndTime(day)
	resEvents, err := s.getListEventsBetweenTwoDateInclude(ctx, dayStart, dayEnd)
	return resEvents, err
}

func (s *Storage) GetListEventsOnWeekByDay(ctx context.Context, day time.Time) ([]storage.Event, error) {
	dayStart := helpers.DateStartTime(day)
	dayEnd := helpers.DateEndTime(day.Add(168 * time.Hour))
	resEvents, err := s.getListEventsBetweenTwoDateInclude(ctx, dayStart, dayEnd)
	return resEvents, err
}

func (s *Storage) GetListEventsOnMonthByDay(ctx context.Context, day time.Time) ([]storage.Event, error) {
	dayStart := helpers.DateStartTime(day)
	dayEnd := helpers.DateEndTime(day.Add(720 * time.Hour))
	resEvents, err := s.getListEventsBetweenTwoDateInclude(ctx, dayStart, dayEnd)
	return resEvents, err
}
