package main

import (
	"context"

	"github.com/AkihikoITOH/protobuf-over-zeromq-concept/client/pb"
	"github.com/golang/protobuf/proto"
	"github.com/google/logger"
	zmq "github.com/pebbe/zmq4"
)

const DefaultListenerEndpoint = "tcp://127.0.0.1:8080"

type Listener struct {
	*zmq.Socket
	messageChannel chan *pb.Message
	errorChannel   chan error
	logger         *logger.Logger
}

func NewListener(logger *logger.Logger) (*Listener, error) {
	socket, err := zmq.NewSocket(zmq.PULL)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	listener := &Listener{
		Socket:         socket,
		messageChannel: make(chan *pb.Message),
		errorChannel:   make(chan error),
		logger:         logger,
	}
	return listener, nil
}

func (listener *Listener) Setup(endpoint string) error {
	err := listener.Bind(endpoint)
	if err != nil {
		listener.logger.Error(err)
		return err
	}
	listener.logger.Info("Connection established.")
	return nil
}

func (listener *Listener) Messages() <-chan *pb.Message {
	return listener.messageChannel
}

func (listener *Listener) Errors() <-chan error {
	return listener.errorChannel
}

func (listener *Listener) Listen(ctx context.Context) {
	listener.logger.Info("Listening to messages.")
	defer close(listener.messageChannel)
	defer close(listener.errorChannel)

	done := false
	for !done {
		select {
		case <-ctx.Done():
			done = true
			break
		default:
			msg, err := listener.RecvBytes(0)
			if err != nil {
				listener.logger.Error(err)
				listener.errorChannel <- err
				done = true
				break
			} else {
				pbmsg := &pb.Message{}
				listener.logger.Info("Received message")
				err = proto.Unmarshal(msg, pbmsg)
				if err != nil {
					listener.logger.Error(err)
					listener.errorChannel <- err
					done = true
					break
				}
				listener.messageChannel <- pbmsg
			}
		}
	}

	return
}

func (listener *Listener) Close() {
	listener.Socket.Close()
}
