FROM golang:1.22

WORKDIR /go/src/app

COPY . .

RUN go build -v ./...

CMD ["./bot_builder_engine"]