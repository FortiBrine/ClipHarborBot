FROM golang:1.25-alpine AS builder

RUN apk update && \
    apk add --no-cache \
    ca-certificates \
    git \
    && rm -rf /var/cache/apk/*

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -installsuffix cgo \
    -ldflags="-w -s -extldflags '-static'" \
    -o /build/clipharbor-bot \
    ./cmd/bot

FROM alpine:3.20

RUN apk add --no-cache \
    ca-certificates \
    yt-dlp \
    ffmpeg \
    tzdata

RUN addgroup -g 1000 -S appgroup && \
    adduser -u 1000 -S appuser -G appgroup -h /app

WORKDIR /app

RUN mkdir -p /app/temp && \
    chown -R appuser:appgroup /app

COPY --from=builder --chown=appuser:appgroup /build/clipharbor-bot /app/clipharbor-bot

USER appuser

EXPOSE 2000

ENV TMPDIR=/app/temp \
    HOME=/app

CMD ["/app/clipharbor-bot"]

