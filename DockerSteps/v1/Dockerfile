FROM golang:1.10

ENV SRC_DIR=/go/src/github.com/mapleFU/DockerSteps/v1

ADD ./ $SRC_DIR

RUN go build $SRC_DIR/main.go

CMD ["./main"]