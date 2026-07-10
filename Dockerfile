FROM docker.io/library/golang:1.26.4-alpine3.24 AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download && go mod verify

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go tool sqlc generate

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux \
    go build -trimpath \
    -ldflags="-w -s" \
    -o /build/clipharborbot \
    ./cmd/bot

FROM docker.io/library/alpine:3.24
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
