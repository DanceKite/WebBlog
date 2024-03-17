FROM golang:1.20-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download

COPY . .

RUN go build -o webblog_app .

# 创建一个小小镜像
FROM debian:buster-slim

COPY ./wait-for.sh /
COPY ./templates /templates
COPY ./static /static
COPY ./conf /conf

COPY --from=builder /build/webblog_app /

RUN set -eux; \
    apt-get update; \
    apt-get install -y \
        --no-install-recommends \
        netcat; \
        chmod 755 wait-for.sh


#ENTRYPOINT ["/webblog_app","conf/config.yaml"]

