package main

import (
	"github.com/Shopify/sarama"
)

func main() {
	producer, _ := sarama.NewSyncProducer([]string{"localhost:9092"}, nil)
	// message := &sarama.ProducerMessage{Topic: "demo", Partition: 0}
	message := &sarama.ProducerMessage{Topic: "demo", Partition: 0}
	producer.SendMessage(message)
}
