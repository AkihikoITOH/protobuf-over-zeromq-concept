# Protocol Buffers & ZeroMQ

This project is a proof of concept for an application using [ZeroMQ](http://zeromq.org/) as the communication protocol and [Protocol Buffers](https://developers.google.com/protocol-buffers/) as the message serialization mechanism.

## Components

- `server` - contains files for a server.
- `client` - contains files for clients which talks to a server.
- `protos` - contains `.proto` files which defines Protocol Buffer messages.

## Setup

To generate Go files which contain structs generated from the `.proto` files, run `./build/generate_pb.sh` in `server` or `client` directory.

- Build server program: `cd server; go build`
- Build client program: `cd client; go build`

## Run

- Start server: `cd server; ./server`
- Start client: `cd client; ./client`
