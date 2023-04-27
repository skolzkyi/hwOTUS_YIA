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
	flag.StringVar(&configFilePath, "config", "./configs/", "Path to config_sheduler.env")
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
	log, err := logger.New(config.Logger.Level)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()
	log.Info("servAddr: " + config.GetServerURL())
	log.Info("kafkaAddr: " + config.GetKafkaAddr())
	sheduler := NewSheduler()
	sheduler.Init(log, config.GetShedulerPeriod())
	notificationAgent := NewAgent()
	notificationAgent.Init("GetListEventsNotificationByDayAgent", config.GetNotificationEventPeriod(), AgentActionGetListEventsNotificationByDayAgent) //nolint:lll
	oldEventCleaningAgent := NewAgent()
	oldEventCleaningAgent.Init("DeleteOldEventsByDayAgent", config.GetCleanOldEventPeriod(), AgentActionDeleteOldEventsByDay) //nolint:lll
	sheduler.AddAgent(notificationAgent)
	sheduler.AddAgent(oldEventCleaningAgent)
	sheduler.RunAgents(ctx, config)
}
