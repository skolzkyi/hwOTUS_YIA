package main

import (
	"context"
	"time"

	"github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/logger"
)

type Sheduler struct {
	agents   []*Agent
	Logger   *logger.LogWrap
	periodic time.Duration
	ticker   *time.Ticker
}

func NewSheduler() *Sheduler {
	return &Sheduler{}
}
func (s *Sheduler) Init(logg *logger.LogWrap, periodic time.Duration) {
	s.Logger = logg
	s.agents = make([]*Agent, 0)
	s.periodic = periodic
}
func (s *Sheduler) AddAgent(agent *Agent) {
	s.agents = append(s.agents, agent)
}
func (s *Sheduler) RunAgents(ctx context.Context, config Config) {
	s.Logger.Info("Sheduler Up")
	s.ticker = time.NewTicker(s.periodic)
	for {
		select {
		case <-ctx.Done():
			s.Logger.Info("Sheduler Down")
			return
		case <-s.ticker.C:
			for _, curAgent := range s.agents {
				curTime := time.Now()
				controlTime := curAgent.lastActionTime.Add(curAgent.periodic)
				if controlTime.Before(curTime) || controlTime.Equal(curTime) || curAgent.lastActionTime.IsZero() || curAgent.firstStart {
					ctxAct, _ := context.WithTimeout(ctx, 10*time.Second)
					s.Logger.Info("ShedulerAgentStarted: " + curAgent.name)
					err := curAgent.action(ctxAct, config, *s.Logger, curAgent.firstStart)
					if err != nil {
						s.Logger.Error("ShedulerAgentError: " + err.Error() + " (" + curAgent.name + ")")
					} else {
						curAgent.firstStart = false
						s.Logger.Info("ShedulerAgentStopped: " + curAgent.name)
					}
					curAgent.lastActionTime = curTime
					curAgent.lastError = err
				}
			}
		}
	}
}

type Agent struct {
	name           string
	periodic       time.Duration
	action         func(ctx context.Context, config Config, log logger.LogWrap, firstStart bool) error
	lastActionTime time.Time
	lastError      error
	firstStart     bool
}

func NewAgent() *Agent {
	return &Agent{}
}
func (a *Agent) Init(name string, periodic time.Duration, action func(ctx context.Context, config Config, log logger.LogWrap, firstStart bool) error) {
	a.name = name
	a.periodic = periodic
	a.action = action
	a.firstStart = true
}
