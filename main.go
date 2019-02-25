package main

import (
	"context"
	"log"
	"os"
	"time"

	kafka "github.com/newrelic-forks/kafka-go"
)

func main() {
	logger := log.New(os.Stderr, "kafka-go: ", log.LstdFlags)

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:         []string{"localhost:9092"},
		Topic:           "test-2",
		Async:           true,
		BatchSize:       100,
		MaxMessageBytes: 1048576,
		RequiredAcks:    -1,
		QueueCapacity:   16384,
		BatchTimeout:    500 * time.Millisecond,
		ErrorLogger:     logger,
	})
	defer writer.Close()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{"localhost:9092"},
		Topic:       "test-2",
		GroupID:     "group-test",
		ErrorLogger: logger,
	})
	defer reader.Close()

	ctx := context.Background()
	k, v := []byte{}, []byte{}

	// Write messages continuously.
	go func() {
		i := 0
		for {
			i++
			err := writer.WriteMessages(ctx, kafka.Message{
				Key:   k,
				Value: v,
			})
			if err != nil {
				panic(err)
			}

			if i%10000 == 0 {
				log.Printf("wrote %d messages", i)
			}
			time.Sleep(time.Millisecond)
		}
	}()

	var last int64
	j := 0
	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			panic(err)
		}
		j++
		if j%100 == 0 {
			log.Printf("read %d messages", j)
		}

		// Log when the offset jumps.
		if msg.Offset-last > 1 {
			log.Println("OFFSET", msg.Offset)
		}
		last = msg.Offset

		time.Sleep(time.Millisecond * 50)
	}
}
