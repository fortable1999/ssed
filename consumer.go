package main

import (
	"fmt"
	"log"
	"regexp"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

type Consumer struct {
	consumer *cluster.Consumer
	zk []string
	topics []string
}

func NewConsumer(zkList []string) (*Consumer, error) {
	config := cluster.NewConfig()
	config.Group.Topics.Whitelist = regexp.MustCompile(`access_log-\w+`)
	// c, err := sarama.NewConsumer(zkList, config)
	c, err := cluster.NewConsumer(zkList, "demo_group", nil, config)
	if err != nil {
		panic(err)
	}
	consumer := Consumer{ consumer:c, zk:zkList, topics:[]string{""} }
	return &consumer, err
}

func (c *Consumer) Start(msgCh chan *sarama.ConsumerMessage,
						 errCh chan error,
						 sigCh chan int,
						 doneCh chan int) {
	defer c.consumer.Close()
	// consumer, err := (*c.consumer).ConsumePartition(c.topics[0], 0, sarama.OffsetOldest)
	// if err != nil {
	// 	log.Fatal(err)
	// 	doneCh <-1
	// 	return
	// }
	for {
		fmt.Println("next event...")
		select {
		// case err := <-c.consumer.Errors():
		// 	log.Println(err)
		// 	errCh <-err
		case msg := <-c.consumer.Messages():
				log.Println("Received messages", string(msg.Key), string(msg.Value))
				msgCh <-msg
		case err := <-c.consumer.Errors():
				log.Println("Received error", string(err.Error()))
				errCh <-err
		case <-sigCh:
			log.Println("Receiving stop signal")
			doneCh <- 0
		}
	}
}

func (c *Consumer) Close() error {
	err := c.Close()
	log.Println("consumer stopped with error", err)
	return err
}

