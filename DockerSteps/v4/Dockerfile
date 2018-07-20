FROM golang:1.10-alpine

ENV SRC_DIR=/go/src/github.com/mapleFU/Docker-Practice/DockerSteps/v3

WORKDIR $SRC_DIR

EXPOSE 8080

COPY ./ $SRC_DIR

RUN go build $SRC_DIR/main.go

CMD ["./main"]