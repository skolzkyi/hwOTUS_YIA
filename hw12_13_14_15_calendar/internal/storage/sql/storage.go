package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	mysql "github.com/go-sql-driver/mysql"
	helpers "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/helpers"
	storage "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/storage/event"
)

var ErrSQLTimeConvert = errors.New("SQL DateTime convertation error")

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
	dsn := helpers.StringBuild(config.GetDbUser(), ":", config.GetDbPassword(), "@/", config.GetDbName(), "?parseTime=true")
	fmt.Println("dsn: ", dsn)
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

func (s *Storage) GetEvent(ctx context.Context, id int) (storage.Event, error) {
	stmt := "SELECT id, title , userID, description , dateStart, dateStop, eventMessageTimeDelta FROM eventsTable WHERE id = ?"

	row := s.DB.QueryRowContext(ctx, stmt, id)

	event := &storage.Event{}

	var ntStart mysql.NullTime
	var ntEnd mysql.NullTime
	var int64Delta int64

	err := row.Scan(&event.ID, &event.Title, &event.UserID, &event.Description, &ntStart, &ntEnd, &int64Delta)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.Event{}, storage.ErrNoRecord
		} else {
			return storage.Event{}, err
		}
	}
	if ntStart.Valid {
		event.DateStart = ntStart.Time
	} else {
		err = ErrSQLTimeConvert
		return storage.Event{}, err
	}
	if ntStart.Valid {
		event.DateStop = ntEnd.Time
	} else {
		err = ErrSQLTimeConvert
		return storage.Event{}, err
	}
	dr := time.Duration(int64Delta) * time.Millisecond
	event.EventMessageTimeDelta = dr

	return *event, nil

}

func (s *Storage) CreateEvent(ctx context.Context, value storage.Event) (int, error) {
	fmt.Println("inSQLcreate")
	ok, err := s.isEventOnThisTimeExcluded(ctx, value)
	if err != nil {
		fmt.Println("busy check error: ", err.Error())
		return 0, err
	}
	if ok {
		return 0, storage.ErrDateBusy
	}
	stmt := "INSERT INTO eventsTable(title , userID, description , dateStart, dateStop, eventMessageTimeDelta) VALUES (?,?,?,?,?,?)"
	res, err := s.DB.ExecContext(ctx, stmt, value.Title, value.UserID, value.Description, value.DateStart, value.DateStop, int64(value.EventMessageTimeDelta))
	if err != nil {
		fmt.Println("CreateEvent error: ", err.Error())
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		fmt.Println("CreateEvent get new id error: ", err.Error())
		return 0, err
	}

	return int(id), nil
}

func (s *Storage) UpdateEvent(ctx context.Context, value storage.Event) error {
	stmt := "UPDATE eventsTable SET title =?,userID=?,description=?,dateStart=?, dateStop=?, eventMessageTimeDelta=? WHERE id=?"

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

func (s *Storage) DeleteEvent(ctx context.Context, id int) error {
	stmt := "DELETE from eventsTable WHERE id=?"

	_, err := s.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) getListEventsBetweenTwoDateInclude(ctx context.Context, StartTime time.Time, EndTime time.Time) ([]storage.Event, error) {
	resEvents := make([]storage.Event, 0)
	startTimeStr := StartTime.Format("2006-01-02 15:04:05")
	endTimeStr := EndTime.Format("2006-01-02 15:04:05")
	stmt := "SELECT id, title , userID, description , dateStart, dateStop, eventMessageTimeDelta FROM eventsTable WHERE CAST(dateStart AS DATE) BETWEEN CAST('" + startTimeStr + "' AS DATE)  AND CAST('" + endTimeStr + "' AS DATE)"

	rows, err := s.DB.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	event := storage.Event{}

	var ntStart mysql.NullTime
	var ntEnd mysql.NullTime
	var int64Delta int64
	for rows.Next() {
		err = rows.Scan(&event.ID, &event.Title, &event.UserID, &event.Description, &ntStart, &ntEnd, &int64Delta)
		if err != nil {
			return nil, err
		}
		if ntStart.Valid {
			event.DateStart = ntStart.Time
		} else {
			err = ErrSQLTimeConvert
			return nil, err
		}
		if ntEnd.Valid {
			event.DateStop = ntEnd.Time
		} else {
			err = ErrSQLTimeConvert
			return nil, err
		}
		dr := time.Duration(int64Delta) * time.Millisecond
		event.EventMessageTimeDelta = dr

		resEvents = append(resEvents, event)
		event = storage.Event{}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

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

func (s *Storage) isEventOnThisTimeExcluded(ctx context.Context, value storage.Event) (bool, error) {
	startTimeStr := value.DateStart.Format("2006-01-02 15:04:05")
	endTimeStr := value.DateStop.Format("2006-01-02 15:04:05")
	stmt := "SELECT id FROM eventsTable WHERE UserID=? AND CAST(dateStart AS DATE) BETWEEN CAST('" + startTimeStr + "' AS DATE) AND  CAST('" + endTimeStr + "' AS DATE) OR CAST(dateStop AS DATE) BETWEEN CAST('" + startTimeStr + "' AS DATE) AND  CAST('" + endTimeStr + "' AS DATE)"
	fmt.Println("stmt: ", stmt)
	//CAST(dateStart AS DATE)
	row := s.DB.QueryRowContext(ctx, stmt, value.UserID)

	var id int

	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		} else {
			fmt.Println("error: ", err.Error())
			os.Exit(1)
			//return false, err
		}
	}
	return true, nil
}
