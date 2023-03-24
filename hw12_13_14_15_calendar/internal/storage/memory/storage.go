package memorystorage

import (
	"context"
	"io"
	"sort"
	"sync"
	"time"

	helpers "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/helpers"
	storage "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/storage/event"
)

type Storage struct {
	mu sync.RWMutex
	m  map[string]storage.Event
	io.Closer
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Init(_ context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m = make(map[string]storage.Event)
	return nil
}

func (s *Storage) Close() error {
	return nil
}

func (s *Storage) GetEvent(ctx context.Context, id string) (storage.Event, error) {
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

func (s *Storage) CreateEvent(ctx context.Context, value storage.Event) error {
	select {
	case <-ctx.Done():
		return storage.ErrStorageTimeout
	default:
		s.mu.Lock()
		defer s.mu.Unlock()
		s.m[value.ID] = value
		return nil
	}
}

func (s *Storage) UpdateEvent(ctx context.Context, value storage.Event) error {
	err := s.CreateEvent(ctx, value)
	return err
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	select {
	case <-ctx.Done():
		return storage.ErrStorageTimeout
	default:
		s.mu.Lock()
		defer s.mu.Unlock()
		delete(s.m, id)
		return nil
	}
}

func (s *Storage) getListEventsBetweenTwoDateInclude(ctx context.Context, StartTime time.Time, EndTime time.Time) ([]storage.Event, error) {
	resEvents := make([]storage.Event, 0)
	select {
	case <-ctx.Done():
		return nil, storage.ErrStorageTimeout
	default:
		s.mu.RLock()
		for _, curEvent := range s.m {
			if helpers.DateBetweenInclude(curEvent.DateStart, StartTime, EndTime) || helpers.DateBetweenInclude(curEvent.DateStop, StartTime, EndTime) {
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
