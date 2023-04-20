package kafka
// pochemuto linter typecheck na servere schitaet, chto  segmentio/kafka-go ne ispolzuetsia, a kafka. - ne opredeleno, no vse compiliruetsia i rabotaet
import (
	"context"
	"errors"
	"time"

	"github.com/segmentio/kafka-go" //nolint:typecheck
)

var ErrWriteMessage = errors.New("failed to write messages")

type Writer struct {
	KWriter *kafka.Writer //nolint:typecheck
}



func NewWriter() Writer {
	return Writer{}
}

func (w *Writer) Init(addr string, port string, topicName string, autoTopicCreation bool) {
	w.KWriter = &kafka.Writer{ //nolint:typecheck
		Addr:     kafka.TCP(addr + ":" + port), //nolint:typecheck
		Topic:    topicName,
		Logger:      kafka.LoggerFunc(logf), //nolint:typecheck
	    ErrorLogger: kafka.LoggerFunc(logf), //nolint:typecheck
		Balancer: &kafka.LeastBytes{}, //nolint:typecheck
		AllowAutoTopicCreation: autoTopicCreation,
		BatchSize:    1,
        BatchTimeout: 10 * time.Millisecond,
		Async: true,
	}
}

func (w *Writer) WriteMessagesPack(ctx context.Context, messagesPack []string) error {
	defer recoveryFunction()

	if len(messagesPack) > 0{
		KMessages := make([]kafka.Message, 0) //nolint:typecheck
		for _, curMes := range messagesPack {
			KMessages = append(KMessages, kafka.Message{Key:[]byte(""),Value: []byte(curMes)}) //nolint:typecheck
		}
		
		retries := 3
		for i := 0; i < retries; i++ {
			err := w.KWriter.WriteMessages(ctx, KMessages...)
		
			if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) { //nolint:typecheck
        		time.Sleep(time.Millisecond * 250)
        		continue
    		} 
			if err == nil {
				break
			} 
			return err
		}
	}
	return nil
}

func (w *Writer) Close() error {
	err := w.KWriter.Close()
	return err
}


