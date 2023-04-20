package kafka
// pochemuto linter typecheck na servere schitaet, chto  segmentio/kafka-go ne ispolzuetsia, a kafka. - ne opredeleno, no vse compiliruetsia i rabotaet
import (
	"context"
	"errors"
	

	"github.com/segmentio/kafka-go" //nolint:typecheck
)

var ErrReadMessage = errors.New("failed to read messages")

type KReader struct {
	KReader *kafka.Reader //nolint:typecheck
}


func NewReader() KReader {
	return KReader{}
}

func (r *KReader) Init(addr string, port string, topicname string, groupID string) {

	r.KReader = kafka.NewReader(kafka.ReaderConfig{ //nolint:typecheck
		Brokers:  []string{addr + ":" + port},
		Topic:    topicname,
		GroupID:  groupID,
		MinBytes: 1,
		MaxBytes: 10e6, // 10MB
		Logger:      kafka.LoggerFunc(logf), //nolint:typecheck
	    ErrorLogger: kafka.LoggerFunc(logf), //nolint:typecheck
		StartOffset: kafka.FirstOffset, //nolint:typecheck
	})
}

func (r *KReader) ReadMessage(ctx context.Context) (string, error) {
	defer recoveryFunction()
	m, err := r.KReader.ReadMessage(ctx)
	if err != nil {
		return "", err
	}

	return string(m.Value), nil
}

func (r *KReader) Close() error {
	err := r.KReader.Close()
	return err
}
