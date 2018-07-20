FROM golang:1.10-alpine

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk add --no-cache git

ENV SRC_DIR=/go/src/github.com/mapleFU/GoSQLServerDemo

WORKDIR $SRC_DIR

COPY ../ $SRC_DIR


# Document that the service listens on port 8080.
EXPOSE 8080

RUN go build $SRC_DIR/persistance.go

CMD ["./persistance"]