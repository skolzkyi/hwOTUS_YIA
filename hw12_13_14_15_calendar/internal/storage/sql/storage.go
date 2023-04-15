package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
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

func (s *Storage) Init(ctx context.Context, logger storage.Logger, config storage.Config) error {
	err := s.Connect(ctx, logger, config)
	if err != nil {
		logger.Error("SQL connect error: " + err.Error())
		return err
	}
	err = s.DB.PingContext(ctx)
	if err != nil {
		logger.Error("SQL DB ping error: " + err.Error())
		return err
	}

	return err
}

func (s *Storage) Connect(ctx context.Context, logger storage.Logger, config storage.Config) error {
	select {
	case <-ctx.Done():
		return storage.ErrStorageTimeout
	default:
		dsn := helpers.StringBuild(config.GetDBUser(), ":", config.GetDBPassword(), "@/", config.GetDBName(), "?parseTime=true") //nolint:lll
		// fmt.Println("dsn: ", dsn)
		var err error
		s.DB, err = sql.Open("mysql", dsn)
		if err != nil {
			logger.Error("SQL open error: " + err.Error())
			return err
		}

		s.DB.SetConnMaxLifetime(config.GetDBConnMaxLifetime())
		s.DB.SetMaxOpenConns(config.GetDBMaxOpenConns())
		s.DB.SetMaxIdleConns(config.GetDBMaxIdleConns())

		return nil
	}
}

func (s *Storage) Close(ctx context.Context, logger storage.Logger) error {
	select {
	case <-ctx.Done():
		return storage.ErrStorageTimeout
	default:
		err := s.DB.Close()
		if err != nil {
			logger.Error("SQL DB close error: " + err.Error())
			return err
		}
	}
	return nil
}

func (s *Storage) GetEvent(ctx context.Context, logger storage.Logger, id int) (storage.Event, error) {
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
		logger.Error("SQL row scan GetEvent error: " + err.Error() + " stmt: " + stmt)
		return storage.Event{}, err
	}
	if ntStart.Valid {
		event.DateStart = ntStart.Time
	} else {
		logger.Error("SQL GetEvent dateTime convert error")
		return storage.Event{}, ErrSQLTimeConvert
	}
	if ntStart.Valid {
		event.DateStop = ntEnd.Time
	} else {
		logger.Error("SQL GetEvent dateTime convert error")
		return storage.Event{}, ErrSQLTimeConvert
	}
	dr := time.Duration(int64Delta) * time.Millisecond
	event.EventMessageTimeDelta = dr

	return *event, nil
}

func (s *Storage) CreateEvent(ctx context.Context, logger storage.Logger, value storage.Event) (int, error) {
	// fmt.Println("inSQLcreate")
	ok, err := s.isEventOnThisTimeExcluded(ctx, logger, value, false)
	if err != nil {
		logger.Error("SQL CreateEvent busy check error" + err.Error())
		return 0, err
	}
	if ok {
		return 0, storage.ErrDateBusy
	}
	stmt := "INSERT INTO eventsTable(title , userID, description , dateStart, dateStop, eventMessageTimeDelta) VALUES (?,?,?,?,?,?)"                           //nolint:lll
	res, err := s.DB.ExecContext(ctx, stmt, value.Title, value.UserID, value.Description, value.DateStart, value.DateStop, int64(value.EventMessageTimeDelta)) //nolint:lll
	if err != nil {
		logger.Error("SQL DB exec stmt CreateEvent error: " + err.Error() + " stmt: " + stmt)
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		logger.Error("SQL CreateEvent get new id error: " + err.Error() + " stmt: " + stmt)
		return 0, err
	}

	return int(id), nil
}

func (s *Storage) UpdateEvent(ctx context.Context, logger storage.Logger, value storage.Event) error {
	ok, err := s.isEventOnThisTimeExcluded(ctx, logger, value, true)
	if err != nil {
		logger.Error("SQL UpdateEvent busy check error" + err.Error())
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
		logger.Error("SQL DB exec stmt UpdateEvent error: " + err.Error() + " stmt: " + stmt)
		return err
	}
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, logger storage.Logger, id int) error {
	stmt := "DELETE from eventsTable WHERE id=?"

	_, err := s.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		logger.Error("SQL DB exec stmt DeleteEent error: " + err.Error() + " stmt: " + stmt)
		return err
	}
	return nil
}

func (s *Storage) GetListEventsNotificationByDay(ctx context.Context,logger storage.Logger, dateTime time.Time) ([]storage.Event, error) {
	resEvents := make([]storage.Event, 0)
	dateTimeStr := dateTime.Format("2006-01-02 15:04:05")
	stmt := "SELECT id, title , userID, description , dateStart, dateStop, eventMessageTimeDelta FROM eventsTable WHERE CAST('" + dateTimeStr + "' AS DATE) BETWEEN DATE_SUB(CAST('" + dateTimeStr + "' AS DATE), INTERVAL eventMessageTimeDelta*1000 MICROSECOND)  AND CAST(dateStart AS DATE) AND CAST('" + dateTimeStr + "' AS DATE) < CAST(dateStart AS DATE) ORDER BY dateStart ASC"

	rows, err := s.DB.QueryContext(ctx, stmt)
	if err != nil {
		logger.Error("SQL QueryContext stmt GetListEventsNotificationByDay error: " + err.Error() + " stmt: " + stmt)
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
			logger.Error("SQL rows scan GetListEventsNotificationByDay error")
			return nil, err
		}
		if ntStart.Valid {
			event.DateStart = ntStart.Time
		} else {
			logger.Error("SQL GetListEventsNotificationByDay dateTime convert error")
			err = ErrSQLTimeConvert
			return nil, err
		}
		if ntEnd.Valid {
			event.DateStop = ntEnd.Time
		} else {
			logger.Error("SQL GetListEventsNotificationByDay dateTime convert error")
			err = ErrSQLTimeConvert
			return nil, err
		}
		dr := time.Duration(int64Delta) * time.Millisecond
		event.EventMessageTimeDelta = dr

		resEvents = append(resEvents, event)
		event = storage.Event{}
	}

	if err = rows.Err(); err != nil {
		logger.Error("SQL GetListEventsNotificationByDay rows error: " + err.Error())
		return nil, err
	}

	return resEvents, nil

}

