# 编译阶段
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.sum .
COPY go.mod .

# 设置国内代理
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download

COPY . .
# 编译
RUN go build -o server cmd/main.go

# 创建运行时镜像
FROM alpine:latest

WORKDIR /app

# 拷贝可执行文件
COPY --from=builder /app/server .

EXPOSE 8081

# 启动时需要将配置文件映射到指定目录
CMD ["./server", "--config", "./config/config.yaml"]