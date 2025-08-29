package mesagging

import (
	"context"
	"encoding/json"
	"fmt"
	"food/model"

	"github.com/rs/zerolog/log"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Messaging struct {
	ChMessaging chan []byte
	Topic       string
	Brokers     []string
}

func NewMessaging(topic string, brokers []string) *Messaging {
	return &Messaging{make(chan []byte), topic, brokers}
}

func (msg *Messaging) ProduceRecords() {
	if msg.Topic == "" {
		panic("invalid topc")
	}

	if len(msg.Brokers) < 1 {
		panic("invalid brokers")
	}

	cl, err := kgo.NewClient(
		kgo.SeedBrokers(msg.Brokers...),

		kgo.RequiredAcks(kgo.AllISRAcks()),
	)
	if err != nil {
		panic(err)
	}

	defer cl.Close()
	ctx := context.Background()
	for message := range msg.ChMessaging {

		record := &kgo.Record{Topic: msg.Topic, Value: message, Key: nil}

		cl.Produce(ctx, record, func(r *kgo.Record, err error) {
			if err != nil {
				fmt.Printf("record had a produce error: %v\n", err)
			}
			order := new(model.ORDER)
			json.Unmarshal(r.Value, order)
			fmt.Println("Producer-->", r.ProducerID, "Topid-->", r.Topic, "Partition:", r.Partition, "Offset:", r.Offset, "Value:", order)
		})

	}
	cl.Flush(ctx)
	log.Print("Closed publishing data")

}
