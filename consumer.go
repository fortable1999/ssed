package main

import (
	// "fmt"
	// "log"
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
	for {
		select {
		case msg := <-c.consumer.Messages():
				msgCh <-msg
		case err := <-c.consumer.Errors():
				errCh <-err
		case <-sigCh:
			doneCh <- 0
		}
	}
}
