package main

import (
	"fmt"
	"github.com/Shopify/sarama"
)

func main() {
	consumer, err := NewConsumer([]string{"localhost:9092"})
	consumer.topics = []string{"demo"}
	if err != nil {
		panic(err)
	}

	sigCh := make(chan int)
	doneCh := make(chan int)

	msgCh := make(chan *sarama.ConsumerMessage)
	errCh := make(chan *sarama.ConsumerError)

	go consumer.Start(msgCh, errCh, sigCh, doneCh)

	for {
		select {
		case msg := <-msgCh:
			fmt.Println(msg)
		case err := <-errCh:
			fmt.Println(err)
		case <-doneCh:
			fmt.Println("done!")
		}
	}
	fmt.Println("end")
}
