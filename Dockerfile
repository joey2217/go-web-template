FROM golang:1.26-alpine AS builder

# 切换 Alpine APK 源到腾讯云镜像（必须早于 apk add）
RUN sed -i 's|dl-cdn.alpinelinux.org|mirrors.cloud.tencent.com|g' /etc/apk/repositories

WORKDIR /src

RUN apk add --no-cache ca-certificates

# Go module：代理走 goproxy.cn（七牛/360 维护），sum 库走谷歌中国镜像
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o /out/server ./cmd/server

# ---- 运行阶段 ----
FROM alpine:3.24

# 运行时同样切换 APK 源
RUN sed -i 's|dl-cdn.alpinelinux.org|mirrors.cloud.tencent.com|g' /etc/apk/repositories

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /out/server /app/server
COPY --from=builder /src/config /app/config
COPY --from=builder /src/migrations /app/migrations

EXPOSE 8080

ENTRYPOINT ["/app/server"]
