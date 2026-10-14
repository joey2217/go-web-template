# ---- 构建阶段 ----
FROM golang:1.26-alpine AS builder

WORKDIR /src

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o /out/server ./cmd/server

# ---- 运行阶段 ----
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /out/server /app/server
COPY --from=builder /src/config /app/config
COPY --from=builder /src/migrations /app/migrations

EXPOSE 8080

ENTRYPOINT ["/app/server"]
