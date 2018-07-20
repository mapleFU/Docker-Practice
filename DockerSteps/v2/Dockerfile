FROM golang:1.10-alpine

ENV SRC_DIR=/go/src/github.com/mapleFU/DockerSteps/v2

WORKDIR $SRC_DIR

COPY ./ $SRC_DIR

RUN go build $SRC_DIR/main.go

CMD ["./main"]