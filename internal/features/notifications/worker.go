package notifications

import (
	"context"
	"fmt"
	"log/slog"
	"nds-go-starter/internal/core/worker"

	"github.com/segmentio/kafka-go"
)

type notificationWorker struct {
	reader *kafka.Reader
}

func NewNotificationWorker(reader *kafka.Reader) worker.Worker {
	return &notificationWorker{
		reader: reader,
	}
}

func (w *notificationWorker) Start(ctx context.Context) {
	slog.Info("Notification Worker (Kafka Consumer) started")

	for {
		m, err := w.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			slog.Error("Notification Worker: failed to read message", "error", err)
			continue
		}

		fmt.Printf("Notification Worker RECEIVED Kafka Message: %s\n", string(m.Value))
		// Di sini nantinya bisa ditambahkan logika kirim email/WA/Push Notification
	}
}
