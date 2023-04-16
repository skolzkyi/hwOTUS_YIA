package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/logger"
)

var configFilePath string

func init() {
	flag.StringVar(&configFilePath, "config", "./configs/config_sheduler.env", "Path to config_sheduler.env")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig()
	err := config.Init(configFilePath)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config: ", config)
	logg, err := logger.New(config.Logger.Level)
	if err != nil {
		fmt.Println(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	sheduler := NewSheduler()
	sheduler.Init(logg, config.GetShedulerPeriod())
	notificationAgent := NewAgent()
	notificationAgent.Init("GetListEventsNotificationByDayAgent", config.GetNotificationEventPeriod(), AgentActionGetListEventsNotificationByDayAgent)
	oldEventCleaningAgent := NewAgent()
	oldEventCleaningAgent.Init("DeleteOldEventsByDayAgent", config.GetCleanOldEventPeriod(), AgentActionDeleteOldEventsByDay)
	sheduler.AddAgent(notificationAgent)
	sheduler.AddAgent(oldEventCleaningAgent)
	sheduler.RunAgents(ctx, config)
}
