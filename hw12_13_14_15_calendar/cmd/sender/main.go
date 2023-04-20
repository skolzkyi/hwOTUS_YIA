package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	helpers "github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/helpers"
	"github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/internal/logger"
	"github.com/skolzkyi/hwOTUS_YIA/hw12_13_14_15_calendar/kafka"
)

type Notification struct {
	ID        int
	Title     string
	DateStart string
	User      string
}

func (n *Notification) String() string {
	return helpers.StringBuild("[ID: ", strconv.Itoa(n.ID), ", Title: ", n.Title, " DateStart: ", n.DateStart, " User: ", n.User) //nolint:lll
}

var configFilePath string

func init() {
	flag.StringVar(&configFilePath, "config", "./configs/config_sender.env", "Path to config_sender.env")
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
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	kafkaReader := kafka.NewReader()
	kafkaReader.Init(config.GetKafkaAddr(), config.GetKafkaPort(), config.GetKafkaTopicName(), "C_sender")
	log.Info("kafkaAddr: " + config.GetKafkaAddr())
	log.Info("Sender up")
	for {
		select {
		case <-ctx.Done():
			log.Info("Sender down")
			err = kafkaReader.Close()
			if err != nil {
				log.Error("KafkaReader close error: " + err.Error())
			}
			os.Exit(1) //nolint:gocritic
		default:
			kafkaMessage, err := kafkaReader.ReadMessage(ctx)
			if err != nil {
				log.Error("Sender crush on read kafka messages: " + err.Error())
				cancel()
			}
			if kafkaMessage != "" {
				notif := Notification{}
				err = json.Unmarshal([]byte(kafkaMessage), &notif)
				if err != nil {
					log.Error("Sender error unmarshalling notification: " + err.Error())
					continue
				}
				sendNotification(notif)
			}
		}
	}
}

func sendNotification(notif Notification) {
	fmt.Println("Message from sender: ", notif.String())
}
