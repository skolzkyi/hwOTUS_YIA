package memorystorage

import (
	"context"
	"fmt"
	"sort"
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

func (s *Storage) Init(_ context.Context, _ storage.Config) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m = make(map[int]storage.Event)
	s.idCounter = 0
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	return nil
}

func (s *Storage) GetEvent(ctx context.Context, id int) (storage.Event, error) {
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

func (s *Storage) CreateEvent(ctx context.Context, value storage.Event) (int, error) {
	select {
	case <-ctx.Done():
		return 0, storage.ErrStorageTimeout
	default:
		ok, err := s.isEventOnThisTimeExcluded(ctx, value, false)
		if err != nil {
			fmt.Println("busy check error: ", err.Error())
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
			return 0, storage.ErrIDNotUnique
		}
		s.m[id] = value
		s.idCounter++
		return id, nil
	}
}

func (s *Storage) UpdateEvent(ctx context.Context, value storage.Event) error {
	select {
	case <-ctx.Done():
		return storage.ErrStorageTimeout
	default:
		ok, err := s.isEventOnThisTimeExcluded(ctx, value, true)
		if err != nil {
			fmt.Println("busy check error: ", err.Error())
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

func (s *Storage) DeleteEvent(ctx context.Context, id int) error {
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
func (s *Storage) getListEventsBetweenTwoDateInclude(ctx context.Context, startTime time.Time, endTime time.Time) ([]storage.Event, error) { //nolint:lll
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
