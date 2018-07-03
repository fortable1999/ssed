// Copyright (c) 2017 Ismael Celis

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"fmt"
	// "strings"
	"log"
	"net/http"
	"github.com/Shopify/sarama"
)

// Example SSE server in Golang.
//     $ go run sse.go

type Broker struct {

	// Events are pushed to this channel by the main events-gathering routine
	Notifier chan string

	// New client connections
	newClients chan chan string

	// Closed client connections
	closingClients chan chan string

	// Client connections registry
	clients map[chan string]bool
}

func NewServer() (broker *Broker) {
	// Instantiate a broker
	broker = &Broker{
		Notifier:       make(chan string, 1),
		newClients:     make(chan chan string),
		closingClients: make(chan chan string),
		clients:        make(map[chan string]bool),
	}

	// Set it running - listening and broadcasting events
	go broker.listen()

	return
}

func (broker *Broker) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	// Make sure that the writer supports flushing.
	//
	flusher, ok := rw.(http.Flusher)

	if !ok {
		http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	q := req.URL.Query()
	topics, ok := q["topics"]
	if !ok {
		rw.Header().Set("Server", "A Go Web Server")
		rw.WriteHeader(400)
		rw.Write([]byte("Bad request"))
		return
	}
	for i, topic := range topics {
		fmt.Println(i, topic)
	}
	// topics := strings.Split(value, ",")

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	// Each connection registers its own message channel with the Broker's connections registry
	messageChan := make(chan string)

	// Signal the broker that we have a new connection
	broker.newClients <- messageChan

	// Remove this client from the map of connected clients
	// when this handler exits.
	defer func() {
		fmt.Println("finished connection")
		broker.closingClients <- messageChan
	}()

	// Listen to connection close and un-register messageChan
	notify := rw.(http.CloseNotifier).CloseNotify()


	for {

		// Write to the ResponseWriter
		// Server Sent Events compatible
		select {
		case msg := <-messageChan:
			fmt.Fprintf(rw, "data: %s\n\n", msg)
		case <-notify:
			broker.closingClients <- messageChan
		}

		// Flush the data immediatly instead of buffering it for later.
		flusher.Flush()
	}

}

func (broker *Broker) listen() {
	for {
		select {
		case s := <-broker.newClients:

			// A new client has connected.
			// Register their message channel
			broker.clients[s] = true
			log.Printf("Client added. %d registered clients", len(broker.clients))
		case s := <-broker.closingClients:

			// A client has dettached and we want to
			// stop sending them messages.
			delete(broker.clients, s)
			log.Printf("Removed client. %d registered clients", len(broker.clients))
		case event := <-broker.Notifier:

			// We got a new event from the outside!
			// Send event to all connected clients
			for clientMessageChan, _ := range broker.clients {
				clientMessageChan <- event
			}
		}
	}

}

func main() {

	broker := NewServer()

	consumer, err := NewConsumer([]string{"localhost:9092"})
	consumer.topics = []string{"demo"}
	if err != nil {
		panic(err)
	}

	sigCh := make(chan int)
	doneCh := make(chan int)

	msgCh := make(chan *sarama.ConsumerMessage)
	errCh := make(chan error)

	go consumer.Start(msgCh, errCh, sigCh, doneCh)

	go func() {
		for {
			select {
			case msg := <-msgCh:
				eventString := string(msg.Value)
				broker.Notifier <- eventString
			case err := <-errCh:
				fmt.Println(err)
			case <-doneCh:
				break
			}

		}
	}()

	log.Fatal("HTTP server error: ", http.ListenAndServe("localhost:3000", broker))

}
