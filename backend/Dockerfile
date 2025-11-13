# ==========================
# CyberKube Backend Dockerfile
# ==========================

FROM golang:1.22-alpine AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
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

