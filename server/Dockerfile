FROM golang:1.10-alpine

WORKDIR /go/src/server
COPY . .

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure
RUN go build

CMD ["server"]
