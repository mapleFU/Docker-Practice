FROM golang:1.10-alpine

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk add --no-cache git

ENV SRC_DIR=/go/src/github.com/mapleFU/Images/dockerfile-learn

WORKDIR $SRC_DIR

COPY ./ $SRC_DIR

EXPOSE 8080 8080

RUN go build $SRC_DIR/main.go

CMD ["./main"]