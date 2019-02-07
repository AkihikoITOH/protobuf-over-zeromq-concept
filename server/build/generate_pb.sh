#!/bin/sh

protoc -I ../protos --go_out=./pb ../protos/*.proto
