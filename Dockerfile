FROM golang:latest AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /gopath/src/faucet

COPY . .
RUN go build -ldflags="-s -w" -o /main.go


FROM ubuntu:20.04

ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/main /app/main

CMD ["./main"]
