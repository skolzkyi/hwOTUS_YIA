package kafka

import (
	"errors"
	"net"
	"strconv"
	"fmt"
	"runtime/debug"

	"github.com/segmentio/kafka-go" //nolint:typecheck
)

var ErrDialLeader = errors.New("failed to dial leader")
var ErrCloseConn = errors.New("failed to close connection")

func logf(msg string, a ...interface{}) {
	fmt.Printf(msg, a...)
	fmt.Println()
}

func recoveryFunction() {
	if recoveryMessage:=recover(); recoveryMessage != nil {
	  fmt.Println("kafka_recovery_message(panic): ",recoveryMessage)
	  fmt.Println("stack: ", string(debug.Stack()))
	}
}

func CreateTopic(topicName string, kafkaURL string, kafkaPort string) error {
	conn, err := kafka.Dial("tcp", kafkaURL+":"+kafkaPort) //nolint:typecheck
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}
	var controllerConn *kafka.Conn //nolint:typecheck
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))) //nolint:typecheck
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{ //nolint:typecheck
		{
			Topic:             topicName,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		return err
	}
	return nil
}
