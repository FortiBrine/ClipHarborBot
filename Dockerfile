FROM golang:1.25-alpine3.23 AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath \
    -ldflags="-w -s" \
    -o /build/clipharborbot \
    ./cmd/bot

FROM alpine:3.23
ARG YTDLP_VERSION=2026.03.13

WORKDIR /app

RUN addgroup -S appgroup \
    && adduser -S appuser -G appgroup -h /app \
    && install -d -o appuser -g appgroup /app/temp

# install yt-dlp
RUN wget -O /usr/local/bin/yt-dlp \
        https://github.com/yt-dlp/yt-dlp/releases/download/${YTDLP_VERSION}/yt-dlp_musllinux \
        && chmod +x /usr/local/bin/yt-dlp

COPY --from=builder --chown=appuser:appgroup /build/clipharborbot .

USER appuser

EXPOSE 2000

ENV TMPDIR=/app/temp \
    HOME=/app

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -q -O- http://localhost:2000/health || exit 1

CMD ["/app/clipharborbot"]
