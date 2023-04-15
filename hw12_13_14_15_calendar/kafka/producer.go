package kafka

import (
	"context"
	"errors"

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

func (w *Writer) Init(addr string, port string, topicname string) {

	w.kWriter = &kafka.Writer{
		Addr:     kafka.TCP(addr + ":" + port),
		Topic:    topicname,
		Balancer: &kafka.LeastBytes{},
	}
}

func (w *Writer) WriteMessagesPack(ctx context.Context, messagesPack []string) error {
	kMessages := make([]kafkago.Message, 0)
	for _, curMes := range messagesPack {
		kMessages = append(kMessages, kafka.Message{Value: []byte(curMes)})
	}
	err := w.kWriter.WriteMessages(ctx, kMessages...)

	return err
}

func (w *Writer) Close() error {
	err := w.kWriter.Close()
	return err
}
