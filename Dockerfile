FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -installsuffix cgo \
    -ldflags="-w -s -extldflags '-static'" \
    -o /build/clipharborbot \
    ./cmd/bot

FROM alpine:3.23.3

RUN apk add --no-cache yt-dlp && \
    addgroup -g 1000 -S appgroup && \
    adduser -u 1000 -S appuser -G appgroup -h /app

WORKDIR /app

RUN mkdir -p /app/temp && \
    chown -R appuser:appgroup /app

COPY --from=builder --chown=appuser:appgroup /build/clipharborbot /app/clipharborbot

USER appuser

EXPOSE 2000

ENV TMPDIR=/app/temp \
    HOME=/app

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD /bin/busybox wget -qO- http://localhost:2000/health || exit 1

CMD ["/app/clipharborbot"]
