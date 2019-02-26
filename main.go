package main

import (
	"context"
	"log"
	"os"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

func main() {
	ctx := context.Background()
	logger := log.New(os.Stderr, "kafka-go: ", log.LstdFlags|log.Lshortfile)
	topic := "test-3"

	p, err := kafka.LookupPartition(ctx, "tcp", "localhost:9092", topic, 0)
	if err != nil {
		panic(err)
	}

	conn, err := kafka.DialPartition(ctx, "tcp", "localhost:9092", p)
	if err != nil {
		panic(err)
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   topic,
		//Async:         true,
		BatchSize:     1,
		RequiredAcks:  -1,
		QueueCapacity: 1,
		BatchTimeout:  10 * time.Millisecond,
		ErrorLogger:   logger,
		//Logger:        logger,
	})
	defer writer.Close()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{"localhost:9092"},
		Topic:       topic,
		GroupID:     "group-test-3",
		ErrorLogger: logger,
		//Logger:          logger,
		AutoOffsetReset: kafka.FirstOffset,
		QueueCapacity:   1,
		MaxWait:         time.Second,
	})
	defer reader.Close()

	msgs := make([]kafka.Message, 100)

	var last int64
	var written int
	var read int

	go func() {
		for {
			err := writer.WriteMessages(ctx, msgs...)
			if err != nil {
				panic(err)
			}

			written += len(msgs)
			if written%100000 == 0 {
				log.Printf("%d messages written", written)
			}
		}
	}()

	go func() {
		for {
			firstOffset, lastOffset, err := conn.ReadOffsets()
			if err != nil {
				panic(err)
			}
			log.Printf("first offset: %d, last offset: %d", firstOffset, lastOffset)
			time.Sleep(5 * time.Second)
		}
	}()

	for {
		firstOffset, lastOffset, err := conn.ReadOffsets()
		if err != nil {
			panic(err)
		}

		msg, _ := reader.ReadMessage(ctx)

		read++
		if read%1000 == 0 {
			log.Printf("%d messages read", read)
		}

		// Log when the offset jumps.
		if msg.Offset-last > 1 {
			log.Printf("first offset: %d, last offset: %d", firstOffset, lastOffset)
			log.Println("OFFSET", msg.Offset)
		}
		last = msg.Offset
		time.Sleep(time.Millisecond)
	}
}
