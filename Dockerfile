FROM golang:latest AS builder

# 为我们的镜像设置必要的环境变量
ENV CGO_ENABLED=0 \
    GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

# 项目根目录, 项目完整路径
WORKDIR /gopath/src/faucet

# 将代码复制到工作目录中
COPY . .
COPY ./views /faucet/views

# 将我们的代码编译成二进制可执行文件  可执行文件名为 faucet
RUN go mod tidy && go build -ldflags="-s -w" -o /faucet/faucet .

FROM ubuntu:20.04

ENV TZ Asia/Shanghai

WORKDIR /faucet

# 把工作目录中的文件复制到容器
COPY --from=builder /faucet/faucet /faucet/faucet
COPY --from=builder /faucet/views /faucet/views

# 声明服务端口
EXPOSE 3003

#
ENTRYPOINT ["/faucet/faucet"]