func (s *Storage) DeleteOldEventsByDay(ctx context.Context,logger storage.Logger, dateTime time.Time) (int, error) {
	controlDateTime := dateTime.Add(-8760 * time.Hour) // -365 days
	controlDateTimeStr := controlDateTime.Format("2006-01-02 15:04:05")
	stmt := "DELETE from eventsTable WHERE CAST(dateStop AS DATE) < CAST('" + controlDateTimeStr + "' AS DATE)"

	res, err := s.DB.ExecContext(ctx, stmt)
	if err != nil {
		logger.Error("SQL DeleteOldEventsByDay DB exec error: " + err.Error())
		return 0, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		logger.Error("SQL DeleteOldEventsByDay rows affected error: " + err.Error())
		return 0, err
	}
	return int(count), nil
}

// small name not informative.
func (s *Storage) getListEventsBetweenTwoDateInclude(ctx context.Context, logger storage.Logger, startTime time.Time, endTime time.Time) ([]storage.Event, error) { //nolint:lll
	resEvents := make([]storage.Event, 0)
	startTimeStr := startTime.Format("2006-01-02 15:04:05")
	endTimeStr := endTime.Format("2006-01-02 15:04:05")
	stmt := "SELECT id, title , userID, description , dateStart, dateStop, eventMessageTimeDelta FROM eventsTable WHERE CAST(dateStart AS DATE) BETWEEN CAST('" + startTimeStr + "' AS DATE)  AND CAST('" + endTimeStr + "' AS DATE)" //nolint:lll

	rows, err := s.DB.QueryContext(ctx, stmt)
	if err != nil {
		logger.Error("SQL getListEventsBetweenTwoDateInclude DB query error: " + err.Error() + " stmt: " + stmt)
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
			logger.Error("SQL rows scan getListEventsBetweenTwoDateInclude error")
			return nil, err
		}
		if ntStart.Valid {
			event.DateStart = ntStart.Time
		} else {
			logger.Error("SQL getListEventsBetweenTwoDateInclude dateTime convert error")
			return nil, ErrSQLTimeConvert
		}
		if ntEnd.Valid {
			event.DateStop = ntEnd.Time
		} else {
			logger.Error("SQL getListEventsBetweenTwoDateInclude dateTime convert error")
			return nil, ErrSQLTimeConvert
		}
		dr := time.Duration(int64Delta) * time.Millisecond
		event.EventMessageTimeDelta = dr

		resEvents = append(resEvents, event)
		event = storage.Event{}
	}

	if err = rows.Err(); err != nil {
		logger.Error("SQL getListEventsBetweenTwoDateInclude rows error: " + err.Error())
		return nil, err
	}

	return resEvents, nil
}

func (s *Storage) GetListEventsonDayByDay(ctx context.Context, logger storage.Logger, day time.Time) ([]storage.Event, error) { //nolint:lll
	dayStart := helpers.DateStartTime(day)
	dayEnd := helpers.DateEndTime(day)
	resEvents, err := s.getListEventsBetweenTwoDateInclude(ctx, logger, dayStart, dayEnd)
	return resEvents, err
}

func (s *Storage) GetListEventsOnWeekByDay(ctx context.Context, logger storage.Logger, day time.Time) ([]storage.Event, error) { //nolint:lll
	weekTime := 168 * time.Hour
	dayStart := helpers.DateStartTime(day)
	dayEnd := helpers.DateEndTime(day.Add(weekTime))
	resEvents, err := s.getListEventsBetweenTwoDateInclude(ctx, logger, dayStart, dayEnd)
	return resEvents, err
}

func (s *Storage) GetListEventsOnMonthByDay(ctx context.Context, logger storage.Logger, day time.Time) ([]storage.Event, error) { //nolint:lll
	monthTime := 720 * time.Hour
	dayStart := helpers.DateStartTime(day)
	dayEnd := helpers.DateEndTime(day.Add(monthTime))
	resEvents, err := s.getListEventsBetweenTwoDateInclude(ctx, logger, dayStart, dayEnd)
	return resEvents, err
}

func (s *Storage) isEventOnThisTimeExcluded(ctx context.Context, logger storage.Logger, value storage.Event, exclID bool) (bool, error) { //nolint:lll
	startTimeStr := value.DateStart.Format("2006-01-02 15:04:05")
	endTimeStr := value.DateStop.Format("2006-01-02 15:04:05")
	var addStmt string
	if exclID {
		addStmt = " AND ID!=? "
	}
	stmt := "SELECT id FROM eventsTable WHERE UserID=?" + addStmt + " AND CAST(dateStart AS DATE) BETWEEN CAST('" + startTimeStr + "' AS DATE) AND  CAST('" + endTimeStr + "' AS DATE) OR CAST(dateStop AS DATE) BETWEEN CAST('" + startTimeStr + "' AS DATE) AND  CAST('" + endTimeStr + "' AS DATE)" //nolint:lll
	// fmt.Println("stmt: ", stmt)

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
		logger.Error("SQL isEventOnThisTimeExcluded rows error: " + err.Error())
		return false, err
	}
	return true, nil
}
