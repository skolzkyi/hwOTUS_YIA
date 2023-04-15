package kafka

import (
	"context"
	"errors"

	"github.com/segmentio/kafka-go"
	kafkago "github.com/segmentio/kafka-go"
)

var ErrReadMessage = errors.New("failed to read messages")

type KReader struct {
	kReader *kafkago.Reader
}

func NewReader() *KReader {
	reader := &kafkago.Reader{}
	return &KReader{
		kReader: reader,
	}

}

func (r *KReader) Init(addr string, port string, topicname string, group_ID string) {

	r.kReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{addr + ":" + port},
		Topic:    topicname,
		GroupID:  group_ID,
		MinBytes: 1,
		MaxBytes: 10e6, // 10MB
	})

}

func (r *KReader) ReadMessage(ctx context.Context) (string, error) {

	m, err := r.kReader.ReadMessage(ctx)
	if err != nil {
		return "", err
	}

	return string(m.Value), nil
}

func (r *KReader) Close() error {
	err := r.kReader.Close()
	return err
}
