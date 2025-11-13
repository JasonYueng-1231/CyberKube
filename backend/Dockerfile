# ==========================
# CyberKube Backend Dockerfile
# ==========================

ARG GO_IMAGE=golang:1.22-alpine
FROM ${GO_IMAGE} AS builder
WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct

# 先复制模块文件，进行依赖预下载（容错）
COPY go.mod ./
RUN go mod download || true

# 再复制全部源码并整理依赖，生成 go.sum
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/cyberkube ./cmd/server

FROM alpine:3.20
WORKDIR /app
RUN adduser -D -H -u 10001 app && apk add --no-cache ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime || true
USER app
COPY --from=builder /out/cyberkube /app/cyberkube

ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["/app/cyberkube"]
