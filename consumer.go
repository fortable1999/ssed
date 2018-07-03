package main

import (
	// "fmt"
	"encoding/hex"
	"log"
	// "regexp"
	// "github.com/Shopify/sarama"
	kafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/bsm/go-guid"
	// sarama "github.com/bsm/sarama"
)

type Consumer struct {
	consumer *kafka.Consumer
	zk []string
	topics []string
}

func NewLogConsumer(zkList []string, topics []string) (*Consumer, error) {
	// config := sarama.NewConfig()

	if topics == nil {
		topics = []string{"^access_log-.*"}
	}
	// c, err := sarama.NewConsumer(zkList, topics,)
	guid := hex.EncodeToString(guid.New128().Bytes())
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id": guid,
		// "auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	c.SubscribeTopics(topics, nil)
	consumer := Consumer{ consumer:c, zk:zkList, topics:topics }
	return &consumer, err
}

func (c *Consumer) Start(msgCh chan *kafka.Message,
						 errCh chan error,
						 sigCh chan struct{},
					     doneCh chan struct{}) {
	defer func () {
		c.consumer.Close()
	}()

	run := true
	for run {
		select {
		case <-sigCh:
			log.Println("Caught signal %v: terminating\n")
			run = false
		default:
			ev := c.consumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				log.Printf("%% Message on %s:\n%s\n",
					e.TopicPartition, string(e.Value))
				if e.Headers != nil {
					log.Printf("%% Headers: %v\n", e.Headers)
				}
				msgCh <-e
			case kafka.PartitionEOF:
				log.Printf("%% Reached %v\n", e)
			case kafka.Error:
				// log.Fprintf(os.Stderr, "%% Error: %v\n", e)
				run = false
				errCh <-e
			default:
				log.Printf("Ignored %v\n", e)
			}
		}
	}

	// for {
	// 	log.Println("start looping")
	// 	select {
	// 	case msg := <-c.consumer.Messages():
	// 		log.Printf("message event %s", string(msg.Value))
	// 		msgCh <-msg
	// 	case err := <-c.consumer.Errors():
	// 		log.Printf("error event %s", err.Error())
	// 		errCh <-err
	// 	case <-sigCh:
	// 		doneCh <-struct{}{}
	// 	}
	// }
}
