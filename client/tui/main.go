package main

import (
	"context"
	"io/ioutil"
	"log"

	"github.com/AkihikoITOH/protobuf-over-zeromq-concept/client/tui/pb"
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

func setupListener() (*Listener, error) {
	listener, err := NewListener(lgr)
	if err != nil {
		return nil, err
	}
	listener.Setup(DefaultListenerEndpoint)

	return listener, nil
}

func setupPublisher() (*Publisher, error) {
	publisher, err := NewPublisher(lgr)
	if err != nil {
		return nil, err
	}

	publisher.Setup(DefaultPublisherEndpoint)

	return publisher, nil
}

func main() {
	listener, err := setupListener()
	if err != nil {
		lgr.Fatal(err.Error())
	}
	defer listener.Close()

	publisher, err := setupPublisher()
	if err != nil {
		lgr.Fatal(err.Error())
	}
	defer publisher.Close()

	listenCtx, cancelListen := context.WithCancel(context.Background())
	defer cancelListen()

	publishCtx, cancelPublish := context.WithCancel(context.Background())
	defer cancelPublish()

	go listener.Listen(listenCtx)
	go publisher.Poll(publishCtx)

	view := NewView(newUser(), listener.Messages(), publisher.Messages())
	if err != nil {
		log.Fatal(err.Error())
	}
	view.BuildUI()
	view.Start()
}
