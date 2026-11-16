# Build
FROM golang:1.26-alpine AS builder

WORKDIR /src

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/velo-bot ./cmd/velo-bot

# Runtime
FROM alpine:3.22

RUN apk add --no-cache ca-certificates \
    && adduser -D -H -u 10001 bot \
    && mkdir -p /app/logger/logs \
    && chown -R bot:bot /app

WORKDIR /app

COPY --from=builder /out/velo-bot /app/velo-bot

USER bot

CMD ["./velo-bot"]
