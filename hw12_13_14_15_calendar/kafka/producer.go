package kafka

import (
	"context"
	"errors"
	//"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	kafkago "github.com/segmentio/kafka-go"
)

var ErrWriteMessage = errors.New("failed to write messages")

type Writer struct {
	kWriter *kafka.Writer
}



func NewWriter() Writer {
	return Writer{}
}

func (w *Writer) Init(addr string, port string, topicName string, autoTopicCreation bool) {

	w.kWriter = &kafka.Writer{
		Addr:     kafka.TCP(addr + ":" + port),
		Topic:    topicName,
		Logger:      kafka.LoggerFunc(logf),
	    ErrorLogger: kafka.LoggerFunc(logf),
		Balancer: &kafka.LeastBytes{},
		AllowAutoTopicCreation: autoTopicCreation,
		BatchSize:    1,
        BatchTimeout: 10 * time.Millisecond,
		Async: true,
	}
}

func (w *Writer) WriteMessagesPack(ctx context.Context, messagesPack []string) error {
	defer recoveryFunction()

	if len(messagesPack) > 0{
		kMessages := make([]kafkago.Message, 0)
		for _, curMes := range messagesPack {
			kMessages = append(kMessages, kafka.Message{Key:[]byte(""),Value: []byte(curMes)})
		}
		
		retries := 3
		for i := 0; i < retries; i++ {
			err := w.kWriter.WriteMessages(ctx, kMessages...)
		
			if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
        		time.Sleep(time.Millisecond * 250)
        		continue
    		} else {
				if err == nil {
					break
				} 
				return err
			}
		}
	}
	return nil
}

func (w *Writer) Close() error {
	err := w.kWriter.Close()
	return err
}


