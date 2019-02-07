package main

import (
	"github.com/AkihikoITOH/protobuf-over-zeromq-concept/client/pb"
	"github.com/golang/protobuf/proto"
	"github.com/google/logger"
	zmq "github.com/pebbe/zmq4"
)

const DefaultPublisherEndpoint = "tcp://127.0.0.1:8080"

type Publisher struct {
	*zmq.Socket
	logger *logger.Logger
}

func NewPublisher(logger *logger.Logger) (*Publisher, error) {
	socket, err := zmq.NewSocket(zmq.PUSH)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &Publisher{socket, logger}, nil
}

func (publisher *Publisher) Setup(endpoint string) error {
	err := publisher.Connect(endpoint)
	if err != nil {
		publisher.logger.Error(err)
		return err
	}
	publisher.logger.Info("Connection established.")
	return nil
}

func (publisher *Publisher) Publish(message *pb.Message) error {
	publisher.logger.Infof("Publishing message: %s", message)
	data, err := proto.Marshal(message)
	if err == nil {
		publisher.SendBytes(data, 0)
	}
	return err
}

func (publisher *Publisher) Close() {
	publisher.Socket.Close()
}
