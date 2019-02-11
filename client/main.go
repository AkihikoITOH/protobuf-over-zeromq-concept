package main

import (
	"context"
	"io/ioutil"
	"log"
	"sync"

	"github.com/AkihikoITOH/protobuf-over-zeromq-concept/client/pb"
	"github.com/Pallinder/sillyname-go"
	"github.com/google/logger"
	uuid "github.com/satori/go.uuid"
)

var (
	lgr = logger.Init("Server", false, false, ioutil.Discard)
)

func newUser() *pb.User {
	return &pb.User{Uuid: uuid.Must(uuid.NewV4()).String(), Name: sillyname.GenerateStupidName()}
}

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

	listenerWg.Add(1)
	go func() {
		defer listenerWg.Done()
		listener.Listen(ctx)
	}()

	outgoingMessages := make(chan *pb.Message)

	go func() {
		for {
			msg, ok := <-outgoingMessages
			if !ok {
				cancel()
				break
			}
			publisher.Publish(msg)
		}
	}()

	view := NewView(newUser(), listener.Messages(), outgoingMessages)
	if err != nil {
		log.Fatal(err.Error())
	}
	view.BuildUI()
	view.Start()
}
