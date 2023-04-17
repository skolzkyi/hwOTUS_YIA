package kafka

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	kafkago "github.com/segmentio/kafka-go"
)

var ErrWriteMessage = errors.New("failed to write messages")

type Writer struct {
	kWriter *kafkago.Writer
}



func NewWriter() *Writer {
	writer := &kafkago.Writer{}
	return &Writer{
		kWriter: writer,
	}
}

func (w *Writer) Init(addr string, port string, topicName string, autoTopicCreation bool) {

	w.kWriter = &kafka.Writer{
		Addr:     kafka.TCP(addr + ":" + port),
		Topic:    topicName,
		Logger:      kafka.LoggerFunc(logf),
	    ErrorLogger: kafka.LoggerFunc(logf),
		Balancer: &kafka.LeastBytes{},
		AllowAutoTopicCreation: autoTopicCreation,
	}
}

func (w *Writer) WriteMessagesPack(ctx context.Context, messagesPack []string) error {
	var err error
	kMessages := make([]kafkago.Message, 0)
	for _, curMes := range messagesPack {
		kMessages = append(kMessages, kafka.Message{Value: []byte(curMes)})
	}
	fmt.Println("messages: ",kMessages)
	retries := 3
	for i := 0; i < retries; i++ {
		err = w.kWriter.WriteMessages(ctx, kMessages...)
    	fmt.Println("WrmP: ",err.Error())
		
		if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
        	time.Sleep(time.Millisecond * 250)
        	continue
    	}
	}
	
	return err
}

func (w *Writer) Close() error {
	err := w.kWriter.Close()
	return err
}
