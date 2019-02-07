package main

import (
	"context"
	"io/ioutil"
	"log"
	"sync"

	"github.com/google/logger"
)

var (
	lgr = logger.Init("Server", true, true, ioutil.Discard)
)

func main() {
	listener, err := NewListener(lgr)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer listener.Close()

	publisher, err := NewPublisher(lgr)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer publisher.Close()

	listener.Setup(DefaultListenerEndpoint)
	publisher.Setup(DefaultPublisherEndpoint)
	ctx, cancel := context.WithCancel(context.Background())

	listenerWg := &sync.WaitGroup{}
	publisherWg := &sync.WaitGroup{}

	listenerWg.Add(1)
	go func() {
		defer listenerWg.Done()
		listener.Listen(ctx)
	}()

	publisherWg.Add(1)
	go func() {
		defer publisherWg.Done()
		done := false
		for !done {
			select {
			case message := <-listener.Messages():
				publisher.Publish(message)
			case <-listener.Errors():
				cancel()
				done = true
			}
		}
	}()

	listenerWg.Wait()
	publisherWg.Wait()
}
