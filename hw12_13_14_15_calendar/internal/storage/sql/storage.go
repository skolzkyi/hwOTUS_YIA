package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // for driver
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
	err = s.DB.PingContext(ctx)

	return err
}

func (s *Storage) Connect(ctx context.Context, config storage.Config) error {
	select {
	case <-ctx.Done():
		return storage.ErrStorageTimeout
	default:
		dsn := helpers.StringBuild(config.GetDBUser(), ":", config.GetDBPassword(), "@/", config.GetDBName(), "?parseTime=true") //nolint:lll
		// fmt.Println("dsn: ", dsn)
		var err error
		s.DB, err = sql.Open("mysql", dsn)
		if err != nil {
			return err
		}

		s.DB.SetConnMaxLifetime(config.GetDBConnMaxLifetime())
		s.DB.SetMaxOpenConns(config.GetDBMaxOpenConns())
		s.DB.SetMaxIdleConns(config.GetDBMaxIdleConns())

		return nil
	}
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
	stmt := "SELECT id, title , userID, description , dateStart, dateStop, eventMessageTimeDelta FROM eventsTable WHERE id = ?" //nolint:lll

	row := s.DB.QueryRowContext(ctx, stmt, id)

	event := &storage.Event{}

	var ntStart sql.NullTime
	var ntEnd sql.NullTime
	var int64Delta int64

	err := row.Scan(&event.ID, &event.Title, &event.UserID, &event.Description, &ntStart, &ntEnd, &int64Delta)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.Event{}, storage.ErrNoRecord
		}
		return storage.Event{}, err
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
	// fmt.Println("inSQLcreate")
	ok, err := s.isEventOnThisTimeExcluded(ctx, value, false)
	if err != nil {
		fmt.Println("busy check error: ", err.Error())
		return 0, err
	}
	if ok {
		return 0, storage.ErrDateBusy
	}
	stmt := "INSERT INTO eventsTable(title , userID, description , dateStart, dateStop, eventMessageTimeDelta) VALUES (?,?,?,?,?,?)"                           //nolint:lll
	res, err := s.DB.ExecContext(ctx, stmt, value.Title, value.UserID, value.Description, value.DateStart, value.DateStop, int64(value.EventMessageTimeDelta)) //nolint:lll
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
	ok, err := s.isEventOnThisTimeExcluded(ctx, value, true)
	if err != nil {
		fmt.Println("busy check error: ", err.Error())
		return err
	}
	if ok {
		return storage.ErrDateBusy
	}
	stmt := "UPDATE eventsTable SET title =?,userID=?,description=?,dateStart=?, dateStop=?, eventMessageTimeDelta=? WHERE id=?" //nolint:lll

	_, err = s.DB.ExecContext(ctx, stmt, value.Title, value.UserID, value.Description, value.DateStart, value.DateStop, value.EventMessageTimeDelta, value.ID) //nolint:lll
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrNoRecord
		}
		return err
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

// small name not informative.
func (s *Storage) getListEventsBetweenTwoDateInclude(ctx context.Context, startTime time.Time, endTime time.Time) ([]storage.Event, error) { //nolint:lll
	resEvents := make([]storage.Event, 0)
	startTimeStr := startTime.Format("2006-01-02 15:04:05")
	endTimeStr := endTime.Format("2006-01-02 15:04:05")
	stmt := "SELECT id, title , userID, description , dateStart, dateStop, eventMessageTimeDelta FROM eventsTable WHERE CAST(dateStart AS DATE) BETWEEN CAST('" + startTimeStr + "' AS DATE)  AND CAST('" + endTimeStr + "' AS DATE)" //nolint:lll

	rows, err := s.DB.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	event := storage.Event{}

	var ntStart sql.NullTime
	var ntEnd sql.NullTime
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

func (s *Storage) isEventOnThisTimeExcluded(ctx context.Context, value storage.Event, exclID bool) (bool, error) {
	startTimeStr := value.DateStart.Format("2006-01-02 15:04:05")
	endTimeStr := value.DateStop.Format("2006-01-02 15:04:05")
	var addStmt string
	if exclID {
		addStmt = " AND ID!=? "
	}
	stmt := "SELECT id FROM eventsTable WHERE UserID=?" + addStmt + " AND CAST(dateStart AS DATE) BETWEEN CAST('" + startTimeStr + "' AS DATE) AND  CAST('" + endTimeStr + "' AS DATE) OR CAST(dateStop AS DATE) BETWEEN CAST('" + startTimeStr + "' AS DATE) AND  CAST('" + endTimeStr + "' AS DATE)" //nolint:lll
	fmt.Println("stmt: ", stmt)

	var row *sql.Row
	if exclID {
		row = s.DB.QueryRowContext(ctx, stmt, value.UserID, value.ID)
	} else {
		row = s.DB.QueryRowContext(ctx, stmt, value.UserID)
	}

	var id int

	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
