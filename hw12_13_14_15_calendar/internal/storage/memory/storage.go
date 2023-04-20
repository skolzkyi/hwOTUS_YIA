package memorystorage

import (
	"context"
	"sort"
	"strconv"
	"sync"
	"time"

	helpers "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/helpers"
	storage "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/storage/event"
)

type Storage struct {
	mu        sync.RWMutex
	m         map[int]storage.Event
	idCounter int
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Init(_ context.Context, _ storage.Logger, _ storage.Config) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m = make(map[int]storage.Event)
	s.idCounter = 0
	return nil
}

func (s *Storage) Close(_ context.Context, _ storage.Logger) error {
	return nil
}

func (s *Storage) GetEvent(ctx context.Context, _ storage.Logger, id int) (storage.Event, error) {
	select {
	case <-ctx.Done():
		return storage.Event{}, storage.ErrStorageTimeout
	default:
		s.mu.RLock()
		defer s.mu.RUnlock()
		var err error
		val, ok := s.m[id]
		if !ok {
			err = storage.ErrNoRecord
		}
		return val, err
	}
}

func (s *Storage) CreateEvent(ctx context.Context, logger storage.Logger, value storage.Event) (int, error) {
	select {
	case <-ctx.Done():
		return 0, storage.ErrStorageTimeout
	default:
		ok, err := s.isEventOnThisTimeExcluded(ctx, logger, value, false)
		if err != nil {
			logger.Error("Memory storage CreateEvent busy check error" + err.Error())
			return 0, err
		}
		if ok {
			return 0, storage.ErrDateBusy
		}
		id := s.idCounter
		value.ID = s.idCounter

		s.mu.Lock()
		defer s.mu.Unlock()
		_, ok = s.m[id]
		if ok {
			logger.Error("Memory storage not unique ID error, ID:" + strconv.Itoa(id))
			return 0, storage.ErrIDNotUnique
		}
		s.m[id] = value
		s.idCounter++
		return id, nil
	}
}

func (s *Storage) UpdateEvent(ctx context.Context, logger storage.Logger, value storage.Event) error {
	select {
	case <-ctx.Done():
		return storage.ErrStorageTimeout
	default:
		ok, err := s.isEventOnThisTimeExcluded(ctx, logger, value, true)
		if err != nil {
			logger.Error("Memory storage UpdateEvent busy check error" + err.Error())
			return err
		}
		if ok {
			return storage.ErrDateBusy
		}
		_, ok = s.m[value.ID]
		if !ok {
			return storage.ErrNoRecord
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		s.m[value.ID] = value
	}
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, _ storage.Logger, id int) error {
	select {
	case <-ctx.Done():
		return storage.ErrStorageTimeout
	default:
		_, ok := s.m[id]
		if !ok {
			return storage.ErrNoRecord
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		delete(s.m, id)
		return nil
	}
}

// small name not informative.
func (s *Storage) getListEventsBetweenTwoDateInclude(ctx context.Context, _ storage.Logger, startTime time.Time, endTime time.Time) ([]storage.Event, error) { //nolint:lll
	resEvents := make([]storage.Event, 0)
	select {
	case <-ctx.Done():
		return nil, storage.ErrStorageTimeout
	default:
		s.mu.RLock()
		for _, curEvent := range s.m {
			startDateBetweenInclude := helpers.DateBetweenInclude(curEvent.DateStart, startTime, endTime)
			endDateBetweenInclude := helpers.DateBetweenInclude(curEvent.DateStop, startTime, endTime)
			if startDateBetweenInclude || endDateBetweenInclude {
				resEvents = append(resEvents, curEvent)
			}
		}
		s.mu.RUnlock()
		sort.SliceStable(resEvents, func(i, j int) bool {
			return resEvents[i].DateStart.Before(resEvents[j].DateStart)
		})
		return resEvents, nil
	}
}

func (s *Storage) GetListEventsNotificationByDay(ctx context.Context,_ storage.Logger, dateTime time.Time) ([]storage.Event, error) {
	resEvents := make([]storage.Event, 0)
	select {
	case <-ctx.Done():
		return nil, storage.ErrStorageTimeout
	default:
		s.mu.RLock()
		for _, curEvent := range s.m {
			controlTime := curEvent.DateStart.Add(-1 * curEvent.EventMessageTimeDelta)
			if (dateTime.After(controlTime) || dateTime.Equal(controlTime)) && curEvent.DateStart.After(dateTime) {
				resEvents = append(resEvents, curEvent)
			}
		}
		s.mu.RUnlock()
		sort.SliceStable(resEvents, func(i, j int) bool {
			return resEvents[i].DateStart.Before(resEvents[j].DateStart)
		})
		return resEvents, nil
	}
}

func (s *Storage) DeleteOldEventsByDay(ctx context.Context,_ storage.Logger, dateTime time.Time) (int, error) {
	var curEvent storage.Event
	var i int
	idEventForDeletion := make(map[int]struct{})
	controlTime := dateTime.Add(-8760 * time.Hour)
	select {
	case <-ctx.Done():
		return i, storage.ErrStorageTimeout
	default:
		s.mu.RLock()
		for _, curEvent = range s.m {
			if curEvent.DateStop.Before(controlTime) { // -365 days
				idEventForDeletion[curEvent.ID] = struct{}{}
			}
		}
		s.mu.RUnlock()
		s.mu.Lock()
		for curID := range idEventForDeletion {
			delete(s.m, curID)
			i++
		}
		s.mu.Unlock()

		return i, nil
	}
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

func (s *Storage) isEventOnThisTimeExcluded(ctx context.Context, _ storage.Logger, value storage.Event, exclID bool) (bool, error) { //nolint:lll
	select {
	case <-ctx.Done():
		return false, storage.ErrStorageTimeout
	default:
		s.mu.RLock()
		for _, curEvent := range s.m {
			startDateBetweenInclude := helpers.DateBetweenInclude(curEvent.DateStart, value.DateStart, value.DateStop)
			endDateBetweenInclude := helpers.DateBetweenInclude(curEvent.DateStop, value.DateStart, value.DateStop)
			if (startDateBetweenInclude || endDateBetweenInclude) && curEvent.UserID == value.UserID {
				if exclID && curEvent.ID == value.ID {
					continue
				}
				s.mu.RUnlock()
				return true, nil
			}
		}
		s.mu.RUnlock()

		return false, nil
	}
}
