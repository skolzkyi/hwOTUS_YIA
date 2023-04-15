package main

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/logger"
	pb "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/server/grpc/pb"
	kafka "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/kafka"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Notification struct {
	ID        int
	Title     string
	DateStart string
	User      string
}

func AgentActionGetListEventsNotificationByDayAgent(ctx context.Context, config Config, log logger.LoggerWrap, firstStart bool) error {
	curDate := time.Now()
	address := config.GetServerURL() + ":" + config.GetGRPSPort()
	conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("grpc connect error: " + err.Error())
		return err
	}

	client := pb.NewCalendarClient(conn)

	inData := &pb.GetEventsOnDayRequest{}
	inData.Date = timestamppb.New(curDate)
	outData, err := client.GetListEventsNotificationByDay(ctx, inData)
	if err != nil {
		log.Error("grpc GetListEventsNotificationByDay error: " + err.Error())
		return err
	}
	events := outData.GetEvents()
	log.Info("pbEventsLen: " + strconv.Itoa(len(events)))
	if len(events) > 0 {
		if firstStart && !config.GetKafkaAutoCreateTopicEnable() {
			err = kafka.CreateTopic(config.GetKafkaTopicName(), config.GetKafkaAddr(), config.kafkaPort)
			if err != nil {
				log.Error("kafka CreateTopic error: " + err.Error())
				return err
			}
		}
		kafkaWriter := kafka.NewWriter()
		kafkaWriter.Init(config.GetKafkaAddr(), config.GetKafkaPort(), config.GetKafkaTopicName())
		kafkaMessages := make([]string, 0)
		for _, curEvent := range events {
			notif := Notification{
				ID:        int(curEvent.GetId()),
				Title:     curEvent.GetTitle(),
				DateStart: curEvent.GetDatestart().AsTime().Format("2006-01-02 15:04:05"),
				User:      curEvent.GetUserid(),
			}
			jsonstring, err := json.Marshal(notif)
			if err != nil {
				log.Error("kafka json.Marshal error: " + err.Error())
				return err
			}
			kafkaMessages = append(kafkaMessages, string(jsonstring))
		}
		err = kafkaWriter.WriteMessagesPack(ctx, kafkaMessages)
		if err != nil {
			log.Error("kafka WriteMessagesPack error: " + err.Error())
			return err
		}
		err = kafkaWriter.Close()
		if err != nil {
			log.Error("kafka kafkaWriter.Close error: " + err.Error())
			return err
		}
	}
	return nil
}

func AgentActionDeleteOldEventsByDay(ctx context.Context, config Config, log logger.LoggerWrap, _ bool) error {
	curDate := time.Now()
	//curDate := time.Date(2024, 4, 19, 0, 0, 0, 1, time.UTC)
	address := config.GetServerURL() + ":" + config.GetGRPSPort()
	conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("grpc connect error: " + err.Error())
		return err
	}

	client := pb.NewCalendarClient(conn)

	inData := &pb.DeleteOldEventsRequest{}
	inData.Date = timestamppb.New(curDate)
	outData, err := client.DeleteOldEvents(ctx, inData)
	if err != nil {
		log.Error("grpc DeleteOldEventsByDay error: " + err.Error())
		return err
	}
	errorFromServer := outData.GetError()
	log.Info("Result: Count delOldEvent -  " + strconv.Itoa(int(outData.GetCount())) + "; Err - " + outData.GetError())

	if errorFromServer != "" {
		log.Error("AgentActionDeleteOldEventsByDay error on server: " + err.Error())
		return errors.New(errorFromServer)
	}

	return nil
}
